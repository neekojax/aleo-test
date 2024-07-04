package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
)

// 定义一个结构体来匹配JSON响应的结构
type BlockResponse struct {
	Success   bool        `json:"success"`
	Count     int         `json:"count"`
	Message   string      `json:"message"`
	BlockData []BlockData `json:"block_data"`
}

type BlockData struct {
	Network               int    `json:"network"`
	Round                 int    `json:"round"`
	Height                int    `json:"height"`
	CumulativeWeight      int    `json:"cumulative_weight"`
	CumulativeProofTarget int    `json:"cumulative_proof_target"`
	CoinbaseTarget        int    `json:"coinbase_target"`
	ProofTarget           int    `json:"proof_target"`
	LastCoinbaseTarget    int    `json:"last_coinbase_target"`
	LastCoinbaseTimestamp int    `json:"last_coinbase_timestamp"`
	Timestamp             int    `json:"timestamp"`
	Power                 int    `json:"power"`
	Reward                int    `json:"reward"`
	BlockReward           int    `json:"block_reward"`
	TargetTotal           int    `json:"TargetTotal"`
	BlockHash             string `json:"block_hash"`
	Epoch                 int    `json:"epoch"`
	Transactions          int    `json:"Transactions"`
	Solutions             int    `json:"Solutions"`
	Time                  string `json:"time"`
}

func main() {
	// 获取当前时间
	now := time.Now()

	// 定义开始时间
	startDate := time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC)

	// 初始化每天的奖励总和变量
	sums := make(map[string]int)

	// 打开Excel文件
	f, err := excelize.OpenFile("rewards.xlsx")
	if err != nil {
		fmt.Println("Error opening Excel file:", err)
		return
	}

	// 指定工作表名称
	sheetName := "Sheet1"

	// 获取工作表中的所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println("Error getting rows from sheet:", err)
		return
	}

	// 遍历所有行
	for i, row := range rows {
		// 跳过标题行（假设第一行是标题行）
		if i == 0 {
			continue
		}

		// 读取日期时间单元格值
		dateTimeValue := row[0] // 假设日期时间在第一列

		dateTimeValue = dateTimeValue[:len(dateTimeValue)-6]
		// 将 'T' 替换为空格
		dateTimeValue = strings.Replace(dateTimeValue, "T", " ", 1)

		// 解析日期时间并转换为时间戳
		dateTime, err := time.Parse("2006-01-02 15:04:05", dateTimeValue)
		if err != nil {
			fmt.Printf("Error parsing date-time from cell A%d: %s\n", i+1, err)
			continue
		}
		timestamp := dateTime.Unix()

		// 检查时间戳是否在指定时间段内
		for d := startDate; d.Before(now) || d.Equal(now); d = d.AddDate(0, 0, 1) {
			dayStart := d.Unix()
			dayEnd := d.AddDate(0, 0, 1).Unix()

			if timestamp >= dayStart && timestamp < dayEnd {
				// 读取奖励单元格值
				rewardValue := row[1] // 假设奖励在第二列

				// 转换奖励单元格值为整数
				intValue, err := strconv.Atoi(rewardValue)
				if err != nil {
					fmt.Printf("Error converting reward cell B%d value to integer: %s\n", i+1, err)
					continue
				}

				// 累加到对应天数的总和中
				key := d.Format("2006-01-02")
				sums[key] += intValue
			}
		}
	}

	totalReward := 0

	// 输出每天的总和结果
	for d := startDate; d.Before(now) || d.Equal(now); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		fmt.Printf("Reward on %s is: %.6f\n", key, float64(sums[key])/1e6)
		totalReward += sums[key]
	}

	fmt.Printf("Total Reward start on 2024/7/2 00:00:00 is: %.6f\n", float64(totalReward)/1e6)
}

// 辅助函数：将时间字符串转换为时间戳
func getTimeStamp(dateTimeStr string) int64 {
	dateTime, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	if err != nil {
		fmt.Println("Error parsing date-time:", err)
		return 0
	}
	return dateTime.Unix()
}
