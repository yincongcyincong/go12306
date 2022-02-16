package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/tools/12306/conf"
	"github.com/tools/12306/module"
	"github.com/tools/12306/utils"
	"net/url"
	"regexp"
	"strings"
)

var (
	TokenRe           = regexp.MustCompile("var globalRepeatSubmitToken = '(.+)';")
	TicketInfoRe      = regexp.MustCompile("var ticketInfoForPassengerForm=(.+);")
	OrderRequestParam = regexp.MustCompile("var orderRequestDTO=(.+);")
)

func GetTrainInfo(searchParam *module.SearchParam) ([]*module.TrainData, error) {

	var err error
	searchRes := new(module.TrainRes)
	err = utils.RequestGetWithCDN(utils.GetCookieStr(), fmt.Sprintf("https://kyfw.12306.cn/otn/%s?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=ADULT",
		conf.QueryUrl, searchParam.TrainDate, searchParam.FromStation, searchParam.ToStation), searchRes, nil, utils.GetCdn())
	if err != nil {
		seelog.Error(err)
		return nil, err
	}

	if searchRes.HTTPStatus != 200 && searchRes.Status {
		seelog.Errorf("get train info fail: %+v", searchRes)
		return nil, errors.New("get train info fail")
	}

	searchDatas := make([]*module.TrainData, len(searchRes.Data.Result))
	for i, res := range searchRes.Data.Result {
		resSlice := strings.Split(res, "|")
		sd := new(module.TrainData)
		sd.Status = resSlice[1]
		sd.TrainNo = resSlice[3]
		sd.FromStationName = searchRes.Data.Map[resSlice[6]]
		sd.ToStationName = searchRes.Data.Map[resSlice[7]]
		sd.FromStation = resSlice[6]
		sd.ToStation = resSlice[7]

		if resSlice[1] == "预订" {
			sd.SecretStr = resSlice[0]
			sd.LeftTicket = resSlice[29]
			sd.StartTime = resSlice[8]
			sd.ArrivalTime = resSlice[9]
			sd.DistanceTime = resSlice[10]

			sd.SeatInfo = make(map[string]string)
			sd.SeatInfo["特等座"] = resSlice[conf.SeatType["特等座"]]
			sd.SeatInfo["商务座"] = resSlice[conf.SeatType["商务座"]]
			sd.SeatInfo["一等座"] = resSlice[conf.SeatType["一等座"]]
			sd.SeatInfo["二等座"] = resSlice[conf.SeatType["二等座"]]
			sd.SeatInfo["软卧"] = resSlice[conf.SeatType["软卧"]]
			sd.SeatInfo["硬卧"] = resSlice[conf.SeatType["硬卧"]]
			sd.SeatInfo["硬座"] = resSlice[conf.SeatType["硬座"]]
			sd.SeatInfo["无座"] = resSlice[conf.SeatType["无座"]]
			sd.SeatInfo["动卧"] = resSlice[conf.SeatType["动卧"]]
			sd.SeatInfo["软座"] = resSlice[conf.SeatType["软座"]]
		}

		searchDatas[i] = sd
	}
	return searchDatas, nil
}

func GetRepeatSubmitToken() (*module.SubmitToken, error) {

	body, err := utils.RequestGetWithoutJson(utils.GetCookieStr(), "https://kyfw.12306.cn/otn/confirmPassenger/initDc", nil)

	matchRes := TokenRe.FindStringSubmatch(string(body))
	submitToken := new(module.SubmitToken)
	if len(matchRes) > 1 {
		submitToken.Token = matchRes[1]
	}

	ticketRes := TicketInfoRe.FindSubmatch(body)
	if len(ticketRes) > 1 {
		ticketRes[1] = bytes.Replace(ticketRes[1], []byte("'"), []byte(`"`), -1)
		err = json.Unmarshal(ticketRes[1], &submitToken.TicketInfo)
		if err != nil {
			seelog.Error(err)
			return nil, err
		}
	}

	orderRes := OrderRequestParam.FindSubmatch(body)
	if len(orderRes) > 1 {
		orderRes[1] = bytes.Replace(orderRes[1], []byte("'"), []byte(`"`), -1)
		err = json.Unmarshal(orderRes[1], &submitToken.OrderRequestParam)
		if err != nil {
			seelog.Error(err)
			return nil, err
		}
	}

	return submitToken, nil
}

func GetPassengers(submitToken *module.SubmitToken) (*module.PassengerRes, error) {

	data := make(url.Values)
	data.Set("_json_att", "")
	data.Set("REPEAT_SUBMIT_TOKEN", submitToken.Token)
	res := new(module.PassengerRes)
	err := utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/confirmPassenger/getPassengerDTOs", res, nil)
	if err != nil {
		seelog.Error(err)
		return nil, err
	}

	if res.Status && res.HTTPStatus != 200 {
		seelog.Error(err)
		return nil, err
	}

	for _, p := range res.Data.NormalPassengers {
		passengerTicketStr := fmt.Sprintf("0,%s,%s,%s,%s,%s,N,%s",
			p.PassengerType, p.PassengerName, p.PassengerIdTypeCode, p.PassengerIdNo, p.MobileNo, p.AllEncStr)
		oldPassengerStr := fmt.Sprintf("%s,%s,%s,%s_",
			p.PassengerName, p.PassengerIdTypeCode, p.PassengerIdNo, p.PassengerType)
		p.PassengerTicketStr = passengerTicketStr
		p.OldPassengerStr = oldPassengerStr
	}

	return res, nil

}

func CheckUser() error {
	data := make(url.Values)
	data.Set("_json_att", "")
	res := new(module.CheckUserRes)
	err := utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/login/checkUser", res, nil)
	if err != nil {
		seelog.Error(err)
		return err
	}

	if res.Status && res.HTTPStatus != 200 {
		seelog.Errorf("checkUser fail:")
		return err
	}

	if !res.Data.Flag {
		seelog.Errorf("check user fail: %+v", res)
		return errors.New("check user fail")
	} else {
		seelog.Info("check user success")
	}
	return nil

}
