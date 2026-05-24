package config

import "time"

const (
	ListenAddr = "0.0.0.0:9090"
	Kubeconfig = "C:\\Users\\tang_\\.kube\\config"
	//日志行数
	PodLogTailLine = 2000
	//数据库配置
	DbType = "mysql"
	DbUser = "root"
	DbPwd  = "root123"
	DbHost = "47.98.44.167"
	DbPort = 3306
	DbName = "k8s_workflow"
	// 打印mysql debug sql 日志
	LogMode = true
	// 连接池配置
	MaxIdleConns = 10  //最大空闲连接
	MaxOpenConns = 100 //最大连接数
	MaxLifeTime  = 30 * time.Second  //最大生存时间
)
