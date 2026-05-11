package mydig

import (
	"math/rand"
	"strings"
)

// 解析二进制响应
func parseResp(b []byte) string {
	var msg string
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
