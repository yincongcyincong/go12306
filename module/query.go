package module

type TrainRes struct {
	HTTPStatus int    `json:"httpstatus"`
	Message    string `json:"message"`
	Status     bool   `json:"status"`
	Data       struct {
		Result []string          `json:"result"`
		Flag   string            `json:"flag"`
		Map    map[string]string `json:"map"`
	} `json:"data"`
}

type TrainData struct {
	SecretStr        string
	TrainNo          string
	FromStationName  string
	FromStation      string
	ToStationName    string
	ToStation        string
	TrainLocation    string
	StationTrainCode string
	LeftTicket       string
	StartTime        string
	ArrivalTime      string
	DistanceTime     string
	Status           string
	IsCanNate        string
	TrainName        string
	SeatInfo         map[string]string
}

type SearchParam struct {
	TrainDate       string `json:"train_date"`
	FromStation     string `json:"from_station"`
	ToStation       string `json:"to_station"`
	FromStationName string `json:"from_station_name"`
	ToStationName   string `json:"to_station_name"`
	SeatType        string `json:"seat_type"`
}

type PassengerRes struct {
	ValidateMessagesShowId string `json:"validateMessagesShowId"`
	Status                 bool   `json:"status"`
	HTTPStatus             int    `json:"httpstatus"`
	Data                   struct {
		NotifyForGat     string       `json:"notify_for_gat"`
		IsExist          bool         `json:"is_exist"`
		ExMsg            string       `json:"exMsg"`
		TwoIsOpenCLick   []string     `json:"two_isOpenClick"`
		OtherIsOpenClick []string     `json:"other_isOpenClick"`
		NormalPassengers []*Passenger `json:"normal_passengers"`
	} `json:"data"`
	Messages []string `json:"messages"`
}

type Passenger struct {
	PassengerName       string `json:"passenger_name"`
	SexCode             string `json:"sex_code"`
	SexName             string `json:"sex_name"`
	BornDate            string `json:"born_date"`
	CountryCode         string `json:"country_code"`
	PassengerIdTypeCode string `json:"passenger_id_type_code"`
	PassengerIdTypeName string `json:"passenger_id_type_name"`
	PassengerIdNo       string `json:"passenger_id_no"`
	PassengerType       string `json:"passenger_type"`
	PassengerTypeName   string `json:"passenger_type_name"`
	MobileNo            string `json:"mobile_no"`
	PhoneNo             string `json:"phone_no"`
	Email               string `json:"email"`
	Address             string `json:"address"`
	Postalcode          string `json:"postalcode"`
	FirstLetter         string `json:"first_letter"`
	RecordCount         string `json:"record_count"`
	TotalTimes          string `json:"total_times"`
	IndexId             string `json:"index_id"`
	AllEncStr           string `json:"AllEncStr"`
	IsAdult             string `json:"IsAdult"`
	IsYongThan10        string `json:"IsYongThan10"`
	IsYongThan14        string `json:"IsYongThan14"`
	IsOldThan60         string `json:"IsOldThan60"`
	IfReceive           string `json:"if_receive"`
	IsActive            string `json:"is_active"`
	IsBuyTicket         string `json:"is_buy_ticket"`
	LastTime            string `json:"last_time"`
	PassengerUuid       string `json:"passenger_uuid"`
	PassengerTicketStr  string
	OldPassengerStr     string
	PassengerInfo       string
	Alias               string
}

type SubmitToken struct {
	Token             string
	TicketInfo        map[string]interface{}
	OrderRequestParam map[string]interface{}
}

type CheckUserRes struct {
	ValidateMessagesShowId string `json:"validateMessagesShowId"`
	Status                 bool   `json:"status"`
	HTTPStatus             int    `json:"httpstatus"`
	Data                   struct {
		Flag bool `json:"flag"`
	} `json:"data"`
	Messages []string `json:"messages"`
}

type SearchInfo struct {
	PassengerType map[string]string
	OrderSeatType map[string]string
	Station       map[string]string
	Passengers    []*Passenger
}

type OrderParam struct {
	TrainData    *TrainData      `json:"train_data"`
	Passengers   []*Passenger    `json:"-"`
	SearchParam  *SearchParam    `json:"search_param"`
	PassengerMap map[string]bool `json:"passenger_name"`
}
