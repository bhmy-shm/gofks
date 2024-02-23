package breaker

import (
	"log"
	"testing"
	"time"
)

// 记录时间
func recordResponseTime(rw *RollingWindow, responseTime float64) {
	rw.Add(responseTime)
}

// 记录错误
func recordError(rw *RollingWindow, occurred bool) {
	if occurred {
		rw.Add(1)
	}
}

// 计算最近10分钟内每分钟的平均响应时间
func TestRollingWindow_AverageResponseTime(t *testing.T) {
	// 创建一个新的滚动窗口，用于10个1分钟的桶。
	//rw := NewRollingWindow(10, time.Minute)
	bucketDuration := time.Duration(int64(windowDuration) / int64(buckets))
	rw := NewRollingWindow(buckets, bucketDuration)

	// 记录10分钟内的响应时间。
	for i := 0; i < 10; i++ {
		recordResponseTime(rw, float64(i*100)) // 假设响应时间递增。
		time.Sleep(time.Second)                // 等待1分钟模拟真实情况。
	}

	// 计算平均响应时间。
	var sum float64
	var count int64
	rw.Reduce(func(b *Bucket) {
		sum += b.Sum
		count += b.Count
	})

	// 由于我们是逐分钟记录的，因此每个桶应该只有一个值。
	average := sum / float64(count)

	// 验证平均响应时间是否符合预期。
	expectedAverage := float64((0 + 100 + 200 + 300 + 400 + 500 + 600 + 700 + 800 + 900) / 10)
	log.Println(expectedAverage, average, "The average response time should be equal to the expected value.")
}

// 计算某个服务的错误率。
func TestRollingWindow_ErrorRate(t *testing.T) {

	bucketDuration := time.Duration(int64(windowDuration) / int64(buckets))
	rw := NewRollingWindow(buckets, bucketDuration)

	// 记录10分钟内的错误情况。
	totalRequests := 200
	totalErrors := 0
	for i := 0; i < totalRequests; i++ {
		occurred := i%12 == 0 // 每10个请求中有1个是错误的。
		recordError(rw, occurred)
		if occurred {
			totalErrors++
		}
		//time.Sleep(time.Second * 1) // 每1秒一个请求。
	}

	// 计算错误率。
	var errorSum float64
	var requestCount int64
	rw.Reduce(func(b *Bucket) {
		errorSum += b.Sum
		requestCount += b.Count
	})

	// 错误率计算为错误总数除以请求总数。
	log.Println("requestCount", requestCount)
	log.Println("errSum:", errorSum)
	errorRate := errorSum / float64(requestCount)

	// 验证错误率是否符合预期。
	log.Println("total:", totalErrors, totalRequests)
	expectedErrorRate := float64(totalErrors) / float64(totalRequests)
	log.Println(expectedErrorRate, errorRate, "The error rate should be equal to the expected value.")
}
