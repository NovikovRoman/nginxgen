package main

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	// symfony    = flag.Bool("symfony", false, "Конфигурация для symfony 4.")

	// Обязательные
	phpVersion = flag.String("php-version", "", "Версия PHP (7.2, 7.4).")
	root       = flag.String("root", "", "Абсолютный путь до public.")
	domain     = flag.String("domain", "", "Домен без www.")
	proxyPort  = flag.Int("proxy", 0, "Порт proxy.")

	// Необязательные
	static = flag.String("static", "",
		"Директории в public со статическим контентом через запятую (css, image, files …).")
	ip    = flag.String("ip", "", "IP сервера.")
	https = flag.Bool("https", false, "С https")
)

func main() {
	flag.Parse()

	if *phpVersion == "" || *root == "" || *domain == "" || *proxyPort <= 0 {
		log.Fatalln(help())
	}

	projectPath := regexp.MustCompile(`(?si)/[^/]+$`).ReplaceAllString(*root, "")

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

	// изменим контент нужных сниппетов и сохраним для проекта
	b, err := symfony4()
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(filepath.Join(snippetPath, "nginxgen-symfony4-"+*phpVersion+".conf"), b, permission)
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

	b, err = createConfig(logPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(filepath.Join(configPath, *domain+".conf"), b, permission)

	fmt.Println(aboutLogFormat())
}

func createConfig(logPath string) (b []byte, err error) {
	if *https {
		b, err = ioutil.ReadFile(filepath.Join(templates, "https.conf"))

	} else {
		b, err = ioutil.ReadFile(filepath.Join(templates, "http.conf"))
	}

	if err != nil {
		return
	}

	b = bytes.ReplaceAll(b, []byte("{domain}"), []byte(*domain))
	b = bytes.ReplaceAll(b, []byte("{root}"), []byte(*root))
	b = bytes.ReplaceAll(b, []byte("{proxy}"), []byte(strconv.Itoa(*proxyPort)))

	listen80 := "80"
	listen443 := "443"
	listenProxy := strconv.Itoa(*proxyPort)
	if *ip != "" {
		listen80 = *ip + ":" + listen80
		listen443 = *ip + ":" + listen443
		listenProxy = *ip + ":" + listenProxy
	}

	b = bytes.ReplaceAll(b, []byte("{listen_80}"), []byte(listen80))
	b = bytes.ReplaceAll(b, []byte("{listen_443}"), []byte(listen443))
	b = bytes.ReplaceAll(b, []byte("{listen_proxy}"), []byte(listenProxy))
	b = bytes.ReplaceAll(b, []byte("{php-version}"), []byte(*phpVersion))

	b = bytes.ReplaceAll(b, []byte("{access_log}"), []byte(logPath+"/access.log main"))
	b = bytes.ReplaceAll(b, []byte("{error_log}"), []byte(logPath+"/error.log"))
	b = bytes.ReplaceAll(b, []byte("{proxy_access_log}"), []byte(logPath+"/access-proxy.log main"))
	b = bytes.ReplaceAll(b, []byte("{proxy_error_log}"), []byte(logPath+"/error-proxy.log"))

	staticDirs := ""
	if *static != "" {
		var bStaticDirs []byte
		bStaticDirs, err = getStaticDirs()
		if err != nil {
			return
		}
		staticDirs = string(bStaticDirs)
	}
	b = bytes.ReplaceAll(b, []byte("{static-dirs}"), []byte(staticDirs))
	return
}

func symfony4() (b []byte, err error) {
	b, err = ioutil.ReadFile(filepath.Join(templates, "symfony4.conf"))
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
