package seslog

import (
	"time"
)

type Options struct {
	CHDSN   string
	Address string
}

type NginxEvent struct {
	nginx_ip         string
	nginx_ip_uint32  uint32
	nginx_event_time time.Time
	nginx_hostname   string
	nginx_tag        string
}

type TimeZoneInfo struct {
	zonename   string
	zoneoffset int32
}

type URLParsed struct {
	scheme   string
	domain   string
	path     string
	arg_keys []string
	arg_vals [][]string
}

type ConnectionInfo struct {
	Connections_active  uint16 `field:"connections_active"`
	Connections_reading uint16 `field:"connections_reading"`
	Connections_waiting uint16 `field:"connections_waiting"`
	Connections_writing uint16 `field:"connections_writing"`
}

type GeoipInfo struct {
	Geoip_country_code string  `field:"geoip_country_code"`
	Geoip_latitude     float64 `field:"geoip_latitude"`
	Geoip_longitude    float64 `field:"geoip_longitude"`
}

type UserAgentInfo struct {
	ua_family        string
	ua_major         string
	ua_minor         string
	ua_patch         string
	ua_os_family     string
	ua_os_major      string
	ua_os_minor      string
	ua_os_patch      string
	ua_os_patchminor string
	ua_device_family string
}

type UpstreamInfo struct {
	Upstream_response_length uint64  `field:"upstream_response_length"`
	Upstream_response_time   float64 `field:"upstream_response_time"`
	Upstream_status          uint16  `field:"upstream_status"`
}

type AccessLogEvent struct {
	NginxEvent
	TimeZoneInfo
	ConnectionInfo
	GeoipInfo
	UpstreamInfo
	UserAgentInfo
	url_parsed     URLParsed
	referer_parsed URLParsed

	Body_bytes_sent          uint64    `field:"body_bytes_sent"`
	Connections_active       uint16    `field:"connections_active"`
	Connections_reading      uint16    `field:"connections_reading"`
	Connections_waiting      uint16    `field:"connections_waiting"`
	Connections_writing      uint16    `field:"connections_writing"`
	Content_length           uint64    `field:"content_length"`
	Geoip_country_code       string    `field:"geoip_country_code"`
	Geoip_latitude           float64   `field:"geoip_latitude"`
	Geoip_longitude          float64   `field:"geoip_longitude"`
	Http_host                string    `field:"http_host"`
	Http_referer             string    `field:"http_referer"`
	Http_user_agent          string    `field:"http_user_agent"`
	Http_x_forwarded_for     string    `field:"http_x_forwarded_for"`
	Remote_addr              string    `field:"remote_addr"`
	Request_method           string    `field:"request_method"`
	Request_time             float64   `field:"request_time"`
	Request_uri              string    `field:"request_uri"`
	Scheme                   string    `field:"scheme"`
	Status                   uint16    `field:"status"`
	Tcpinfo_rtt              uint64    `field:"tcpinfo_rtt"`
	Tcpinfo_rttvar           uint64    `field:"tcpinfo_rttvar"`
	Time_local               time.Time `field:"time_local"`
	Upstream_cache_status    string    `field:"upstream_cache_status"`
	Upstream_response_length uint64    `field:"upstream_response_length"`
	Upstream_response_time   float64   `field:"upstream_response_time"`
	Upstream_status          uint16    `field:"upstream_status"`
	Uri                      string    `field:"uri"`

	remote_addr_uint32 uint32
}
