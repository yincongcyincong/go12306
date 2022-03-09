package action

import (
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/yincongcyincong/go12306/module"
	"github.com/yincongcyincong/go12306/utils"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

func SubmitOrder(trainData *module.TrainData, searchParam *module.SearchParam) error {
	var err error
	data := make(url.Values)
	data.Set("train_date", searchParam.TrainDate)
	data.Set("back_train_date", time.Now().Format("2006-01-02"))
	data.Set("tour_flag", "dc")
	data.Set("purpose_codes", "ADULT")
	data.Set("query_from_station_name", searchParam.FromStationName)
	data.Set("query_to_station_name", searchParam.ToStationName)
	secretStr, err := url.QueryUnescape(trainData.SecretStr)
	if err != nil {
		return err
	}
	data.Set("secretStr", secretStr)

	submitOrder := new(module.SubmitOrderRes)
	err = utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/leftTicket/submitOrderRequest", submitOrder, map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc"})
	if err != nil {
		return err
	}

	if submitOrder.Data != "0" {
		return errors.New(fmt.Sprintf("提交订单失败: %+v", submitOrder))
	}

	return nil
}

func CheckOrder(passengers []*module.Passenger, submitToken *module.SubmitToken, searchParam *module.SearchParam) error {
	//passengerTicketStr : 座位编号,0,票类型,乘客名,证件类型,证件号,手机号码,保存常用联系人(Y或N)
	//oldPassengersStr: 乘客名,证件类型,证件号,乘客类型
	passengerTicketStr := ""
	oldPassengerStr := ""
	for _, p := range passengers {
		passengerTicketStr = fmt.Sprintf("%s%s,%s_", passengerTicketStr, searchParam.SeatType, p.PassengerTicketStr)
		oldPassengerStr = oldPassengerStr + p.OldPassengerStr
	}

	data := make(url.Values)
	data.Set("bed_level_order_num", "000000000000000000000000000000")
	data.Set("passengerTicketStr", passengerTicketStr)
	data.Set("oldPassengerStr", oldPassengerStr)
	data.Set("tour_flag", "dc")
	data.Set("randCode", "")
	data.Set("sessionId", "")
	data.Set("sig", "")
	data.Set("cancel_flag", "2")
	data.Set("_json_att", "")
	data.Set("whatsSelect", "1")
	data.Set("scene", "nc_login")
	data.Set("REPEAT_SUBMIT_TOKEN", submitToken.Token)

	checkOrderRes := new(module.CheckOrderRes)
	err := utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo", checkOrderRes, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return err
	}

	if !checkOrderRes.Status || !checkOrderRes.Data.SubmitStatus {
		return errors.New(fmt.Sprintf("检查订单失败: %+v", checkOrderRes))
	}

	return nil
}

func GetQueueCount(submitToken *module.SubmitToken, searchParam *module.SearchParam) error {
	var err error
	startTime, err := time.Parse("2006-01-02", searchParam.TrainDate)
	if err != nil {
		return err
	}

	data := make(url.Values)
	data.Set("train_location", submitToken.TicketInfo["train_location"].(string))
	data.Set("purpose_codes", submitToken.TicketInfo["purpose_codes"].(string))
	data.Set("_json_att", "")
	data.Set("leftTicket", submitToken.TicketInfo["leftTicketStr"].(string))
	data.Set("toStationTelecode", submitToken.TicketInfo["queryLeftTicketRequestDTO"].(map[string]interface{})["to_station"].(string))
	data.Set("fromStationTelecode", submitToken.TicketInfo["queryLeftTicketRequestDTO"].(map[string]interface{})["from_station"].(string))
	data.Set("REPEAT_SUBMIT_TOKEN", submitToken.Token)
	data.Set("train_no", submitToken.TicketInfo["queryLeftTicketRequestDTO"].(map[string]interface{})["train_no"].(string))
	data.Set("stationTrainCode", submitToken.TicketInfo["queryLeftTicketRequestDTO"].(map[string]interface{})["station_train_code"].(string))
	data.Set("train_date", fmt.Sprintf("%s 00:00:00 GMT+0800 (中国标准时间)", startTime.Format("Mon Jan 02 2006")))
	data.Set("seatType", searchParam.SeatType)

	queueRes := new(module.QueueCountRes)
	err = utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/confirmPassenger/getQueueCount", queueRes, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return err
	}

	if !queueRes.Status {
		return errors.New("购买失败，排队状态有误")
	}

	ticketNum, _ := strconv.Atoi(queueRes.Data.Ticket)
	if queueRes.Data.Ticket != "充足" && ticketNum <= 0 {
		return errors.New("购买失败，余票不足")
	}

	if queueRes.Data.Op2 == "true" {
		return errors.New("排队人数超过票数")
	}

	return nil
}

func ConfirmQueue(passengers []*module.Passenger, submitToken *module.SubmitToken, searchParam *module.SearchParam) error {
	
	passengerTicketStr := ""
	oldPassengerStr := ""
	for _, p := range passengers {
		passengerTicketStr = fmt.Sprintf("%s%s,%s_", passengerTicketStr, searchParam.SeatType, p.PassengerTicketStr)
		oldPassengerStr = oldPassengerStr + p.OldPassengerStr
	}

	data := make(url.Values)
	data.Set("passengerTicketStr", passengerTicketStr)
	data.Set("oldPassengerStr", oldPassengerStr)
	data.Set("randCode", "")
	data.Set("purpose_codes", submitToken.TicketInfo["purpose_codes"].(string))
	data.Set("key_check_isChange", submitToken.TicketInfo["key_check_isChange"].(string))
	data.Set("leftTicketStr", submitToken.TicketInfo["leftTicketStr"].(string))
	data.Set("train_location", submitToken.TicketInfo["train_location"].(string))
	data.Set("choose_seats", "")
	data.Set("seatDetailType", "000")
	data.Set("is_jy", "N")
	data.Set("is_cy", "Y")
	data.Set("REPEAT_SUBMIT_TOKEN", submitToken.Token)
	data.Set("encryptedData", strconv.Itoa(rand.Intn(math.MaxInt64)))
	data.Set("whatsSelect", "1")
	data.Set("roomType", "00")
	data.Set("dwAll", "N")
	data.Set("_json_att", "")

	confirmQueue := new(module.ConfirmQueueRes)
	err := utils.Request(utils.ReplaceSpecailChar(data.Encode()), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueue", confirmQueue, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return err
	}

	switch data := confirmQueue.Data.(type) {
	case string:
		return errors.New(data)
	case module.ConfirmData:
		if !data.SubmitStatus {
			return errors.New(fmt.Sprintf("确认排队信息失败: %+v", confirmQueue.Data))
		}
	}

	return nil
}

func OrderWait(submitToken *module.SubmitToken) (*module.OrderWaitRes, error) {

	var err error
	orderWaitUrl := fmt.Sprintf("https://kyfw.12306.cn/otn/confirmPassenger/queryOrderWaitTime?random=%s&tourFlag=dc&_json_att=&REPEAT_SUBMIT_TOKEN=%s", strconv.Itoa(rand.Intn(math.MaxInt64)), submitToken.Token)
	orderWaitRes := new(module.OrderWaitRes)
	err = utils.RequestGet(utils.GetCookieStr(), orderWaitUrl, orderWaitRes, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return nil, err
	}

	if orderWaitRes.Data.OrderId != "" {
		return orderWaitRes, nil
	} else {
		switch orderWaitRes.Data.WaitTime {
		case -100:
			seelog.Info("重新获取订单号")
		case -2, -3:
			seelog.Errorf("订单失败获取消")
		default:
			seelog.Infof("等待时间:%d,等待人数：%d", orderWaitRes.Data.WaitTime, orderWaitRes.Data.WaitCount)
		}
		return nil, errors.New("需要继续等待")
	}
}

func OrderResult(submitToken *module.SubmitToken, orderNo string) error {

	var err error
	data := make(url.Values)
	data.Set("orderSequence_no", orderNo)
	data.Set("REPEAT_SUBMIT_TOKEN", submitToken.Token)
	data.Set("json_att", "")

	orderRes := new(module.OrderResultRes)
	err = utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/confirmPassenger/resultOrderForDcQueue", orderRes, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
	if err != nil {
		return err
	}

	if !orderRes.Data.SubmitStatus {
		return errors.New(fmt.Sprintf("获取订单信息失败: %+v", orderRes))
	}

	return nil
}

//func AutoBuy(passenger *module.Passenger, trainData *module.TrainData, submitToken *module.SubmitToken) {
//	//passengerTicketStr : 座位编号,0,票类型,乘客名,证件类型,证件号,手机号码,保存常用联系人(Y或N)
//	//oldPassengersStr: 乘客名,证件类型,证件号,乘客类型
//	var err error
//
//	data := make(url.Values)
//	data.Set("bed_level_order_num", "000000000000000000000000000000")
//	data.Set("passengerTicketStr", "O,"+passenger.PassengerTicketStr)
//	data.Set("oldPassengerStr", passenger.OldPassengerStr)
//	data.Set("tour_flag", "dc")
//	data.Set("cancel_flag", "2")
//	data.Set("purpose_codes", "ADULT")
//	data.Set("REPEAT_SUBMIT_TOKEN", submitToken.Token)
//	data.Set("query_from_station_name", trainData.FromStation)
//	data.Set("query_to_station_name", trainData.ToStation)
//	data.Set("train_date", trainData.StartTime)
//
//	trainData.SecretStr, err = url.QueryUnescape(trainData.SecretStr)
//	if err != nil {
//		log.Panicln(err)
//	}
//	data.Set("secretStr", trainData.SecretStr)
//
//	qrImage := new(module.AutoBuyRes)
//	err = utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/confirmPassenger/autoSubmitOrderRequest", qrImage, map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"})
//	if err != nil {
//		log.Panicln(err)
//	}
//}
