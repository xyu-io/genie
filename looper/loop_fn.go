package looper

import "time"

// TimeLoopThen 循环执行某个函数，等待周期到达再执行
func TimeLoopThen(interval time.Duration, runNow bool, fn func(time.Time)) func() {
	if fn == nil {
		return nil
	}
	if runNow {
		fn(time.Now())
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			now := <-ticker.C
			fn(now)
		}
	}()
	return func() {
		ticker.Stop()
	}
}
