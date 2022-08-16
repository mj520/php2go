package php2go

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type address struct {
	ip   uint32
	port uint16
}

func (a address) Ip() uint32 {
	return a.ip
}

func (a address) Port() uint16 {
	return a.port
}

func (a address) GetIp() string {
	return Long2ip(a.ip)
}

func (a address) GetPort() int {
	return int(a.port)
}

func (a address) GetAddress() string {
	return a.GetIp() + ":" + strconv.Itoa(a.GetPort())
}

type Server struct {
	address
}

var server Server

//NewServer 注入共享服务的ip 和 port
func NewServer(ip string, port interface{}) Server {
	if ip == "" || ip == "0.0.0.0" {
		ip = GetOutBoundIP() //内网
	}
	server = Server{
		address: address{
			ip:   IP2long(ip),
			port: uint16(GetInterfaceToInt(port)),
		},
	}
	return server
}

// GetServer 共享 必须先NewServer
func GetServer() Server {
	return server
}

type Client struct {
	address
}

// NewClient client
func NewClient(ip string, port int) Client {
	client := Client{
		address: address{
			ip:   IP2long(ip),
			port: uint16(port),
		},
	}
	return client
}

//Connection 24字符串16进制 连接协议
type Connection struct {
	Server
	Client
	id uint32
}

var connectionId uint32
var connectionClients sync.Map
var connectionFormat = []string{"N4", "N2", "N4", "N2", "N4"}
var generateIdMutex sync.Mutex

const minUint uint32 = 1000000000

// GenerateId 连接满 死循环 调用应当放到连接处判断 基本不肯 锁
func generateId() uint32 {
	generateIdMutex.Lock()
	for {
		if connectionId <= minUint {
			connectionId = minUint
		}
		connectionId++
		if _, ok := connectionClients.Load(connectionId); !ok {
			break
		}
	}
	generateIdMutex.Unlock()
	return connectionId
}

func (c *Connection) Id() uint32 {
	return c.id
}

//NewConnection 连接传入
func NewConnection(clientIp string, clientPort int) *Connection {
	client := NewClient(clientIp, clientPort)
	c := &Connection{
		Server: GetServer(),
		Client: client,
		id:     generateId(),
	}
	//connectionClients.Store(c.Id, c)
	return c
}

//Pack 编码
func (c *Connection) Pack() string {
	p := &Protocol{}
	p.Format = connectionFormat
	clientId := p.Pack16(
		int64(c.Server.Ip()),
		int64(c.Server.Port()),
		int64(c.Client.Ip()),
		int64(c.Client.Port()),
		int64(c.id),
	)
	return clientId
}

//UnPack 编码
func (c *Connection) UnPack(id string) {
	p := &Protocol{}
	p.Format = connectionFormat
	s := p.UnPack16(id)
	c.Server.ip = uint32(s[0])
	c.Server.port = uint16(s[1])
	c.Client.ip = uint32(s[2])
	c.Client.port = uint16(s[3])
	c.id = uint32(s[4])
}

// ReleaseId 释放id
func (c *Connection) ReleaseId(id uint32) {
	connectionClients.Delete(id)
}

//GetOutBoundIP 获取出口ip 内网
func GetOutBoundIP() string {
	conn, err := net.Dial("udp", "223.5.5.5:53")
	if err != nil {
		log.Fatal(err)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	return ip
}

//GetInBoundIP 获取入口 外网ip
func GetInBoundIP() string {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	ip := string(body)
	return ip
}
