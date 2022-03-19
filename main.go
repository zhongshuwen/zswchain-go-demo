package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"math/rand"

	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/zswitems"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var apis = []*zsw.API{
	zsw.New("https://node3.tn1.chao7.cn"),
	zsw.New("https://node4.tn1.chao7.cn"),
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var version = "dev"

var kexinJiedianAWalletName = zsw.AccountName("kxjdtestgm1a")

//var userAZhongShuWenUuid = "017f5d8a-f70d-4602-b85f-b24751953e4d"
var userAWalletName = zsw.AccountName("usertestgm1c")

//var userBZhongShuWenUuid = "017f5d8a-f6f3-4594-833c-9a877e7af54b"
var userBWalletName = zsw.AccountName("usertestgm1d")

func Quit(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	os.Exit(1)
}

func NoError(err error, message string, args ...interface{}) {
	if err != nil {
		Quit(message+": "+err.Error(), args...)
	}
}

func toJson(v interface{}) string {
	out, err := json.MarshalIndent(v, "", "  ")
	NoError(err, "unable to marshal json")

	return string(out)
}

func runTxBasic(ctx context.Context, api *zsw.API, actions []*zsw.Action) string {

	txOpts := &zsw.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		panic(fmt.Errorf("filling tx opts: %w", err))
	}

	tx := zsw.NewTransaction(actions, txOpts)
	signedTx, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, zsw.CompressionNone)
	if err != nil {
		panic(fmt.Errorf("sign transaction: %w", err))
	}

	content, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		panic(fmt.Errorf("json marshalling transaction: %w", err))
	}

	fmt.Println(string(content))
	fmt.Println()

	response, err := api.PushTransaction(context.Background(), packedTx)
	if err != nil {
		panic(fmt.Errorf("push transaction: %w", err))
	}

	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
	return hex.EncodeToString(response.Processed.ID)
}
func UuidToUint128OrQuit(uuidString string) zsw.Uint128 {
	var x zsw.Uint128
	NoError(x.FromUuidString(uuidString), "Invalid uuid: '%s'", uuidString)
	return x

}

type ItemBalanceTableRow struct {
	ItemId           uint64 `json:"item_id"`
	Status           uint32 `json:"status"`
	Balance          uint64 `json:"balance"`
	BalanceInCustody uint64 `json:"balance_in_custody"`
	BalanceFrozen    uint64 `json:"balance_frozen"`
}

func QueryUserCangpin(ctx context.Context, api *zsw.API, account zsw.AccountName) (out *[]ItemBalanceTableRow, errOut error) {
	var rowReq = zsw.GetTableRowsRequest{
		Code:       "zsw.items",
		Scope:      string(account),
		Table:      "itembalances",
		LowerBound: "", //use this to paginate with the last result's id
		Limit:      10, //results to fetch
		JSON:       true,
	}
	var resp, err = api.GetTableRows(ctx, rowReq)
	if err != nil {
		errOut = err
		return
	}
	var x []ItemBalanceTableRow
	resp.JSONToStructs(&x)
	return &x, nil
}
func MintOneRand(api *zsw.API, ch chan<- string) {
	ch <- runTxBasic(context.Background(), api, []*zsw.Action{
		zswitems.NewItemMint( //mint数字藏品
			kexinJiedianAWalletName, // minter/平台Issuer
			userAWalletName,         // receiver
			kexinJiedianAWalletName, //custodian，如果是用户自己用："nullnullnull"
			0,                       // T+X秒 （用户得到藏品之后需要等多少秒才可以转移/交易）

			// mint 52藏品A1给用户，mint 95藏品A2给用户
			[]uint64{5622}, //数字藏品ids
			[]uint64{1},    //量

			"random: "+randSeq(10), //memo
		),
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	ch := make(chan string)
	/*

	   {
	     "account": "zsw.items",
	     "name": "mint",
	     "data": {
	       "minter": "kxjdtestgm1a",
	       "to": "usertestgm1d",
	       "to_custodian": "kxjdtestgm1a",
	       "item_ids": [5622],
	       "amounts": [2000],
	       "memo":"here are your products",
	       "freeze_time": 0
	     }
	   }
	*/

	keyBag := &zsw.KeyBag{}
	NoError(
		keyBag.ImportPrivateKeyFromEnv(context.Background(), "KEXIN_JIEDIAN_A_PRIVATE_KEY"),
		"missing KEXIN_JIEDIAN_A_PRIVATE_KEY",
	)

	NoError(
		keyBag.ImportPrivateKeyFromEnv(context.Background(), "USER_A_PRIVATE_KEY"),
		"missing USER_A_PRIVATE_KEY",
	)
	NoError(
		keyBag.ImportPrivateKeyFromEnv(context.Background(), "USER_B_PRIVATE_KEY"),
		"missing USER_B_PRIVATE_KEY",
	)

	lenApis := len(apis)
	for _, api := range apis {
		api.SetSigner(keyBag)
	}

	fmt.Println("mint item")
	for i := 0; i < 100; i++ {
		go MintOneRand(apis[i%lenApis], ch)
		time.Sleep(50 * time.Millisecond)
	}
	for i := 0; i < 100; i++ {
		fmt.Printf("result %d %s", i, <-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
