package main

import (
	"flag"
	"fmt"
	"quant/config"
	"quant/stock"
	"time"
)

var dir string

func init() {
	flag.StringVar(&dir, "dir", "etc", "config dir")
}

func main() {
	config.InitGlobalConfig(dir)

	// 1.获取股票列表
	// 2.订阅实时股价数据
	// 3.并发处理每只股票的数据

	// 策略列表也是会变得，所以需要定时更新

	GlobalStrategyList, err := stock.GetGlobalStrategyList()

	for {
		fmt.Println(GlobalStrategyList, err)
		time.Sleep(time.Second * 10)
	}
}
