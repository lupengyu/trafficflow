package constant

const (
	LongitudeMin     float64 = 117.7737  // 经度最小值
	LongitudeMax     float64 = 118.63037 // 经度最大值
	LatitudeMin      float64 = 24.08784  // 纬度最小值
	LatitudeMax      float64 = 24.691    // 纬度最大值
	BigShipLength    int     = 100       // 大型船舶判断条件
	HalfNauticalMile float64 = 926       // 0.5海里
	NauticalMile     float64 = 1852      // 1海里
	StaticShip       float64 = 0.03      // 静止船舶判断条件
)

type ShipMeta struct {
	MMSI int // 用户ID
}

type PositionMeta struct {
	ID                  int     `json:"ID"`
	MessageType         int     `json:"Message_Type"`          // 消息类型
	RepeatIndicator     int     `json:"Repeat_Indicator"`      // 转发指示符
	MMSI                int     `json:"MMSI"`                  // 用户ID
	NavigationStatus    int     `json:"Navigation_Status"`     // 导航状态
	ROT                 int     `json:"ROT"`                   // 旋转速率
	SOG                 float64 `json:"SOG"`                   // 地面航速
	PositionAccuracy    int     `json:"Position_Accuracy"`     // 位置精度
	Longitude           float64 `json:"Longitude"`             // 经度
	Latitude            float64 `json:"Latitude"`              // 纬度
	COG                 float64 `json:"COG"`                   // 地面航向
	HDG                 int     `json:"HDG"`                   // 实际航向
	TimeStamp           int     `json:"Time_stamp"`            // 时戳
	ReservedForRegional int     `json:"Reserved_for_regional"` // 特定操作指示符
	RAIMFlag            int     `json:"RAIM_flag"`             // RAIM标志
	Year                int     `json:"Year"`                  // 年
	Month               int     `json:"Month"`                 // 月
	Day                 int     `json:"Day"`                   // 日
	Hour                int     `json:"Hour"`                  // 时
	Minute              int     `json:"Minute"`                // 分
	Second              int     `json:"Second"`                // 秒
}

type InfoMeta struct {
	ID               int     `json:"ID"`
	NavigationStatus int     `json:"Navigation_Status"` // 转发指示符
	MMSI             int     `json:"MMSI"`              // 用户ID
	AIS              int     `json:"AIS"`               // AIS版本指示符
	IMO              int     `json:"IMO"`               // IMO编号
	CallSign         string  `json:"Call_Sign"`         // 呼号
	Name             string  `json:"Name"`              // 名称
	ShipType         int     `json:"Ship_Type"`         // 船舶和货物类型
	A                int     `json:"A"`                 // 船舱距船头距离
	B                int     `json:"B"`                 // 船舱距船尾距离
	C                int     `json:"C"`                 // 船舱距左侧距离
	D                int     `json:"D"`                 // 船舱距右侧距离
	Length           int     `json:"Length"`            // 船舶长度
	Width            int     `json:"Width"`             // 船舶宽度
	PositionType     int     `json:"Position_Type"`     // 位置精度
	ETAMonth         int     `json:"ETA_Month"`         // 预估到达时间-月
	ETADay           int     `json:"ETA_Day"`           // 预估到达时间-日
	ETAHour          int     `json:"ETA_Hour"`          // 预估到达时间-时
	ETAMinute        int     `json:"ETA_Minute"`        // 预估到达时间-分
	Draft            float64 `json:"Draft"`             // 目前最大静态吃水
	Destination      string  `json:"Destination"`       // 目的地
	Year             int     `json:"Year"`              // 年
	Month            int     `json:"Month"`             // 月
	Day              int     `json:"Day"`               // 日
	Hour             int     `json:"Hour"`              // 时
	Minute           int     `json:"Minute"`            // 分
	Second           int     `json:"Second"`            // 秒
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
	StartTime *Data // 起始时间
	EndTime   *Data // 结束时间
	LotDivide int   // 经度划分数目
	LatDivide int   // 纬度划分数目
}

type CulTrafficResponse struct {
	AreaTraffics [][]AreaTraffic // 区域交通量统计
	TrafficData  *TrafficData    // 综合小时交通量统计
}

type AreaTraffic struct {
	ShipMap               map[int]int   // 区域内船舶map		->	日流量
	Traffic               int           // 区域内流量			->	日流量和
	HourShipMap           []map[int]int // 区域内各小时船舶map	->	小时流量
	HourTraffic           []int         // 区域内小时流量		->	小时流量和
	SmallShipMap          map[int]int   // 区域内小船map		->	日小船流量
	SmallShipTraffic      int           // 区域内小船流量		->	日小船流量和
	HourSmallShipMap      []map[int]int // 区域内各小时小船map	->	小时小船流量
	HourSmallShipTraffic  []int         // 区域内小时小船流量	->	小时小船流量和
	BigShipMap            map[int]int   // 区域内大船map		->	日大船流量
	BigShipTraffic        int           // 区域内大船流量		->	日大船流量和
	HourBigShipMap        []map[int]int // 区域内各小时大船map	->	小时大船流量
	HourBigShipTraffic    []int         // 区域内小时大船流量	->	小时大船流量和
	Type0ShipMap          map[int]int   // 渔船map
	Type0ShipTraffic      int           // 渔船交通量
	HourType0ShipMap      []map[int]int // 渔船小时map
	HourType0ShipTraffic  []int         // 渔船小时交通量
	Type6xShipMap         map[int]int   // 客船map
	Type6xShipTraffic     int           // 客船交通量
	HourType6xShipMap     []map[int]int // 客船小时map
	HourType6xShipTraffic []int         // 客船小时交通量
	Type7xShipMap         map[int]int   // 货轮map
	Type7xShipTraffic     int           // 货轮交通量
	HourType7xShipMap     []map[int]int // 货轮小时map
	HourType7xShipTraffic []int         // 货轮小时交通量
	Type8xShipMap         map[int]int   // 油轮map
	Type8xShipTraffic     int           // 油轮交通量
	HourType8xShipMap     []map[int]int // 油轮小时map
	HourType8xShipTraffic []int         // 油轮小时交通量
}

type TrafficData struct {
	HourTrafficSum           []int // 各区域小时流量相和		->	总小时流量和
	HourSmallShipTrafficSum  []int // 各区域小时小船流量相和		->	总小时小船流量和
	HourBigShipTrafficSum    []int // 各区域小时大船流量相和		->	总小时大船流量和
	HourType0ShipTrafficSum  []int // 各区域小时渔船流量相和
	HourType6xShipTrafficSum []int // 各区域小时客船流量相和
	HourType7xShipTrafficSum []int // 各区域小时货轮流量相和
	HourType8xShipTrafficSum []int // 各区域小时油轮流量相和
}

type CulDensityRequest struct {
	Time      *Data // 时间
	DeltaT    *Data // 时间范围
	LotDivide int   // 经度划分数目
	LatDivide int   // 纬度划分数目
}

type CulDensityResponse struct {
	DensityData *DensityData    // 总体密度数据
	AreaDensity [][]AreaDensity // 区域密度数据
}

type DensityData struct {
	ShipDensity      int // 船舶密度
	SmallShipDensity int // 小船密度
	BigShipDensity   int // 大船密度
	Type0Density     int // 渔船密度
	Type6xDensity    int // 客船密度
	Type7xDensity    int // 货轮密度
	Type8xDensity    int // 油轮密度
}

type AreaDensity struct {
	ShipMap          map[int]int // 区域内船舶map
	Density          int         // 区域内船舶密度
	SmallShipShipMap map[int]int // 区域内小船map
	SmallShipDensity int         // 区域内小船密度
	BigShipShipMap   map[int]int // 区域内大船map
	BigShipDensity   int         // 区域内大船密度
	Type0ShipMap     map[int]int // 区域内渔船map
	Type0Density     int         // 区域内渔船密度
	Type6xShipMap    map[int]int // 区域内客船map
	Type6xDensity    int         // 区域内客船密度
	Type7xShipMap    map[int]int // 区域内货轮map
	Type7xDensity    int         // 区域内货轮密度
	Type8xShipMap    map[int]int // 区域内油轮map
	Type8xDensity    int         // 区域内油轮密度
}

type CulSpeedRequest struct {
	Time      *Data // 时间
	DeltaT    *Data // 时间范围
	LotDivide int   // 经度划分数目
	LatDivide int   // 纬度划分数目
}

type CulSpeedResponse struct {
	SpeedData *SpeedData    // 总体速度数据
	AreaSpeed [][]AreaSpeed // 区域速度数据
}

type SpeedData struct {
	ShipSpeed      float64 // 平均速度
	ShipCnt        int     // 船舶总数
	ShipSpeedRange []int   // 船舶速度区间
}

type AreaSpeed struct {
	ShipSpeed       float64         // 区域平均速度
	ShipCnt         int             // 区域船舶总数
	ShipSpeedSumMap map[int]float64 // 区域单一船舶数据项速度取值总和
	ShipSpeedCnt    map[int]int     // 区域单一船舶数据项数目
}

type Position struct {
	Longitude float64 // 经度
	Latitude  float64 // 纬度
}

type Track struct {
	PrePosition      *Position // 之前位置
	DeWeightDoorLine bool      // 判断是否过门线
	Time             *Data     // 当前时间
	Deviation        int64     // 当前时间精度
	COG              float64   // 船舶COG
	SOG              float64   // 船舶SOG
}

type CulDoorLineRequest struct {
	StartPosition *Position // 门线起点
	EndPosition   *Position // 门线终点
	StartTime     *Data     // 起始时间
	EndTime       *Data     // 终止时间
}

type CulDoorLineResponse struct {
	Cnt            int // 过门线总数
	DeWeightingCnt int // 去重后过门线总数
}

type CulSpacingRequest struct {
	Time   *Data // 时间
	DeltaT *Data // 时间范围
}

type CulSpacingResponse struct {
	MinSpacing   float64                 // 最短相距距离
	MinSpaceA    int                     // 最短相距距离-船A
	MinSpaceB    int                     // 最短相距距离-船B
	APosition    *Position               // A船位置
	BPosition    *Position               // B船位置
	SpacingMap   map[int]float64         // 船舶最短距离map
	SpacingRange []int                   // 船舶最短距离分布区间
	ShipSpacing  map[int]map[int]float64 // 船舶间距map
	TrackMap     map[int]*Track          // 船舶当前信息
}

type CulMeetingRequest struct {
	StartTime *Data // 起始时间
	EndTime   *Data // 终止时间
	DeltaT    *Data // 时间范围
	TimeRange *Data // 时间精度
}

type CulMeetingResponse struct {
	SimpleMeeting              int   // 简单会遇数目
	ComplexMeeting             int   // 复杂会遇数目
	SimpleDamageMeeting        int   // 简单危险会遇数目
	ComplexDamageMeeting       int   // 复杂危险会遇数目
	ForecastDamageMeeting      int   // 预测危险会遇数目
	DamageMeetingAvoid         int   // 危险会遇规避数目
	AngleForecastDamageMeeting []int // 各个角度预测危险会遇数目
	AngleDamageMeetingAvoid    []int // 各个角度危险会遇规避数目
}

type MeetingIntersection struct {
	DCPA    float64 // 最近会遇距离
	TCPA    float64 // 最近会遇时间
	Azimuth float64 // 会遇方向
	VR      float64 // 相对速度
}

type EarlyWarningRequest struct {
	StartTime *Data // 起始时间
	EndTime   *Data // 终止时间
	DeltaT    *Data // 时间范围
	TimeRange *Data // 时间精度
	MMSI      int   // 用户ID
}

type EarlyWarningResponse struct {
	Warning []*Warning // 预警信息
}

type Alert struct {
	MMSI                int                  // 预警对象船只
	ShipTrack           *Track               // 预警对象船只数据
	IsEmergency         bool                 // 是否是紧急情况
	Distance            float64              // 相对距离
	Azimuth             float64              // 相对方位
	MeetingIntersection *MeetingIntersection // 危险会遇预测数据
	UDCPA               float64              // DCPA预警值
	UTCPA               float64              // TCPA预警值
	UB                  float64              // 相对角度预警值
	UD                  float64              // 相对距离预警
	UV                  float64              // 本船速度预警值
	Danger              float64              // 预警评分
}

type Warning struct {
	MasterShipTrack *Track   // 主船只数据
	Alerts          []*Alert // 警告列表
	Time            *Data    // 警告时间
}
