package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	templates           = "templates"
	permission          = 0600
	permissionDirectory = 0755
)

var (
	// Обязательные
	phpVersion = kingpin.Flag("php", "Версия PHP (7.2, 7.4 и тп).").Required().String()
	root       = kingpin.Flag("root", "Абсолютный путь до директории проекта.").Required().String()
	domain     = kingpin.Flag("domain", "Домен без www.").Required().String()

	// Необязательные
	proxyPort = kingpin.Flag("proxy-port", "Порт proxy.").Short('p').Default("80").Int()

	public = kingpin.Flag("public", "Относительный путь до public директории от директории проекта.").
		Default("/public").String()

	static = kingpin.Flag("static",
		"Директории в public со статическим контентом через запятую (css, image, files …).").String()

	ip = kingpin.Flag("ip", "IP сервера.").String()

	ssl = kingpin.Flag("ssl", "С SSL-серитификатом.").Short('s').Bool()

	publicPath  string
	proxyDomain string
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	projectPath := regexp.MustCompile(`(?si)/\s*$`).ReplaceAllString(*root, "")
	publicPath = projectPath + "/" + regexp.MustCompile(`(?si)(^\s*/|/\s*$)`).ReplaceAllString(*public, "")

	configPath := filepath.Join(projectPath, "nginxgen")
	if err := os.MkdirAll(configPath, permissionDirectory); err != nil {
		log.Fatalln(err)
	}

	snippetPath := filepath.Join(configPath, "snippets")
	if err := os.MkdirAll(snippetPath, permissionDirectory); err != nil {
		log.Fatalln(err)
	}

	// сохраним snippets
	err := saveSnippets(snippetPath,
		"nginxgen-basic-security.conf",
		"nginxgen-gzip.conf",
		"nginxgen-cut-index-php.conf",
		"nginxgen-fastcgi-php.conf",
		"nginxgen-proxy.conf",
	)
	if err != nil {
		log.Fatalln(err)
	}

	var b []byte
	// изменим контент нужных сниппетов и сохраним для проекта
	if b, err = symfony(); err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(filepath.Join(snippetPath, "nginxgen-symfony-"+*phpVersion+".conf"), b, permission)
	if err != nil {
		log.Fatalln(err)
	}

	b, err = fastcgi()
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(filepath.Join(snippetPath, "nginxgen-fastcgi"+*phpVersion+".conf"), b, permission)
	if err != nil {
		log.Fatalln(err)
	}

	logPath := projectPath + "/logs"
	if err := os.MkdirAll(logPath, permissionDirectory); err != nil {
		log.Fatalln(err)
	}

	proxyDomain = "proxy_" + regexp.MustCompile(`\.`).ReplaceAllString(*domain, "_")

	b, err = createConfig(logPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(filepath.Join(configPath, *domain+".conf"), b, permission)

	fmt.Println(aboutLogFormat( *domain,proxyDomain, *root))
}

func createConfig(logPath string) (b []byte, err error) {
	if *ssl {
		b, err = ioutil.ReadFile(filepath.Join(templates, "https.conf"))

	} else {
		b, err = ioutil.ReadFile(filepath.Join(templates, "http.conf"))
	}

	if err != nil {
		return
	}

	b = bytes.ReplaceAll(b, []byte("{domain}"), []byte(*domain))
	b = bytes.ReplaceAll(b, []byte("{public}"), []byte(publicPath))
	b = bytes.ReplaceAll(b, []byte("{proxy}"), []byte(strconv.Itoa(*proxyPort)))

	listen80 := "80"
	listen443 := "443"
	listenProxy := strconv.Itoa(*proxyPort)
	if *ip != "" {
		listen80 = *ip + ":" + listen80
		listen443 = *ip + ":" + listen443
		listenProxy = *ip + ":" + listenProxy
	}

	b = bytes.ReplaceAll(b, []byte("{proxy_domain}"), []byte(proxyDomain))
	b = bytes.ReplaceAll(b, []byte("{listen_80}"), []byte(listen80))
	b = bytes.ReplaceAll(b, []byte("{listen_443}"), []byte(listen443))
	b = bytes.ReplaceAll(b, []byte("{listen_proxy}"), []byte(listenProxy))
	b = bytes.ReplaceAll(b, []byte("{php-version}"), []byte(*phpVersion))

	b = bytes.ReplaceAll(b, []byte("{access_log}"), []byte(logPath+"/access.log"))
	b = bytes.ReplaceAll(b, []byte("{error_log}"), []byte(logPath+"/error.log"))
	b = bytes.ReplaceAll(b, []byte("{proxy_access_log}"), []byte(logPath+"/access-proxy.log"))
	b = bytes.ReplaceAll(b, []byte("{proxy_error_log}"), []byte(logPath+"/error-proxy.log"))

	staticDirs := []byte("")
	if *static != "" {
		staticDirs, err = getStaticDirs()
		if err != nil {
			return
		}
		s := strings.Split(*static, ",")
		for i := range s {
			s[i] = strings.TrimSpace(s[i])
		}
		staticDirs = bytes.ReplaceAll(staticDirs, []byte("{dirs}"), []byte(strings.Join(s, "|")))
		staticDirs = bytes.ReplaceAll(staticDirs, []byte("{public}"), []byte(publicPath))
	}

	b = bytes.ReplaceAll(b, []byte("{static-dirs}"), staticDirs)
	return
}

func symfony() (b []byte, err error) {
	b, err = ioutil.ReadFile(filepath.Join(templates, "symfony.conf"))
	if err != nil {
		return
	}

	b = bytes.ReplaceAll(b, []byte("{php-version}"), []byte(*phpVersion))
	return
}

func fastcgi() (b []byte, err error) {
	b, err = ioutil.ReadFile(filepath.Join(templates, "fastcgi.conf"))
	if err != nil {
		return
	}

	b = bytes.ReplaceAll(b, []byte("{php-version}"), []byte(*phpVersion))
	return
}

func getStaticDirs() (b []byte, err error) {
	b, err = ioutil.ReadFile(filepath.Join(templates, "static-dirs.conf"))
	if err != nil {
		return
	}

	dirs := strings.Split(*static, ",")
	for i := range dirs {
		dirs[i] = strings.TrimSpace(dirs[i])
	}

	b = bytes.ReplaceAll(b, []byte("{static-dirs}"), []byte(strings.Join(dirs, "|")))
	return
}
