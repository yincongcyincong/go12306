package action

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/yincongcyincong/go12306/module"
	"github.com/yincongcyincong/go12306/utils"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func CreateImage() (*module.QrImage, error) {
	_, err := utils.RequestGetWithoutJson(utils.GetCookieStr(), "https://kyfw.12306.cn/otn/login/init", nil)
	if err != nil {
		return nil, err
	}

	data := make(url.Values)
	data.Set("appid", "otn")
	qrImage := new(module.QrImage)
	err = utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/passport/web/create-qr64", qrImage, nil)
	if err != nil {
		return nil, err
	}
	if qrImage.ResultCode != "0" {
		return nil, errors.New(fmt.Sprintf("获取二维码失败: %+v", qrImage))
	}

	image, err := base64.StdEncoding.DecodeString(qrImage.Image)
	if err != nil {
		return nil, err
	}

	err = createQrCode(image)
	if err != nil {
		return nil, err
	}

	return qrImage, nil
}

func QrLogin(qrImage *module.QrImage) error {

	// 扫描二维码
	var err error
	data := make(url.Values)
	data.Set("appid", "otn")
	data.Set("uuid", qrImage.Uuid)
	data.Set("RAIL_DEVICEID", utils.GetCookieVal("RAIL_DEVICEID"))
	data.Set("RAIL_EXPIRATION", utils.GetCookieVal("RAIL_EXPIRATION"))
	qrRes := new(module.QrRes)
	for i := 0; i < 60; i++ {
		err = utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/passport/web/checkqr", qrRes, nil)
		if err == nil && qrRes.ResultCode == "2" {
			break
		} else {
			seelog.Infof("请在'./conf/'查看二维码并用12306扫描登陆，二维码暂未登陆，继续查看二维码状态")
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return err
	}
	utils.AddCookie(map[string]string{"uamtk": qrRes.Uamtk})

	err = GetLoginData()
	if err != nil {
		return err
	}

	return nil
}

func GetLoginData() error {

	// 验证信息，获取tk
	data := make(url.Values)
	data.Set("appid", "otn")

	tk := new(module.TkRes)
	err := utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/passport/web/auth/uamtk", tk, nil)
	if err != nil {
		return err
	}
	if tk.ResultCode != 0 {
		return errors.New(fmt.Sprintf("登陆失败: %+v", tk))
	}
	utils.AddCookie(map[string]string{"tk": tk.Newapptk})

	// 通过tk校验
	data.Set("tk", utils.GetCookieVal("tk"))
	userRes := new(module.UserRes)
	err = utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/uamauthclient", userRes, nil)
	if err != nil {
		return err
	}
	if userRes.ResultCode != 0 {
		seelog.Error(userRes.ResultMessage)
		return errors.New(userRes.ResultMessage)
	}

	// 初始化api
	apiRes := new(module.ApiRes)
	err = utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/index/initMy12306Api", apiRes, nil)
	if err != nil {
		return err
	}
	if !apiRes.Status || apiRes.HTTPStatus != 200 {
		return errors.New(fmt.Sprintf("初始化12306 api失败: %+v", apiRes))
	}
	seelog.Infof("%s 登陆成功", apiRes.Data["user_name"])

	// 获取特殊cookie字段
	staticTk := new(module.TkRes)
	err = utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/passport/web/auth/uamtk-static", staticTk, nil)
	if err != nil {
		return err
	}
	utils.AddCookie(map[string]string{"tk": staticTk.Newapptk})

	if !CheckLogin() {
		return errors.New("自动登陆失败！")
	}

	return nil
}

func CheckLogin() bool {
	// 获取查询或者的query url
	data := make(url.Values)
	data.Set("_json_att", "")
	confRes := new(module.InitConfRes)
	err := utils.Request(data.Encode(), utils.GetCookieStr(), "https://kyfw.12306.cn/otn/login/conf", confRes, nil)
	if err != nil {
		return false
	}
	utils.QueryUrl = confRes.Data.QueryUrl

	return confRes.Data.IsLogin == "Y"
}

func LoginOut() error {
	req, err := http.NewRequest("GET", "https://kyfw.12306.cn/otn/login/loginOut", strings.NewReader(""))
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Set("Cookie", utils.GetCookieStr())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	_, err = utils.GetClient().Do(req)
	if err != nil {
		seelog.Error(err)
	}
	return err
}

func createQrCode(captchBody []byte) error {
	_, err := os.Stat("./conf")
	if err != nil {
		seelog.Warn(err)
		err = os.Mkdir("./conf", os.ModePerm)
		if err != nil {
			return err
		}
	}

	imgPath := "./conf/qrcode.png"
	err = ioutil.WriteFile(imgPath, captchBody, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
