package main

import (
	"encoding/json"
	"fmt"
)

type Job struct {
	Name string `json:"name"`
	CronExpr string `json:"cronExpr"`
	Msg string `json:"msg"`
}
func main() {
	var (
		str string
		job Job
		err error
		jsonbytes []byte
	)
	str = `{"name":"test", "cronExpr":"* * * * *", "msg":"success"}`
	if err = json.Unmarshal([]byte(str), &job); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(job)

	if jsonbytes, err = json.Marshal(job); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonbytes))
}
