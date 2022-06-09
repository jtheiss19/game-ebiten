package network

import (
	"bytes"
	"encoding/gob"
	"net"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Message struct {
	Type string
	Data interface{}
}

type Packet struct {
	Parts   []int
	Message Message
}

var (
	typeMap = map[string]reflect.Type{
		"join request": reflect.TypeOf(""),
		"update comp":  reflect.TypeOf(""),
		"new comp":     reflect.TypeOf(""),
		"id":           reflect.TypeOf(""),
	}

	host    = "localhost:"
	tcpPort = "8081"
	udpPort = "8080"

	HandleUDPRequestFunc  = handleUDPRequest
	HandleTCPRequestFunc  = handleTCPRequest
	HandleTCPResponseFunc = handleTCPRequest
)

func RegisterType(thingToRegister interface{}) {
	gob.Register(thingToRegister)
}

func ReadPacketTCP(dec *gob.Decoder) (*Packet, error) {

	// Decode
	p := &Packet{}
	err := dec.Decode(&p)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Check
	_, ok := typeMap[p.Message.Type]
	if !ok {
		return nil, err
	}

	return p, nil
}

func SendPacketTCP(enc *gob.Encoder, data interface{}, typeOfRequest string) error {
	// Create Packet
	responsePacket := &Packet{
		Parts: []int{1},
		Message: Message{
			Type: typeOfRequest,
			Data: data,
		},
	}

	// Encode Data
	err := enc.Encode(responsePacket)
	if err != nil {
		return err
	}

	return nil
}

func ReadPacketUDP(conn *net.UDPConn) (*Packet, error) {
	// Read
	buf := make([]byte, 4096)
	n, _, err := conn.ReadFromUDP(buf[:])
	logrus.Debug("just read from UDP connection")
	if err != nil {
		logrus.Error("UDP Bad Connection Read")
		logrus.Error(err)
		return nil, err
	}

	// Decode
	logrus.Debug("going to decode buffer from UDP read")
	dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
	p := &Packet{}
	err = dec.Decode(&p)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	logrus.Debug("just decoded buffer from UDP read")

	// Check
	_, ok := typeMap[p.Message.Type]
	if !ok {
		logrus.Error("UDP class lookup bad")
		return nil, err
	}

	return p, nil
}

func SendPacketUDP(conn net.PacketConn, dst *net.UDPAddr, data interface{}, typeOfRequest string) error {
	// Create Packet
	packet := &Packet{
		Parts: []int{1},
		Message: Message{
			Type: typeOfRequest,
			Data: data,
		},
	}

	// Encode
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(packet)
	if err != nil {
		return err
	}

	// Send
	_, err = conn.WriteTo(buf.Bytes(), dst)
	if err != nil {
		return err
	}

	return nil
}
