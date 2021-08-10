package seslog

import (
	"bytes"
	"compress/gzip"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/lib/data"

	"github.com/golang/glog"
)

const INSERT_SQL = "INSERT INTO seslog.access_log (" +
	"nginx_tag, " +
	"nginx_event_time, " +
	"nginx_ip_uint32, " +
	"nginx_hostname, " +

	"zonename, " +
	"zoneoffset, " +

	"body_bytes_sent, " +
	"connections_active, " +
	"connections_reading, " +
	"connections_waiting, " +
	"connections_writing, " +
	"content_length, " +

	"http_scheme, " +
	"http_domain," +
	"http_path," +
	"http_arg_keys," +
	"http_arg_vals," +

	"http_referer_scheme, " +
	"http_referer_domain," +
	"http_referer_path," +
	"http_referer_arg_keys," +
	"http_referer_arg_vals," +

	"http_location_scheme, " +
	"http_location_domain," +
	"http_location_path," +
	"http_location_arg_keys," +
	"http_location_arg_vals," +

	"ua_family," +
	"ua_major," +
	"ua_minor," +
	"ua_patch," +
	"ua_os_family," +
	"ua_os_major," +
	"ua_os_minor," +
	"ua_os_patch," +
	"ua_os_patchminor," +
	"ua_device_family," +

	"http_x_forwarded_for, " +
	"remote_addr_uint32, " +
	"request_method, " +
	"request_time, " +
	"status, " +
	"tcpinfo_rtt, " +
	"tcpinfo_rttvar, " +

	"upstream_response_length, " +
	"upstream_response_time, " +
	"upstream_status, " +

	"uri" +
	") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

type CHWriter struct {
	sync.Mutex
	Options
	connect           clickhouse.Clickhouse
	events            AccessLogEvents
	block             *data.Block
	skip_write_backup bool
}

func checkFatal(err error, msgFormat string) {
	if err != nil {
		glog.Fatalf(msgFormat, err)
	}
}

func NewCHWriter(options Options) *CHWriter {
	connect, err := clickhouse.OpenDirect(options.CHDSN)
	checkFatal(err, "ClickHouse connect problem: %s")
	this := &CHWriter{
		Options: options,
		connect: connect,
		events:  AccessLogEvents{},
	}
	return this
}

func (this *CHWriter) AddEvent(event AccessLogEvent) {
	this.Lock()
	this.events = append(this.events, event)
	this.Unlock()
}

func (this *CHWriter) startWatcher() {
	dur, err := time.ParseDuration(this.Options.FlushInterval)
	if this.notOk(err) {
		return
	}
	for range time.Tick(dur) {
		this.writeEvents()
	}
}

func (this *CHWriter) reconnect() error {
	glog.Info("Try to reconnect...")
	connect, err := clickhouse.OpenDirect(this.Options.CHDSN)
	if err == nil {
		this.connect = connect
		glog.Info("Reconnected!")
		return nil
	}
	for range time.After(30 * time.Second) {
		glog.Info("Try to reconnect after 30 sec...")
		connect, err := clickhouse.OpenDirect(this.Options.CHDSN)
		if err == nil {
			this.connect = connect
			glog.Info("Reconnected!")
			return nil
		}
		break
	}
	glog.Infof("Reconnect failed: %s", err)
	return err
}

func (this *CHWriter) notOk(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == "driver: bad connection" {
		go this.reconnect()
	}

	glog.Warningf("ClickHouse error: %s", err)
	return true
}

func (this *CHWriter) writeEvents() {
	this.Lock()
	defer this.Unlock()
	if len(this.events) == 0 {
		return
	}
	events := make(AccessLogEvents, len(this.events))
	copy(events, this.events)
	this.events = AccessLogEvents{}
	go this.makeEventsTx(events)
}

func (this *CHWriter) toJsonGZ(events AccessLogEvents) ([]byte, error) {
	var gzBuffer bytes.Buffer
	json_byte_arr, err := json.Marshal(events)
	if this.notOk(err) {
		return nil, err
	}
	zipWriter, err := gzip.NewWriterLevel(&gzBuffer, gzip.BestCompression)
	if this.notOk(err) {
		return nil, err
	}
	_, err = zipWriter.Write(json_byte_arr)
	if this.notOk(err) {
		return nil, err
	}
	if err := zipWriter.Close(); this.notOk(err) {
		return nil, err
	}
	return gzBuffer.Bytes(), nil
}

func (this *CHWriter) writeEventsToFile(events AccessLogEvents) {
	byte_arr, err := this.toJsonGZ(events)

	if this.notOk(err) {
		return
	}
	filename := fmt.Sprintf("/tmp/seslog.events.%s.json.gz", time.Now().Format("20060102150405"))
	if err := ioutil.WriteFile(filename, byte_arr, 0644); this.notOk(err) {
		return
	}
	glog.Infof("Written %d bytes to %s", len(byte_arr), filename)
}

func ReadGzFile(filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (this *CHWriter) fromJsonGz(filename string) (AccessLogEvents, error) {
	byte_arr, err := ReadGzFile(filename)
	if this.notOk(err) {
		return nil, err
	}

	events := AccessLogEvents{}
	if err := json.Unmarshal(byte_arr, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (this *CHWriter) FromFilesToCH() {
	this.skip_write_backup = true
	pattern := "/tmp/seslog.events.*.json.gz"
	filenames, err := filepath.Glob(pattern)
	if this.notOk(err) {
		return
	}
	for _, filename := range filenames {
		glog.Infof("ReadGzFile: %s", filename)
		events, err := this.fromJsonGz(filename)
		if this.notOk(err) {
			continue
		}
		if err := this.makeEventsTx(events); this.notOk(err) {
			continue
		}
		err = os.Remove(filename)
		this.notOk(err)
	}
}

func (this *CHWriter) makeEventsTx(events AccessLogEvents) error {
	commited := false

	//fallback
	defer func() {
		if !commited && !this.skip_write_backup {
			go this.writeEventsToFile(events)
		}
	}()

	tx, err := this.connect.Begin()
	if this.notOk(err) {
		return err
	}

	//tx defer
	defer func() {
		if !commited {
			err := tx.Rollback()
			this.notOk(err)
		}
	}()

	stmt, err := this.connect.Prepare(INSERT_SQL)
	//stmt fail
	if this.notOk(err) {
		return err
	}

	//stmt defer
	defer func() {
		err := stmt.Close()
		this.notOk(err)
	}()

	block, err := this.connect.Block()
	if this.notOk(err) {
		return err
	}

	for _, event := range events {
		err = block.AppendRow([]driver.Value{
			event.Nginx_tag,
			event.Nginx_event_time,
			event.Nginx_ip_uint32,
			event.Nginx_hostname,

			event.Zonename,
			event.Zoneoffset,

			event.Body_bytes_sent,
			event.Connections_active,
			event.Connections_reading,
			event.Connections_waiting,
			event.Connections_writing,
			event.Content_length,

			event.Url_parsed.Scheme,
			event.Url_parsed.Domain,
			event.Url_parsed.Path,
			clickhouse.Array(event.Url_parsed.Arg_keys),
			clickhouse.Array(clickhouse.Array(event.Url_parsed.Arg_vals)),

			event.Referer_parsed.Scheme,
			event.Referer_parsed.Domain,
			event.Referer_parsed.Path,
			clickhouse.Array(event.Referer_parsed.Arg_keys),
			clickhouse.Array(clickhouse.Array(event.Referer_parsed.Arg_vals)),

			event.Location_parsed.Scheme,
			event.Location_parsed.Domain,
			event.Location_parsed.Path,
			clickhouse.Array(event.Location_parsed.Arg_keys),
			clickhouse.Array(clickhouse.Array(event.Location_parsed.Arg_vals)),

			event.Ua_family,
			event.Ua_major,
			event.Ua_minor,
			event.Ua_patch,
			event.Ua_os_family,
			event.Ua_os_major,
			event.Ua_os_minor,
			event.Ua_os_patch,
			event.Ua_os_patchminor,
			event.Ua_device_family,

			event.Http_x_forwarded_for,
			event.Remote_addr_uint32,
			event.Request_method,
			event.Request_time,
			event.Status,

			event.Tcpinfo_rtt,
			event.Tcpinfo_rttvar,

			event.Upstream_response_length,
			event.Upstream_response_time,
			event.Upstream_status,

			event.Uri,
		})
		if this.notOk(err) {
			continue
		}
	}

	err = this.connect.WriteBlock(block)
	if this.notOk(err) {
		return err
	}

	err = this.connect.Commit()
	if this.notOk(err) {
		return err
	}
	commited = true
	glog.Infof("CHWriter: written %d events", len(events))

	return nil
}
