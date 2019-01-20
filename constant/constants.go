package constant

const (
	LongitudeMin  float64 = 117.7737
	LongitudeMax  float64 = 118.63037
	LatitudeMin   float64 = 24.08784
	LatitudeMax   float64 = 24.691
	BigShipLength int     = 100
)

type ShipMeta struct {
	MMSI int
}

type PositionMeta struct {
	ID                  int     `json:"ID"`
	MessageType         int     `json:"Message_Type"`
	RepeatIndicator     int     `json:"Repeat_Indicator"`
	MMSI                int     `json:"MMSI"`
	NavigationStatus    int     `json:"Navigation_Status"`
	ROT                 int     `json:"ROT"`
	SOG                 float64 `json:"SOG"`
	PositionAccuracy    int     `json:"Position_Accuracy"`
	Longitude           float64 `json:"Longitude"`
	Latitude            float64 `json:"Latitude"`
	COG                 float64 `json:"COG"`
	HDG                 int     `json:"HDG"`
	TimeStamp           int     `json:"Time_stamp"`
	ReservedForRegional int     `json:"Reserved_for_regional"`
	RAIMFlag            int     `json:"RAIM_flag"`
	Year                int     `json:"Year"`
	Month               int     `json:"Month"`
	Day                 int     `json:"Day"`
	Hour                int     `json:"Hour"`
	Minute              int     `json:"Minute"`
	Second              int     `json:"Second"`
}

type InfoMeta struct {
	ID               int     `json:"ID"`
	NavigationStatus int     `json:"Navigation_Status"`
	MMSI             int     `json:"MMSI"`
	AIS              int     `json:"AIS"`
	IMO              int     `json:"IMO"`
	CallSign         string  `json:"Call_Sign"`
	Name             string  `json:"Name"`
	ShipType         int     `json:"Ship_Type"`
	A                int     `json:"A"`
	B                int     `json:"B"`
	C                int     `json:"C"`
	D                int     `json:"D"`
	Length           int     `json:"Length"`
	Width            int     `json:"Width"`
	PositionType     int     `json:"Position_Type"`
	ETAMonth         int     `json:"ETA_Month"`
	ETADay           int     `json:"ETA_Day"`
	ETAHour          int     `json:"ETA_Hour"`
	ETAMinute        int     `json:"ETA_Minute"`
	Draft            float64 `json:"Draft"`
	Destination      string  `json:"Destination"`
	Year             int     `json:"Year"`
	Month            int     `json:"Month"`
	Day              int     `json:"Day"`
	Hour             int     `json:"Hour"`
	Minute           int     `json:"Minute"`
	Second           int     `json:"Second"`
}

type GetShipResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []int  `json:"data"`
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

type Data struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

type CulTrafficRequest struct {
	StartTime *Data
	EndTime   *Data
	LotDivide int
	LatDivide int
}

type CulTrafficResponse struct {
	AreaTraffics [][]AreaTraffic
	TrafficData  *TrafficData
}

type AreaTraffic struct {
	ShipMap               map[int]int   //区域内船舶map		->	日流量
	Traffic               int           //区域内流量			->	日流量和
	HourShipMap           []map[int]int //区域内各小时船舶map	->	小时流量
	HourTraffic           []int         //区域内小时流量		->	小时流量和
	SmallShipMap          map[int]int   //区域内小船map		->	日小船流量
	SmallShipTraffic      int           //区域内小船流量		->	日小船流量和
	HourSmallShipMap      []map[int]int //区域内各小时小船map	->	小时小船流量
	HourSmallShipTraffic  []int         //区域内小时小船流量	->	小时小船流量和
	BigShipMap            map[int]int   //区域内大船map		->	日大船流量
	BigShipTraffic        int           //区域内大船流量		->	日大船流量和
	HourBigShipMap        []map[int]int //区域内各小时大船map	->	小时大船流量
	HourBigShipTraffic    []int         //区域内小时大船流量	->	小时大船流量和
	Type0ShipMap          map[int]int
	Type0ShipTraffic      int
	HourType0ShipMap      []map[int]int
	HourType0ShipTraffic  []int
	Type6xShipMap         map[int]int
	Type6xShipTraffic     int
	HourType6xShipMap     []map[int]int
	HourType6xShipTraffic []int
	Type7xShipMap         map[int]int
	Type7xShipTraffic     int
	HourType7xShipMap     []map[int]int
	HourType7xShipTraffic []int
	Type8xShipMap         map[int]int
	Type8xShipTraffic     int
	HourType8xShipMap     []map[int]int
	HourType8xShipTraffic []int
}

type TrafficData struct {
	HourTrafficSum           []int //各区域小时流量相和			->	总小时流量和
	HourSmallShipTrafficSum  []int //各区域小时小船流量相和		->	总小时小船流量和
	HourBigShipTrafficSum    []int //各区域小时大船流量相和		->	总小时大船流量和
	HourType0ShipTrafficSum  []int
	HourType6xShipTrafficSum []int
	HourType7xShipTrafficSum []int
	HourType8xShipTrafficSum []int
}
