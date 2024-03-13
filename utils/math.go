package utils

func Min(vals ...int64) int64 {
	var min int64
	for _, val := range vals {
		if min == 0 || val <= min {
			min = val
		}
	}
	return min
}
