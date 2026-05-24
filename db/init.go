package db

import (
	"dashboard-api/config"
	"dashboard-api/model"
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
)



var (
	isInit 		bool
	GORM		*gorm.DB
	err			error
)


func Init()  {
	if isInit {
		return
	}

	//判断数据库是否存在, 不存在就创建
	db, errs := sql.Open(config.DbType, fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
		config.DbUser, config.DbPwd, config.DbHost, config.DbPort))

	if errs != nil {
		panic("数据库连接失败：" + errs.Error())
	}

	defer db.Close()

	//(创建数据库)
	if _, errs := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET UTF8", config.DbName)); errs != nil {
		panic("创建数据库失败：" + errs.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
	config.DbUser,
	config.DbPwd,
	config.DbHost,
	config.DbPort,
	config.DbName)

	GORM, err = gorm.Open(config.DbType, dsn)
	if err != nil {
		panic("数据库连接失败. " + err.Error())
	}

	GORM.LogMode(config.LogMode)

	//迁移数据表(创建表结构)
	GORM.Set("gorm:table_options", "CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ENGINE=InnoDB").AutoMigrate(&model.Workflow{})
	logger.Info("自动迁移数据库表成功")

	//开启连接池
	//连接池最大允许的空闲连接数
	GORM.DB().SetConnMaxIdleTime(config.MaxIdleConns)
	//设置连接可复用的最大连接时间
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(time.Duration(config.MaxLifeTime))

	isInit = true
	logger.Info("连接数据库成功...")
}