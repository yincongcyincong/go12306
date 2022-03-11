package action

import (
	"fmt"
	"github.com/yincongcyincong/go12306/module"
	"testing"
	"time"
)

func Test_Login(t *testing.T) {
	qrImage, err := CreateImage()
	if err != nil {
		t.Fatal("还未login")
		return
	}
	qrImage.Image = ""

	err = QrLogin(qrImage)
	if err != nil {
		t.Fatal("还未login")
		return
	}
}

func Test_Hb(t *testing.T) {
	if !CheckLogin() {
		t.Fatal("还未login")
	}

	searchParam := &module.SearchParam{
		TrainDate:       time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		FromStation:     "BJP",
		ToStation:       "TJP",
		FromStationName: "北京",
		ToStationName:   "天津",
		SeatType:        "M",
	}

	trains, err := GetTrainInfo(searchParam)
	if err != nil {
		t.Fatal(err)
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

	err = AfterNateChechFace(trainData, searchParam)
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = AfterNateSuccRate(trainData, searchParam)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = CheckUser()
	if err != nil {
		t.Fatal(err)
		return
	}

	err = AfterNateSubmitOrder(trainData, searchParam)
	if err != nil {
		t.Fatal(err)
		return
	}
	submitToken := &module.SubmitToken{
		Token: "",
	}
	passengers, err := GetPassengers(submitToken)
	if err != nil {
		t.Fatal(err)
		return
	}
	pgs := passengers.Data.NormalPassengers[:1]
	err = PassengerInit()
	if err != nil {
		t.Fatal(err)
		return
	}
	err = AfterNateGetQueueNum()
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = AfterNateConfirmHB(pgs, searchParam, trainData)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func Test_Order(t *testing.T) {
	if !CheckLogin() {
		t.Fatal("还未login")
	}

	searchParam := &module.SearchParam{
		TrainDate:       time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		FromStation:     "BJP",
		ToStation:       "TJP",
		FromStationName: "北京",
		ToStationName:   "天津",
		SeatType:        "M",
	}

	trains, err := GetTrainInfo(searchParam)
	if err != nil {
		t.Fatal(err)
		return
	}

	var trainData *module.TrainData
	for _, train := range trains {
		if train.SeatInfo["一等座"] == "有" {
			trainData = train
			break
		}
	}
	fmt.Println(fmt.Sprintf("%+v", trainData))

	err = CheckUser()
	if err != nil {
		t.Fatal(err)
		return
	}

	err = SubmitOrder(trainData, searchParam)
	if err != nil {
		t.Fatal(err)
		return
	}

	submitToken, err := GetRepeatSubmitToken()
	if err != nil {
		t.Fatal(err)
		return
	}

	passengers, err := GetPassengers(submitToken)
	if err != nil {
		t.Fatal(err)
		return
	}
	pgs := passengers.Data.NormalPassengers[:1]

	err = CheckOrder(pgs, submitToken, searchParam)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = GetQueueCount(submitToken, searchParam)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = ConfirmQueue(pgs, submitToken, searchParam)
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
		return
	}
}
