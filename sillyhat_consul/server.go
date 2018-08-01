package sillyhat_consul

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	consulapi "github.com/hashicorp/consul/api"
	//consulconnect "github.com/hashicorp/consul/connect"
	"github.com/gin-gonic/gin"
	"fmt"
	"net"
)

const(
	registration_port = 19527
	recv_buf_len = 1024
	error_ret = -1;
	success_ret = 1;
	timeout = "3s"
	interval = "10s"
	deregister_critical_service_after = "30s"
	//interval = "10s"
	//deregister_critical_service_after = "30s"
)


type ConsulServer struct {
	Router *gin.Engine
	ID      string
	Name    string
	Tags    []string
	LocalIP    string
	Port    int
	ConsulAddress string
	IsHTTPS bool
	IsClient bool
}

func NewConsulServer(Router *gin.Engine,ModuleName string,ModulePort int,ConsulAddress string) *ConsulServer {
	localIP := localIP()
	log.Info("localIP : ",localIP)
	return &ConsulServer{Router:Router,Name:ModuleName,Tags:[]string{ModuleName},LocalIP:localIP,Port:ModulePort,ConsulAddress:ConsulAddress,IsHTTPS:false}
	//return &ConsulServer{Router:Router,ID:localIP + ":" + strconv.Itoa(ModulePort),Name:ModuleName,Tags:[]string{ModuleName},LocalIP:localIP,Port:ModulePort,ConsulAddress:ConsulAddress,IsHTTPS:false}
}

func (c ConsulServer) String() string {
	return fmt.Sprintf("%b", c)
}


func acceptLoop(l net.Listener) {
	log.Info(l)
}


func (c ConsulServer) RegisterServer() {
	log.Info("ConsulServer : ",c.String())
	healthAPI :="/health"
	config := consulapi.DefaultConfig()
	config.Address = c.ConsulAddress

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Panic("consul client error : ", err)
	}
	//if c.IsClient {
	//	go func() {
	//		//svc, _ := consulconnect.NewService(c.Name, client)
	//		//listen, _ := tls.Listen("tcp", ":8080", svc.ServerTLSConfig())
	//		//acceptLoop(listen)
	//
	//		// Create an instance representing this service. "my-service" is the
	//		// name of _this_ service. The service should be cleaned up via Close.
	//		svc, _ := consulconnect.NewService(c.Name, client)
	//		defer svc.Close()
	//		// Creating an HTTP server that serves via Connect
	//		listener, err := tls.Listen("tcp", ":8888", svc.ServerTLSConfig())
	//		if err != nil{
	//			log.Panic(err)
	//		}
	//		defer listener.Close()
	//		go acceptLoop(listener)
	//		// Create an instance representing this service. "my-service" is the
	//		// name of _this_ service. The service should be cleaned up via Close.
	//		//svc, _ := consulconnect.NewService(c.Name, client)
	//		//defer svc.Close()
	//		//// Creating an HTTP server that serves via Connect
	//		//server := &http.Server{
	//		//	Addr:      ":8888",
	//		//	TLSConfig: svc.ServerTLSConfig(),
	//		//	// ... other standard fields
	//		//}
	//		//// Serve!
	//		////server.ListenAndServeTLS("","")
	//		//server.ListenAndServe()
	//	}()
	//}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = c.ID
	registration.Name = c.Name
	//registration.Port = consulServer.Port
	//registration.Port = registration_port
	registration.Port = c.Port
	registration.Tags = c.Tags
	//registration.Address = c.Address
	registration.Address = c.LocalIP
	//registration.ID = module
	//registration.Name = module
	//registration.Port = registration_port
	//registration.Tags = []string{module}
	//registration.Address = "127.0.0.1"
	var http string
	if c.IsHTTPS {
		http = fmt.Sprintf("https://%s:%d%s", registration.Address, c.Port, healthAPI)
	}else{
		http = fmt.Sprintf("http://%s:%d%s", registration.Address, c.Port, healthAPI)
	}
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           http,
		Timeout:                        timeout,
		Interval:                       interval,			//健康检查间隔
		GRPC:							fmt.Sprintf("%v:%v/%v", c.LocalIP, c.Port, c.Name),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
		DeregisterCriticalServiceAfter: deregister_critical_service_after, //注销时间，相当于过期时间,check失败后30秒删除本服务
		//DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		//Interval:                       time.Duration(10) * time.Second,
	}
	//err = client.Agent().ServiceRegister(registration)
	//if err != nil {
	//	log.Panic("register server error : ", err)
	//}
	if err := client.Agent().ServiceRegister(registration); err != nil {
		log.Panic("register server error : ", err)
	}

	c.Router.GET(healthAPI, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Description":"Golang project client status",
			"Status":"UP",
		})
	})

}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}