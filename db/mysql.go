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
	List[define.DSN_MASTER] = createInstance(define.DSN_MASTER)
	List[define.DSN_SLAVE] = createInstance(define.DSN_SLAVE)
	List[define.DSN_LOG1] = createInstance(define.DSN_LOG1)
	List[define.DSN_LOG2] = createInstance(define.DSN_LOG2)
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
			// LogLevel:                  logger.Silent, // Log level
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
	doc := "SELECT EXISTS(SELECT 1 FROM " + table + " WHERE " + field + " = ?) AS found"
	orm.Raw(doc, value).Scan(&o)
	return o.Found
}

func RunMySql() {
	connectMySQLPooling()
}
