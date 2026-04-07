package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
type Userinfo struct{
	gorm.Model
	// ID int   `gorm:"autoIncrement:false"`
	Name string
	Gender string
	Hobby string
}

type User struct{

}

func main() {
	dsn :="root:123456@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 创建表自动迁移
	db.AutoMigrate(&Userinfo{})


	// 创建数据行
	 u1:=Userinfo{Name: "阿贝尔", Gender: "男", Hobby: "吉他"}
	// u2:=Userinfo{0,"小伊","女","看书"}
	// u3:=Userinfo{13,"d","男","ajij"}
	db.Create(&u1)
	// db.Create(&u2)
	// db.Delete(&u1)
	// db.Create(&u3)
	// db.Delete(&Userinfo{}, 3)
	// db.Delete(&Userinfo{}, 4)
	// var x Userinfo
	// db.Find(&x,"id = ?", 1)
	// fmt.Printf("%#v",x)
	// db.Delete(&Userinfo{},[]int{13,123})
}

