package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1. Receive data (GET)")
		fmt.Println("2. Send data (POST)")
		fmt.Printf(">> ")
		scanner.Scan()
		choo := scanner.Text()
		if choo == "1"{
			receiveData()
		} else if choo == "2"{
			sendData()
		} else {
			fmt.Println("Please select 1 or 2 :)")
		}
	}
}

var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
	},
	Timeout: 7 * time.Second,
}

func receiveData(){
	ctx, cancel:= context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	req, err:= http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:5678", nil)
	if err != nil{
		return
	}

	resp, err:= client.Do(req)
	if err != nil{
		return
	}
	defer resp.Body.Close()

	dataRe, err:= io.ReadAll(resp.Body)
	if err != nil{
		return
	}

	fmt.Println("Info from server:", string(dataRe))
}

func sendData(){
	data := map[string]string{
		"Name":  "Budi",
		"Age":   "50",
		"Email": "BB@gmail.com",
	}

	jsonData, err:= json.Marshal(data)
	if err!=nil{
		return
	}

	contx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err:= http.NewRequestWithContext(contx, http.MethodPost, "http://localhost:5678/post", bytes.NewBuffer(jsonData))
	if err!=nil{
		return
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err!=nil{
		return
	}
	defer response.Body.Close()

	content, err:= io.ReadAll(response.Body)
	if err!=nil{
		return
	}

	fmt.Println("Server response:", string(content))
}
