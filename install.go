package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
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

func getIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	log.Println("Starting / Начало работы")

	centos := flag.Bool("centos", false, "On CentOS/RHEL")

	portStats := flag.String("p", "", "Is the local port for stats")
	port := flag.String("H", "443", "Is the port, used by clients to connect to the proxy")
	secret := flag.String("S", randomHex(16), "Secret")
	tag := flag.String("P", "", "Ad tag get here @MTProxybot")
	domain := flag.String("D", "www.google.com", "Domain with TLS 1.3 support")

	uninstall := flag.String("uninstall", "", "Removing a server")

	flag.Parse()

	if *uninstall != "" {
		cmd("systemctl stop MTProxy-" + *uninstall + ".service")
		cmd("systemctl disable MTProxy-" + *uninstall + ".service")
		cmd("rm /etc/systemd/system/MTProxy-" + *uninstall + ".service")

		log.Println("Uninstall complete / Удаление завершено")
		return
	}

	log.Println("Dependency check / Проверка зависимостей")

	if *centos {
		cmd("yum -y install openssl-devel zlib-devel")
		cmd("yum -y groupinstall \"Development Tools\"")
	} else {
		cmd("apt -y install git make build-essential libssl-dev zlib1g-dev")
	}

	log.Println("Installing / Установка")
	cmd("git clone https://github.com/TelegramMessenger/MTProxy && cd MTProxy && make && cd objs/bin && " +
		"curl -s https://core.telegram.org/getProxySecret -o proxy-secret && " +
		"curl -s https://core.telegram.org/getProxyConfig -o proxy-multi.conf")

	cmd("cd /opt && mkdir mtproxy")
	cmd("cp MTProxy/objs/bin/mtproto-proxy /opt/mtproxy/mtproto-proxy")
	cmd("cp MTProxy/objs/bin/proxy-multi.conf /opt/mtproxy/proxy-multi.conf")
	cmd("cp MTProxy/objs/bin/proxy-secret /opt/mtproxy/proxy-secret")

	cmd("rm -r MTProxy")

	log.Println("Creating a service / Создание службы")
	cmd("touch /etc/systemd/system/MTProxy-" + *port + ".service")

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

	systemd := `[Unit]
Description=MTProxy
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/mtproxy
ExecStart=/opt/mtproxy/mtproto-proxy -u nobody` + first + " -H " + *port + " -S " + *secret + second + third + ` --aes-pwd proxy-secret proxy-multi.conf
Restart=on-failure

[Install]
WantedBy=multi-user.target`

	cmd("echo \"" + systemd + "\" >> /etc/systemd/system/MTProxy-" + *port + ".service")

	cmd("systemctl daemon-reload")
	cmd("systemctl restart MTProxy-" + *port + ".service")
	cmd("systemctl enable MTProxy-" + *port + ".service")

	src := []byte(*domain)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	fmt.Println("\n\n\ntg://proxy?server=" + getIP().String() + "&port=" + *port + "&secret=ee" + *secret + fmt.Sprintf("%s\n", dst)+"\n\n\n")
}
