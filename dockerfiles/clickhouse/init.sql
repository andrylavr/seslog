CREATE DATABASE IF NOT EXISTS seslog;

CREATE TABLE IF NOT EXISTS seslog.access_log(
  nginx_event_date Date DEFAULT toDate(nginx_event_time),
  nginx_tag String,
  nginx_event_time DateTime DEFAULT now(),
  nginx_hostname String,
  nginx_ip_uint32 UInt32,

  zonename String,
  zoneoffset Int32,

	body_bytes_sent UInt64,
	connections_active UInt16,
	connections_reading UInt16,
	connections_waiting UInt16,
	connections_writing UInt16,
	content_length UInt64,

	geoip_country_code FixedString(2),
	geoip_latitude Float64,
	geoip_longitude Float64,

	http_scheme String,
  http_domain String,
  http_path String,
  http_arg_keys Array(String),
  http_arg_vals Array(Array(String)),

  http_referer_scheme String,
  http_referer_domain String,
  http_referer_path String,
  http_referer_arg_keys Array(String),
  http_referer_arg_vals Array(Array(String)),

  ua_family String,
  ua_major String,
  ua_minor String,
  ua_patch String,
  ua_os_family String,
  ua_os_major String,
  ua_os_minor String,
  ua_os_patch String,
  ua_os_patchminor String,
  ua_device_family String,

	http_x_forwarded_for String,
	remote_addr_uint32 UInt32,
	request_method String,
	request_time Float64,

	status UInt16,
	tcpinfo_rtt UInt64,
	tcpinfo_rttvar UInt64,

	upstream_response_length UInt64,
	upstream_response_time Float64,
	upstream_status UInt16,
	uri String
) ENGINE = MergeTree
        PARTITION BY toYYYYMM(nginx_event_date)
        ORDER BY (nginx_event_date, nginx_tag, nginx_event_time, nginx_hostname);
