# seslog

`seslog` is nginx syslog server.
It collects nginx's access logs and write them to ClickHouse DB.

## Install
Add new `log_format` to your nginx config (http section)
```
log_format seslog_format '$body_bytes_sent\t$connections_active\t$connections_reading\t$connections_waiting\t$connections_writing\t$content_length\t$http_host\t$http_referer\t$http_user_agent\t$http_x_forwarded_for\t$remote_addr\t$request_method\t$request_time\t$request_uri\t$scheme\t$status\t$tcpinfo_rtt\t$tcpinfo_rttvar\t$time_local\t$upstream_cache_status\t$upstream_response_length\t$upstream_response_time\t$upstream_status\t$uri\t$sent_http_location';
```

Then add access_log (preferred section)  
```
access_log syslog:server=<YOUR_SESLOG_IP>:5514,tag=<YOUR_PROJECT_NAME> seslog_format;
```
For example:
```
access_log syslog:server=127.0.0.1:5514,tag=php_admin_panel seslog_format if=$sesloggable;
```

## Tips
Please use `$loggable` (or anything like that) variable to avoid useless logging

e.g. (http context):
```
map $request_uri $sesloggable {
    default                                             1;
    ~*\.(ico|css|js|gif|jpg|jpeg|png|svg|woff|ttf|eot)$ 0;
}
```