package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hungtrd/bill-calculator/pkg/gemini"
	"github.com/joho/godotenv"
)

var prompt = `
Bạn là một thuật toán OCR để chuyển chữ viết tay thành văn bản có độ chính xác cao.
Đây là bức ảnh một hoá đơn viết bằng Tiếng Việt, bao gồm tên mặt hàng, số lượng và đơn giá được viết bên cạnh.
Hãy đọc thông tin từ bức ảnh và đưa ra dữ liệu dưới dạng bảng gồm 3 cột: tên mặt hàng, số lượng và đơn giá.
Xuất dữ liệu dưới định dạng CSV với dấu phẩy (,) là ký tự phân cách. Nếu số thập phân, sử dụng dấu chấm (.) để ngăn cách phần nguyên và phần thập phân. đơn giá có thể là số thập phân
Ví dụ về định dạng CSV:
    Tên mặt hàng,Số lượng,Đơn giá
    Mặt hàng 1,1.00,100.00
    Mặt hàng 2,2.00,200.50
    Mặt hàng 3,3.00,300.75
`

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := gemini.NewClient(gemini.Model1_5Pro)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	resp, err := client.GenerateContentWithFile(context.Background(), "./tmp/files/bill_4.JPG", prompt)
	if err != nil {
		log.Fatal(err)
	}

	bill := ParseBillData(resp)
	fmt.Println(bill.Total)
}

type Bill struct {
	Name     string
	Customer string
	Items    []Item
	Total    float64
}

type Item struct {
	Name     string
	Quantity float64
	Price    float64
}

func ParseBillData(csvData string) Bill {
	fmt.Println(csvData)
	csvData = normalizeData(csvData)
	strReader := strings.NewReader(csvData)
	reader := csv.NewReader(strReader)

	records, err := reader.ReadAll()
	// Checks for the error
	if err != nil {
		fmt.Println("Error reading records")
	}

	items := []Item{}

	var total float64 = 0
	for _, row := range records {
		if len(row) >= 3 {
			quantity, err := strconv.ParseFloat(row[1], 64)
			if err != nil {
				fmt.Println("Error parsing quantity value:", row[1])
				continue
			}
			price, err := strconv.ParseFloat(row[2], 64)
			if err != nil {
				fmt.Println("Error parsing price value:", row[2])
				continue
			}
			item := Item{
				Name:     row[0],
				Quantity: quantity,
				Price:    price,
			}
			items = append(items, item)
			total += handleRowsData(item)
		}
	}

	return Bill{Items: items, Total: total}
}

func handleRowsData(item Item) float64 {
	total := item.Quantity * item.Price
	fmt.Printf("%s: %.2f x %.2f = %.2f\n", item.Name, item.Quantity, item.Price, total)

	return total
}

func normalizeData(data string) string {
	data = strings.ReplaceAll(data, "```csv", "")
	data = strings.ReplaceAll(data, "```", "")
	return data
}
