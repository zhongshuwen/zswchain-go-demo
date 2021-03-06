package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/ecc"
	"github.com/zhongshuwen/zswchain-go/zswitems"
)

func RunDebugScenarioA(ctx context.Context, api *zsw.API, creator zsw.AccountName, newKexinJiedian zsw.AccountName, newKexinJiedianPublicKey string, newKexinJiedianZswId string) (error, string) {
	actions := []*zsw.Action{}
	actions = append(actions, GetActionsCreateUserWithResources(
		creator,
		newKexinJiedian,
		newKexinJiedianPublicKey,
		1000000,
		zsw.NewZSWAsset(10000),
		zsw.NewZSWAsset(2000),
	)...)
	actions = append(actions, GetActionsSetupKexinJiedianPermissions(
		creator,
		newKexinJiedian,
		newKexinJiedianZswId,
	)...)

	return nil, runTxBasic(context.Background(), api, actions)
}

func RunDebugScenarioB(ctx context.Context, api *zsw.API, authorizer zsw.AccountName, kexinJiedian zsw.AccountName, kexinJiedianPublicKey string, kexinJiedianZswId string, recipientUser zsw.AccountName) (error, string) {
	collectionZswId := uuid.New().String()
	itemTemplateZswId := uuid.New().String()
	itemZswId := uuid.New().String()

	actions := GetActionsCreateExampleCollectionItemFlow1155(
		authorizer,
		kexinJiedian,
		collectionZswId,
		itemTemplateZswId,
		itemZswId,
	)
	mintActions := []*zsw.Action{
		zswitems.NewItemMint(
			kexinJiedian,
			recipientUser,
			kexinJiedian,
			0,
			[]uint64{UuidToUint128OrQuit(itemZswId).Get40BitId()},
			[]uint64{1},
			"An item for you!",
		),
	}

	actions = append(actions, mintActions...)

	return nil, runTxBasic(context.Background(), api, actions)
}
func getRandUUID() string {
	b := make([]byte, 4)

	rand.Read(b)
	return fmt.Sprintf("00000000-0000-0000-0000-0000%s", hex.EncodeToString(b))
	//rand.Intn(100)

}

func RunDebugScenarioC(ctx context.Context, authorizer zsw.AccountName, newKexinJiedian zsw.AccountName) (error, string) {

	api := zsw.New(os.Getenv("ZSW_API_URL"))
	api.Debug = true

	keyBag := &zsw.KeyBag{}

	NoError(
		keyBag.ImportPrivateKeyFromEnv(context.Background(), "ZSW_CONTENT_REVIEW_PRIVATE_KEY"),
		"missing ZSW_CONTENT_REVIEW_PRIVATE_KEY",
	)
	kexinJiedianZswId := uuid.New().String()
	collectionZswId := uuid.New().String()
	itemTemplateZswId := uuid.New().String()
	itemZswId := uuid.New().String()

	kxjdPrivateKey, err := ecc.NewPrivateKey(os.Getenv("KEXIN_JIEDIAN_A_PRIVATE_KEY"))
	if err != nil {
		return err, ""
	}
	userAPrivateKey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return err, ""
	}
	userBPrivateKey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return err, ""
	}
	keyBag.Append(kxjdPrivateKey)

	userAName := zsw.AccountName(fmt.Sprintf("usra1%s", RandomLowercaseStringAZ(7)))
	userBName := zsw.AccountName(fmt.Sprintf("usrb1%s", RandomLowercaseStringAZ(7)))

	fmt.Printf("-- ??????????????? --\n?????????%s\n?????????%s\n?????????%s\n--------------------------------------------------\n", newKexinJiedian, kxjdPrivateKey.PublicKey().String(), kxjdPrivateKey.String())
	fmt.Printf("-- ??????A --\n?????????%s\n?????????%s\n?????????%s\n--------------------------------------------------\n", userAName, userAPrivateKey.PublicKey().String(), userAPrivateKey.String())
	fmt.Printf("-- ??????B --\n?????????%s\n?????????%s\n?????????%s\n--------------------------------------------------\n", userBName, userBPrivateKey.PublicKey().String(), userBPrivateKey.String())

	api.SetSigner(keyBag)
	actions := []*zsw.Action{}
	actions = append(actions, GetActionsCreateUserWithResources(
		authorizer,
		newKexinJiedian,
		kxjdPrivateKey.PublicKey().String(),
		1000000,
		zsw.NewZSWAsset(10000),
		zsw.NewZSWAsset(2000),
	)...)
	actions = append(actions, GetActionsCreateUserWithResources(
		authorizer,
		userAName,
		userAPrivateKey.PublicKey().String(),
		3000,
		zsw.NewZSWAsset(0),
		zsw.NewZSWAsset(0),
	)...)
	actions = append(actions, GetActionsCreateUserWithResources(
		authorizer,
		userBName,
		userBPrivateKey.PublicKey().String(),
		3000,
		zsw.NewZSWAsset(0),
		zsw.NewZSWAsset(0),
	)...)
	runTxBasic(context.Background(), api, actions)

	time.Sleep(time.Second * 2)
	actions2 := []*zsw.Action{}
	actions2 = append(actions2, GetActionsSetupKexinJiedianPermissions(
		authorizer,
		newKexinJiedian,
		kexinJiedianZswId,
	)...)
	runTxBasic(context.Background(), api, actions2)
	time.Sleep(time.Second * 2)
	actions2 = GetCreateExampleCollection(
		authorizer,
		newKexinJiedian,
		collectionZswId,
	)
	actions2 = append(actions2, GetActionsCreateExampleCollectionItemFlow1155(
		authorizer,
		newKexinJiedian,
		collectionZswId,
		itemTemplateZswId,
		itemZswId,
	)...)
	mintActions := []*zsw.Action{

		zswitems.NewItemMint(
			newKexinJiedian,
			userAName,
			newKexinJiedian,
			0,
			[]uint64{UuidToUint128OrQuit(itemZswId).Get40BitId()},
			[]uint64{100},
			"An item for you!",
		),
		zswitems.NewItemMint(
			newKexinJiedian,
			userBName,
			newKexinJiedian,
			0,
			[]uint64{UuidToUint128OrQuit(itemZswId).Get40BitId()},
			[]uint64{200},
			"An item for you 2!",
		),
	}
	actions2 = append(actions2, mintActions...)

	//actions = append(actions, mintActions...)

	return nil, runTxBasic(context.Background(), api, actions2)

}
