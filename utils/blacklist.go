package utils

import (
	"github.com/cihub/seelog"
	"sync"
	"time"
)

type blacklist struct {
	list map[string]int64
	lock sync.Mutex
}

var bl  = &blacklist{
	list: make(map[string]int64),
	lock: sync.Mutex{},
}

func InitBlacklist() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error(err)
				seelog.Flush()
			}
		}()
		ticker := time.NewTicker(1 * time.Second)
		for {
			<- ticker.C
			deleteBlackList()
		}
	}()
}

func InBlackList(trainNo string) bool {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	_, ok := bl.list[trainNo]
	return ok
}

func AddBlackList(trainNo string) {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	bl.list[trainNo] = time.Now().Unix() + 60
}

func deleteBlackList() {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	timeNow := time.Now().Unix()
	for k, ts := range bl.list {
		if ts <= timeNow {
			delete(bl.list, k)
		}
	}
}