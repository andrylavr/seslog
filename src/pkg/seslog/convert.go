package seslog

import (
	"errors"
	"github.com/golang/glog"
	"github.com/satyrius/gonx"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

const NginxTimeLayout = "02/Jan/2006:15:04:05 -0700"

type VarType uint

const (
	Invalid VarType = iota
	Uint64
	Uint32
	Uint16
	Uint8
	Bool
	String
	Float64
	Time
)

// map of nginx's variable types
var typemap = map[string]VarType{
	"args":                     String,
	"body_bytes_sent":          Uint64,
	"connections_active":       Uint16,
	"connections_reading":      Uint16,
	"connections_waiting":      Uint16,
	"connections_writing":      Uint16,
	"content_length":           Uint64,
	"geoip_country_code":       String,
	"geoip_latitude":           Float64,
	"geoip_longitude":          Float64,
	"host":                     String,
	"http_host":                String,
	"http_referer":             String,
	"http_user_agent":          String,
	"http_x_forwarded_for":     String,
	"remote_addr":              String,
	"request_method":           String,
	"request_time":             Float64,
	"request_uri":              String,
	"scheme":                   String,
	"status":                   Uint16,
	"tcpinfo_rtt":              Uint64,
	"tcpinfo_rttvar":           Uint64,
	"time_local":               Time,
	"upstream_cache_status":    String,
	"upstream_response_length": Uint64,
	"upstream_response_time":   Float64,
	"upstream_status":          Uint16,
	"uri":                      String,
}

func convert(key string, strval string) (interface{}, error) {
	vartype, ok := typemap[key]
	if !ok {
		return nil, errors.New("Field converter not found (" + key + ")")
	}

	switch vartype {

	case String:
		if strval == "-" {
			return "", nil
		}
		return strval, nil

	case Time:
		if strval == "-" {
			return time.Unix(0, 0), nil
		}
		return time.Parse(NginxTimeLayout, strval)

	case Float64:
		if strval == "-" {
			return float64(0), nil
		}
		return strconv.ParseFloat(strval, 64)

	case Uint64:
		if strval == "-" {
			return uint64(0), nil
		}
		return strconv.ParseUint(strval, 10, 64)

	case Uint16:
		if strval == "-" {
			return uint16(0), nil
		}
		res, err := strconv.ParseUint(strval, 10, 64)
		if err == nil {
			return uint16(res), err
		}
		return nil, err
	}

	return nil, errors.New("Field converter not found (" + key + ")")
}

var struct_map = make(map[string]int)

func init() {
	t := reflect.TypeOf(AccessLogEvent{})
	for fieldNum := 0; fieldNum < t.NumField(); fieldNum++ {
		typeField := t.Field(fieldNum)
		typeField.Type.Kind()
		tag := typeField.Tag
		mapkey := tag.Get("field")
		struct_map[mapkey] = fieldNum
	}
}

func (this *AccessLogServer) parseURL(rawurl string, urlParsed *URLParsed) error {
	info, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	urlParsed.scheme = info.Scheme
	urlParsed.domain = info.Host
	urlParsed.path = info.Path
	query := info.Query()
	for key, val := range query {
		urlParsed.arg_keys = append(urlParsed.arg_keys, key)
		urlParsed.arg_vals = append(urlParsed.arg_vals, val)
	}
	return nil
}

func (this *AccessLogServer) parseUserAgent(uastring string, uainfo *UserAgentInfo) {
	client := this.uaparser.Parse(uastring)

	uainfo.ua_family = client.UserAgent.Family
	uainfo.ua_major = client.UserAgent.Major
	uainfo.ua_minor = client.UserAgent.Minor
	uainfo.ua_patch = client.UserAgent.Patch
	uainfo.ua_os_family = client.Os.Family
	uainfo.ua_os_major = client.Os.Major
	uainfo.ua_os_minor = client.Os.Minor
	uainfo.ua_os_patch = client.Os.Patch
	uainfo.ua_os_patchminor = client.Os.PatchMinor
	uainfo.ua_device_family = client.Device.Family
}

func (this *AccessLogServer) parseEventURL(output *AccessLogEvent) error {
	url_parsed := &output.url_parsed
	url_parsed.scheme = output.Scheme
	url_parsed.domain = output.Http_host
	path, querystring := splitRequestURI(output.Request_uri)
	url_parsed.path = path
	query, err := url.ParseQuery(querystring)
	if err != nil {
		return err
	}
	for key, val := range query {
		url_parsed.arg_keys = append(url_parsed.arg_keys, key)
		url_parsed.arg_vals = append(url_parsed.arg_vals, val)
	}

	return nil
}

func (this *AccessLogServer) fields2event(fields gonx.Fields, output *AccessLogEvent) {
	elem := reflect.ValueOf(output).Elem()
	for field_key, field_val := range fields {
		fieldNum, ok := struct_map[field_key]
		if !ok {
			continue
		}
		converted, err := convert(field_key, field_val)
		if err != nil {
			glog.Warningln(err)
			continue
		}
		if converted == nil {
			continue
		}

		elem.Field(fieldNum).Set(reflect.ValueOf(converted))
	}

	zonename, zoneoffset := output.Time_local.Zone()
	output.zonename = zonename
	output.zoneoffset = int32(zoneoffset)

	_ = this.parseURL(output.Http_referer, &output.referer_parsed)
	if err := this.parseEventURL(output); err != nil {
		glog.Warningln(err)
	}
	this.parseUserAgent(output.Http_user_agent, &output.UserAgentInfo)

	output.nginx_ip_uint32 = IP2UInt32(net.ParseIP(output.nginx_ip))
	output.remote_addr_uint32 = IP2UInt32(net.ParseIP(output.Remote_addr))
}
