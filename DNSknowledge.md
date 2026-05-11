# DNS + dig + Wireshark 学习笔记

---

## 一、DNS 是什么

DNS（Domain Name System）是互联网的"电话簿"，负责把域名翻译成 IP 地址。

浏览器不认识 `en.wikiversity.org`，只认识 IP 地址（如 `103.102.166.224`）。DNS 的工作就是帮你完成这个翻译。

### 1.1 查询过程（递归查询）

电脑——DNS recursor——由它进行递归：
1. Root nameserver
  .com(which book) - return TLD server
2. TLD server
  which page
3. Domain nameserver
  which row - return IP Address

### 1.2 常见 DNS 记录类型

| A | 域名 → IPv4 地址 | `dyna.wikimedia.org → 103.102.166.224` |
| AAAA | 域名 → IPv6 地址 | `en.wikiversity.org → 2001:df2:e500:ed1a::1` |
| CNAME | 别名 → 真实域名 | `en.wikiversity.org → dyna.wikimedia.org` |
| MX | 邮件服务器 | `gmail.com → smtp.google.com` |
| PTR | IP → 域名（反向查询） | `172.26.192.1 → LAPTOP-xxx.mshome.net` |
| NS | 域名 → 权威 DNS 服务器 | `wikiversity.org → ns1.wikimedia.org` |
| SOA | Zone 的起始授权记录 | 包含主 NS、管理员邮箱、序列号等 |

### 1.3 CNAME 别名链

CNAME 永远指向另一个域名，**不会直接给 IP**。顺着链往下找，直到 A 记录才是终点：

```
en.wikiversity.org   CNAME  →  dyna.wikimedia.org
dyna.wikimedia.org   A      →  103.102.166.224  ← 真实 IP 在这里
```

### 1.4 TTL（生存时间）

每条 DNS 记录都有 TTL（秒），表示这条记录可以被缓存多久。

- TTL 大（如 86400 = 1天）→ 稳定服务，很少变动
- TTL 小（如 10 秒）→ 频繁切换 IP，通常用于负载均衡（如 Wikimedia）

### 1.5 DNS 服务器 ≠ Web 服务器

`114.114.114.114` 是 DNS 服务器，只监听 **UDP 53 端口**，不提供网页。

在浏览器里输入它没有响应，因为它不监听 HTTP（80/443）端口——就像图书馆查询台只告诉你书在哪，不给你送书。

---

## 二、dig 命令

格式： dig [@server可以指定DNS服务器为你解析] [name域名] [type查询什么信息] [+queryoptions查询/返回方式要求]

### 2.1 基本用法

```bash
# 查询 A 记录（默认）
dig google.com
# 指定 DNS 服务器查询
dig  google.com @8.8.8.8

# 查询指定记录类型
dig google.com MX
dig google.com AAAA
dig google.com NS

# 反向查询（IP → 域名）
dig -x 8.8.8.8

# 精简输出，只看结果，比较常用
dig google.com +short

# 追踪完整递归过程（从根服务器开始）
dig google.com +trace
```



### 2.2 读懂 dig 输出

```
; <<>> DiG 9.x <<>> en.wikiversity.org
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 2

;; QUESTION SECTION:        ← 你问了什么
;en.wikiversity.org.    IN  A

;; ANSWER SECTION:          ← 答案
en.wikiversity.org. 3600 IN CNAME dyna.wikimedia.org.
dyna.wikimedia.org.   10 IN A     103.102.166.224

;; Query time: 44 msec      ← DNS 查询耗时
;; SERVER: 114.114.114.114  ← 回答你的 DNS 服务器
```

各 Section 含义：
- **QUESTION**：发出的查询内容
- **ANSWER**：直接回答（A / CNAME 等记录）
- **AUTHORITY**：哪个权威服务器负责这个域
- **ADDITIONAL**：附加信息（权威服务器的 IP 等）


### 2.3 遇到的trouble：
```
1. trace不成功
原因：trace是要直接接触RootDNS，但是我的环境（WSL + 本地 DNS代理）不能自由访问公网 root DNS
解决（不好）：```dig +trace google.com @8.8.8.8```recursor还是参与查询，其优化等特性导致输出丢失信息比较多
2. DNESEC不出现RRSIG
原因：依旧是被recursor隐藏
解决：```dig +trace google.com @8.8.8.8```能返回RRSIG
```

---

## 三、Wireshark 抓包

### 3.1 界面三个区域

```
┌──────────────────────────────────────────┐
│  包列表区（Packet List）                  │  ← 单击选包
│  No. Time  Source  Destination  Info     │
├──────────────────────────────────────────┤
│  包详情区（Packet Details）               │  ← 展开字段
│  ▶ Frame / ▶ Ethernet / ▶ IP / ▶ DNS    │
├──────────────────────────────────────────┤
│  原始字节区（Packet Bytes）               │  ← 点字段会高亮
│  0000  00 74 9c 92 27 42 ...             │
└──────────────────────────────────────────┘
```

**注意**：用**单击**选包（下方详情区展开，支持链接跳转）；**双击**会弹出独立窗口，不支持 [Response In: xxxx] 链接跳转。

### 3.2 抓 DNS 包的步骤

```
1. 清除 DNS 缓存（确保发出真实请求）
   Windows: ipconfig /flushdns

2. 双击 WLAN（或有波形的网卡）开始抓包

3. 过滤栏输入 dns 回车，只显示 DNS 包

4. 在另一个窗口触发 DNS 查询
   nslookup en.wikiversity.org
   或直接打开浏览器访问任意网站

5. 回到 Wireshark，找到对应的包
```

### 3.3 常用过滤语法

```
dns                          只看 DNS 包
udp.port == 53               同上
dns && ip.addr == 114.114.114.114   只看和指定 DNS 服务器的通信
frame.number == 2326         精确定位某一帧
http                         只看 HTTP
ip.addr == 10.16.203.92      只看某 IP 的通信
```

### 3.4 选择哪张网卡

| 网卡名 | 说明 |
|--------|------|
| WLAN | Wi-Fi，有波形就选这个 |
| 以太网 | 有线（插网线时用） |
| Adapter for loopback | 本机自己通信（127.0.0.1） |
| 本地连接* 3/4/5 | 虚拟网卡（WSL/VPN），一般不用 |

---

## 四、DNS 报文结构（RFC 1035）

DNS 报文分 5 个 Section，对应 Wireshark 里的字段：

```
┌────────────────────────────────────┐
│  Header（固定 12 字节）             │
│  Transaction ID / Flags / 计数字段  │
├────────────────────────────────────┤
│  Question Section                  │
│  就是Queries，问的域名 + 类型）      |
├────────────────────────────────────┤
│  Answer Section                    │
│  回答的资源记录（RR）               │
├────────────────────────────────────┤
│  Authority Section                 │
│  权威 NS 服务器信息（常为空）        │
├────────────────────────────────────┤
│  Additional Section                │
│  附加信息（常为空）                 │
└────────────────────────────────────┘
```

### 4.1 Header Flags 字段详解

query 包 `Flags: 0x0100` vs response 包 `Flags: 0x8180`：

| 标志位 | 名称 | query | response | 含义 |
|--------|------|-------|----------|------|
| QR | Query/Response | 0 | 1 | 0=查询，1=响应 |
| Opcode | 操作码 | 0 | 0 | 0=标准查询 |
| AA | Authoritative Answer | 0 | 0/1 | 1=权威服务器直接回答 |
| TC | Truncated | 0 | 0 | 1=响应被截断，需改用 TCP |
| RD | Recursion Desired | 1 | 1 | 1=请求递归查询 |
| RA | Recursion Available | 0 | 1 | 1=服务器支持递归 |
| RCODE | Reply Code | - | 0 | 0=No error，3=NXDOMAIN |

### 4.2 实际抓包对照（ en.wikiversity.org 查询）

**Query 包（Frame 2325）：**
```
Transaction ID: 0x0002
Flags: 0x0100  （QR=0 查询，RD=1 要递归）
Questions: 1
Answer RRs: 0   ← query 包没有答案
Queries:
  en.wikiversity.org  type A  class IN
[Response In: 2326]  ← Wireshark 帮你标注答案在哪帧
```

**Response 包（Frame 2326）：**
```
Transaction ID: 0x0002  ← 和 query 一样，靠它配对
Flags: 0x8180  （QR=1 响应，RA=1 支持递归，RCODE=0 无错误）
Questions: 1
Answer RRs: 2   ← 回答了 2 条记录
Answers:
  en.wikiversity.org  CNAME  dyna.wikimedia.org  TTL=26082s
  dyna.wikimedia.org  A      103.102.166.224     TTL=10s
[Request In: 2325]
[Time: 44.494 ms]  ← 这次 DNS 查询总耗时
```

### 4.3 底部十六进制怎么看

```
0030  00 00 00 00 00 00 02 65  6e 0b 77 69 6b 69 76 65
0040  72 73 69 74 79 03 6f 72  67 00 00 01 00 01
                               ↑
              6e=n 77=w 69=i 6b=k ... 就是 "en.wikiversity.org"
```

点击详情区的任意字段，底部对应字节高亮，可以直观看到每个字段在二进制里占了哪几个字节。

---

## 五、常见问题

**Q: 为什么浏览器输入 114.114.114.114 无响应？**
A: 114.114.114.114 是 DNS 服务器，只监听 UDP 53 端口响应 DNS 查询，不提供网页服务。

**Q: 怎么区分 CNAME（别名）和真实主机名？**
A: 看记录类型。CNAME 的值是另一个域名；A 记录的值是 IP 地址，A 记录才是终点。

**Q: ipconfig /flushdns 之后为什么缓存里还有记录？**
A: 刷新后残留的是本机自身的记录（如 `LAPTOP-xxx.mshome.net`），这是正常的本地解析条目，不影响实验。


---

