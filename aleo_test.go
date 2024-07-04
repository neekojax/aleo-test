package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestFetchRewardData(t *testing.T) {
	// Create a new Excel file
	f := excelize.NewFile()
	// Create a new sheet
	sheetName := "Sheet1"
	// Set the active sheet of the workbook
	//f.SetActiveSheet(index)

	// Write headers
	f.SetCellValue(sheetName, "A1", "Time")
	f.SetCellValue(sheetName, "B1", "Reward")

	pageCount := 1
	setPageCount := false

	// Iterate over the first 100 pages
	for page := 0; page <= pageCount; page++ {
		// Create the URL with the current page number
		url := fmt.Sprintf("https://testnetbeta.aleo123.io/api/v5/mainnetv0/blocks/list?page=%d&page_size=100", page)

		fmt.Printf("get %s\n", url)

		// Send the GET request
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			continue
		}

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		// Parse the JSON response
		var blockResponse BlockResponse
		err = json.Unmarshal(body, &blockResponse)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}

		if !setPageCount {

			pageCount = (blockResponse.Count-57700)/100 + 1
			fmt.Printf("count %d, page count: %d\n", blockResponse.Count, pageCount)
			setPageCount = true
		}

		// Write the time and reward to the Excel file
		for i, block := range blockResponse.BlockData {
			row := strconv.Itoa(page*100 + i + 2) // +2 because we start from row 2 and i is 0-indexed
			f.SetCellValue(sheetName, "A"+row, block.Time)
			f.SetCellValue(sheetName, "B"+row, block.Reward)
		}
	}

	// Save the Excel file
	if err := f.SaveAs("rewards.xlsx"); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}
