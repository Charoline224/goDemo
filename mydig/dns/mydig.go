package dns

import (
	"fmt"
	"net"
)

func typeToString(t uint16) string {
	switch t {
	case 1:
		return "A"
	case 2:
		return "NS"
	case 5:
		return "CNAME"
	case 28:
		return "AAAA"
	default:
		return fmt.Sprintf("TYPE%d", t)
	}
}

func sendQuery(query []byte, server string) ([]byte, error) {
	//建立UDP连接
	conn, err := net.Dial("udp", server+":53")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//传输请求
	_, err = conn.Write(query)
	if err != nil {
		return nil, err
	}
	//返回响应
	resp := make([]byte, 512)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, err
	}
	return resp[:n], nil
}

func Resolve(domain, server string, qtype uint16, trace bool) ([]string, error) {

	var results []string
	query := buildQuery(domain, qtype)
	if trace {
		fmt.Printf(";; 查询 %s 类型%d 从服务器 %s\n", domain, qtype, server)
	}
	resp, err := sendQuery(query, server)
	if err != nil {
		return nil, err
	}
	res := parseResp(resp)
	fmt.Printf(";; ANCount:%d NSCount:%d ARCount:%d\n", res.ANCount, res.NSCount, res.ARCount)
	var cname string
	for _, answer := range res.Anwers {
		if trace {
			fmt.Printf("%s\t%d\t%s\t%s\n", answer.Name, answer.TTL, typeToString(answer.Type), answer.Data)
		}
		if qtype == answer.Type {
			results = append(results, answer.Data)
		}
		if answer.Type == 5 {
			// return resolve(answer.Data, server , answer.Type, tarce)
			cname = answer.Data
		}
	}
	//有A了直接返回
	if len(results) != 0 {
		return results, nil
	}
	//否则查查cname
	if cname != "" {
		result, err := Resolve(cname, server, 1, trace)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	//NS查
	for _, nsdo := range res.Autorities {
		if trace {
			fmt.Printf("%s\t%d\t%s\t%s\n", nsdo.Name, nsdo.TTL, typeToString(nsdo.Type), nsdo.Data)
		}
		if nsdo.Type == 2 {
			for _, nsip := range res.Additionals {
				//debug出现问题：
				// NS指令运行错误。
				//原因:
				// DNS 服务器返回的 Additional Section 里，对于每个 NS 服务器，同时给了 IPv4（A）和 IPv6（AAAA）两条记录
				if nsip.Name == nsdo.Data && nsip.Type == 1 {
					return Resolve(domain, nsip.Data, qtype, trace)
				}
			}
			nsIPs, err := Resolve(nsdo.Data, "8.8.8.8", 1, trace)
			if err != nil || len(nsIPs) == 0 {
				continue
			}
			return Resolve(domain, nsIPs[0], qtype, trace)
		}
	}
	return results, nil
}

//记录trouble：
//是不是results里面有了东西就要返回？这样会不会使部分NS无法保存？
//在查TLD的阶段，answer一直都是空的，一层层的NS一直存在authorities里面。
//直到有一个NS到达权威服务器，权威服务器一次性会返回所有NS到answer，就可以了。
