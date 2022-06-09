package network

import (
	"encoding/gob"
	"net"

	"github.com/sirupsen/logrus"
)

func handleTCPRequest(enc *gob.Encoder, packet *Packet) {
	logrus.Warn("using default TCP request handler")
}

func StartTCPConnection() (net.Conn, error) {
	// Connect
	conn, err := net.Dial("tcp", host+tcpPort)
	if err != nil {
		return &net.TCPConn{}, err
	}

	go tcpHandleConnectionClient(conn)

	return conn, nil
}

func ListenTCP() {
	// Create Listener
	listener, err := net.Listen("tcp", host+tcpPort)
	if err != nil {
		logrus.Error("TCP Bad Listener Start")
		logrus.Error(err)
	}
	defer listener.Close()

	// Start Listener
	logrus.Info("starting TCP Listener")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Error("TCP Bad Connection Acception")
			logrus.Error(err)
			continue
		}
		logrus.Info("TCP new connection")

		// Handle Every New Connection
		go tcpHandleConnectionServer(conn)
	}
}

func tcpHandleConnectionServer(conn net.Conn) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	for {
		// Read
		p, err := ReadPacketTCP(dec)
		if err != nil {
			logrus.Debug(err)
			break
		}

		// Handle
		go HandleTCPRequestFunc(enc, p)
	}
	logrus.Info("closing connection")
}

func tcpHandleConnectionClient(conn net.Conn) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	for {
		// Read
		p, err := ReadPacketTCP(dec)
		if err != nil {
			logrus.Debug(err)
			break
		}

		// Handle
		HandleTCPResponseFunc(enc, p)
	}
	logrus.Info("closing connection")
}
