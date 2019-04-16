package seslog

import (
	"time"
)

type Options struct {
	CHDSN         string
	Address       string
	FlushInterval string
}

type NginxEvent struct {
	Nginx_ip         string
	Nginx_ip_uint32  uint32
	Nginx_event_time time.Time
	Nginx_hostname   string
	Nginx_tag        string
}

type TimeZoneInfo struct {
	Zonename   string
	Zoneoffset int32
}

type URLParsed struct {
	Scheme   string
	Domain   string
	Path     string
	Arg_keys []string
	Arg_vals [][]string
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
	Ua_family        string
	Ua_major         string
	Ua_minor         string
	Ua_patch         string
	Ua_os_family     string
	Ua_os_major      string
	Ua_os_minor      string
	Ua_os_patch      string
	Ua_os_patchminor string
	Ua_device_family string
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
	Url_parsed     URLParsed
	Referer_parsed URLParsed

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

	Remote_addr_uint32 uint32
}

type AccessLogEvents []AccessLogEvent
