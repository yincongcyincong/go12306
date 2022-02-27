package main

import (
	"flag"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/tools/12306/conf"
	"github.com/tools/12306/module"
	"github.com/tools/12306/notice"
	"github.com/tools/12306/utils"
	"github.com/tools/12306/view"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	runType   = flag.String("run_type", "command", "command: 命令行模式，web：网页模式")
	wxrobot   = flag.String("wxrobot", "", "企业微信机器人通知")
	mustDevice = flag.String("must_device", "0", "强制生成设备信息")
)

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
        <format id="main" format="%UTCDate %UTCTime [%LEV] %RelFile:%Line - %Msg%n"></format>
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

}

func initHttp() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
			seelog.Flush()
		}
	}()

	http.HandleFunc("/create-image", CreateImageReq)
	http.HandleFunc("/login", QrLoginReq)
	http.HandleFunc("/logout", UserLogoutReq)
	http.HandleFunc("/search-train", SearchTrain)
	http.HandleFunc("/search-info", SearchInfo)
	http.HandleFunc("/order-view", OrderView)
	http.HandleFunc("/order", IsLogin(StartOrderReq))
	http.HandleFunc("/re-login", ReLogin)
	http.HandleFunc("/", LoginView)
	http.HandleFunc("/send-msg", SendMsg)
	if err := http.ListenAndServe(":28178", nil); err != nil {
		log.Panicln(err)
	}
}

func main() {

	flag.Parse()

	initLog()
	conf.InitConf()
	initCookieInfo()

	go utils.InitBlacklist()
	go utils.InitAvailableCDN()
	go initHttp()

	if *runType == "command" {
		go CommandStart()
	}

	//fmt.Println(utils.CreateLogDeviceParam().Encode())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	select {
	case <-sigs:
		seelog.Info("用户登出")
		utils.WriteCookieToFile()
		LoginOut()
	}
}

func IsLogin(reqFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// 判断是否login
		err := GetLoginData()
		if err != nil {
			utils.HTTPFailResp(w, http.StatusInternalServerError, 2, "not login", "")
			return
		}

		reqFunc(w, r)
	}
}

func CreateImageReq(w http.ResponseWriter, r *http.Request) {
	qrImage, err := CreateImage()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	utils.HTTPSuccResp(w, qrImage)
}

func QrLoginReq(w http.ResponseWriter, r *http.Request) {
	qrImage := new(module.QrImage)
	err := utils.EncodeParam(r, qrImage)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	err = QrLogin(qrImage)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	utils.HTTPSuccResp(w, "")
}

func UserLogoutReq(w http.ResponseWriter, r *http.Request) {
	// cookie写入文件
	utils.WriteCookieToFile()
	err := LoginOut()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	utils.HTTPSuccResp(w, "")
}

func SearchInfo(w http.ResponseWriter, r *http.Request) {

	res := new(module.SearchInfo)
	submitToken, err := GetRepeatSubmitToken()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	passengers, err := GetPassengers(submitToken)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	res.Passengers = passengers.Data.NormalPassengers
	res.Station = conf.Station
	res.PassengerType = conf.PassengerType
	res.OrderSeatType = conf.OrderSeatType

	utils.HTTPSuccResp(w, res)
}

func SearchTrain(w http.ResponseWriter, r *http.Request) {
	searchParam := new(module.SearchParam)
	err := utils.EncodeParam(r, searchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	res, err := GetTrainInfo(searchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	utils.HTTPSuccResp(w, res)
}

func OrderView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, view.ViewHtml)
}

func StartOrderReq(w http.ResponseWriter, r *http.Request) {
	orderParam := new(module.OrderParam)
	err := utils.EncodeParam(r, orderParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = CheckUser()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = SubmitOrder(orderParam.TrainData, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	submitToken, err := GetRepeatSubmitToken()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	passengers, err := GetPassengers(submitToken)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	orderParam.Passengers = make([]*module.Passenger, 0)
	for _, p := range passengers.Data.NormalPassengers {
		if _, ok := orderParam.PassengerMap[p.PassengerName]; ok {
			orderParam.Passengers = append(orderParam.Passengers, p)
		}
	}

	err = CheckOrder(orderParam.Passengers, submitToken, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = GetQueueCount(submitToken, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = ConfirmQueue(orderParam.Passengers, submitToken, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	var orderWaitRes *module.OrderWaitRes
	for i := 0; i < 10; i++ {
		orderWaitRes, err = OrderWait(submitToken)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		if orderWaitRes.Data.OrderId != "" {
			break
		}
	}

	err = OrderResult(submitToken, orderWaitRes.Data.OrderId)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	utils.HTTPSuccResp(w, "购票成功")
}

func ReLogin(w http.ResponseWriter, r *http.Request) {
	err := GetLoginData()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	utils.HTTPSuccResp(w, "重新登陆成功")
}

func LoginView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, view.ViewHtml)
}

func SendMsg(w http.ResponseWriter, r *http.Request) {
	err := notice.SendWxrootMessage(*wxrobot, "车票购买成功")
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	utils.HTTPSuccResp(w, "")
}

