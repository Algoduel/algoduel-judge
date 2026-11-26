// Package ladder — лестница ступеней: размеры тестов, очки и оценка роста.
package ladder

import "math"

// Config — параметры лестницы задачи: ступени n = 2^From … 2^To.
type Config struct {
	From int // показатель нижней ступени
	To   int // показатель верхней ступени
}

// Exps возвращает показатели всех ступеней снизу вверх.
func (c Config) Exps() []int {
	exps := make([]int, 0, c.To-c.From+1)
	for e := c.From; e <= c.To; e++ {
		exps = append(exps, e)
	}
	return exps
}

// Score — очки за эффективность: доля пройденных ступеней.
func Score(passed, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(passed) / float64(total)
}

// GrowthExponents оценивает, во сколько раз растут затраты при удвоении n,
// в виде показателя степени: ~1 для O(n), ~2 для O(n²), чуть больше 1 для O(n log n).
// costs — число инструкций на подряд идущих ступенях (уже без стоимости старта).
func GrowthExponents(costs []uint64) []float64 {
	var out []float64
	for i := 1; i < len(costs); i++ {
		if costs[i-1] == 0 || costs[i] == 0 {
			continue
		}
		out = append(out, math.Log2(float64(costs[i])/float64(costs[i-1])))
	}
	return out
}
