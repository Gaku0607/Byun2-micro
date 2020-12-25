package repository

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func repositoryConfig() (string, error) {
	return fmt.Sprintf("root:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"dader487", "127.0.0.1:3306", "Byun2micro"), nil
}

func repositoryConnect(conf string) (*gorm.DB, error) {

	db, err := gorm.Open("mysql", conf)
	if err != nil {
		return nil, err
	}

	db.SingularTable(true) //默認s去除
	db.LogMode(true)       //開啟日誌

	//設置連接持參數
	db.DB().SetMaxOpenConns(16)  // 最大閒置數
	db.DB().SetMaxIdleConns(100) // 最大連結束
	db.DB().SetConnMaxLifetime(time.Second * time.Duration(100))

	return db, nil
}
