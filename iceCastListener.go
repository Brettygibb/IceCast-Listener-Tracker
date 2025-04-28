package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Stats struct {
	Icecast struct {
		Source struct {
			Listeners int `json:"listeners"`
		} `json:"source"`
	} `json:"icestats"`
}

var client = &http.Client{
	Timeout: 5 * time.Second,
}

const csvFile = "C:/Users/user/Desktop/BGRadio Fredericton/Go FIles/icecast.csv"

func main() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	currentDate := time.Now().Format("2006-01-02")
	hourlyData := make([][]int, 24) // 24 hours, each with []int of samples

	for range ticker.C {
		now := time.Now()
		hour := now.Hour()
		date := now.Format("2006-01-02")

		// Fetch Icecast data
		resp, err := client.Get("http://localhost:8000/status-json.xsl")
		if err != nil {
			fmt.Println(now.Format(time.RFC3339), "Error fetching data:", err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var stats Stats
		if err := json.Unmarshal(body, &stats); err != nil {
			fmt.Println(now.Format(time.RFC3339), "Error parsing JSON:", err)
			continue
		}

		listeners := stats.Icecast.Source.Listeners
		fmt.Printf("%s - Listeners: %d\n", now.Format("15:04"), listeners)

		// Store the listener count
		hourlyData[hour] = append(hourlyData[hour], listeners)

		// If it's midnight, flush yesterday's data to file
		if hour == 0 && now.Minute() == 0 && currentDate != date {
			writeCSVRow(currentDate, hourlyData)
			hourlyData = make([][]int, 24) // reset
			currentDate = date
		}
	}
}

func writeCSVRow(date string, data [][]int) {
	// Load or create file
	file, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Read existing rows
	reader := csv.NewReader(file)
	existingRows, _ := reader.ReadAll()

	// Create header if needed
	if len(existingRows) == 0 {
		header := []string{"Date"}
		for i := 0; i < 24; i++ {
			header = append(header, fmt.Sprintf("%02d:00", i))
		}
		existingRows = append(existingRows, header)
	}

	// Build new row
	newRow := []string{date}
	for _, hourData := range data {
		if len(hourData) == 0 {
			newRow = append(newRow, "")
			continue
		}
		sum := 0
		for _, v := range hourData {
			sum += v
		}
		avg := float64(sum) / float64(len(hourData))
		newRow = append(newRow, strconv.FormatFloat(avg, 'f', 2, 64))
	}

	existingRows = append(existingRows, newRow)

	// Reset file, write all
	file.Truncate(0)
	file.Seek(0, 0)
	writer := csv.NewWriter(file)
	writer.WriteAll(existingRows)
	writer.Flush()

	fmt.Println("✅ Wrote row for:", date)
}
