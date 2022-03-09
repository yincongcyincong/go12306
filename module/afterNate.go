package module

type ChechFace struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		LoginFlag    bool `json:"login_flag"`
		IsShowQrcode bool `json:"is_show_qrcode"`
		FaceFlag     bool `json:"face_flag"`
	} `json:"data"`
}

type SuccRate struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		Flag []*SuccRateFlag `json:"flag"`
	} `json:"data"`
}

type SuccRateFlag struct {
	Level          string `json:"level"`
	SeatTypeCode   string `json:"seat_type_code"`
	TrainNo        string `json:"train_no"`
	StartTrainDate string `json:"start_train_date"`
	Info           string `json:"info"`
}

type AfterNateSubmit struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

type AfterNateQueueNum struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		Flag     bool        `json:"flag"`
		QueueNum []*QueueNum `json:"queueNum"`
	} `json:"data"`
}

type QueueNum struct {
	QueueInfo        string `json:"queue_info"`
	QueueLevel       string `json:"queue_level"`
	SeatTypeCode     string `json:"seat_type_code"`
	TrainNo          string `json:"train_no"`
	TrainDate        string `json:"train_date"`
	StationTrainCode string `json:"station_train_code"`
}

type PassengerInit struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
}

type AfterNatConfirm struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		Flag      bool   `json:"flag"`
		ReserveNo string `json:"reserve_no"`
	} `json:"data"`
}
