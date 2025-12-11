package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type LoadTestResult struct {
	StatusCode int
	Duration   time.Duration
	Timestamp  time.Time
}

type Report struct {
	TotalRequests     int64
	TotalTime         time.Duration
	SuccessRequests   int64 // HTTP 200
	StatusDistribution map[int]int64
	MinDuration       time.Duration
	MaxDuration       time.Duration
	AvgDuration       time.Duration
}

func main() {
	url := flag.String("url", "", "URL to test")
	requests := flag.Int64("requests", 0, "Total number of requests")
	concurrency := flag.Int("concurrency", 1, "Number of concurrent requests")
	flag.Parse()

	if *url == "" || *requests == 0 || *concurrency == 0 {
		fmt.Println("Error: --url, --requests, and --concurrency parameters are required")
		flag.PrintDefaults()
		return
	}

	fmt.Printf("Starting load test\n")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Total Requests: %d\n", *requests)
	fmt.Printf("Concurrency: %d\n\n", *concurrency)

	report := performLoadTest(*url, *requests, *concurrency)
	printReport(report)
}

func performLoadTest(url string, totalRequests int64, concurrency int) *Report {
	startTime := time.Now()
	resultChan := make(chan LoadTestResult, concurrency)
	
	var wg sync.WaitGroup
	var requestCounter int64 = 0
	var totalDuration int64 = 0
	var minDuration time.Duration = time.Duration(1<<63 - 1) // max int64
	var maxDuration time.Duration = 0

	// Worker function
	worker := func() {
		defer wg.Done()
		client := &http.Client{
			Timeout: time.Second * 30,
		}

		for {
			current := atomic.AddInt64(&requestCounter, 1)
			if current > totalRequests {
				break
			}

			start := time.Now()
			resp, err := client.Get(url)
			duration := time.Since(start)

			if err != nil {
				resultChan <- LoadTestResult{
					StatusCode: 0,
					Duration:   duration,
					Timestamp:  time.Now(),
				}
			} else {
				resp.Body.Close()
				resultChan <- LoadTestResult{
					StatusCode: resp.StatusCode,
					Duration:   duration,
					Timestamp:  time.Now(),
				}
			}

			// Track min/max/total duration
			durNano := duration.Nanoseconds()
			atomic.AddInt64(&totalDuration, durNano)
			
			for {
				currentMin := atomic.LoadInt64((*int64)(&minDuration))
				if durNano < currentMin {
					if atomic.CompareAndSwapInt64((*int64)(&minDuration), currentMin, durNano) {
						break
					}
				} else {
					break
				}
			}

			for {
				currentMax := atomic.LoadInt64((*int64)(&maxDuration))
				if durNano > currentMax {
					if atomic.CompareAndSwapInt64((*int64)(&maxDuration), currentMax, durNano) {
						break
					}
				} else {
					break
				}
			}
		}
	}

	// Start workers
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go worker()
	}

	// Collect results in background
	statusDistribution := make(map[int]int64)
	var successCount int64 = 0
	var resultWg sync.WaitGroup
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for result := range resultChan {
			statusDistribution[result.StatusCode]++
			if result.StatusCode == 200 {
				atomic.AddInt64(&successCount, 1)
			}
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)
	resultWg.Wait()

	totalTime := time.Since(startTime)
	finalCount := atomic.LoadInt64(&requestCounter)
	if finalCount > totalRequests {
		finalCount = totalRequests
	}

	avgDuration := time.Duration(0)
	if finalCount > 0 {
		avgDuration = time.Duration(atomic.LoadInt64(&totalDuration) / finalCount)
	}

	return &Report{
		TotalRequests:      finalCount,
		TotalTime:          totalTime,
		SuccessRequests:    atomic.LoadInt64(&successCount),
		StatusDistribution: statusDistribution,
		MinDuration:        minDuration,
		MaxDuration:        maxDuration,
		AvgDuration:        avgDuration,
	}
}

func printReport(report *Report) {
	fmt.Println("\n========== LOAD TEST REPORT ==========")
	fmt.Printf("Total Time:        %s\n", report.TotalTime)
	fmt.Printf("Total Requests:    %d\n", report.TotalRequests)
	fmt.Printf("Successful (200):  %d\n", report.SuccessRequests)
	fmt.Printf("Min Duration:      %v\n", report.MinDuration)
	fmt.Printf("Max Duration:      %v\n", report.MaxDuration)
	fmt.Printf("Avg Duration:      %v\n", report.AvgDuration)
	
	if report.TotalTime > 0 {
		requestsPerSecond := float64(report.TotalRequests) / report.TotalTime.Seconds()
		fmt.Printf("Requests/second:   %.2f\n", requestsPerSecond)
	}

	fmt.Println("\nStatus Code Distribution:")
	for status, count := range report.StatusDistribution {
		if status == 0 {
			fmt.Printf("  Error (timeout/failed): %d\n", count)
		} else {
			fmt.Printf("  HTTP %d: %d\n", status, count)
		}
	}
	fmt.Println("=======================================")
}
