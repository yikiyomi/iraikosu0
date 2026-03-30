package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播的channnel
	Message chan string
}

// NewServer 创建一个服务器接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message 广播消息channel的goroutine，有消息就发送给全体user
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将msg发送给全体user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()

	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//当前连接的业务...
	user := NewUser(conn,this)
	//用户上线
	user.Online()
	//接收用户是否活跃
	isLive :=make(chan bool)
	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				return
			}
			//提取用户消息
			msg := string(buf[:n-1])
			//用户针对msg进行消息处理
			user.DoMessage(msg)
			//用户活跃
			isLive <- true
		}
	}()
		//当前handler阻塞
	for{
	select{
		case <-isLive:
			//活跃，重置定时器
		case <- time.After(time.Second *300):
		user.SendMsg("滚出去！")
		//若超市将当前客户端关闭
		//关闭连接
		conn.Close()
		//退出当前handler
		return

	}
		}
}

// 启动服务器端口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Printf("net.listen err: %v\n", err)
		fmt.Println("创建监听失败....")
	}
	//close listen socket
	defer listener.Close()
	fmt.Printf("监听端口：%d\n", this.Port)

	//启动监听message的go
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			fmt.Println("接受失败...")
			continue
		}
		//do handler
		go this.Handler(conn)
	}

}
