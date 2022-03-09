package http

import (
	"fmt"
	"github.com/yincongcyincong/go12306/action"
	"github.com/yincongcyincong/go12306/module"
	"github.com/yincongcyincong/go12306/notice"
	"github.com/yincongcyincong/go12306/utils"
	"github.com/yincongcyincong/go12306/view"
	"net/http"
	"time"
)

func IsLogin(reqFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// 判断是否login
		err := action.GetLoginData()
		if err != nil {
			utils.HTTPFailResp(w, http.StatusInternalServerError, 2, "not login", "")
			return
		}

		reqFunc(w, r)
	}
}

func CreateImageReq(w http.ResponseWriter, r *http.Request) {
	qrImage, err := action.CreateImage()
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
	err = action.QrLogin(qrImage)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	utils.HTTPSuccResp(w, "")
}

func UserLogoutReq(w http.ResponseWriter, r *http.Request) {
	// cookie写入文件
	utils.SaveConf()
	err := action.LoginOut()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	utils.HTTPSuccResp(w, "")
}

func SearchInfo(w http.ResponseWriter, r *http.Request) {

	res := new(module.SearchInfo)
	submitToken, err := action.GetRepeatSubmitToken()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	passengers, err := action.GetPassengers(submitToken)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	res.Passengers = passengers.Data.NormalPassengers
	res.Station = utils.Station
	res.PassengerType = utils.PassengerType
	res.OrderSeatType = utils.OrderSeatType

	utils.HTTPSuccResp(w, res)
}

func SearchTrain(w http.ResponseWriter, r *http.Request) {
	searchParam := new(module.SearchParam)
	err := utils.EncodeParam(r, searchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	res, err := action.GetTrainInfo(searchParam)
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

	err = action.CheckUser()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = action.SubmitOrder(orderParam.TrainData, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	submitToken, err := action.GetRepeatSubmitToken()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	passengers, err := action.GetPassengers(submitToken)
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

	err = action.CheckOrder(orderParam.Passengers, submitToken, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = action.GetQueueCount(submitToken, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = action.ConfirmQueue(orderParam.Passengers, submitToken, orderParam.SearchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	var orderWaitRes *module.OrderWaitRes
	for i := 0; i < 10; i++ {
		orderWaitRes, err = action.OrderWait(submitToken)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		if orderWaitRes.Data.OrderId != "" {
			break
		}
	}

	err = action.OrderResult(submitToken, orderWaitRes.Data.OrderId)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	utils.HTTPSuccResp(w, "购票成功")
}

func ReLogin(w http.ResponseWriter, r *http.Request) {
	err := action.GetLoginData()
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
	err := notice.SendWxrootMessage(utils.WxRobot, "车票购买成功")
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	utils.HTTPSuccResp(w, "")
}

func StartHBReq(w http.ResponseWriter, r *http.Request) {
	searchParam := &module.SearchParam{
		TrainDate:       "2022-03-11",
		FromStation:     "BJP",
		ToStation:       "TJP",
		FromStationName: "北京",
		ToStationName:   "天津",
		SeatType:        "M",
	}

	trains, err := action.GetTrainInfo(searchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	var trainData *module.TrainData
	for _, train := range trains {
		if train.SeatInfo["一等座"] == "无" && train.IsCanNate == "1" {
			trainData = train
			break
		}
	}
	fmt.Println(fmt.Sprintf("%+v", trainData))

	err = action.AfterNateChechFace(trainData, searchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	_, err = action.AfterNateSuccRate(trainData, searchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = action.CheckUser()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = action.AfterNateSubmitOrder(trainData, searchParam)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	submitToken := &module.SubmitToken {
		Token: "",
	}
	passengers, err := action.GetPassengers(submitToken)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	pgs := passengers.Data.NormalPassengers[:1]

	err = action.PassengerInit()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	err = action.AfterNateGetQueueNum()
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	_, err = action.AfterNateConfirmHB(pgs, searchParam, trainData)
	if err != nil {
		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	utils.HTTPSuccResp(w, "候补成功")
}
