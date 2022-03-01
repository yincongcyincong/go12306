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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type cookieInfo struct {
	cookie map[string]string
	lock   sync.Mutex
}

type parIndex struct {
	par   string
	index int
}

const (
	slitPar    = "parseInt(c/3)+1"
	sha256Par  = "SHA256("
	zaPar      = "za("
	qaPar      = "Qa("
	qa2Par     = "127==="
	raPar      = "Ra("
	reversePar = "for(d=a.length-1;0<=d;d--)c+=a.charAt(d)"
	changePar  = "parseInt(c/2)"
	encodePar  = "length%2"
)

var (
	cookie = &cookieInfo{
		cookie: make(map[string]string),
		lock:   sync.Mutex{},
	}
	AlgIDRe   = regexp.MustCompile("algID(.*?)x26")
	hashAlgRe = regexp.MustCompile(`(?s),hashAlg:function\(a,b,c\)\{(.*?)},(?s)`)
)

func GetDeviceInfo() {

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
	body, err := RequestGetWithoutJson("", "https://kyfw.12306.cn/otn/HttpZF/GetJS", nil)
	if err != nil {
		seelog.Error(err)
		return nil
	}

	matchData := AlgIDRe.FindSubmatch(body)
	if len(matchData) < 2 {
		seelog.Error("get algID fail")
		return nil
	}
	algId := strings.TrimLeft(string(matchData[1]), `\x3d`)
	algId = strings.TrimRight(algId, `\`)

	token := ""
	data := url.Values{}
	token += "adblock0"
	data.Set(getDeviceParam("adblock"), "0")
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
	data.Set(getDeviceParam("plugins"), "1412399caf7126b9506fee481dd0a407")
	width := strconv.Itoa(GetRand(500, 1000))
	token += "scrAvailSize" + width + "x1440"
	data.Set(getDeviceParam("scrAvailSize"), width+"x1440")
	token += "srcScreenSize30xx900x1440"
	data.Set(getDeviceParam("srcScreenSize"), "30xx900x1440")
	token += "storeDbi1l1o1s1"
	data.Set(getDeviceParam("storeDb"), "i1l1o1s1")
	token += "timeZone-8"
	data.Set(getDeviceParam("timeZone"), "-8")
	token += "touchSupport99115dfb07133750ba677d055874de87"
	data.Set(getDeviceParam("touchSupport"), "99115dfb07133750ba677d055874de87")
	webNo := strconv.Itoa(GetRand(5000, 7000))
	token += "userAgentMozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0." + webNo + ".109 Safari/537.36"
	data.Set(getDeviceParam("userAgent"), "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0."+webNo+".109 Safari/537.36")
	token += "webSmartID74a173cc6a9e7335c27eddd372be213a"
	data.Set(getDeviceParam("webSmartID"), "74a173cc6a9e7335c27eddd372be213a")

	data.Set("hashCode", createHashCode(token, string(body)))
	data.Set("timestamp", strconv.Itoa(int(time.Now().Unix()*1000)))
	data.Set("algID", algId)
	return data
}

func getDeviceParam(param string) string {
	if paramRef, ok := conf.LogDeviceMap[param]; ok {
		return paramRef
	}

	return param
}

func createHashCode(token, body string) string {
	body = strings.Replace(body, "\n", "", -1)
	body = strings.Replace(body, "\r", "", -1)
	matchFunc := hashAlgRe.FindSubmatch([]byte(body))
	fmt.Println(string(matchFunc[1]))

	parIdx := findParIndex(string(matchFunc[1]))
	sort.Slice(parIdx, func(i, j int) bool {
		return parIdx[i].index < parIdx[j].index
	})

	for _, par := range parIdx {
		fmt.Println(fmt.Sprintf("%+v, %s", par, token))
		switch par.par {
		case slitPar:
			token = slitToken(token)
		case qaPar, qa2Par:
			token = qa(token)
		case sha256Par, zaPar:
			token = sha256Token(token)
		case reversePar:
			token = reverse(token)
		case changePar:
			token = changeStr(token)
		case encodePar:
			token = encodeToken(token)
		case raPar:
		default:
			seelog.Error("par is not found")
		}
	}

	return token
}

func sha256Token(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	token = base64.StdEncoding.EncodeToString(h.Sum(nil))
	token = strings.Replace(token, "+", "-", -1)
	token = strings.Replace(token, "/", "_", -1)
	token = strings.Replace(token, "=", "", -1)
	return token
}

func findParIndex(body string) []*parIndex {
	pars := []string{slitPar, sha256Par, qaPar, changePar, reversePar, zaPar, raPar, qa2Par, encodePar}
	res := make([]*parIndex, 0)
	for _, par := range pars {
		startIdx := 0

		for {
			idx := strings.Index(body[startIdx:], par)
			if idx == -1 {
				break
			}
			parIdx := &parIndex{
				par:   par,
				index: startIdx + idx,
			}
			startIdx += idx + len(par)
			res = append(res, parIdx)
		}
	}
	return res
}

func encodeToken(token string) string {
	tLen := len(token)
	if tLen%2 == 0 {
		return token[tLen/2:tLen] + token[0:tLen/2]
	}

	return token[tLen/2+1:tLen] + token[tLen/2:tLen/2+1] + token[0:tLen/2]
}

func qa(token string) string {
	tokenRune := []rune(token)
	for i := 0; i < len(tokenRune); i++ {
		if tokenRune[i] != 127 {
			tokenRune[i] = tokenRune[i] + 1
		} else {
			tokenRune[i] = 0
		}
	}

	return string(tokenRune)
}

func slitToken(token string) string {
	tLen := len(token)
	tf := tLen / 3
	if tLen%3 != 0 {
		tf = tLen/3 + 1
	}
	if tLen >= 3 {
		a := token[tf*2 : tLen]
		b := token[0:tf]
		c := token[tf : 2*tf]
		token = a + b + c
	}
	return token
}

func changeStr(token string) string {
	tokenRune := []byte(token)
	for i := 0; i < len(tokenRune)/2; i++ {
		if i%2 == 0 {
			tokenRune[i], tokenRune[len(tokenRune)-1-i] = tokenRune[len(tokenRune)-1-i], tokenRune[i]
		}
	}

	return string(tokenRune)
}

func reverse(token string) string {
	a := []rune(token)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return string(a)
}
