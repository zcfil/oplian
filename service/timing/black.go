package timing

import (
	"log"
	"oplian/global"
	"time"
)

func DelBlackCache() {
	now := time.Now()
	zero := time.After(time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Sub(now))
	atime := time.NewTimer(time.Hour * 24)
	go func() {
		for {
			select {
			case <-atime.C:
			case <-zero:
			}
			global.BlackLock.Lock()
			for token, exp := range global.BlackCache {
				//Remove expired tokens
				if time.Now().Unix() > exp {
					delete(global.BlackCache, token)
					log.Panicln("移除：", token)
				}
			}
			global.BlackLock.Unlock()
			atime = time.NewTimer(time.Hour * 24)
		}
	}()
}
