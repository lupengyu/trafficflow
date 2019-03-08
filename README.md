# AIS交通流特性统计分析

## 功能
### 已经实现功能
* 统计船舶交通量
    * 时间段计算（统计一个时间段内的数据）
    * 区分船舶大小
    * 区分船舶类型
    * 区分区域
* 统计船舶密度
    * 瞬时计算（时间段很短）
    * 区分船舶大小
    * 区分船舶类型
    * 区分区域
* 统计船舶航速
    * 瞬时计算
    * 统计平均航速
    * 区分区域统计平均航速
    * 航速区分为 [0, 5] (5, 10] (10, 15] (15, 20] (20, +∞)
* 统计船舶航迹
    * 时间段计算
    * 统计航迹过门线的次数（去重与不去重）
* 统计船舶间距
    * 瞬时计算
    * 统计每艘船与别的船的最短间距
    * 统计所有船中的最短间距
    * 间距区分为 [0, 50) [50, 300] (300, +∞)
* 统计船舶会遇
    * 时间段内瞬时计算
    * 自定义计算精度(默认1分钟)
    * 分为简单会遇(一次会遇仅有两艘船参与)与复杂会遇(一次会遇有多艘船参与)
    * 分为普通会遇(距离小于0.5海里)与危险会遇(外船介入本船船舶领域)

### TODO
* (P0)预测普通会遇中的会遇点
* (P0)计算会遇中的回避率(即在普通会遇中出现预测危险会遇, 但是通过回避以避免)与回避行为发生地
* (P1)实现对危险行为的实时预警
* (P2)记录普通会遇发生地
* (P3)可视化展示

## 目录结构  
│  Gopkg.lock 依赖库管理项  
│  Gopkg.toml 依赖库管理项  
│  main.go 程序入口  
│  
│--client 数据流客户端  
│  │--http http客户端  
│  │      config.go http服务器定义  
│  │      get_info.go  
│  │      get_info_test.go  
│  │      get_position.go  
│  │      get_position_test.go  
│  │      get_ship.go  
│  │      get_ship_test.go  
│  │      request.go http请求聚类  
│  │  
│  │--sql mysql客户端  
│          table_info.go info表  
│          table_info_test.go  
│          table_position.go position表  
│          table_position_test.go  
│          table_ship.go ship表  
│          table_ship_test.go  
│  
│--constant 类型与常数声明  
│      constants.go  
│  
│--dal 驱动  
│  │--mysql mysql-for-go驱动  
│  │      drive.go  驱动程序  
│  │--cache 程序缓存  
│          drive.go  驱动程序  
│          shipinfo.go  
│  
│--handler 方法类  
│      density.go 统计船密度流  
│      doorline.go 统计轨迹门线  
│      meeting.go 统计会遇  
│      spacing.go 统计船间距  
│      speed.go 统计船速  
│      traffic.go 统计交通流  
│  
│--helper 帮助类  
│      helper.go  通用helper类  
│      math.go  数学helper类  
│      helpers_test.go  
│  
│--vendor 依赖项  
