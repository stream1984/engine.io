package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

type ppp struct {
}

var Payload *ppp

func init() {
	Payload = new(ppp)
}

func writePacket(bf *bytes.Buffer, packet *Packet) error {
	var data []byte
	var err error
	var length int
	if packet.Option&BINARY != BINARY {
		data, err = stringEncoder.encode(packet)
		if err != nil {
			return err
		}
		length = utf8.RuneCount(data)
	} else {
		data, err = base64Encoder.encode(packet)
		if err != nil {
			return err
		}
		length = len(data)
	}
	_, err = bf.WriteString(fmt.Sprintf("%d:", length))
	if err != nil {
		return err
	}
	_, err = bf.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (p *ppp) Encode(packets ...*Packet) ([]byte, error) {
	if len(packets) < 1 {
		return nil, errors.New("input packets is empty")
	}
	bf := new(bytes.Buffer)
	for _, it := range packets {
		if err := writePacket(bf, it); err != nil {
			return nil, err
		}
	}
	return bf.Bytes(), nil
}

func (p *ppp) DecodeString(str string) ([]*Packet, error) {
	return p.Decode([]byte(str))
}

func (p *ppp) Decode(input []byte) ([]*Packet, error) {
	var size int
	var err error
	var rest, content []byte = input, nil
	var packets = make([]*Packet, 0)
	var packet *Packet
	for len(rest) > 0 {
		if size == 0 {
			size, rest, err = readPacketLength(rest)
		} else {
			content, rest, err = readPacketString(rest, size)
			packet, err = readPacket(content)
			if err != nil {
				return nil, err
			}
			packets = append(packets, packet)
			size = 0
		}
	}
	return packets, err
}

func readPacket(input []byte) (*Packet, error) {
	if input[0] != 'b' {
		return stringEncoder.decode(input)
	} else {
		return base64Encoder.decode(input)
	}
}

func readPacketLength(input []byte) (int, []byte, error) {
	for i := 0; i < len(input); i++ {
		if input[i] != ':' {
			continue
		}
		size, err := strconv.Atoi(string(input[:i]))
		if err != nil {
			return 0, nil, err
		}
		return size, input[i+1:], nil
	}
	return 0, nil, errors.New("invalid payload string")
}

func readPacketString(input []byte, size int) ([]byte, []byte, error) {
	var i, w int
	for i, w = 0, 0; i < len(input) && size > 0; {
		_, width := utf8.DecodeRune(input[i:])
		w = width
		size--
		i += w
	}
	return input[:i], input[i:], nil
}