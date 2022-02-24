package module

type QrImage struct {
	Image         string `json:"image"`
	ResultCode    string `json:"result_code"`
	ResultMessage string `json:"result_message"`
	Uuid          string `json:"uuid"`
}

type QrRes struct {
	ResultMessage string `json:"result_message"`
	ResultCode    string `json:"result_code"`
	Uamtk         string `json:"uamtk"`
}

type TkRes struct {
	ResultMessage string `json:"result_message"`
	ResultCode    int    `json:"result_code"`
	Newapptk      string `json:"newapptk"`
}

type UserRes struct {
	ResultMessage string `json:"result_message"`
	ResultCode    int    `json:"result_code"`
	Apptk         string `json:"apptk"`
	Username      string `json:"username"`
}

type LoginRes struct {
	QrRes   *QrRes
	TkRes   *TkRes
	UserRes *UserRes
}

type ApiRes struct {
	ValidateMessagesShowId string                 `json:"validateMessagesShowId"`
	Status                 bool                   `json:"status"`
	HTTPStatus             int                    `json:"httpstatus"`
	Data                   map[string]interface{} `json:"data"`
	Messages               []string               `json:"messages"`
}

type LoginUser struct {
	QrRes       *QrRes
	TkRes       *TkRes
	UserRes     *UserRes
	ApiRes      *ApiRes
	SubmitToken *SubmitToken
	Passenger   *Passenger
	TrainData   *TrainData
	BuyStatus   int
}

type InitConfRes struct {
	ValidateMessagesShowId string `json:"validateMessagesShowId"`
	Status                 bool   `json:"status"`
	HTTPStatus             int    `json:"httpstatus"`
	Data                   struct {
		IsstudentData     bool     `json:"isstudentDate"`
		IsMessagePassCode string   `json:"is_message_passCode"`
		StudentDate       []string `json:"studentDate"`
		IsUamLogin        string   `json:"is_uam_login"`
		IsLoginPassCode   string   `json:"is_login_passCode"`
		IsSweepLogin      string   `json:"is_sweep_login"`
		QueryUrl          string   `json:"queryUrl"`
		Now               int      `json:"now"`
		IsLogin           string   `json:"is_login"`
		LoginUrl          string   `json:"login_url"`
		StuControl        int      `json:"stu_control"`
		OtherControl      int      `json:"other_control"`
	} `json:"data"`
	Messages []string `json:"messages"`
}

type DeviceInfo struct {
	Exp        string `json:"exp"`
	Dfp        string `json:"dfp"`
	CookieCode string `json:"cookieCode"`
}
