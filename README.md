# AIS交通流特性统计分析

## 目录结构  
│  Gopkg.lock  
│  Gopkg.toml  
│  main.go  
│  
├─client  
│  ├─http  
│  │      config.go  
│  │      get_info.go  
│  │      get_info_test.go  
│  │      get_position.go  
│  │      get_position_test.go  
│  │      get_ship.go  
│  │      get_ship_test.go  
│  │      request.go  
│  │  
│  └─sql  
│          table_info.go  
│          table_info_test.go  
│          table_position.go  
│          table_position_test.go  
│          table_ship.go  
│          table_ship_test.go  
│  
├─constant  
│      constants.go  
│  
├─dal  
│  └─mysql  
│          drive.go  
│  
├─handler  
│      density.go  
│      doorline.go  
│      spacing.go  
│      speed.go  
│      traffic.go  
│  
├─helper  
│      helpers.go  
│      helpers_test.go  
│  
└─vendor  
