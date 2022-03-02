package action

//type LoginInfo struct {
//	Cookie string
//}
//
//func Login() string {
//	client := http.Client{
//		Transport: &http.Transport{
//			DialContext: (&net.Dialer{
//				Timeout:   5 * time.Second,
//				KeepAlive: 5 * time.Second,
//			}).DialContext,
//			TLSHandshakeTimeout:   5 * time.Second,
//			ResponseHeaderTimeout: 5 * time.Second,
//			ExpectContinueTimeout: 1 * time.Second,
//			MaxIdleConnsPerHost:   -1,
//		},
//	}
//
//	initReq, err := http.NewRequest("GET", "https://kyfw.12306.cn/otn/login/init", strings.NewReader(""))
//	if err != nil {
//		log.Panicln(err)
//	}
//	initReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
//	resp, err := client.Do(initReq)
//	if err != nil {
//		log.Panicln(err)
//	}
//
//
//	tmpCookie := "JSESSIONID=5999F3C48D622554BBF930476136B63C; BIGipServerpool_index=770703882.43286.0000; route=6f50b51faa11b987e576cdb301e545c4; BIGipServerotn=653263370.64545.0000; guidesStatus=off; highContrastMode=defaltMode; cursorStatus=off; RAIL_EXPIRATION=1644399124143; RAIL_DEVICEID=HxVKRYFybjjce3j_YFUSj3YCSikCtGnQMTRB_ivkGogYJI_Zub0z5XAjSE6mis4hAeHrm0b9WIr8rpCwIpTP3wfFa2PUE67-RmNB25iPrKTA1XFxiQk4PywZh0czQHuGNifLJpeXTUzMDC7fRpMy5qH0kWuIktLB"
//
//	captchReq, err := http.NewRequest("GET", "https://kyfw.12306.cn/passport/captcha/captcha-image?login_site=E&module=login&rand=sjrand&0.31745546375395106", strings.NewReader(""))
//	if err != nil {
//		log.Panicln(err)
//		return ""
//	}
//	captchReq.Header.Set("Cookie", tmpCookie)
//	resp, err := client.Do(captchReq)
//	if err != nil {
//		log.Panicln(err)
//		return ""
//	}
//
//	captchBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Panicln(err)
//		return ""
//	}
//	createCaptcha(captchBody)
//
//	return tmpCookie
//
//}
//
//func createCaptcha(captchBody []byte) {
//	imgPath := "./image/tkcode.png"
//	err := ioutil.WriteFile(imgPath, captchBody, fs.ModePerm)
//	if err != nil {
//		log.Panicln(err)
//		return
//	}
//}
//
//func CaptchaCheck(tmpCookie string, captcha string) {
//	captchaResList := strings.Split(captcha, ",")
//	randCode := ""
//	for _, captchaRes := range captchaResList {
//		switch captchaRes {
//		case "1":
//			randCode += ",40,77"
//		case "2":
//			randCode += ",112,77"
//		case "3":
//			randCode += ",184,77"
//		case "4":
//			randCode += ",256,77"
//		case "5":
//			randCode += ",40,149"
//		case "6":
//			randCode += ",112,149"
//		case "7":
//			randCode += ",184,149"
//		case "8":
//			randCode += ",256,149"
//		}
//	}
//
//	client := http.Client{
//		Transport: &http.Transport{
//			DialContext: (&net.Dialer{
//				Timeout:   5 * time.Second,
//				KeepAlive: 5 * time.Second,
//			}).DialContext,
//			TLSHandshakeTimeout:   5 * time.Second,
//			ResponseHeaderTimeout: 5 * time.Second,
//			ExpectContinueTimeout: 1 * time.Second,
//			MaxIdleConnsPerHost:   -1,
//		},
//	}
//	randCode = strings.TrimLeft(randCode, ",")
//	fmt.Println(randCode)
//	microTime := time.Now().UnixNano() / 1000000
//	microTimeStr := strconv.FormatInt(microTime, 10)
//
//	var clusterinfo = url.Values{}
//	//var clusterinfo = map[string]string{}
//	clusterinfo.Add("answer", randCode)
//	clusterinfo.Add("rand", "sjrand")
//	clusterinfo.Add("login_site", "E")
//	data := clusterinfo.Encode()
//
//	checkReq, err := http.NewRequest("POST", "https://kyfw.12306.cn/passport/captcha/captcha-check?answer="+randCode+"&rand=sjrand&login_site=E&_="+microTimeStr, strings.NewReader(data))
//	if err != nil {
//		log.Panicln(err)
//		return
//	}
//
//	checkReq.Header.Set("Cookie", tmpCookie)
//	checkReq.Header.Set("Referer", "https://kyfw.12306.cn/otn/resources/login.html")
//	checkReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
//	checkReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
//	checkReq.Header.Set("Origin", "https://kyfw.12306.cn")
//
//	resp, err := client.Do(checkReq)
//	if err != nil {
//		log.Panicln(err)
//		return
//	}
//
//	checkBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Panicln(err)
//		return
//	}
//	fmt.Println(string(checkBody))
//}


//func BuyProcess(w http.ResponseWriter, r *http.Request) {
//	qrImage, err := CreateImage()
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//	qrImage.Image = ""
//
//	err = QrLogin(qrImage)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	//submitToken, err := GetRepeatSubmitToken()
//	//if err != nil {
//	//	utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//	//	return
//	//}
//	//
//	//passengers, err := GetPassengers(submitToken)
//	//if err != nil {
//	//	utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//	//	return
//	//}
//
//	searchParam := &module.SearchParam{
//		TrainDate:       "2022-02-17",
//		FromStation:     "BJP",
//		ToStation:       "TJP",
//		FromStationName: "北京",
//		ToStationName:   "天津",
//		SeatType:        "O",
//	}
//	res, err := GetTrainInfo(searchParam)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//	fmt.Println(fmt.Sprintf("%+v", res[2]))
//
//	orderParam := &module.OrderParam{
//		TrainData:   res[2],
//		SearchParam: searchParam,
//		PassengerMap: map[string]bool{
//			"尹聪": true,
//		},
//	}
//	d, _ := json.Marshal(orderParam)
//	fmt.Println(string(d))
//
//	err = CheckUser()
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	err = SubmitOrder(orderParam.TrainData, orderParam.SearchParam)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	submitToken, err := GetRepeatSubmitToken()
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	passengers, err := GetPassengers(submitToken)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//	orderParam.Passengers = passengers.Data.NormalPassengers[:1]
//
//	err = CheckOrder(orderParam.Passengers, submitToken, orderParam.SearchParam)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	err = GetQueueCount(submitToken, orderParam.SearchParam)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	err = ConfirmQueue(orderParam.Passengers, submitToken, orderParam.SearchParam)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	var orderWaitRes *module.OrderWaitRes
//	for i := 0; i < 20; i++ {
//		orderWaitRes, err = OrderWait(submitToken)
//		if err != nil {
//			time.Sleep(3 * time.Second)
//			continue
//		}
//		if orderWaitRes.Data.OrderId != "" {
//			break
//		}
//	}
//
//	err = OrderResult(submitToken, orderWaitRes.Data.OrderId)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	utils.HTTPSuccResp(w, "")
//}


//func LoginProcess(w http.ResponseWriter, r *http.Request) {
//	qrImage, err := CreateImage()
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//	qrImage.Image = ""
//
//	err = QrLogin(qrImage)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	submitToken, err := GetRepeatSubmitToken()
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//
//	passengers, err := GetPassengers(submitToken)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//	fmt.Println(passengers)
//
//	searchParam := &module.SearchParam{
//		TrainDate:       "2022-02-17",
//		FromStation:     "BJP",
//		ToStation:       "TJP",
//		FromStationName: "北京",
//		ToStationName:   "天津",
//		SeatType:        "O",
//	}
//	res, err := GetTrainInfo(searchParam)
//	if err != nil {
//		utils.HTTPFailResp(w, http.StatusInternalServerError, 1, err.Error(), "")
//		return
//	}
//	fmt.Println(res)
//
//	utils.HTTPSuccResp(w, "")
//}