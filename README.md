# Nginxgen

## Начало работы

Настроить nginx.

```
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
```

# Использование

```
usage: nginxgen --php=PHP --root=ROOT --domain=DOMAIN [<flags>]

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
      --php=PHP           Версия PHP (7.2, 7.4 и тп).
      --root=ROOT         Абсолютный путь до директории проекта.
      --domain=DOMAIN     Домен без www.
  -p, --proxy-port=80     Порт proxy.
      --public="/public"  Относительный путь до public директории от директории проекта.
      --static=STATIC     Директории в public со статическим контентом через запятую (css, image, files …).
      --ip=IP             IP сервера.
  -s, --ssl               С SSL-серитификатом.
```

Пример:

```
nginxgen --php=7.4 --root=/home/user/project/ --domain=site.ru --public=/public
```

# SSL

## Установка

Установить `dehydrated`:

```
sudo apt install dehydrated
```

## Домены

Прописываем домены:

```
sudo nano /etc/dehydrated/domains.txt
```

Пример:

```
# comment site.ru
example.com www.example.com
example2.com www.example2.com test.example2.com
```

# Первый запуск

```
sudo mkdir -p /etc/nginx/ssl
sudo openssl dhparam -dsaparam -out /etc/nginx/ssl/dhparams.pem 2048

sudo dehydrated --register --accept-terms
```

# Получение/обновление сертификата

Получить/обновить сертификаты вручную:

```
sudo dehydrated -c -g && sudo nginx -s reload
```

Cron-задача `root` для обновления:

```
0 4 6 * * dehydrated -c -g && nginx -s reload
```