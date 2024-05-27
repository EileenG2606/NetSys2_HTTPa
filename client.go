package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main(){
	scanner:= bufio.NewScanner(os.Stdin)
	for{
		fmt.Println("1. Receive data (GET)")
		fmt.Println("2. Send data (POST)")
		fmt.Printf(">> ")
		scanner.Scan()
		choo := scanner.Text()
		if choo == "1"{
			receiveData()
		} else if choo == "2"{
			sendData()
		} else{
			fmt.Println("Please select 1 or 2 :)")
		}
	}
}

func loadTLSConfig() (*tls.Config, error){
	certPool:= x509.NewCertPool()
	certP,err:= os.ReadFile("cert.pem")
	if err != nil{
		return nil, err
	}
	certPool.AppendCertsFromPEM(certP)
	return &tls.Config{
		RootCAs: certPool,
	}, nil
}

var client= &http.Client{
	Timeout: 7*time.Second,
}

func receiveData(){
	ctx, cancel:= context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	tlsConf,err:= loadTLSConfig()
	if err!=nil{
		fmt.Println(err)
		return
	}

	// konfigurasi tls diimplementasikan dalan http
	client.Transport= &http.Transport{
		TLSClientConfig: tlsConf,
	}

	req,err:= http.NewRequestWithContext(ctx, http.MethodGet, "https://localhost:5678", nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	resp, err:= client.Do(req)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	dataRe,err:= io.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println(err)
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

	jsonData,err:= json.Marshal(data)
	if err!=nil{
		fmt.Println(err)
		return
	}

	contx,cancel:= context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tlsConfig, err := loadTLSConfig()
	if err!=nil{
		fmt.Println(err)
		return
	}

	client.Transport= &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	req, err:= http.NewRequestWithContext(contx, http.MethodPost, "https://localhost:5678/post", bytes.NewBuffer(jsonData))
	if err!=nil{
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	response,err:= client.Do(req)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	content,err:= io.ReadAll(response.Body)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println("Server response:", string(content))
}
