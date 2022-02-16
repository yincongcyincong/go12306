package module

type AutoBuyRes struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		IfShowPassCode string `json:"ifShowPassCode"`
	} `json:"data"`
}

type BuyRes struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
}

type QueueCountRes struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		Count  string `json:"count"`
		Ticket string `json:"ticket"`
		Op2    string `json:"op_2"`
		CountT string `json:"countT"`
		Op1    string `json:"op_1"`
	} `json:"data"`
}

type CheckOrderRes struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		CanChooseBeds      string `json:"canChooseBeds"`
		CanChooseSeats     string `json:"canChooseSeats"`
		ChooseSeats        string `json:"choose_Seats"`
		IsCanChooseMid     string `json:"isCanChooseMid"`
		IfShowPassCodeTime string `json:"ifShowPassCodeTime"`
		SubmitStatus       bool   `json:"submitStatus"`
		SmokeStr           string `json:"smokeStr"`
	} `json:"data"`
}

type ConfirmQueueRes struct {
	ValidateMessagesShowId string      `json:"validateMessagesShowId"`
	Status                 bool        `json:"status"`
	HTTPStatus             int         `json:"httpstatus"`
	Messages               []string    `json:"messages"`
	Data                   interface{} `json:"data"`
}

type ConfirmData struct {
	IsAsync      string `json:"isAsync"`
	SubmitStatus bool   `json:"submitStatus"`
}

type OrderWaitRes struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		QueryOrderWaitTimeStatus bool   `json:"queryOrderWaitTimeStatus"`
		Count                    int    `json:"count"`
		WaitTime                 int    `json:"waitTime"`
		RequestId                int    `json:"requestId"`
		WaitCount                int    `json:"WaitCount"`
		TourFlag                 string `json:"tourFlag"`
		OrderId                  string `json:"orderId"`
	} `json:"data"`
}

type OrderResultRes struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		SubmitStatus bool `json:"submitStatus"`
	} `json:"data"`
}

type SubmitOrderRes struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   string   `json:"data"`
}
