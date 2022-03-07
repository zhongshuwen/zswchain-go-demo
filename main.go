package main

import (
	"fmt"
	"os"
	"context"
	"encoding/json"

	zsw "github.com/zhongshuwen/zswchain-go"
)

var version = "dev"



func NoError(err error, message string, args ...interface{}) {
	if err != nil {
		Quit(message+": "+err.Error(), args...)
	}
}

func Quit(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	os.Exit(1)
}


func main() {
	api := zsw.New("http://localhost:3031")
	ctx := context.Background()

	infoResp, err := api.GetInfo(ctx)
	NoError(err, "unable to get chain info")

	fmt.Println("Chain Info", toJson(infoResp))

	accountResp, _ := api.GetAccount(ctx, "zsw.admin")
	fmt.Println("Account Info", toJson(accountResp))
}

func toJson(v interface{}) string {
	out, err := json.MarshalIndent(v, "", "  ")
	NoError(err, "unable to marshal json")

	return string(out)
}