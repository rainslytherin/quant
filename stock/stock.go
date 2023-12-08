package stock

import (
	"quant/database"
	"time"
)

var GlobalStrategyList *[]Strategy

type Strategy struct {
	ID         int
	StrikeCode string
	StrikeName string
	Title      string
	Enable     int
	Multiple   int
	MultipleCn int
	Action     int
	EffectTime string
	UpperLimit float64
	LowerLimit float64
	DelayTime  int
	IsDeleted  int
	CreateTime string
	Remark     string
}

func GetGlobalStrategyList() (*[]Strategy, error) {
	GlobalStrategyList, err := GetStrategyList()
	if err != nil {
		return nil, err
	}
	go flushGlobalStrategyList()
	return GlobalStrategyList, nil
}

func flushGlobalStrategyList() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			strategyList, _ := GetStrategyList()
			GlobalStrategyList = strategyList
		}

	}
}

// get strategy list from database
func GetStrategyList() (*[]Strategy, error) {
	// get database connetion from global_db
	// query the strategy list
	// return the strategy list

	db, err := database.GetGlobalDB()
	if err != nil {
		return nil, err
	}

	sql := "select * from strategy where enable = 1"

	var strategyList []Strategy
	err = db.Select(&strategyList, sql)

	return &strategyList, err
}
