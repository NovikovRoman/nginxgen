package main

func help() string {
	// -symfony=…		- для Symfony 4.
	return `
-php-version=…	- Версия PHP (7.2, 7.4).
-root=…		- Абсолютный путь до директории проекта.
-domain=…	- Домен без www.
-proxy=…	- Порт proxy.

-public=…   - Относительный путь до public директории от директории проекта. По-умолчанию «public». Необязательно.
-static=…	- Директории в public со статическим контентом через запятую (css, image, files …). Необязательно.
-ip=…		- IP сервера. Необязательно.
-https		- C https. Необязательно.
`
}

func aboutLogFormat() string {
	return `
Добавьте в nginx.conf:
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
}
