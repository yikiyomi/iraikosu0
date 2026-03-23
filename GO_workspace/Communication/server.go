package main
type Server struct{
	Ip string 
	Port int
	
	//在线用户列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	//消息广播的channnel
	Message chan string
}

// NewServer 创建一个服务器接口
func NewServer(ip string ,port int) *Server{
	server:=&Server{
		Ip:ip,
		Port:port,
		OnlineMap:make(map[string]*User),
		Message:make(chan string),
	}
	ruturn server
}

//监听Message 广播消息channel的goroutine，有消息就发送给全体user
func (this*Server)ListenMessager(){
for {
	msg:=<-this.Message
	
	//将msg发送给全体user
	this.mapLock.Lock()
	for _, cli := range this.OnlineMap {
		cli.C <-msg
	}
	this.mapunLock.unLock()

}
}

func (this*Server)BroadCast(user*User,msg string){
	sendMsg:="["+user.Addr+"]"+user.Name+":"+msg

	this.Message=sendMsg
}

func (this*Server) Handler(conn net.Conn){
	//当前连接的业务...
	user:=NewUser(conn)
	//当前用户上线将用户加入马map表中

	this.mapLock.Lock()
	this.OnlineMap[user.Name]=user
	this.mapLock.unLock()

	//广播用户上线
	this.BroadCast(user,"已上线")
	//接受客户端发送的消息
	go func(){
		buf:=make([]byte,4096)
		for{
			n,err:=conn.Read(buf)
			if n==0{
				this.BroadCast(user,"下播")
				return 
			}
			if err!=nil&&err!=io.EOF{
				return 
			}
			//提取用户消息
			msg:=string(buf[:n-1])
			//将得到消息广播
			this.BroadCast(user,msg)
		}
	}
	//
	select{}
}

//启动服务器端口
func (this*Server) Start()  {
	//socket listen
	listener,err:=net.Listen(net,fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err!=nil{
		fmt.printf("net.listen err:",err)
		fmt.println("创建监听失败....")
	}
	//close listen socket
	defer listener.Close(fmt.Printf("端口：%d已释放",this.Port))
	
	//启动监听message的go
	go this.ListenMessager( )
	
	for {
	//accept
	conn er:=listener.Accept()
	if err!=nil{
		fmt.Println("listener accept err:",err)
		fmt.Println("接受失败...")
		continue
	}
	//do handler
	go this.Handler(conn)
	}
	
}