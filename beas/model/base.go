package model

import (
	"mixed/beas/file"
	"mixed/beas/setting"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var ShopEngine *xorm.Engine

// initial DB engine
func InitDBEngine() {
	if m, err := xorm.NewEngine(setting.DB_DRIVER, setting.DB_DSN); err != nil {
		panic("Init xorm engine err :" + err.Error())
	} else {
		tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "cc_")
		m.SetTableMapper(tbMapper)
		m.SetMaxIdleConns(setting.DB_MAX_IDLE_CONNS)
		m.SetMaxOpenConns(setting.DB_MAX_OPEN_CONNS)
		if logFd, err := file.CreateFile(setting.DB_LOGS_FILE); err != nil {
			panic("Create DB Log Error :" + err.Error())
		} else {
			m.SetLogger(xorm.NewSimpleLogger(logFd))
		}
		// debug mode
		if setting.IsDebug {
			m.ShowSQL(true)
		}
		ShopEngine = m
	}

	// keep connect active
	// go func() {
	// 	var ticker *time.Ticker = time.NewTicker(10 * time.Minute)
	// 	var repeatTimes int
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			if err := ShopEngine.Ping(); err != nil {
	// 				if repeatTimes >= 10 {
	// 					panic("DB is offline" + err.Error())
	// 				}
	// 			} else {
	// 				repeatTimes = 0
	// 			}
	// 		default:
	// 		}
	// 	}
	// }()
}
