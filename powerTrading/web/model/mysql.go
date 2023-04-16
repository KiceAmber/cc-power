package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Engine *sqlx.DB

func InitDB() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/cc-power?charset=uft8mb4&parseTime=True"
	Engine, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("初始化连接数据库失败,Error: ", err)
		return
	}
	Engine.SetMaxIdleConns(10)
	Engine.SetMaxOpenConns(20)
	return
}
