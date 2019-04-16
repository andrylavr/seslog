# seslog

`seslog` is nginx syslog server.
It collects nginx's access logs and write them to ClickHouse DB.

## Install
Add new `log_format` to your nginx config (http section)
```
log_format seslog_format '$body_bytes_sent	$connections_active	$connections_reading	$connections_waiting	$connections_writing	$content_length	$geoip_country_code	$geoip_latitude	$geoip_longitude	$http_host	$http_referer	$http_user_agent	$http_x_forwarded_for	$remote_addr	$request_method	$request_time	$request_uri	$scheme	$status	$tcpinfo_rtt	$tcpinfo_rttvar	$time_local	$upstream_cache_status	$upstream_response_length	$upstream_response_time	$upstream_status	$uri';
```

Then add access_log (preferred section)  
```
access_log syslog:server=<YOUR_SESLOG_IP>:5514,tag=<YOUR_PROJECT_NAME> seslog_format;
```
For example:
```
access_log syslog:server=127.0.0.1:5514,tag=php_admin_panel seslog_format if=$loggable;
```

## Tips
Please use `$loggable` (or anything like that) variable to avoid useless logging

e.g. (http context):
```
map $request_uri $loggable {
    default                                             1;
    ~*\.(ico|css|js|gif|jpg|jpeg|png|svg|woff|ttf|eot)$ 0;
}
```