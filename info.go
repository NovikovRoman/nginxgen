package main

func aboutLogFormat(domain, proxyDomain, root string) (res string) {
	domainConf := domain + ".conf"
	hr := "---------------------------------------"
	res =  "\n\033[0;31mНастройте `nginx.conf`. Для SSL установите `dehydrated`. Подробнее в README.md.\nЕсли это первая генерация, то перенесите сниппеты в nginx:\033[0m\nsudo mv " + root + "nginxgen/snippets/* /etc/nginx/snippets/\n" + hr
	res += "\n\n\033[0;32mДобавить в /etc/hosts:\033[0m\n127.0.0.1    " + proxyDomain + "\nили заменить на свой backend host\n" + hr
	res += "\n\n\033[0;32mСкопировать конфиг:\033[0m\nsudo cp " + root + "nginxgen/" + domainConf + " /etc/nginx/sites-available/\n\n"
	res += "\033[0;32mСоздать символическую ссылку:\033[0m\nsudo ln -s /etc/nginx/sites-available/" + domainConf + " /etc/nginx/sites-enabled/" + domain + "\n"
	return res + hr
}
