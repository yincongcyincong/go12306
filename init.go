package main

import (
	"github.com/cihub/seelog"
	"github.com/yincongcyincong/go12306/conf"
	http12306 "github.com/yincongcyincong/go12306/http"
	"github.com/yincongcyincong/go12306/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Init() {
	initLog()
	conf.InitConf(*wxrobot)
	initCookieInfo()

	go utils.InitBlacklist()
	go utils.InitAvailableCDN()
	go initHttp()

	if *runType == "command" {
		go CommandStart()
	}
}

func initLog() {
	logType := `<console/>`
	if *runType == "web" {
		logType = `<file path="log/log.log"/>`
	}

	logger, err := seelog.LoggerFromConfigAsString(`<seelog type="sync" minlevel="info">
    <outputs formatid="main">
        ` + logType + `
    </outputs>
	<formats>
        <format id="main" format="%Date %Time [%LEV] %RelFile:%Line - %Msg%n"></format>
    </formats>
</seelog>`)
	if err != nil {
		log.Panicln(err)
	}
	err = seelog.ReplaceLogger(logger)
	if err != nil {
		log.Panicln(err)
	}
}

func initCookieInfo() {
	// 用户自己设置设置device信息
	err := utils.ReadCookieFromFile()
	if err != nil {
		seelog.Error("read cookie file fail: ", err)
	}

	railExpStr := utils.GetCookieVal("RAIL_EXPIRATION")
	railExp, _ := strconv.Atoi(railExpStr)
	if railExp <= int(time.Now().Unix()*1000) || *mustDevice == "1" {
		seelog.Info("开始重新获取设备信息")
		utils.GetDeviceInfo()
	}

	if utils.GetCookieVal("RAIL_DEVICEID") == "" || utils.GetCookieVal("RAIL_EXPIRATION") == "" {
		panic("获取设备信息失败")
	}

	conf.UserAgent = utils.GetCookieVal("User-Agent")
}

func initHttp() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
			seelog.Flush()
		}
	}()

	http.HandleFunc("/create-image", http12306.CreateImageReq)
	http.HandleFunc("/login", http12306.QrLoginReq)
	http.HandleFunc("/logout", http12306.UserLogoutReq)
	http.HandleFunc("/search-train", http12306.SearchTrain)
	http.HandleFunc("/search-info", http12306.SearchInfo)
	http.HandleFunc("/order-view", http12306.OrderView)
	http.HandleFunc("/order", http12306.IsLogin(http12306.StartOrderReq))
	http.HandleFunc("/re-login", http12306.ReLogin)
	http.HandleFunc("/", http12306.LoginView)
	http.HandleFunc("/send-msg", http12306.SendMsg)
	if err := http.ListenAndServe(":28178", nil); err != nil {
		log.Panicln(err)
	}
}