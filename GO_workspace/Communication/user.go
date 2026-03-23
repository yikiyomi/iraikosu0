package main 
type User struct{
	Name string 
	Addr string
	C chan string 
	conn net.Conn
}
//创建一个用户api
func NewUser(conn net.Conn) *User{
	userAddr:=conn.RemoteAddr().String()

	user:=&User{
		Name:userAddr,
		Addr:userAddr,
		c:make(chan string)
		conn: conn
	}
	go user.ListenMessage()

	ruturn user
}
//监听当前user channel的方法，有消息就发送客户端
func (this*User) ListenMessage(){
	for {
		msg:=<-this.C

		this.conn.Write([]byte(msg+"\n"))
	}
}