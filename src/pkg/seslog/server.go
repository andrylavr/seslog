package seslog

import (
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/patrickmn/go-cache"
	"github.com/satyrius/gonx"
	"github.com/ua-parser/uap-go/uaparser"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func getStringFromLogParts(logParts format.LogParts, key string) string {
	str := ""
	part, ok := logParts[key]
	if ok {
		switch v := part.(type) {
		case string:
			str = v
		}
	}

	return str
}

func getNginxHostname(logParts format.LogParts) string {
	return getStringFromLogParts(logParts, "hostname")
}

func getNginxTag(logParts format.LogParts) string {
	return getStringFromLogParts(logParts, "tag")
}

func getNginxEventTimestamp(logParts format.LogParts) time.Time {
	timestamp := time.Now()
	timestampPart, ok := logParts["timestamp"]
	if ok {
		switch v := timestampPart.(type) {
		case time.Time:
			timestamp = v
		}
	}

	return timestamp
}

func getNginxIP(logParts format.LogParts) string {
	hostport := getStringFromLogParts(logParts, "client")
	host, _, err := net.SplitHostPort(hostport)
	if err == nil {
		return host
	}
	return ""
}

func getEventContent(logParts format.LogParts) string {
	content := ""
	contentPart, ok := logParts["content"]
	if ok {
		switch v := contentPart.(type) {
		case string:
			content = v
		}
	}

	return content
}

func (this *AccessLogServer) handleLogParts() {
	for logParts := range this.channel {
		access_log_event := AccessLogEvent{}
		access_log_event.Nginx_ip = getNginxIP(logParts)
		access_log_event.Nginx_event_time = getNginxEventTimestamp(logParts)
		access_log_event.Nginx_hostname = getNginxHostname(logParts)
		access_log_event.Nginx_tag = getNginxTag(logParts)

		content := getEventContent(logParts)
		entry, err := this.nginx_parser.ParseString(content)
		if err != nil {
			glog.Warningf("ParseString fail: %s", err)
			continue
		}

		fields := entry.Fields()
		this.fields2event(fields, &access_log_event)
		this.chwriter.AddEvent(access_log_event)
	}
}

type AccessLogServer struct {
	options      Options
	syslogServer *syslog.Server
	handler      *syslog.ChannelHandler
	chwriter     *CHWriter
	channel      chan format.LogParts
	nginx_parser *gonx.Parser
	uaparser     *uaparser.Parser
	uacache      *cache.Cache
}

const log_format = `$body_bytes_sent	$connections_active	$connections_reading	$connections_waiting	$connections_writing	$content_length	$http_host	$http_referer	$http_user_agent	$http_x_forwarded_for	$remote_addr	$request_method	$request_time	$request_uri	$scheme	$status	$tcpinfo_rtt	$tcpinfo_rttvar	$time_local	$upstream_cache_status	$upstream_response_length	$upstream_response_time	$upstream_status	$uri	$sent_http_location`

func NewAccessLogServer(options Options) (*AccessLogServer, error) {
	channel := make(syslog.LogPartsChannel)
	uaparser_inst, err := uaparser.New("./resources/regexes.yaml")
	if err != nil {
		glog.Fatalf("uaparser fail: %s", err)
	}
	server := &AccessLogServer{
		options:      options,
		syslogServer: syslog.NewServer(),
		handler:      syslog.NewChannelHandler(channel),
		chwriter:     NewCHWriter(options),
		channel:      channel,
		nginx_parser: gonx.NewParser(log_format),
		uaparser:     uaparser_inst,
		uacache:      cache.New(1*time.Hour, 2*time.Hour),
	}

	server.syslogServer.SetFormat(syslog.RFC3164)
	server.syslogServer.SetHandler(server.handler)

	return server, nil
}

func (this *AccessLogServer) RunServer() error {
	if err := this.syslogServer.ListenUDP(this.options.Address); err != nil {
		return err
	}
	glog.Infof("Seslog server listen UDP [%s]", this.options.Address)
	if err := this.syslogServer.Boot(); err != nil {
		return err
	}

	go this.handleLogParts()
	go this.chwriter.startWatcher()

	this.syslogServer.Wait()

	return nil
}
