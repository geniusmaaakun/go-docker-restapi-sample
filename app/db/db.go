package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
 	err error
)

/*
* DBの接続設定を実施
*/
func Init() {
	// .envを読み込む
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	// MySQLへの接続情報を定義
	dsn := os.Getenv(("MYSQL_USER")) +":"+os.Getenv(("MYSQL_PASSWORD")) +"@tcp("+ os.Getenv(("MYSQL_HOST")) +":" +os.Getenv(("MYSQL_PORT"))+ ")/"+ os.Getenv(("MYSQL_DATABASE")) +"?charset=utf8mb4&parseTime=True&loc=Local"

	// DBインスタンスを生成
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// グローバル変数に代入する必要あり
	DB = db
}

/*
* DBの接続情報を取得
*/
func GetDB() *gorm.DB {
	return DB
}

/*
* DBを閉じる
*/
func CloseDB() {
	sqlDB, _ := DB.DB()
	if err = sqlDB.Close(); err != nil {
		panic(err)
	}
}