package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	zsw "github.com/zhongshuwen/zswchain-go"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandomLowercaseStringAZ(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {

	//zsw.EnableDebugLogging(zsw.NewLogger(false))

	rand.Seed(time.Now().UnixNano())

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "fullprocess",
			Aliases: []string{"full"},
			Usage:   "创建新用户，创建schema（链上metadata）创建collection，创建itemtpl创建item，mintingitem",
			Action: func(c *cli.Context) error {
				err, _ := RunDebugScenarioC(c.Context, "zsw.admin", zsw.AccountName(RandomLowercaseStringAZ(12)))
				return err
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
