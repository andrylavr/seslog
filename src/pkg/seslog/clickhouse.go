package seslog

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kshvakov/clickhouse"
)

type CHWriter struct {
	sync.Mutex
	connect *sql.DB
	events  []AccessLogEvent
}

func NewCHWriter(options Options) *CHWriter {
	connect, err := sql.Open("clickhouse", options.CHDSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		glog.Fatalf("ClickHouse connect problem: %s", err)
	}

	this := &CHWriter{
		connect: connect,
		events:  []AccessLogEvent{},
	}

	return this
}

func (this *CHWriter) AddEvent(event AccessLogEvent) {
	this.Lock()
	this.events = append(this.events, event)
	this.Unlock()
}

func (this *CHWriter) startWatcher() {
	for range time.Tick(5 * time.Second) {
		this.WriteEvents()
	}
}

func notOk(err error) bool {
	if err != nil {
		glog.Warningf("ClickHouse error: %s", err)
		return true
	}
	return false
}

func (this *CHWriter) WriteEvents() {
	if len(this.events) == 0 {
		return
	}
	this.Lock()
	events := make([]AccessLogEvent, len(this.events))
	copy(events, this.events)
	this.events = []AccessLogEvent{}
	this.Unlock()

	go func(events []AccessLogEvent) {
		tx, err := this.connect.Begin()
		if notOk(err) {
			return
		}
		const insertsql = "INSERT INTO seslog.access_log (" +
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

			"geoip_country_code, " +
			"geoip_latitude, " +
			"geoip_longitude, " +

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
		stmt, err := tx.Prepare(insertsql)
		if notOk(err) {
			return
		}
		for _, event := range events {
			_, err := stmt.Exec(
				event.nginx_tag,
				event.nginx_event_time,
				event.nginx_ip_uint32,
				event.nginx_hostname,

				event.zonename,
				event.zoneoffset,

				event.Body_bytes_sent,
				event.Connections_active,
				event.Connections_reading,
				event.Connections_waiting,
				event.Connections_writing,
				event.Content_length,

				event.Geoip_country_code,
				event.Geoip_latitude,
				event.Geoip_longitude,

				event.url_parsed.scheme,
				event.url_parsed.domain,
				event.url_parsed.path,
				clickhouse.Array(event.url_parsed.arg_keys),
				clickhouse.Array(clickhouse.Array(event.url_parsed.arg_vals)),

				event.referer_parsed.scheme,
				event.referer_parsed.domain,
				event.referer_parsed.path,
				clickhouse.Array(event.referer_parsed.arg_keys),
				clickhouse.Array(clickhouse.Array(event.referer_parsed.arg_vals)),

				event.ua_family,
				event.ua_major,
				event.ua_minor,
				event.ua_patch,
				event.ua_os_family,
				event.ua_os_major,
				event.ua_os_minor,
				event.ua_os_patch,
				event.ua_os_patchminor,
				event.ua_device_family,

				event.Http_x_forwarded_for,
				event.remote_addr_uint32,
				event.Request_method,
				event.Request_time,
				event.Status,

				event.Tcpinfo_rtt,
				event.Tcpinfo_rttvar,

				event.Upstream_response_length,
				event.Upstream_response_time,
				event.Upstream_status,

				event.Uri,
			)
			if notOk(err) {
				continue
			}
		}
		err = tx.Commit()
		notOk(err)
		err = stmt.Close()
		notOk(err)
	}(events)

}
