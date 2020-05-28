package main

import (
	"encoding/json"
	"fmt"
)

type TaskResponse struct {
	Status int `json:"status"`
	Msg string `json:"msg"`
	Data []Data `json:"data"`
}

type Data struct{
	Task int `json:"task"`
	Interval string `json:"interval"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	Data string `json:"data"`
	More string `json:"more"`
	From string `json:"from"`
	Status int `json:"status"`
	Msg string `json:"msg"`
}

func main(){
	Json := "{\"status\":200,\"msg\":\"ok\",\"data\":[\"{\\\"task\\\":553389713283,\\\"interval\\\":\\\"\\\",\\\"start_time\\\":\\\"\\\",\\\"end_time\\\":\\\"\\\",\\\"data\\\":\\\"\\\",\\\"more\\\":\\\"\\\",\\\"from\\\":\\\"\\\",\\\"status\\\":200,\\\"msg\\\":\\\"OK\\\"}\"]}\n"

	a := &TaskResponse{}

	json.Unmarshal([]byte(Json),a)

	fmt.Println(a)

	//data := a.Data[0]
	//
	//b := &Data{}
	//
	//json.Unmarshal([]byte(data),b)
	//
	//fmt.Println(b)
}
