package dns

import (
	"fmt"
	"math/rand"
	"strings"
)

type RR struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RdLength uint16
	Data     string
}
type Message struct {
	ID          uint16
	Flag        uint16
	QDCount     uint16
	ANCount     uint16
	NSCount     uint16
	ARCount     uint16
	Anwers      []RR
	Autorities  []RR
	Additionals []RR
}

// 解析域名
func parseName(b []byte, index int) (string, int) {
	var name []string
	for b[index] != 0x00 {
		if b[index] >= 0xC0 {
			origin := index + 2
			high := b[index]
			low := b[index+1]
			index = int(high&0x3F)<<8 | int(low)
			s, _ := parseName(b, index)
			name = append(name, s)
			return strings.Join(name, "."), origin
		}
		length := int(b[index])
		index++
		name = append(name, string(b[index:length+index]))
		index += length
	}
	index++
	return strings.Join(name, "."), index
}
func parseRR(b []byte, index int) (RR, int) {
	// for range cnt {
	var answer RR
	answer.Name, index = parseName(b, index)
	answer.Type = uint16(b[index])<<8 | uint16(b[index+1])
	index += 2
	answer.Class = uint16(b[index])<<8 | uint16(b[index+1])
	index += 2
	answer.TTL = uint32(b[index])<<24 | uint32(b[index+1])<<16 | uint32(b[index+2])<<8 | uint32(b[index+3])
	index += 4
	answer.RdLength = uint16(b[index])<<8 | uint16(b[index+1])
	index += 2
	switch answer.Type {
	case 1:
		answer.Data = fmt.Sprintf("%d.%d.%d.%d", b[index], b[index+1], b[index+2], b[index+3])
		index += 4
	case 2, 5:
		name, idx := parseName(b, index)
		answer.Data = fmt.Sprintf("%s", name)
		index = idx
	default:
		index += int(answer.RdLength)
	}
	return answer, index
	//}
}

// 解析二进制响应
func parseResp(b []byte) Message {
	var msg Message
	msg.ID = uint16(b[0])<<8 | uint16(b[1])
	msg.Flag = uint16(b[2])<<8 | uint16(b[3])
	msg.QDCount = uint16(b[4])<<8 | uint16(b[5])
	msg.ANCount = uint16(b[6])<<8 | uint16(b[7])
	msg.NSCount = uint16(b[8])<<8 | uint16(b[9])
	msg.ARCount = uint16(b[10])<<8 | uint16(b[11])
	index := 12
	//跳过questionsection
	for range msg.QDCount {
		//跳过域名
		for b[index] != 0x00 {
			length := int(b[index])
			index += length + 1
		}
		//跳过0x00，type，class
		index += 5
	}
	var answer RR
	for range msg.ANCount {
		answer, index = parseRR(b, index)
		msg.Anwers = append(msg.Anwers, answer)
	}
	for range msg.NSCount {
		answer, index = parseRR(b, index)
		msg.Autorities = append(msg.Autorities, answer)
	}
	for range msg.ARCount {
		answer, index = parseRR(b, index)
		msg.Additionals = append(msg.Additionals, answer)
	}
	return msg
}

// 构造二进制报文
func buildQuery(s string, qtype uint16) []byte {
	buf := make([]byte, 0)
	//Header:
	//transID
	id := uint16(rand.Intn(65535))
	buf = append(buf, byte(id>>8), byte(id&0xFF))
	//flags
	buf = append(buf, 0x01, 0x00)
	//QDCount: 1
	buf = append(buf, 0x00, 0x01)
	//ANCount, NSCount, ARCount: 全0
	buf = append(buf, 0x00, 0x00)
	buf = append(buf, 0x00, 0x00)
	buf = append(buf, 0x00, 0x00)
	//域名
	dn := strings.Split(s, ".")
	for _, name := range dn {
		buf = append(buf, byte(len(name)))
		buf = append(buf, []byte(name)...)
	}
	buf = append(buf, 0x00)
	// Type A = 0x0001
	buf = append(buf, byte(qtype>>8), byte(qtype&0xFF))
	// Class IN = 0x0001
	buf = append(buf, 0x00, 0x01)
	return buf
}
