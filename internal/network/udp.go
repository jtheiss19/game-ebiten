package network

import (
	"net"

	"github.com/sirupsen/logrus"
)

func handleUDPRequest(conn net.Conn, packet *Packet) {
	logrus.Warn("using default UDP request handler")
}

func StartListen() {
	go listenUDP()
	go ListenTCP()
}

func listenUDP() {

	dst, err := net.ResolveUDPAddr("udp", host+udpPort)
	if err != nil {
		logrus.Error("UDP Bad creation")
		logrus.Error(err)
	}
	conn, err := net.ListenUDP("udp", dst)
	if err != nil {
		logrus.Error("UDP Bad Connection")
		logrus.Error(err)
	}
	defer conn.Close()

	logrus.Info("starting UDP Listener")
	for {

		p, err := ReadPacketUDP(conn)
		if err != nil {
			logrus.Error(err)
			continue
		}

		// Handle
		logrus.Debug("handling udp")
		go HandleUDPRequestFunc(conn, p)
	}
}

func SendUDP(data interface{}, typeOfRequest string) error {
	// Conenct
	dst, err := net.ResolveUDPAddr("udp", host+udpPort)
	if err != nil {
		return err
	}
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return err
	}
	defer conn.Close()

	SendPacketUDP(conn, dst, data, typeOfRequest)

	return nil
}
