package main
import(
	"fmt"
	//"reflect"
)
func  ranger(m*map[string]int){
for key,value:=range *m{
		fmt.Printf("%s怎么样:%d\n",key,value)
	}
}
func  deleter(m*map[string]int){
	delete(*m,"google")
	fmt.Println("删除google：\n")
}
func main()  {
    siteMao:=make(map[string]int)
	siteMao["google"]=1
	siteMao["getup"]=123
	siteMao["xiaoyi"]=211
	siteMao["hao"]=212121
	fmt.Println(siteMao)
	fmt.Println("浩怎么样：\n",siteMao["hao"])
	ranger(&siteMao)
	siteMao["beru"]=65
	name,ok:=siteMao["beru"]
	if(ok){
		fmt.Println("beru怎么样：",name)
	}else{
		fmt.Println(name,"出错",ok)
	}
	deleter(&siteMao)
	ranger(&siteMao)
}





/*package main 
import"fmt"
func main(){
fmt.Println("hello world！")


type A interface{
	call()
}
type B struct{
}
type C struct{

}
func (c C) call(){
	fmt.Println("C.call()....")
}
func (b B) call(){
	fmt.Println("B.call()....")
}var a A
	a=new(B)
	a.call()
	a=new(C)
	a.call()


func vttest(yi interface{}){
fmt.Printf("value:%v\ntepe:%T\n",yi,yi)
}
vttest(123)
	vttest("hello world")
	vttest(true)

func test(input interface{}){
	inputtype:=reflect.TypeOf(input)
	fmt.Println("type:",inputtype)
	inputvalue:=reflect.ValueOf(input)
	fmt.Println("value:",inputvalue)
}
func abc(arg interface{}){
fmt.Println("abc()....	")
fmt.Println("type:",reflect.TypeOf(arg))
fmt.Println("value:",reflect.ValueOf(arg))
}
type User struct{
	name string
	id int
}
	var num float64=1.2345
	abc(num)
	u:=User{"hao",13651}
	abc(u)




	type Man struct{
		name string 
		id int
	}
	type Superman struct{
		Man 

		lever int
	}
	func (this*Man) eat(){
		fmt.Println("Man.eat()....")
	}
	func (this*Superman) fly(){
		fmt.Println("Superman.fly()....")
	} 
	func (this*Man) out(){
		fmt.Println(this.name)
		fmt.Println(this.id)
	}
	func (this*Superman) out(){
		fmt.Println(this.name)
		fmt.Println(this.id)
		fmt.Println(this.lever)
	}
	m:=Man{"putong",136}
	s:=Superman{Man{"chaoren",101},123}
	s.eat()
	s.fly()
	m.eat()
	m.out()
	s.out()




	type Books struct{
	name string
	book_id int
	}
	var hao Books
	hao.name="阿贝尔历险记"
	hao.book_id=59826314
	put(&hao)
	}
func put(book *Books){
	fmt.Println("书名："+book.name)
	book.book_id=12345678
	fmt.Println("书籍ID："+fmt.Sprint(book.book_id))


."GO_workspace/lib1"
	."GO_workspace/lib2"
Testlib1()
   Testlib2()


import(
//	"fmt"
	."GO_workspace/lib1"
	."GO_workspace/lib2"
)

func main()  {
   Testlib1()
   Testlib2()
}







*/


	/*var test string
	fmt.Scanf("%s",&test)
	if(test=="true"){
		fmt.Println(test)
	}
	if(test=="false"){
		fmt.Println(test)
	}*/

	/*switch(test){
	case "true":{
		fmt.Println("this is "+test)
	}
	case "false":{
		fmt.Println("this is "+test)
	}*/


	/*//100素数
	func sushu(){
	a:=2
	for ;a<100;a++{
		if a%2==0{
			fmt.Println(a-1)
		}else{
			continue
		}
	}
	}*/

	/*func test(num1*string,num2*string) (){
	*num1,*num2=*num2,*num1
	fmt.Println(*num1,*num2)
}num1,num2:="googel","runoob"
	fmt.Println(num1,num2)
	test(&num1,&num2)
	fmt.Println(num1,num2)


	const(
		a=iota
		b 
		c =iota+2
		d
		e=iota+3
		f
	)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(e)
	fmt.Println(f)
	

	
	func fool(a string,b int) int{
	fmt.Println("fool")
	fmt.Println("a=",a)
	fmt.Println("b=",b)
	c:=100
	return c
}
func fool1(a string,b int) (int,int){
	fmt.Println("fool1")
	fmt.Println("a=",a)
	fmt.Println("b=",b)
	return 100,200
}
func fool2(a string,b int) (r1 int,r2 int){
	fmt.Println("fool2")
	fmt.Println("a=",a)
	fmt.Println("b=",b)
	r1=771
	r2=77771
	return 
}
c:=fool("hhh",555)
	fmt.Println("c=",c)
	x,y:=fool1("yyy",777)
	fmt.Println(x,y)
	ret1,ret2:=fool2("nnn",2121)
	fmt.Println(ret1,ret2)

	
	
	*/