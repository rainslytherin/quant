package game

import (
	"encoding/json"
	"fmt"
	"quant/database"
)

// Get the game table for a given stock and date.
func GetGameTable(stock string, date string) (GameTable, error) {
	return GameTable{}, nil
}

// generate the definition of the GameTable struct
type GameTable struct {
	Stock  string
	Date   string
	Points map[string]string
	Params map[string]string
}

type Point struct {
	Key   string
	Value string
	Type  int
}

func get_game_params(strike_code string, play_date string) (map[string]string, error) {
	sql := fmt.Sprintf("select start_date, params from strike_play where strike_code='%s' and time_key='%s'", strike_code, play_date)

	db, err := database.GetGlobalDB()
	if err != nil {
		return nil, err
	}

	var resp struct {
		StartDate string `db:"start_date"`
		Params    string `db:"params"`
	}

	err = db.Get(&resp, sql)
	if err != nil {
		return nil, err
	}

	if resp.StartDate == "" {
		return nil, fmt.Errorf("no data")
	}

	var params map[string]string
	err = json.Unmarshal([]byte(resp.Params), &params)
	if err != nil {
		return nil, err
	}

	params["start_date"] = resp.StartDate

	return params, nil
}

func get_game_points(strike_code string, play_date string, point_type int) ([]Point, error) {
	sql := fmt.Sprintf("select point_key, point_val, point_type from strike_play_detail where strike_code='%s' and time_key='%s'", strike_code, play_date)
	if point_type != 0 {
		sql += fmt.Sprintf(" and point_type=%d", point_type)
	}

	db, err := database.GetGlobalDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	var points []Point

	for rows.Next() {
		var point Point
		err = rows.Scan(&point.Key, &point.Value, &point.Type)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

func get_global_game_points(strike_code string) ([]Point, error) {
	sql := fmt.Sprintf("select point_key, point_val, point_type from strike_play_detail where point_type=2 ")
	if strike_code != "" {
		sql += fmt.Sprintf(" and strike_code='%s'", strike_code)
	}

	db, err := database.GetGlobalDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	var points []Point

	for rows.Next() {
		var point Point
		err = rows.Scan(&point.Key, &point.Value, &point.Type)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

func get_game_table(strike_code string, play_date string) (GameTable, error) {
	var game_table GameTable
	game_table.Stock = strike_code
	game_table.Date = play_date
	game_table.Points = make(map[string]string)

	stock_points, err := get_game_points(strike_code, play_date, 0)
	if err != nil {
		return game_table, err
	}

	global_points, err := get_global_game_points(strike_code)
	if err != nil {
		return game_table, err
	}

	stock_params, err := get_game_params(strike_code, play_date)
	if err != nil {
		return game_table, err
	}

	game_table.Params = stock_params

	// merge stock points and global points, stock points has higher priority
	for _, global_point := range global_points {
		game_table.Points[global_point.Key] = global_point.Value
	}

	for _, stock_point := range stock_points {
		game_table.Points[stock_point.Key] = stock_point.Value
	}

	return game_table, nil
}
