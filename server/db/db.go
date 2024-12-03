/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
package db

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"openeuler.org/PilotGo/prometheus-plugin/server/config"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
)

var MySQL *gorm.DB

var Url string

type MysqlManager struct {
	ip       string
	port     int
	userName string
	passWord string
	dbName   string
	db       *gorm.DB
}

func MysqldbInit(conf *config.MysqlDBInfo) error {
	err := ensureDatabase(conf)
	if err != nil {
		return err
	}
	_, err = mysqlInit(
		conf.HostName,
		conf.UserName,
		conf.Password,
		conf.DataBase,
		conf.Port)
	if err != nil {
		return err
	}

	MySQL.AutoMigrate(&model.PrometheusTarget{})
	MySQL.AutoMigrate(&model.Rule{})
	MySQL.AutoMigrate(&model.RuleTarget{})
	MySQL.AutoMigrate(&model.Alert{})

	return nil
}

func mysqlInit(ip, username, password, dbname string, port int) (*MysqlManager, error) {
	m := &MysqlManager{
		ip:       ip,
		port:     port,
		userName: username,
		passWord: password,
		dbName:   dbname,
	}
	Url = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		m.userName,
		m.passWord,
		m.ip,
		m.port,
		m.dbName)

	var err error
	m.db, err = gorm.Open(mysql.Open(Url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	MySQL = m.db

	var db *sql.DB
	if db, err = m.db.DB(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return m, nil
}
func ensureDatabase(conf *config.MysqlDBInfo) error {
	Url := fmt.Sprintf("%s:%s@(%s:%d)/?charset=utf8mb4&parseTime=true",
		conf.UserName,
		conf.Password,
		conf.HostName,
		conf.Port)
	db, err := gorm.Open(mysql.Open(Url))
	if err != nil {
		return err
	}

	creatDataBase := "CREATE DATABASE IF NOT EXISTS " + conf.DataBase + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci"
	db.Exec(creatDataBase)

	d, err := db.DB()
	if err != nil {
		return err
	}
	if err = d.Close(); err != nil {
		return err
	}
	return nil
}
