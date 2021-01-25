package mysql

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var db *gorm.DB

func Init() {

	dbLink := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=%s",
		os.Getenv("DATABASE_MYSQL_USER"),
		os.Getenv("DATABASE_MYSQL_PASSWORD"),
		os.Getenv("DATABASE_MYSQL_HOST"),
		os.Getenv("DATABASE_MYSQL_DBNAME"),
		os.Getenv("DATABASE_MYSQL_CHARSET"),
		os.Getenv("DATABASE_MYSQL_TIMEZONE"),
	)

	var err error
	db, err = gorm.Open(os.Getenv("DATABASE_DIALECT"), dbLink)
	if err != nil {
		log.Fatalf("Mysql connection error")
		return
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(viper.GetInt("mysql.maxIdle"))
	db.DB().SetMaxOpenConns(viper.GetInt("mysql.maxConn"))
}

func Orm() *gorm.DB {
	return db
}
