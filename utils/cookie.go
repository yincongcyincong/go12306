package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/tools/12306/conf"
	"github.com/tools/12306/module"
	"io/fs"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type cookieInfo struct {
	cookie map[string]string
	lock   sync.Mutex
}

var (
	cookie = &cookieInfo{
		cookie: make(map[string]string),
		lock:   sync.Mutex{},
	}
	AlgIDRe = regexp.MustCompile("algID(.*?)x26")
)

func GetDeviceInfo() {

	// 动态获取设备信息
	//body, err := RequestGetWithoutJson("", "https://kyfw.12306.cn/otn/HttpZF/GetJS", nil, nil)
	//if err != nil {
	//	seelog.Error(err)
	//	return
	//}
	//
	//matchData := AlgIDRe.FindSubmatch(body)
	//if len(matchData) < 2 {
	//	seelog.Error("get algID fail")
	//	return
	//}
	//algId := strings.TrimLeft(string(matchData[1]), `\x3d`)
	//algId = strings.TrimRight(algId, `\`)
	//
	//data := url.Values{}
	//data.Set("adblock", "0")
	//data.Set("cookieEnabled", "1")
	//data.Set("custID", "133")
	//data.Set("doNotTrack", "unknown")
	//data.Set("flashVersion", "0")
	//data.Set("javaEnabled", "0")
	//data.Set("jsFonts", "c227b88b01f5c513710d4b9f16a5ce52")
	//data.Set("localCode", "3232236206")
	//data.Set("mimeTypes", "52d67b2a5aa5e031084733d5006cc664")
	//data.Set("os", "MacIntel")
	//data.Set("platform", "WEB")
	//data.Set("plugins", "d22ca0b81584fbea62237b14bd04c866")
	//data.Set("scrAvailSize", strconv.Itoa(rand.Intn(1000))+"x1920")
	//data.Set("srcScreenSize", "24xx1080x1920")
	//data.Set("storeDb", "i1l1o1s1")
	//data.Set("timeZone", "-8")
	//data.Set("touchSupport", "99115dfb07133750ba677d055874de87")
	//data.Set("userAgent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
	//data.Set("webSmartID", "f4e3b7b14cc647e30a6267028ad54c56")
	//data.Set("timestamp", strconv.Itoa(int(time.Now().Unix()*1000)))
	//data.Set("algID", algId)

	body, err := RequestGetWithoutJson("", "https://kyfw.12306.cn/otn/HttpZF/logdevice?algID=kPKMh0v14R&hashCode=QVOG5ISqfjYPBVzD1nZK3tWxd-vJH_lBCmgRi1DJYsU&FMQw=0&q4f3=zh-CN&VPIf=1&custID=133&VEek=unknown&dzuS=0&yD16=0&EOQP=c227b88b01f5c513710d4b9f16a5ce52&jp76=fe9c964a38174deb6891b6523b8e4518&hAqN=MacIntel&platform=WEB&ks0Q=1412399caf7126b9506fee481dd0a407&TeRS=794x1440&tOHY=30xx900x1440&Fvje=i1l1o1s1&q5aJ=-8&wNLf=99115dfb07133750ba677d055874de87&0aew=Mozilla/5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/98.0.4758.102%20Safari/537.36&E3gR=6830b7871c4d4d53e2c64935d267dda8&timestamp="+strconv.Itoa(int(time.Now().Unix()*1000)), nil)
	if err != nil {
		seelog.Error(err)
		return
	}
	if bytes.Contains(body, []byte("callbackFunction")) {
		body = bytes.TrimLeft(body, "callbackFunction('")
		body = bytes.TrimRight(body, "')")
		deviceInfo := new(module.DeviceInfo)
		err = json.Unmarshal(body, deviceInfo)
		if err != nil {
			seelog.Error(err)
			return
		}
		if deviceInfo.CookieCode == "" {
			seelog.Error("生成device信息失败, 请通过启动参数手动设置device信息")
			return
		}

		cookie.cookie["RAIL_DEVICEID"] = deviceInfo.Dfp
		cookie.cookie["RAIL_EXPIRATION"] = deviceInfo.Exp
	}

}

func AddCookie(kv map[string]string) {
	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	for k, v := range kv {
		cookie.cookie[k] = v
	}
}

func GetCookieVal(key string) string {
	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	return cookie.cookie[key]
}

func AddCookieStr(setCookies []string) {

	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	for _, setCookie := range setCookies {
		cookieKVs := strings.Split(setCookie, ";")
		for _, cookieKV := range cookieKVs {
			cookieKV = strings.TrimSpace(cookieKV)
			cookieSlice := strings.SplitN(cookieKV, "=", 2)
			if len(cookieSlice) >= 2 {
				cookie.cookie[cookieSlice[0]] = cookieSlice[1]
			}
		}
	}
}

func GetCookieStr() string {
	res := ""
	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	for k, v := range cookie.cookie {
		res = fmt.Sprintf("%s%s=%s; ", res, k, v)
	}
	return res
}

func WriteCookieToFile() {
	cookieStr := GetCookieStr()
	cookiePath := "./conf/cookie"
	err := ioutil.WriteFile(cookiePath, []byte(cookieStr), fs.ModePerm)
	if err != nil {
		seelog.Error(err)
	}
}

func ReadCookieFromFile() error {

	cookiePath := "./conf/cookie"
	cookieByte, err := ioutil.ReadFile(cookiePath)
	if err != nil {
		seelog.Error(err)
		return err
	}

	AddCookieStr([]string{string(cookieByte)})
	return nil
}

func CreateLogDeviceParam() url.Values {
	//body, err := RequestGetWithoutJson("", "https://kyfw.12306.cn/otn/HttpZF/GetJS", nil)
	//if err != nil {
	//	seelog.Error(err)
	//	return nil
	//}
	//
	//matchData := AlgIDRe.FindSubmatch(body)
	//if len(matchData) < 2 {
	//	seelog.Error("get algID fail")
	//	return nil
	//}
	//algId := strings.TrimLeft(string(matchData[1]), `\x3d`)
	//algId = strings.TrimRight(algId, `\`)

	token := ""
	data := url.Values{}
	token += "browserLanguagezh-CN"
	data.Set(getDeviceParam("browserLanguage"), "zh-CN")
	token += "cookieEnabled1"
	data.Set(getDeviceParam("cookieEnabled"), "1")
	token += "custID133"
	data.Set(getDeviceParam("custID"), "133")
	token += "doNotTrackunknown"
	data.Set(getDeviceParam("doNotTrack"), "unknown")
	token += "flashVersion0"
	data.Set(getDeviceParam("flashVersion"), "0")
	token += "javaEnabled0"
	data.Set(getDeviceParam("javaEnabled"), "0")
	token += "jsFontsc227b88b01f5c513710d4b9f16a5ce52"
	data.Set(getDeviceParam("jsFonts"), "c227b88b01f5c513710d4b9f16a5ce52")
	token += "mimeTypesfe9c964a38174deb6891b6523b8e4518"
	data.Set(getDeviceParam("mimeTypes"), "fe9c964a38174deb6891b6523b8e4518")
	token += "osMacIntel"
	data.Set(getDeviceParam("os"), "MacIntel")
	token += "platformWEB"
	data.Set(getDeviceParam("platform"), "WEB")
	token += "plugins1412399caf7126b9506fee481dd0a407"
	data.Set(getDeviceParam("plugins"), "d22ca0b81584fbea62237b14bd04c866")

	// todo 需要改为变量
	token += "scrAvailSize794x1440"
	data.Set(getDeviceParam("scrAvailSize"), "794x1440")
	token += "srcScreenSize30xx900x1440"
	data.Set(getDeviceParam("srcScreenSize"), "30xx900x1440")
	token += "storeDbi1l1o1s1"
	data.Set(getDeviceParam("storeDb"), "i1l1o1s1")
	token += "timeZone-8"
	data.Set(getDeviceParam("timeZone"), "-8")
	token += "touchSupport99115dfb07133750ba677d055874de87"
	data.Set(getDeviceParam("touchSupport"), "99115dfb07133750ba677d055874de87")

	// todo 需要改为变量
	token += "userAgentMozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36"
	data.Set(getDeviceParam("userAgent"), "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	token += "webSmartID6830b7871c4d4d53e2c64935d267dda8"
	data.Set(getDeviceParam("webSmartID"), "6830b7871c4d4d53e2c64935d267dda8")

	data.Set("hashCode", createHashCode(token))
	data.Set("timestamp", strconv.Itoa(int(time.Now().Unix()*1000)))
	//data.Set("algID", algId)
	return data
}

func getDeviceParam(param string) string {
	if paramRef, ok := conf.LogDeviceMap[param]; ok {
		return paramRef
	}

	return param
}

func createHashCode(token string) string {
	tLen := len(token)
	tf := tLen / 3
	if tLen%3 != 0 {
		tf = tLen/3 + 1
	}
	if tLen >= 3 {
		token = token[tf:tf*2] + token[tf*2:tLen] + token[0:tf]
	}

	// 在处理一次
	tLen = len(token)
	tf = tLen / 3
	if tLen%3 != 0 {
		tf = tLen/3 + 1
	}
	if tLen >= 3 {
		token = token[tf*2:tLen] + token[0:tf] + token[tf:2*tf]
	}

	token = encodeToken(token)
	token = encodeToken(token)
	token = encodeToken(token)


	h := sha256.New()
	h.Write([]byte(token))
	base64Token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	base64Token = strings.Replace(base64Token, "+", "-", -1)
	base64Token = strings.Replace(base64Token, "/", "_", -1)
	base64Token = strings.Replace(base64Token, "=", "", -1)
	fmt.Println(base64Token)
	return base64Token
}

func encodeToken(token string)string {
	tLen := len(token)
	if tLen%2 == 0 {
		return token[tLen/2:tLen] + token[0:tLen/2]
	}

	return token[tLen/2+1:tLen] + token[tLen/2:tLen/2+1] +  token[0:tLen/2]
}
