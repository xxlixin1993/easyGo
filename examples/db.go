package main

import (
	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/orm/mysql"
	"fmt"
)

type TestModel struct {
	ID    uint32 `gorm:"primary_key:id"`
	UID   uint32 `gorm:"column:uid"`
	IsDel uint8  `gorm:"column:is_del"`
	Info  string `gorm:"column:info"`
	CTime string `gorm:"column:ctime"`
	MTime string `gorm:"column:mtime"`
}

func (t TestModel) TableName() string {
	return "testmysql"
}

func main() {
	easyGo.InitFrame()
	easyGo.InitMysql()
	go testSql()
	easyGo.WaitSignal()
}

func testSql(){
	db, _ := mysql.GetMasterConn("mysql_first")
	//db.Create(&TestModel{
	//	UID:   2,
	//	IsDel: 1,
	//	Info:  "asdsad",
	//	CTime: "2018-09-12 12:12:12",
	//	MTime: "2018-09-12 12:12:12",
	//})

	tm := &TestModel{}
	for i:=0;i<10 ;i++  {
		db.Where(&TestModel{UID:2}).Find(&tm)
		fmt.Println(tm)
	}

}

