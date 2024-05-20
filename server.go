package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main(){
	mux:= http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Yes, This is from server")
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		data, err:= io.ReadAll(r.Body)
		if err!=nil{
			return
		}
		var requestData map[string]string
		err = json.Unmarshal(data, &requestData)
		if err!=nil{
			return
		}
		fmt.Println("Received data: ", requestData)
		fmt.Fprint(w, "Data successfully received by server")
	})

	server := http.Server{
		Addr:	"localhost:5678",
		Handler: mux,
	}
	
	err := server.ListenAndServe()
	if err!=nil{
		return
	}
}
