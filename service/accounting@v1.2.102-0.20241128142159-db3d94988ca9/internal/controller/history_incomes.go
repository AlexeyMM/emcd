package controller

import "time"

type Complexity struct {
	Diff          float64 `json:"difficulty"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	EventTime     int64   `json:"event_time"`
}

var complexitiesMap = map[string][]Complexity{
	"btc": {},
	"ltc": {},
}

func GetComplexityByTimeList(coin string, times []float64) map[float64]*Complexity {
	incomeComplexities := make(map[float64]*Complexity)

	complexities, ok := complexitiesMap[coin]
	lenCom := len(complexities)
	if !ok || lenCom == 0 {
		return incomeComplexities
	}

	var i int

	for _, t := range times {
		for i < lenCom {
			if t >= float64(bod(complexities[i].EventTime)) {
				if i+1 >= lenCom {
					incomeComplexities[t] = &complexities[i]
					i++
					break
				}

				if t < float64(bod(complexities[i+1].EventTime)) {
					if bod(int64(t)) < complexities[i].EventTime {
						incomeComplexities[t] = &complexities[i]
					}
					i++
					break
				}

				i++
			} else {
				break
			}
		}
	}

	return incomeComplexities
}

func bod(t int64) int64 {
	tm := time.Unix(t, 0)
	year, month, day := tm.In(time.Local).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
}
