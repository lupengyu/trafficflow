package constant

type ShipMeta struct {
	MMSI 	int
}

type PositionMeta struct {
	ID 					int 	`json:"ID"`
	MessageType 		int 	`json:"Message_Type"`
	RepeatIndicator 	int 	`json:"Repeat_Indicator"`
	MMSI 				int 	`json:"MMSI"`
	NavigationStatus 	int 	`json:"Navigation_Status"`
	ROT 				int 	`json:"ROT"`
	SOG 				float64 `json:"SOG"`
	PositionAccuracy	int 	`json:"Position_Accuracy"`
	Longitude 			float64 `json:"Longitude"`
	Latitude 			float64 `json:"Latitude"`
	COG 				float64 `json:"COG"`
	HDG 				int 	`json:"HDG"`
	TimeStamp 			int 	`json:"Time_stamp"`
	ReservedForRegional int 	`json:"Reserved_for_regional"`
	RAIMFlag 			int 	`json:"RAIM_flag"`
	Year 				int 	`json:"Year"`
	Month 				int 	`json:"Month"`
	Day 				int 	`json:"Day"`
	Hour				int 	`json:"Hour"`
	Minute 				int		`json:"Minute"`
	Second				int 	`json:"Second"`
}

type InfoMeta struct {
	ID 					int 	`json:"ID"`
	NavigationStatus 	int 	`json:"Navigation_Status"`
	MMSI 				int 	`json:"MMSI"`
	AIS 				int 	`json:"AIS"`
	IMO 				int 	`json:"IMO"`
	CallSign 			string 	`json:"Call_Sign"`
	Name 				string 	`json:"Name"`
	ShipType			int 	`json:"Ship_Type"`
	A 					int 	`json:"A"`
	B 					int 	`json:"B"`
	C 					int 	`json:"C"`
	D 					int 	`json:"D"`
	Length 				int 	`json:"Length"`
	Width 				int 	`json:"Width"`
	PositionType		int 	`json:"Position_Type"`
	ETAMonth 			int 	`json:"ETA_Month"`
	ETADay 				int 	`json:"ETA_Day"`
	ETAHour				int 	`json:"ETA_Hour"`
	ETAMinute			int		`json:"ETA_Minute"`
	Draft				float64	`json:"Draft"`
	Destination			string	`json:"Destination"`
	Year 				int 	`json:"Year"`
	Month 				int 	`json:"Month"`
	Day 				int 	`json:"Day"`
	Hour				int 	`json:"Hour"`
	Minute 				int		`json:"Minute"`
	Second				int 	`json:"Second"`
}

type GetShipResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    []int 			  `json:"data"`
}

type GetPositionWithShipIDResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []PositionMeta `json:"data"`
}

type GetInfoWithShipIDResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []InfoMeta `json:"data"`
}