package db

import (
	"ddaom/define"
	"log"
	"os"
	_ "os"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var List = make(map[string]*gorm.DB)

func connectMySQLPooling() {
	List[define.Mconn.DsnMaster] = createInstance(define.Mconn.DsnMaster)
	List[define.Mconn.DsnSlave] = createInstance(define.Mconn.DsnSlave)
	List[define.Mconn.DsnLog1Master] = createInstance(define.Mconn.DsnLog1Master)
	List[define.Mconn.DsnLog1Slave] = createInstance(define.Mconn.DsnLog1Slave)
	List[define.Mconn.DsnLog2Master] = createInstance(define.Mconn.DsnLog2Master)
	List[define.Mconn.DsnLog2Slave] = createInstance(define.Mconn.DsnLog2Slave)
}

func createInstance(dsn string) *gorm.DB {
	sqlDB, err := sql.Open("mysql", dsn)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute)
	if err != nil {
		panic(err)
	}
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			// LogLevel:      logger.Info, // Log level
			// LogLevel: logger.Silent, // Log level
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func ExistRow(orm *gorm.DB, table string, field string, value string) bool {
	var o struct {
		Found bool
	}
	doc := "SELECT EXISTS(SELECT 1 FROM " + table + " WHERE " + field + " = ? AND deleted_yn = false) AS found"
	orm.Raw(doc, value).Scan(&o)
	return o.Found
}

func RunMySql() {
	connectMySQLPooling()
}
