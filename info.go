package main

func aboutLogFormat(proxyDomain string) (res string) {
	var hr = "---------------------------------------"
	res = hr + "\nДобавить в /etc/hosts:\n127.0.0.1    " + proxyDomain + "\nили заменить на свой backend host\n" + hr
	res += `
Добавить в nginx.conf:
http {
	…
	proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=all:32m max_size=1g inactive=60m use_temp_path=off;
	…
	##
	# Logging Settings
	##
	log_format nginxgen_proxylog '$http_x_forwarded_for - $remote_addr [$time_local] '
            '"$request" $status $body_bytes_sent "$http_referer" '
            '"$http_user_agent"';
	…
}

Перезапустите nginx:
service nginx restart
`
	return res + hr
}
