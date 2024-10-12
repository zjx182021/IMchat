package utils

import "time"

type TimerFunc func(interface{}) bool

func Timer(delay, tick time.Duration, f TimerFunc, param interface{}) {
	go func() {
		if f == nil {
			return
		}
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				if f(param) {
					t.Reset(tick)
				} else {
					return
				}
			}
		}
	}()
}
