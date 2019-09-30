package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

func cmd(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil && err.Error() != "exit status 1" && err.Error() != "exit status 2" {
		log.Println("[error]", err, "("+cmd+")")
	}

	return string(out)
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println("[error]", err, "(Failed to determine IP | Не удалось определить IP)")
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("[error]", err)
		}
	}()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func main() {
	log.Println("          Starting | Начало работы")

	centos := flag.Bool("centos", false, "On CentOS/RHEL")

	portStats := flag.String("p", "", "Is the local port for stats")
	port := flag.String("H", "443", "Is the port, used by clients to connect to the proxy")
	secret := flag.String("S", randomHex(16), "Secret")
	tag := flag.String("P", "", "Ad tag get here @MTProxybot")
	domain := flag.String("D", "www.google.com", "Domain with TLS 1.3 support")

	uninstall := flag.String("uninstall", "", "Removing a server")

	flag.Parse()

	path := "/etc/systemd/system/MTProxy-" + *port + ".service"

	if *uninstall != "" {
		cmd("systemctl stop MTProxy-" + *uninstall + ".service")
		cmd("systemctl disable MTProxy-" + *uninstall + ".service")
		cmd("rm /etc/systemd/system/MTProxy-" + *uninstall + ".service")

		log.Println("Uninstall complete | Удаление завершено")
		log.Println(" Program completed | Программа завершена")
		return
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Println("A server with such a port has already been created, overwrite it?")
		log.Println("Сервер с таким портом уже создан, перезаписать его?")

		answer := ""
		fmt.Print("\n\ny/N: ")
		_, _ = fmt.Scan(&answer)

		switch answer {
		case "y":
			cmd("systemctl stop MTProxy-" + *port + ".service")
			cmd("systemctl disable MTProxy-" + *port + ".service")
			cmd("rm " + path)
		case "Y":
			cmd("systemctl stop MTProxy-" + *port + ".service")
			cmd("systemctl disable MTProxy-" + *port + ".service")
			cmd("rm " + path)
		case "n":
			log.Println("Program completed | Программа завершена")
			return
		case "N":
			log.Println("Program completed | Программа завершена")
			return
		default:
			log.Println("Invalid input | Некорректный ввод")
			log.Println("Program completed | Программа завершена")
			return
		}
	}

	log.Println("  Dependency check | Проверка зависимостей")

	if *centos {
		cmd("yum update")
		cmd("yum -y install openssl-devel zlib-devel")
		cmd("yum -y groupinstall \"Development Tools\"")
	} else {
		cmd("apt update")
		cmd("apt -y install git make build-essential libssl-dev zlib1g-dev")
	}

	log.Println("        Installing | Установка")
	cmd("git clone https://github.com/hookzof/MTProxy && cd MTProxy && make && cd objs/bin && " +
		"curl -s https://core.telegram.org/getProxySecret -o proxy-secret && " +
		"curl -s https://core.telegram.org/getProxyConfig -o proxy-multi.conf")

	cmd("cd /opt && mkdir mtproxy")
	cmd("cp MTProxy/objs/bin/mtproto-proxy /opt/mtproxy/mtproto-proxy")
	cmd("cp MTProxy/objs/bin/proxy-multi.conf /opt/mtproxy/proxy-multi.conf")
	cmd("cp MTProxy/objs/bin/proxy-secret /opt/mtproxy/proxy-secret")

	cmd("rm -r MTProxy")

	log.Println("Creating a service | Создание службы")
	cmd("touch " + path)

	first := ""
	if *portStats != "" {
		first = " -p " + *portStats
	}

	second := ""
	if *tag != "" {
		second = " -P " + *tag
	}

	third := " -D www.google.com"
	if *domain != "www.google.com" {
		third = " -D " + *domain
	}

	config := `[Unit]
Description=MTProxy
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/mtproxy
ExecStart=/opt/mtproxy/mtproto-proxy -u nobody` + first + " -H " + *port + " -S " + *secret + second + third + ` --aes-pwd proxy-secret proxy-multi.conf
Restart=on-failure
LimitNOFILE=infinity
LimitMEMLOCK=infinity

[Install]
WantedBy=multi-user.target`

	cmd("echo \"" + config + "\" >> " + path)

	cmd("systemctl daemon-reload")
	cmd("systemctl restart MTProxy-" + *port + ".service")
	cmd("systemctl enable MTProxy-" + *port + ".service")

	src := []byte(*domain)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	log.Println(" Program completed | Программа завершена")

	fmt.Println("\n\n\nServer file path | Путь файлов сервера     — /opt/mtproxy/")
	fmt.Println("Config file path | Путь файла конфигурации — " + path)
	fmt.Println("\ntg://proxy?server=" + getIP() + "&port=" + *port + "&secret=ee" + *secret + fmt.Sprintf("%s", dst) + "\n\n\n")
}
