package internal

var DefaultStockfishParams = map[string]interface{}{
    "Debug Log File":         "",
    "Contempt":               0,
    "Min Split Depth":        0,
    "Threads":                1,
    "Ponder":                 false,
    "Hash":                   16,
    "MultiPV":                1,
    "Skill Level":            20,
    "Move Overhead":          10,
    "Minimum Thinking Time":  20,
    "Slow Mover":             100,
    "UCI_Chess960":           false,
    "UCI_LimitStrength":      false,
    "UCI_Elo":                1350,
}

var _delete_counter int

var ParamRestrictions = map[string]ParamRestriction{
    "Debug Log File":       {Type: "string", Min: nil, Max: nil},
    "Threads":              {Type: "int", Min: 1, Max: 1024},
    "Hash":                 {Type: "int", Min: 1, Max: 2048},
    "Ponder":               {Type: "bool", Min: nil, Max: nil},
    "MultiPV":              {Type: "int", Min: 1, Max: 500},
    "Skill Level":          {Type: "int", Min: 0, Max: 20},
    "Move Overhead":        {Type: "int", Min: 0, Max: 5000},
    "Slow Mover":           {Type: "int", Min: 10, Max: 1000},
    "UCI_Chess960":         {Type: "bool", Min: nil, Max: nil},
    "UCI_LimitStrength":    {Type: "bool", Min: nil, Max: nil},
    "UCI_Elo":              {Type: "int", Min: 1320, Max: 3190},
    "Contempt":             {Type: "int", Min: -100, Max: 100},
    "Min Split Depth":      {Type: "int", Min: 0, Max: 12},
    "Minimum Thinking Time":{Type: "int", Min: 0, Max: 5000},
    "UCI_ShowWDL":          {Type: "bool", Min: nil, Max: nil},
}

var Releases = map[string]string{
    "16.0": "2023-06-30",
    "15.1": "2022-12-04",
    "15.0": "2022-04-18",
    "14.1": "2021-10-28",
    "14.0": "2021-07-02",
    "13.0": "2021-02-19",
    "12.0": "2020-09-02",
    "11.0": "2020-01-18",
    "10.0": "2018-11-29",
}
