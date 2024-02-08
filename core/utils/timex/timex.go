package timex

import (
	"fmt"
	"time"
)

// 初始化为一个足够长的过去时间，以防timex.Now() - lastTime的结果等于0。这个过去时间被用作起始时间点
var initTime = time.Now().AddDate(-1, -1, -1)

// Now 函数返回当前时间相对于initTime的时间间隔（即相对时间）。
// 这个函数的返回值类型是time.Duration，表示从initTime到当前时间的时间间隔
func Now() time.Duration {
	return time.Since(initTime)
}

// Since 函数接收一个时间间隔d作为参数，并返回从给定时间间隔d到当前时间的时间差
func Since(d time.Duration) time.Duration {
	return time.Since(initTime) - d
}

// ReprOfDuration 将给定的时间间隔转换为毫秒为单位的字符串
// 例子：1500000000 纳秒（等于 1.5 秒）@return 1500.0ms
func ReprOfDuration(duration time.Duration) string {
	return fmt.Sprintf("%.1fms", float32(duration)/float32(time.Millisecond))
}
