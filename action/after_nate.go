package action

import (
	"errors"
	"fmt"
	"github.com/yincongcyincong/go12306/module"
	"github.com/yincongcyincong/go12306/utils"
	"math"
	"math/rand"
	"net/url"
	"strconv"
)

func AfterNateChechFace(trainData *module.TrainData, searchParam *module.SearchParam) error {
	var err error
	data := make(url.Values)
	data.Set("secretList", trainData.SecretStr + "#" + searchParam.SeatType + "|")
	data.Set("_json_att", "")

	chechFaceRes := new(module.ChechFace)
	err = utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/afterNate/chechFace", chechFaceRes, map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc"})
	if err != nil {
		return err
	}

	if !chechFaceRes.Data.FaceFlag {
		return errors.New(fmt.Sprintf("人脸校验失败，请登陆12306进行人脸核验: %+v", chechFaceRes))
	}

	return nil
}

func AfterNateSuccRate(trainData *module.TrainData, searchParam *module.SearchParam) (*module.SuccRate, error) {
	var err error
	data := make(url.Values)
	data.Set("successSecret", trainData.SecretStr + "#" + searchParam.SeatType)
	data.Set("_json_att", "")

	successRate := new(module.SuccRate)
	err = utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/afterNate/getSuccessRate", successRate, map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc"})
	if err != nil {
		return nil, err
	}

	if !successRate.Status {
		return nil, errors.New(fmt.Sprintf("获取候补信息失败: %+v", successRate))
	}

	return successRate, nil
}


func AfterNateSubmitOrder(trainData *module.TrainData, searchParam *module.SearchParam) error {
	var err error
	data := make(url.Values)
	data.Set("secretList", trainData.SecretStr + "#" + searchParam.SeatType)
	data.Set("_json_att", "")

	submitOrder := new(module.AfterNateSubmit)
	err = utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/afterNate/submitOrderRequest", submitOrder, map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc"})
	if err != nil {
		return err
	}

	if !submitOrder.Status && !submitOrder.Data.Flag  {
		return errors.New(fmt.Sprintf("提交订单失败: %+v", submitOrder))
	}

	return nil
}

func PassengerInit() error {
	var err error
	data := make(url.Values)
	queueRes := new(module.PassengerInit)
	err = utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/afterNate/passengerInitApi", queueRes, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return err
	}

	if !queueRes.Status {
		return errors.New(fmt.Sprintf("初始化用户信息失败：%+v", queueRes))
	}

	return nil
}

func AfterNateGetQueueNum() error {
	var err error
	data := make(url.Values)
	queueRes := new(module.AfterNateQueueNum)
	err = utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/afterNate/getQueueNum", queueRes, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return err
	}

	if !queueRes.Status || !queueRes.Data.Flag {
		return errors.New(fmt.Sprintf("候补失败：%+v", queueRes))
	}

	return nil
}

func AfterNateConfirmHB(passengers []*module.Passenger, searchParam *module.SearchParam, trainData *module.TrainData) (*module.AfterNatConfirm, error) {

	passengerInfo := ""
	for _, p := range passengers {
		passengerInfo = passengerInfo + p.PassengerInfo
	}

	data := make(url.Values)
	data.Set("passengerInfo", passengerInfo)
	data.Set("jzParam", "")
	data.Set("hbTrain", fmt.Sprintf("%s,%s#", trainData.TrainName, searchParam.SeatType))
	data.Set("lkParam", "")
	data.Set("sessionId", "")
	data.Set("sig", "")
	data.Set("scene", "")
	data.Set("encryptedData", strconv.Itoa(rand.Intn(math.MaxInt64)))
	data.Set("if_receive_wseat", "Y")
	data.Set("realize_limit_time_diff", "360") // 候补票距离开车前的截止兑换时间，单位: 分钟，默认: 360

	confirmHB := new(module.AfterNatConfirm)
	err := utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/afterNate/confirmHB", confirmHB, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return nil, err
	}

	if !confirmHB.Status || !confirmHB.Data.Flag {
		return nil, errors.New(fmt.Sprintf("候补车票失败：%+v", confirmHB))
	}

	return confirmHB, nil
}


