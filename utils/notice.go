package utils

import (
	"encoding/json"
	"errors"
)

type WxRobotMsg struct {
	MsgType string `json:"msgtype"`
	Text struct{
		Content string `json:"content"`
	} `json:"text"`
}

type WxRobotRes struct {
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
}

func SendWxrootMessage(robot string, message string) error {
	msg := new(WxRobotMsg)
	msg.Text.Content = message
	msg.MsgType = "text"

	msgBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	res := new(WxRobotRes)
	err = Request(string(msgBody),"", robot, res, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}

	if res.Errcode != 0 {
		return errors.New(res.Errmsg)
	}

	return nil
}
