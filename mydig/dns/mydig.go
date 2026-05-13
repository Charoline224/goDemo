package dns

import (
	"net"
)

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
	resp, err := sendQuery(query, server)
	if err != nil {
		return nil, err
	}
	res := parseResp(resp)
	var cname string
	for _, answer := range res.Anwers {
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
	//NS查询
	for _, nsdo := range res.Autorities {
		if nsdo.Type == 2 {
			for _, nsip := range res.Additionals {
				if nsip.Name == nsdo.Data {
					return Resolve(domain, nsip.Data, qtype, trace)
				}
			}
			nsIPs, err := Resolve(nsdo.Data, "198.41.0.4", 1, trace)
			if err != nil {
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
