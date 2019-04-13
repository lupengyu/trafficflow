# AIS交通流特性统计分析

## 功能
### 已经实现功能
* AIS数据清洗与补全
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
* 统计船舶会遇危险度

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
│--data 输出数据保存  
│  │--clean AIS清洗数据文件夹  
│  │--clean_repair AIS清洗修复数据文件夹  
│  │--doorline 门线数据文件夹  
│  │--meeting 会遇危险度数据文件夹  
│  │--segmentation 船舶航迹航迹分片数据文件夹  
│  │--trajectory 船舶航迹数据文件夹  
│  
│--handler 方法类  
│      clean_repair.go 数据清洗和补全  
│      density.go 统计船密度流  
│      doorline.go 统计轨迹门线  
│      earlywarning.go 船舶实时预警   
│      meeting.go 统计会遇  
│      segmentation.go 获得船舶航迹分片  
│      spacing.go 统计船间距  
│      speed.go 统计船速  
│      traffic.go 统计交通流  
│      trajectory.go 获得船舶航迹  
│  
│--helper 帮助类  
│      helpers_test.go  helper测试类  
│      helper.go  通用helper类  
│      math.go  数学helper类  
│      sort.go  排序helper类  
│  
│--vendor 依赖项  
