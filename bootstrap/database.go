package bootstrap

import (
	"cloud-api-go/pkg/config"
	"cloud-api-go/pkg/database"
	"cloud-api-go/pkg/logger"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

// SetupDB 初始化数据库和 ORM
func SetupDB() {
	var dbConfig gorm.Dialector
	switch config.Get("database.connection") {
	case "mysql":
		// 构建 DSN 信息
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
			config.Get("database.mysql.username"),
			config.Get("database.mysql.password"),
			config.Get("database.mysql.host"),
			config.Get("database.mysql.port"),
			config.Get("database.mysql.database"),
			config.Get("database.mysql.charset"),
		)
		dbConfig = mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         191,   // string 类型字段的默认长度
			SkipInitializeWithVersion: false, // 根据版本自动配置
		})
	case "sqlite":
		// 初始化 sqlite
		db := config.Get("database.sqlite.database")
		dbConfig = sqlite.Open(db)
	default:
		panic(errors.New("database connection not supported"))
	}
	database.Connect(dbConfig, logger.NewGormLogger())
	// 设置最大连接数
	database.SQLDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	// 设置最大空闲连接数
	database.SQLDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))
	// 设置每个连接的过期时间
	database.SQLDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second)
}
