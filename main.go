package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	zsw "github.com/zhongshuwen/zswchain-go"
	zswhq "github.com/zhongshuwen/zswchain-go/system"
	"github.com/zhongshuwen/zswchain-go/zswitems"
	"github.com/zhongshuwen/zswchain-go/zswperms"
)

var version = "dev"

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
func main() {

	var zswContentReviewTeamWalletName = zsw.AccountName("zsw.admin")

	var kexinJiedianAZhongShuWenUuid = "017f64cd-29e2-4957-b3ba-bd094e8ed233"
	var kexinJiedianAWalletName = zsw.AccountName("kxjdtest111a")

	//var userAZhongShuWenUuid = "017f5d8a-f70d-4602-b85f-b24751953e4d"
	var userAWalletName = zsw.AccountName("usertest111a")

	//var userBZhongShuWenUuid = "017f5d8a-f6f3-4594-833c-9a877e7af54b"
	var userBWalletName = zsw.AccountName("usertest111b")

	var collectionAUuid = "017a57ed-9dd4-42e1-aed4-0a8ced8a96d3"
	var itemA1Uuid = "017a57ed-b963-401e-9c2b-288a9dce6206"
	var itemA2Uuid = "017a57ed-b963-401e-9c2b-288a9dce6206"

	api := zsw.New("http://localhost:3031")

	keyBag := &zsw.KeyBag{}

	NoError(
		keyBag.ImportPrivateKeyFromEnv(context.Background(), "ZSW_CONTENT_REVIEW_PRIVATE_KEY"),
		"missing ZSW_CONTENT_REVIEW_PRIVATE_KEY",
	)

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

	var publicKeys, err = keyBag.AvailableKeys(context.Background())
	NoError(err, "Error getting public keys!")

	//var zswContentReviewTeamPublicKey = publicKeys[0]
	var kexinJiedianAPublicKey = publicKeys[1]
	var userAPublicKey = publicKeys[2]
	var userBPublicKey = publicKeys[3]

	api.SetSigner(keyBag)

	fmt.Println("创建可信节点账号...")
	runTxBasic(context.Background(), api, []*zsw.Action{

		//创建可信节点联盟链账号
		zswhq.NewNewAccount(
			zswContentReviewTeamWalletName, //中数文的内容审核管理账号
			kexinJiedianAWalletName,        //可信节点联盟链账号
			kexinJiedianAPublicKey,         //可信节点公钥
		),
		// 给可信节点Minting权限
		zswperms.NewSetZswPerms(
			zswContentReviewTeamWalletName, //中数文的内容审核管理账号
			kexinJiedianAWalletName,        //可信节点联盟链账号
			zsw.NewUint128FromUint64(
				uint64(zsw.ZSW_CORE_PERMS_CONFIRM_AUTHORIZE_USER_TX)| // 此权限赋予客户用户授权交易的权力
					uint64(zsw.ZSW_CORE_PERMS_CONFIRM_AUTHORIZE_USER_TRANSFER_ITEM), //允许可信节点赋予C2C基本数字藏品转移
			),
		),
		// 给可信节点自愿监护权限
		zswitems.NewMakeCustodian(
			zswContentReviewTeamWalletName,                    //中数文的内容审核管理账号
			kexinJiedianAWalletName,                           //平台生成的walletName
			UuidToUint128OrQuit(kexinJiedianAZhongShuWenUuid), //中数文平台的“userId”（登录借口获取的）
			zsw.NewUint128FromUint64(0),                       //现在0，没有用
			zsw.NewUint128FromUint64(
				uint64(zsw.CUSTODIAN_PERMS_ENABLED)| //开通Custodian功能
					uint64(zsw.CUSTODIAN_PERMS_TX_TO_SELF_CUSTODIAN)| //可以authorize用户在自己的
					uint64(zsw.CUSTODIAN_PERMS_SEND_TO_NULL_CUSTODIAN)| //can send from self custodianship to another custodian
					uint64(zsw.CUSTODIAN_PERMS_SEND_TO_ZSW_CUSTODIAN), //can send from self custodianship to a non-custodial null custodian
			),
			0, //0是征程
			0, //其他的可信节点用户要使用你的平台的时候，数字藏品要冻多久（秒）
			[]zsw.AccountName{
				kexinJiedianAWalletName, //为了查看历史方便，可以设置logevent账号，未来也可以加handler
			},
		),
		zswitems.NewMakeIssuer(
			zswContentReviewTeamWalletName,
			kexinJiedianAWalletName,
			UuidToUint128OrQuit(kexinJiedianAZhongShuWenUuid),
			zsw.NewUint128FromUint64(0),
			zsw.NewUint128FromUint64(
				uint64(zsw.ZSW_ITEMS_PERMS_AUTHORIZE_MINT_ITEM)| //允许基本minting的功能
					uint64(zsw.ZSW_ITEMS_PERMS_AUTHORIZE_MINT_TO_NULL_CUSTODIAN), //可以mint到需要用户公钥权限的custodian
			),
			0, //0==正常
		),
	})
	fmt.Println("创建用户账号...")
	runTxBasic(context.Background(), api, []*zsw.Action{
		//创建用户A账号
		zswhq.NewNewAccount(
			zswContentReviewTeamWalletName, //中数文的内容审核管理账号
			userAWalletName,                //用户A联盟链账号
			userAPublicKey,                 //用户A联盟链公钥
		),
		// 创建用户B账号
		zswhq.NewNewAccount(
			zswContentReviewTeamWalletName, //中数文的内容审核管理账号
			userBWalletName,                //用户A联盟链账号
			userBPublicKey,                 //用户A联盟链公钥
		),
	})
	fmt.Println("登记新版税接受者（%s）+ 开通“平台发行方”类权限", kexinJiedianAWalletName)
	runTxBasic(context.Background(), api, []*zsw.Action{
		zswitems.NewMakeRoyaltyUser( //登记谁是版税接受者
			zswContentReviewTeamWalletName,
			kexinJiedianAWalletName,
			UuidToUint128OrQuit(kexinJiedianAZhongShuWenUuid),
			zsw.NewUint128FromUint64(0),
			0,
		),
	})
	fmt.Println("登记版税接受者", kexinJiedianAWalletName)
	runTxBasic(context.Background(), api, []*zsw.Action{
		zswitems.NewMakeRoyaltyUser( //登记谁是版税接受者
			zswContentReviewTeamWalletName,
			kexinJiedianAWalletName,
			UuidToUint128OrQuit(kexinJiedianAZhongShuWenUuid),
			zsw.NewUint128FromUint64(0),
			0,
		),
	})
	fmt.Println("创建新的Collection", kexinJiedianAWalletName)
	runTxBasic(context.Background(), api, []*zsw.Action{
		zswitems.NewMakeCollection( //登记谁是版税接受者
			zswContentReviewTeamWalletName,
			kexinJiedianAWalletName,
			kexinJiedianAWalletName,
			UuidToUint128OrQuit(collectionAUuid).GetTypeACode(),
			UuidToUint128OrQuit(collectionAUuid).GetTypeBCode(),
			1,   // 正常 == 1
			11,  // 正常 == 11
			350, // 一级市场分润：10000分之多少，比如350==3.5%，525==5.25%， 900==9%，等
			525, // 二级市场分润：10000分之多少，比如350==3.5%，525==5.25%， 900==9%，等
			"",  //链上metadata schema name, 未来会支持
			"https://metadata.demo.zhongshuwen.com/metadata/collections/my-metadata.json", //collection metadata url
			kexinJiedianAWalletName, //royalty receiver，链上记录为了透明化
			[]zsw.AccountName{},     // notify账号/合约（为了历史或者你们内部系统方便）
			zsw.ZswItemsMetadata{},  //目前不支持链上metadata，很快会出方案
		),
	})
	fmt.Println("创建Collection与两个数字藏品模版", kexinJiedianAWalletName)
	runTxBasic(context.Background(), api, []*zsw.Action{
		zswitems.NewMakeCollection( //创建新的collection
			zswContentReviewTeamWalletName,
			kexinJiedianAWalletName,                             //创建者
			kexinJiedianAWalletName,                             //平台Issuer
			UuidToUint128OrQuit(collectionAUuid).GetTypeACode(), // API获取的Collecion UUID的A类码
			UuidToUint128OrQuit(collectionAUuid).GetTypeBCode(), // API获取的Collecion UUID的B类码
			1,   // 正常 == 1
			9,   // 正常 == 9 (ITEM_CONFIG_TRANSFERABLE | ITEM_CONFIG_ALLOW_NOTIFY)
			350, // 一级市场分润：10000分之多少，比如350==3.5%，525==5.25%， 900==9%，等
			525, // 二级市场分润：10000分之多少，比如350==3.5%，525==5.25%， 900==9%，等
			"",  //链上metadata schema name, 未来会支持
			"https://metadata.demo.zhongshuwen.com/metadata/collections/my-metadata.json", //collection metadata url
			kexinJiedianAWalletName, //royalty receiver，链上记录为了透明化
			[]zsw.AccountName{},     // notify账号/合约（为了历史或者你们内部系统方便）
			zsw.ZswItemsMetadata{},  //目前不支持链上metadata，很快会出方案
		),
		zswitems.NewMakeItem(
			zswContentReviewTeamWalletName,
			kexinJiedianAWalletName,                      //创建者
			kexinJiedianAWalletName,                      //平台Issuer
			UuidToUint128OrQuit(itemA1Uuid).Get40BitId(), //api获取的item metadata uuid
			UuidToUint128OrQuit(itemA1Uuid),              //api获取的item metadata uuid
			1,                                            // 正常 == 1
			9,                                            // 正常 == 9 (ITEM_CONFIG_TRANSFERABLE | ITEM_CONFIG_ALLOW_NOTIFY)

			UuidToUint128OrQuit(collectionAUuid).GetTypeACode(), // API获取的Collecion UUID的A类码
			5000, //供应上限=5000
			"https://metadata.demo.zhongshuwen.com/metadata/items/1.json", //item metadata url
			"",                     //custom链上metadata schema，未来会支持
			zsw.ZswItemsMetadata{}, //定制链上metadata，未来会支持

		),
		zswitems.NewMakeItem(
			zswContentReviewTeamWalletName,
			kexinJiedianAWalletName,                      //创建者
			kexinJiedianAWalletName,                      //平台Issuer
			UuidToUint128OrQuit(itemA2Uuid).Get40BitId(), //api获取的item metadata uuid
			UuidToUint128OrQuit(itemA2Uuid),              //api获取的item metadata uuid
			1,                                            // 正常 == 1
			9,                                            // 正常 == 9 (ITEM_CONFIG_TRANSFERABLE | ITEM_CONFIG_ALLOW_NOTIFY)

			UuidToUint128OrQuit(collectionAUuid).GetTypeACode(), // API获取的Collecion UUID的A类码
			300000, //供应上限=300000
			"https://metadata.demo.zhongshuwen.com/metadata/items/2.json", //collection metadata url
			"",                     //custom链上metadata schema，未来会支持
			zsw.ZswItemsMetadata{}, //定制链上metadata，未来会支持

		),
	})

	fmt.Println("mint数字藏品给用户A", kexinJiedianAWalletName)
	runTxBasic(context.Background(), api, []*zsw.Action{
		zswitems.NewItemMint( //创建新的collection
			kexinJiedianAWalletName, // minter/平台Issuer
			userAWalletName,         // receiver
			kexinJiedianAWalletName, //custodian，如果是用户自己用："nullnullnull"
			0,                       // T+X秒 （用户得到藏品之后需要等多少秒才可以转移/交易）

			// mint 52藏品A1给用户，mint 95藏品A2给用户
			[]uint64{UuidToUint128OrQuit(itemA1Uuid).Get40BitId(), UuidToUint128OrQuit(itemA2Uuid).Get40BitId()}, //数字藏品ids
			[]uint64{52, 95}, //量

			"memo string", //memo
		),
	})
}
