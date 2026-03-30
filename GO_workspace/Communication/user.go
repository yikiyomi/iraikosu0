package main

import (
	"fmt"
	"net"
	"strings"
)
type User struct{
	Name string 
	Addr string
	C chan string 
	conn net.Conn
	server *Server
}
//创建一个用户api
func NewUser(conn net.Conn,server *Server) *User{
	userAddr:=conn.RemoteAddr().String()

	user:=&User{
		Name:userAddr,
		Addr:userAddr,
		C:make(chan string),
		conn: conn,
		server: server,
	}
	go user.ListenMessage()

	return user;
}

//用户上线功能
func (this *User) Online(){
	//当前用户上线将用户加入马map表中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name]=this
	this.server.mapLock.Unlock()
	//发消息提示上线
	this.server.BroadCast(this,"已上线")
}
//用户下线功能
func (this *User) Offline(){
	//当前用户下线将用户从map表中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	//发消息提示下线
	this.server.BroadCast(this,"已下线")
}
//像当前客户端用户发消息
func (this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}


//用户消息处理功能
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线。。\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()

	     }else if len(msg) >7 && msg[:7] == "rename|"{
		//消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]
		//判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok{
			this.SendMsg("当前用户名已被使用")
		}else{
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName]=this
			this.server.mapLock.Unlock()
			
			this.Name=newName
			this.SendMsg("更新成功"+this.Name+"\n")
		}
	}else if len(msg)>4 &&msg[:3] == "to|"{
		//消息格式to|用户|消息内容
		remoteName:=strings.Split(msg,"|")[1]
		if remoteName==""{
			fmt.Printf("消息格式错误请使用、to|用户|消息，格式\n")
			return
		}
		//获取对方用户名
		remoteUser,ok:=this.server.OnlineMap[remoteName]
		if !ok{
			this.SendMsg("当前用户不存在\n")
			return 
		}
		//获取消息内容并发出消息
		content:=strings.Split(msg,"|")[2]
		if content==""{
			fmt.Printf("无消息内容，请重新发送")
			return	
		}
		remoteUser.SendMsg(this.Name+"对你说"+content)
	    } else {
		this.server.BroadCast(this, msg)
	}
}      
//监听当前user channel的方法，有消息就发送客户端
func (this*User) ListenMessage(){
	for {
		msg:=<-this.C
		this.conn.Write([]byte(msg+"\n"))
	}
}