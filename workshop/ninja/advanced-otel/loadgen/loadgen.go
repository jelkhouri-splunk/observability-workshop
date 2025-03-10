package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Function to generate a random trace ID
func generateTraceID() string {
	return randomHex(16)
}

// Function to generate a random span ID
func generateSpanID() string {
	return randomHex(8)
}

// Helper function to generate a random hex string of a given length
func randomHex(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// Function to get the current timestamp in nanoseconds
func getCurrentTime() int64 {
	return time.Now().UnixNano()
}

// Function to send a base trace
func sendBaseTrace(traceID, spanID string, startTime, endTime int64) {
	spanJSON := map[string]interface{}{
		"resourceSpans": []interface{}{
			map[string]interface{}{
				"resource": map[string]interface{}{
					"attributes": []interface{}{
						map[string]interface{}{
							"key": "service.name",
							"value": map[string]interface{}{
								"stringValue": "cinema-service",
							},
						},
						map[string]interface{}{
							"key": "deployment.environment",
							"value": map[string]interface{}{
								"stringValue": "production",
							},
						},
					},
				},
				"scopeSpans": []interface{}{
					map[string]interface{}{
						"scope": map[string]interface{}{
							"name":    "cinema.library",
							"version": "1.0.0",
							"attributes": []interface{}{
								map[string]interface{}{
									"key": "fintest.scope.attribute",
									"value": map[string]interface{}{
										"stringValue": "Starwars, LOTR",
									},
								},
							},
						},
						"spans": []interface{}{
							map[string]interface{}{
								"traceId":           traceID,
								"spanId":            spanID,
								"name":              "/movie-validator",
								"startTimeUnixNano": fmt.Sprintf("%d", startTime),
								"endTimeUnixNano":   fmt.Sprintf("%d", endTime),
								"kind":              2,
								"status": map[string]interface{}{
									"code":    1,
									"message": "Success",
								},
								"attributes": []interface{}{
									map[string]interface{}{
										"key": "user.name",
										"value": map[string]interface{}{
											"stringValue": "George Lucas",
										},
									},
									map[string]interface{}{
										"key": "user.phone_number",
										"value": map[string]interface{}{
											"stringValue": "+1555-867-5309",
										},
									},
									map[string]interface{}{
										"key": "user.email",
										"value": map[string]interface{}{
											"stringValue": "george@deathstar.email",
										},
									},
									map[string]interface{}{
										"key": "user.account_password",
										"value": map[string]interface{}{
											"stringValue": "LOTR>StarWars1-2-3",
										},
									},
									map[string]interface{}{
										"key": "user.visa",
										"value": map[string]interface{}{
											"stringValue": "4111 1111 1111 1111",
										},
									},
									map[string]interface{}{
										"key": "user.amex",
										"value": map[string]interface{}{
											"stringValue": "3782 822463 10005",
										},
									},
									map[string]interface{}{
										"key": "user.mastercard",
										"value": map[string]interface{}{
											"stringValue": "5555 5555 5555 4444",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	sendJSON("http://localhost:4318/v1/traces", spanJSON)
	fmt.Printf("\nBase trace sent with traceId: %s and spanId: %s\n", traceID, spanID)
}

// Function to send a security trace
func sendSecurityTrace(traceID, spanID string, startTime, endTime int64) {
	securityJSON := map[string]interface{}{
		"resourceSpans": []interface{}{
			map[string]interface{}{
				"resource": map[string]interface{}{
					"attributes": []interface{}{
						map[string]interface{}{
							"key": "service.name",
							"value": map[string]interface{}{
								"stringValue": "password-check",
							},
						},
						map[string]interface{}{
							"key": "deployment.environment",
							"value": map[string]interface{}{
								"stringValue": "security-applications",
							},
						},
					},
				},
				"scopeSpans": []interface{}{
					map[string]interface{}{
						"scope": map[string]interface{}{
							"name":    "movie.library",
							"version": "1.0.0",
						},
						"spans": []interface{}{
							map[string]interface{}{
								"traceId":           traceID,
								"spanId":            spanID,
								"parentSpanId":      generateSpanID(),
								"name":              "password-validation",
								"startTimeUnixNano": fmt.Sprintf("%d", startTime),
								"endTimeUnixNano":   fmt.Sprintf("%d", endTime),
								"kind":              2,
								"status": map[string]interface{}{
									"code":    1,
									"message": "Success",
								},
								"attributes": []interface{}{
									map[string]interface{}{
										"key": "user.name",
										"value": map[string]interface{}{
											"stringValue": "George Lucas",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	sendJSON("http://localhost:4318/v1/traces", securityJSON)
	fmt.Printf("\nSecurity trace sent with traceId: %s and spanId: %s\n", traceID, spanID)
}

// Function to send a health trace
func sendHealthTrace(traceID, spanID string, startTime, endTime int64) {
	healthJSON := map[string]interface{}{
		"resourceSpans": []interface{}{
			map[string]interface{}{
				"resource": map[string]interface{}{
					"attributes": []interface{}{
						map[string]interface{}{
							"key": "service.name",
							"value": map[string]interface{}{
								"stringValue": "frontend-service",
							},
						},
						map[string]interface{}{
							"key": "deployment.environment",
							"value": map[string]interface{}{
								"stringValue": "production",
							},
						},
					},
				},
				"scopeSpans": []interface{}{
					map[string]interface{}{
						"scope": map[string]interface{}{
							"name":    "healthz",
							"version": "1.0.0",
						},
						"spans": []interface{}{
							map[string]interface{}{
								"traceId":           traceID,
								"spanId":            spanID,
								"name":              "/_healthz",
								"startTimeUnixNano": fmt.Sprintf("%d", startTime),
								"endTimeUnixNano":   fmt.Sprintf("%d", endTime),
								"kind":              2,
								"status": map[string]interface{}{
									"code":    1,
									"message": "Success",
								},
							},
						},
					},
				},
			},
		},
	}

	sendJSON("http://localhost:4318/v1/traces", healthJSON)
	fmt.Printf("\nHealth trace sent with traceId: %s and spanId: %s\n", traceID, spanID)
}

// Helper function to send JSON data via HTTP POST
func sendJSON(url string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response:", string(body))
}

// Function to generate a random quote
func getRandomQuote() string {
	lotrQuotes := []string{
		"One does not simply walk into Mordor.",
		"Even the smallest person can change the course of the future.",
		"All we have to decide is what to do with the time that is given us.",
		"There is some good in this world, and it's worth fighting for.",
	}

	starWarsQuotes := []string{
		"Do or do not, there is no try.",
		"The Force will be with you. Always.",
		"I find your lack of faith disturbing.",
		"In my experience, there is no such thing as luck.",
	}

	if rand.Intn(2) == 0 {
		return lotrQuotes[rand.Intn(len(lotrQuotes))]
	}
	return starWarsQuotes[rand.Intn(len(starWarsQuotes))]
}

// Function to generate a random log level
func getRandomLogLevel() string {
	logLevels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	return logLevels[rand.Intn(len(logLevels))]
}

// Function to generate a log entry
func generateLogEntry(jsonOutput bool) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level := getRandomLogLevel()
	message := getRandomQuote()

	if jsonOutput {
		logEntry := map[string]string{
			"timestamp": timestamp,
			"level":     level,
			"message":   message,
		}
		jsonData, _ := json.Marshal(logEntry)
		return string(jsonData)
	}
	return fmt.Sprintf("%s [%s] - %s", timestamp, level, message)
}

// Function to write logs to a file
func writeLogs(jsonOutput bool) {
	logFile := "quotes.log"
	fmt.Printf("Writing logs to %s. Press Ctrl+C to stop.\n", logFile)

	for {
		logEntry := generateLogEntry(jsonOutput)
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		file.WriteString(logEntry + "\n")
		file.Close()
		time.Sleep(1 * time.Second)
	}
}

// Display usage instructions
func printHelp() {
	fmt.Println(`Usage: trace_sender [OPTIONS]
Options:
  -base       Send base traces (enabled by default)
  -health     Send health traces
  -security   Send security traces
  -logs       Enable logging of random quotes to quotes.log
  -json       Output logs in JSON format (only applicable with -logs)
  -h, --help  Display this help message

Example:
  loadgen -health -security   Send health and security traces
  loadgen -logs -json         Write random quotes in JSON format to quotes.log`)
}

func main() {
	// Define flags
	baseFlag := flag.Bool("base", true, "Send base traces")
	healthFlag := flag.Bool("health", false, "Send health traces")
	securityFlag := flag.Bool("security", false, "Send security traces")
	logsFlag := flag.Bool("logs", false, "Enable logging of random quotes to quotes.log")
	jsonFlag := flag.Bool("json", false, "Output logs in JSON format (only applicable with -logs)")
	helpFlag := flag.Bool("h", false, "Display help message")
	helpFlagLong := flag.Bool("help", false, "Display help message")

	flag.Parse()

	// Display help and exit if -h or --help is provided
	if *helpFlag || *helpFlagLong {
		printHelp()
		os.Exit(0)
	}

	// Start logging if -logs flag is provided
	if *logsFlag {
		writeLogs(*jsonFlag)
		return
	}

	fmt.Println("Sending traces every 5 seconds. Use Ctrl-C to stop.")

	for {
		traceID := generateTraceID()
		spanID := generateSpanID()
		currentTime := getCurrentTime()
		endTime := currentTime + int64(time.Second)

		if *baseFlag {
			sendBaseTrace(traceID, spanID, currentTime, endTime)
		}

		if *healthFlag {
			time.Sleep(5 * time.Second)
			sendHealthTrace(traceID, generateSpanID(), getCurrentTime(), getCurrentTime()+int64(time.Second))
		}

		if *securityFlag {
			time.Sleep(5 * time.Second)
			sendSecurityTrace(traceID, generateSpanID(), getCurrentTime(), getCurrentTime()+int64(time.Second))
		}

		time.Sleep(5 * time.Second)
	}
}
