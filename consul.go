package php2go

import (
	"context"
	"fmt"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/consts"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"os"
	"strconv"
	"strings"
)

type ConsulApi struct {
	client *api.Client
}

// NewConsul 连接至consul服务返回一个ConsulApi对象
func NewConsul(addr string, auth string) (*ConsulApi, error) {
	if auth != "" {
		os.Setenv(api.HTTPAuthEnvName, auth)
	}
	cfg := api.DefaultConfig()
	if addr != "" {
		cfg.Address = addr
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ConsulApi{client: client}, nil
}

func (c *ConsulApi) GetClient() *api.Client {
	return c.client
}

func (c *ConsulApi) GetLock(key string, value string, ttl string) (lock *api.Lock, err error) {
	lock, err = c.client.LockOpts(&api.LockOptions{
		Key:        key,
		Value:      []byte(value),
		SessionTTL: ttl,
	})
	if err != nil {
		log.Println("consul lock error ", err)
		return nil, err
	}
	_, err = lock.Lock(nil)
	if err != nil {
		log.Println("consul lock failed ", err)
		return nil, err
	}
	return lock, nil
}

// ServiceWatch 服务监控
func (c *ConsulApi) ServiceWatch(service string, handle watch.HandlerFunc) {
	params := make(map[string]interface{})
	params["type"] = "service"
	params["service"] = service
	w, err := watch.Parse(params)
	if err != nil {
		log.Fatalln(err)
	}
	w.Handler = handle
	err = w.RunWithClientAndHclog(c.GetClient(), nil)
	if err != nil {
		log.Fatalln("consul ServiceWatch error", err)
	}
}

func (c *ConsulApi) GetRegisterService(tag string, ip string, port int) *api.AgentServiceRegistration {
	id := fmt.Sprintf("%s-%s:%d", tag, ip, port)
	service := &api.AgentServiceRegistration{
		ID:      id,            // 服务唯一ID
		Name:    tag,           // 服务名称
		Tags:    []string{tag}, // 为服务打标签
		Address: ip,
		Port:    port,
	}
	return service
}

// RegisterService 将服务注册到consul 健康检查 服务名称相同的才会返回多个
func (c *ConsulApi) RegisterService(service *api.AgentServiceRegistration, check *api.AgentServiceCheck) error {
	if check != nil {
		service.Check = check
	}
	return c.client.Agent().ServiceRegister(service)
}

// ServiceList 服务列表
func (c *ConsulApi) ServiceList(tag string) (map[string]*api.AgentService, error) {
	return c.client.Agent().ServicesWithFilter(fmt.Sprintf(`"%s" in Tags`, tag))
}

// ServiceDeregister 注销服务
func (c *ConsulApi) ServiceDeregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

// Service 服务发现 服务名称相同的才会返回多个
func (c *ConsulApi) Service(service string) ([]*api.ServiceEntry, error) {
	servers, _, err := c.client.Health().Service(service, "", true, nil)
	if err != nil {
		log.Printf("get Health server fail : %s", err)
		return nil, err
	}
	return servers, nil
}

// WatchKeyToPath watch key or keyprefix  to path
// t 默认 目录监控 path 目录 + key 取文件部分
// t=key path 目录 + key 取文件部分
// t=source 转 t=file path=key
// t=file path是空=key 必须是文件(含路径)
func (c *ConsulApi) WatchKeyToPath(key string, path string, t string) {
	if t == "source" || (t == "file" && path == "") {
		t = "file"
		path = key
	}
	key = strings.TrimLeft(key, "/")    //去除/开始
	path = strings.TrimRight(path, "/") //去除/结束
	params := make(map[string]interface{})
	params["type"] = "keyprefix"
	if t == "key" || t == "file" {
		params["type"] = "key"
		params["key"] = key
	} else {
		t = "keyprefix"
		params["prefix"] = key
	}
	w, err := watch.Parse(params)
	if err != nil {
		log.Fatalln(err)
	}
	w.Handler = func(u uint64, i interface{}) {
		if i == nil {
			return //不存在key
		}
		var kvs = make(api.KVPairs, 0)
		if s, ok := i.(*api.KVPair); ok {
			kvs = append(kvs, s)
		} else {
			kvs = i.(api.KVPairs)
		}
		for _, kv := range kvs {
			if kv.Value == nil {
				continue //目录页返回
			}
			var file string
			last := Strrpos(kv.Key, "/", 0)
			if last > 0 {
				file = Substr(kv.Key, uint(last)+1, -1)
			} else {
				file = kv.Key
			}
			var dir string
			if t == "keyprefix" {
				sub := strings.TrimPrefix(kv.Key, key) //去除前缀
				sub = strings.TrimSuffix(sub, file)    //去除后缀
				sub = strings.Trim(sub, "/")           //去除前后
				//log.Println("consul PathWatch kv sub:", sub, "file:", file)
				if sub != "" && path != "" {
					dir = path + "/" + sub
				} else if sub == "" && path != "" {
					dir = path
				} else if sub != "" && path == "" {
					dir = sub
				}
			} else {
				dir = path
				if t == "file" && path != "" {
					last := Strrpos(path, "/", 0)
					if last > 0 {
						dir = Substr(path, 0, last)
						file = Substr(path, uint(last)+1, -1)
					}
				}
			}
			if dir != "" {
				os.MkdirAll(dir, os.ModePerm) //可能出错
				file = dir + "/" + file
			}
			os.WriteFile(file, kv.Value, os.ModePerm)
		}
	}
	err = w.RunWithClientAndHclog(c.GetClient(), nil)
	if err != nil {
		log.Fatalln("consul WatchKeyToPath error", err)
	}
}

func (c *ConsulApi) WatchPathKey(path string, key string, t string) {
	//todo https://github.com/fsnotify/fsnotify
}

// GetAgentServiceCheck 健康检查 checkPath!=tcp 走http 为空时 默认/health
// CONSUL_FRP=1 特殊使用 http 检查路径为 /{checkPath}/{service.ID} 注意绑定 http.HandleFunc("/health/", Health) 注意斜杠兼容 增加的server.ID
func GetAgentServiceCheck(service *api.AgentServiceRegistration, remotePort int, checkPath string) *api.AgentServiceCheck {
	check := &api.AgentServiceCheck{
		Name:     service.ID,
		Timeout:  "5s",
		Interval: "10s",
	}
	if checkPath == "" {
		checkPath = "/health"
	}
	address := service.Address
	if IsFrp() {
		cfg := GetFrpConfig()
		address = cfg.ServerAddr
		checkPath += "/" + service.ID
	}
	if checkPath != consts.TCPProxy {
		var port string
		if remotePort != 80 {
			port = ":" + strconv.Itoa(remotePort)
		}
		check.HTTP = fmt.Sprintf(`http://%s%s%s`, address, port, checkPath)
	} else {
		check.TCP = fmt.Sprintf(`%s:%d`, address, remotePort)
	}
	return check
}

func IsFrp() bool {
	if frp := os.Getenv("CONSUL_FRP"); frp != "" {
		return true
	}
	return false
}
func GetFrpConfig() config.ClientCommonConf {
	cfg := config.GetDefaultClientConf()
	if addr := os.Getenv("CONSUL_FRP_ADDR"); addr != "" {
		cfg.ServerAddr = addr
	}
	if port := os.Getenv("CONSUL_FRP_PORT"); port != "" {
		cfg.ServerPort, _ = strconv.Atoi(port)
	}
	if token := os.Getenv("CONSUL_FRP_TOKEN"); token != "" {
		cfg.Token = token
	}
	return cfg
}

// SetAgentServiceProxyFrp CONSUL_FRP_ADDR=必须子域名;CONSUL_FRP_PORT=7000;CONSUL_FRP_TOKEN=;CONSUL_FRP=1
func SetAgentServiceProxyFrp(service *api.AgentServiceRegistration, remotePort int, checkPath string) {
	cfg := GetFrpConfig()
	if checkPath == "" {
		checkPath = "/health"
	}
	var pxyCfg config.ProxyConf
	if checkPath != consts.TCPProxy {
		subDomain := Substr(cfg.ServerAddr, 0, Strpos(cfg.ServerAddr, ".", 0))
		pxyCfg = &config.HTTPProxyConf{
			BaseProxyConf: config.BaseProxyConf{
				ProxyName: service.ID,
				ProxyType: consts.HTTPProxy,
				LocalSvrConf: config.LocalSvrConf{
					LocalPort: service.Port,
					LocalIP:   service.Address,
				},
			},
			DomainConf: config.DomainConf{
				SubDomain: subDomain,
			},
			Locations: []string{checkPath + "/" + service.ID},
		}
	} else {
		pxyCfg = &config.TCPProxyConf{
			BaseProxyConf: config.BaseProxyConf{
				ProxyName: service.ID,
				ProxyType: consts.TCPProxy,
				LocalSvrConf: config.LocalSvrConf{
					LocalPort: service.Port,
					LocalIP:   service.Address,
				},
			},
			RemotePort: remotePort,
		}
	}
	visitorCfg := config.DefaultVisitorConf(consts.STCPProxy)
	var pxyCfgs = map[string]config.ProxyConf{service.ID: pxyCfg}
	var visitorCfgs = map[string]config.VisitorConf{service.ID: visitorCfg}
	svc, err := client.NewService(cfg, pxyCfgs, visitorCfgs, "")
	ctx := context.Background()
	err = svc.Run(ctx)
	if err != nil {
		panic(err)
	}
}
