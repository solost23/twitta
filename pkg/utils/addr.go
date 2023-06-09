package utils

import (
	"errors"
	"net"
	"strings"
)

// GetFreePort 获取当前未被占用的端口
// return: port, error
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func(l *net.TCPListener) {
		err = l.Close()
		if err != nil {
			panic(err)
		}
	}(listener)

	return listener.Addr().(*net.TCPAddr).Port, nil
}

// GetInternalIP 获取内网 IP
// return: ip, error
func GetInternalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.New("internal IP fetch failed, detail:" + err.Error())
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	res := conn.LocalAddr().String()
	res = strings.Split(res, ":")[0]
	return res, nil
}
