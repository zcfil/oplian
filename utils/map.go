package utils

import "log"

func ArrayNotRepeat(arr []uint64, mp map[uint64]struct{}) int {
	not := 0
	log.Println(arr, mp)
	for i := 0; i < len(arr); i++ {
		if _, ok := mp[arr[i]]; !ok {
			not++
		}
	}
	return not
}
