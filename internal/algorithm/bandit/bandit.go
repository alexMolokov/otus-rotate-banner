package bandit

import "math"

type Stat struct {
	ID     int
	Trials int
	Reward int
}

// Choice выбрать самый подходящий из возможных вариантов (статистик) возвращает ID самого подходящего варинта.
func Choice(stat []Stat, allTrials int) int {
	index := 0
	var statistic float64 = -1
	for i, v := range stat {
		reward := float64(v.Reward) / float64(v.Trials)
		currentStatistic := reward + math.Sqrt(2*math.Log(float64(allTrials))/float64(v.Trials))
		if currentStatistic > statistic {
			index, statistic = i, currentStatistic
		}
	}
	return stat[index].ID
}
