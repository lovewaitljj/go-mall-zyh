package dao

import (
	"context"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _DbMaster, _DbSlave *gorm.DB

// DB 返回只读实例
func DB() *gorm.DB {
	return _DbSlave
}

// DBMaster 返回主库实例
func DBMaster() *gorm.DB {
	return _DbMaster
}

func init() {
	_DbMaster = initDB(config.Database.Master)
	_DbSlave = initDB(config.Database.Slave)
}

// 支持多种数据库类型
func getDialector(dType, dsn string) gorm.Dialector {
	//switch dType {
	//case "postgres":
	//	return postgres.Open(dsn)
	//default:
	//	return mysql.Open(dsn)
	//}
	return mysql.Open(dsn)
}

func initDB(option *config.DbConnectOption) *gorm.DB {
	logger.Info(context.TODO(), "database config", "db", option)
	db, err := gorm.Open(getDialector(option.Type, option.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(option.MaxOpenConn)
	sqlDb.SetMaxIdleConns(option.MaxIdleConn)
	sqlDb.SetConnMaxLifetime(option.MaxLifeTime)
	if err = sqlDb.Ping(); err != nil {
		panic(err)
	}
	return db
}
