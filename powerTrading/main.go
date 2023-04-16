package main

import (
	"fmt"
	"os"
	"powerTrading/sdkInit"
	"powerTrading/service"
	"powerTrading/web"
	"powerTrading/web/controller"
	"powerTrading/web/cryptoCode"
	"powerTrading/web/global"
	"powerTrading/web/groupSignature"
	"powerTrading/web/model"
	"powerTrading/web/tools"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

func main() {
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/powerTrading/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/powerTrading/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/powerTrading/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    os.Getenv("GOPATH") + "/src/powerTrading/chaincode",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	//invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")
	//producer := service.User{
	//	Id:       "431122200103110052",
	//	Purchases:   nil,
	//	Sells:    nil,
	//}
	//producer.Surplus= cryptoCode.HomoEncryptData(strconv.FormatInt(0,10))
	//user1:=&model.User{
	//	UserId: producer.Id,
	//}
	//user1.GroupPrivateKey=groupSignature.GenerateMemberPrivateKey()
	//consumer := service.User{
	//	Id:       "431122200203110052",
	//	Purchases:   nil,
	//	Sells: nil,
	//}
	//user2:=&model.User{
	//	UserId: consumer.Id,
	//}
	//user2.GroupPrivateKey=groupSignature.GenerateMemberPrivateKey()
	//consumer.Surplus= cryptoCode.HomoEncryptData("100")
	//asset1 := service.Electricity{
	//	Id:              "123",
	//	Type:            "光能发电",
	//	State:           true                                                                                                                                                                                                                                                  ,
	//	CurrentOwnerId:  "431122200103110052",
	//	Date:            "2023-2-2 19:55",
	//}
	//asset1.Scope= cryptoCode.HomoEncryptData("100")
	//asset1.Price= cryptoCode.HomoEncryptData("2")
	//asset2 := service.Electricity{
	//	Id:              "456",
	//	Type:            "光能发电",
	//	State:           true                                                                                                                                                                                                                                                  ,
	//	CurrentOwnerId:  "431122200103110052",
	//	Date:            "2023-2-3 19:55",
	//}
	//asset2.Scope= cryptoCode.HomoEncryptData("200")
	//asset2.Price= cryptoCode.HomoEncryptData("3")
	//order1:=service.Electricity{
	//	Id: "abc",
	//	Date:            "2023-2-1 19:55",
	//	State:           true,
	//	CurrentOwnerId:  "431122200203110052",
	//}
	//order1.Scope= cryptoCode.HomoEncryptData("350")
	//order1.Price= cryptoCode.HomoEncryptData("1")
	//order2:=service.Electricity{
	//	Id: "def",
	//	Date:            "2023-2-3 20:55",
	//	State:           true,
	//	CurrentOwnerId:  "431122200203110052",
	//}
	//order2.Scope= cryptoCode.HomoEncryptData("100")
	//order2.Price= cryptoCode.HomoEncryptData("1")
	Org1ServiceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	Org2ServiceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[1], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	// controller.ConsumerApp = controller.Application{
	// 	Setup: Org1ServiceSetup,
	// }
	// controller.ProducerApp = controller.Application{
	// 	Setup: Org2ServiceSetup,
	// }

	controller.InitApp(Org1ServiceSetup, Org2ServiceSetup)

	producer1 := service.User{
		Id:        tools.GenUUID(),
		Purchases: nil,
		Sells:     nil,
	}
	producer1.Surplus = cryptoCode.HomoEncryptData("0")
	pUser1 := &model.User{
		UserId:    producer1.Id,
		Username:  "b1",
		Password:  "1",
		CreatedAt: "2022-01-01 12:30:11",
		Role:      "供电商",
	}
	pUser1.GroupPrivateKey = groupSignature.GenerateMemberPrivateKey()

	producer2 := service.User{
		Id:        tools.GenUUID(),
		Purchases: nil,
		Sells:     nil,
	}
	producer2.Surplus = cryptoCode.HomoEncryptData("0")
	pUser2 := &model.User{
		UserId:    producer2.Id,
		Username:  "b2",
		Password:  "1",
		CreatedAt: "2022-06-01 12:30:11",
		Role:      "供电商",
	}
	pUser2.GroupPrivateKey = groupSignature.GenerateMemberPrivateKey()

	producer3 := service.User{
		Id:        tools.GenUUID(),
		Purchases: nil,
		Sells:     nil,
	}
	producer3.Surplus = cryptoCode.HomoEncryptData("0")
	pUser3 := &model.User{
		UserId:    producer3.Id,
		Username:  "b3",
		Password:  "1",
		CreatedAt: "2022-06-02 12:30:11",
		Role:      "供电商",
	}
	pUser3.GroupPrivateKey = groupSignature.GenerateMemberPrivateKey()

	//初始化购电用户
	consumer1 := service.User{
		Id:        tools.GenUUID(),
		Purchases: nil,
		Sells:     nil,
	}
	cUser1 := &model.User{
		UserId:    consumer1.Id,
		Username:  "a1",
		Password:  "1",
		CreatedAt: "2022-06-04 12:30:11",
		Role:      "购电用户",
	}
	consumer1.Surplus = cryptoCode.HomoEncryptData("1000")
	cUser1.GroupPrivateKey = groupSignature.GenerateMemberPrivateKey()

	//初始化发布的资产
	asset1 := service.Electricity{
		Id:             tools.GenUUID(),
		Type:           "光伏发电",
		State:          true,
		CurrentOwnerId: producer1.Id,
		Date:           "2023-01-02 19:55:00",
	}

	asset2 := service.Electricity{
		Id:             tools.GenUUID(),
		Type:           "光伏发电",
		State:          true,
		CurrentOwnerId: producer2.Id,
		Date:           "2023-01-03 19:55:00",
	}

	asset3 := service.Electricity{
		Id:             tools.GenUUID(),
		Type:           "光伏发电",
		State:          true,
		CurrentOwnerId: producer3.Id,
		Date:           "2023-01-04 19:55:00",
	}

	//初始化发布的订单
	order1 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-01-01 19:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	order2 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-01-02 20:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	order3 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-01-03 20:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	order4 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-02-01 19:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	order5 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-02-02 20:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	order6 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-02-03 20:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	order7 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-03-02 20:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	order8 := service.Electricity{
		Id:             tools.GenUUID(),
		Date:           "2023-03-03 20:55:00",
		State:          true,
		CurrentOwnerId: consumer1.Id,
	}

	fmt.Println("发电用户注册...")
	msg, err := Org2ServiceSetup.ProducerRegister(producer1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}
	fmt.Println("注册成功！")
	global.InsertUser(pUser1)
	fmt.Println(pUser1)

	msg, err = Org2ServiceSetup.ProducerRegister(producer2)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}
	fmt.Println("注册成功！")
	global.InsertUser(pUser2)
	fmt.Println(pUser2)

	msg, err = Org2ServiceSetup.ProducerRegister(producer3)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}
	fmt.Println("注册成功！")
	global.InsertUser(pUser3)
	fmt.Println(pUser3)

	fmt.Println("用电用户注册...")
	msg, err = Org1ServiceSetup.ConsumerRegister(consumer1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}
	fmt.Println("注册成功！")
	global.InsertUser(cUser1)
	fmt.Println(cUser1)

	fmt.Println("发电用户发布第一笔电力资产...")
	_, err = controller.InitLaunchAsset(controller.ProducerApp, producer1.Id, asset1.Id, "1000", "5", asset1.Date)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("发布成功...")

	fmt.Println("发电用户发布第二笔电力资产...")
	_, err = controller.InitLaunchAsset(controller.ProducerApp, producer2.Id, asset2.Id, "800", "7", asset2.Date)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("发布成功...")

	fmt.Println("发电用户发布第三笔电力资产...")
	_, err = controller.InitLaunchAsset(controller.ProducerApp, producer3.Id, asset3.Id, "600", "9", asset3.Date)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("发布成功...")

	fmt.Println("用电用户发布第一笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order1.Id, "100", "7", order1.Date)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("发布成功...")
	fmt.Println("用电用户发布第二笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order2.Id, "50", "9", order2.Date)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("发布成功...")

	fmt.Println("用电用户发布第三笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order3.Id, "150", "8", order3.Date)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("发布成功...")
	}

	fmt.Println("用电用户发布第四笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order4.Id, "60", "5", order4.Date)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("发布成功...")
	}
	fmt.Println("用电用户发布第五笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order5.Id, "70", "3", order5.Date)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("发布成功...")
	}

	fmt.Println("用电用户发布第六笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order6.Id, "110", "6", order6.Date)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("发布成功...")
	}

	fmt.Println("用电用户发布第七笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order7.Id, "130", "8", order7.Date)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("发布成功...")
	}

	fmt.Println("用电用户发布第八笔订单...")
	_, err = controller.InitLaunchOrder(controller.ConsumerApp, consumer1.Id, order8.Id, "90", "5", order8.Date)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("发布成功...")
	}
	fmt.Println("进入市场匹配...")
	controller.StartMarketMatch(controller.ProducerApp, controller.ConsumerApp)
	//fmt.Println("发电用户注册...")
	//msg, err := Org2ServiceSetup.ProducerRegister(producer)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println("信息发布成功, 交易编号为: " + msg)
	//}
	//fmt.Println("注册成功！")
	//global.InsertUser(user1)
	//fmt.Println(user1)
	//fmt.Println("发电用户查询个人信息...")
	//result, err := Org2ServiceSetup.QueryUser("431122200103110052")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &producer)
	//	fmt.Println("根据身份证号码查询信息成功：")
	//	fmt.Println(producer)
	//}
	//fmt.Println("用电用户注册...")
	//msg, err = Org1ServiceSetup.ConsumerRegister(consumer)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println("信息发布成功, 交易编号为: " + msg)
	//}
	//fmt.Println("注册成功！")
	//global.InsertUser(user2)
	//fmt.Println(user2)
	//fmt.Println("用电用户查询个人信息...")
	//result, err = Org1ServiceSetup.QueryUser("431122200203110052")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &consumer)
	//	fmt.Println("根据身份证号码查询信息成功：")
	//	fmt.Println(consumer)
	//}
	//fmt.Println("发电用户发布第一笔电力资产...")
	//judge,err:=LaunchAsset(controller.ProducerApp,user1.UserId,asset1.Id,"2000","2")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	if judge{
	//		fmt.Println("电力资产发布成功")
	//	}else{
	//		fmt.Println("您在一小时以内已经发布过了")
	//	}
	//}
	//fmt.Println("查询电力资产...")
	//result, err = Org2ServiceSetup.QueryAsset(asset1.CurrentOwnerId,asset1.Date)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &asset1)
	//	fmt.Println("根据发布者的身份证件号和时间组合健名查询信息成功：")
	//	fmt.Println(asset1)
	//}
	//fmt.Println("发电用户发布第二笔电力资产...")
	//judge,err=LaunchAsset(controller.ProducerApp,user1.UserId,"124","100","2")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	if judge{
	//		fmt.Println("电力资产发布成功")
	//	}else{
	//		fmt.Println("您在一小时以内已经发布过了")
	//	}
	//}
	//fmt.Println("查询电力资产...")
	//result, err = Org2ServiceSetup.QueryAsset(asset2.CurrentOwnerId,asset2.Date)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &asset2)
	//	fmt.Println("根据发布者的身份证件号和时间组合健名查询信息成功：")
	//	fmt.Println(asset2)
	//}
	//fmt.Println("发电用户用户查询个人信息...")
	//result, err = Org2ServiceSetup.QueryUser("431122200103110052")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &producer)
	//	fmt.Println("根据身份证号码查询信息成功：")
	//	fmt.Println(producer)
	//}
	//
	//fmt.Println("用电用户发布第一笔求购订单...")
	//judge,err=LaunchOrder(controller.ConsumerApp,user2.UserId,order1.Id,"100","1")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	if judge{
	//		fmt.Println("电力资产发布成功")
	//	}else{
	//		fmt.Println("您在一小时以内已经发布过了")
	//	}
	//}
	//fmt.Println("订单发布成功！")
	//fmt.Println("用电用户发布第二笔求购订单...")
	//msg, err = Org1ServiceSetup.AssetPurchase(order2)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println("信息发布成功, 交易编号为: " + msg)
	//}
	//fmt.Println("订单发布成功！")
	//fmt.Println("用电用户查询个人信息...")
	//result, err = Org1ServiceSetup.QueryUser("431122200203110052")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &consumer)
	//	fmt.Println("根据身份证号码查询信息成功：")
	//	fmt.Println(consumer)
	//}
	//fmt.Println("查询区块链所有资产...")
	//var allAssets []service.Electricity
	//result,err =Org2ServiceSetup.QueryAllAssets()
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &allAssets)
	//	fmt.Println(allAssets)
	//}
	//fmt.Println("查询区块链所有订单...")
	//result,err =Org2ServiceSetup.QueryAllOrders()
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &allAssets)
	//	fmt.Println(allAssets)
	//}
	//
	//fmt.Println("进入市场匹配...")
	//err= marketMatch.StartMarketMatch(controller.ProducerApp,controller.ConsumerApp)
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}else{
	//	fmt.Println("市场匹配成功！")
	//}
	//fmt.Println("查询第一笔电力资产所匹配订单...")
	//var records []service.Electricity
	//result, err = Org2ServiceSetup.QueryAssetMatchRecords(asset1.CurrentOwnerId,asset1.Id)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &records)
	//	fmt.Println("根据发布者的Id和时间组合健名查询信息成功：")
	//	fmt.Println(records)
	//}
	//fmt.Println("查询第二笔电力资产所匹配订单...")
	//result, err = Org2ServiceSetup.QueryAssetMatchRecords(asset2.CurrentOwnerId,asset2.Id)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &records)
	//	fmt.Println("根据发布者的Id和时间组合健名查询信息成功：")
	//	fmt.Println(records)
	//}
	//fmt.Println("查询第一笔订单所匹配电力资产...")
	//result, err = Org2ServiceSetup.QueryOrderMatchRecords(order1.CurrentOwnerId,order1.Id)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &records)
	//	fmt.Println("根据发布者的Id和订单Id组合健名查询信息成功：")
	//	fmt.Println(records)
	//}
	//fmt.Println("查询第二笔订单所匹配电力资产...")
	//result, err = Org2ServiceSetup.QueryOrderMatchRecords(order2.CurrentOwnerId,order2.Id)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &records)
	//	fmt.Println("根据发布者的Id和订单Id组合健名查询信息成功：")
	//	fmt.Println(records)
	//}
	//fmt.Println("查询第一笔电力资产...")
	//result, err = Org2ServiceSetup.QueryAsset(asset1.CurrentOwnerId,asset1.Date)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &asset1)
	//	fmt.Println("根据发布者的身份证件号和时间组合健名查询信息成功：")
	//	fmt.Println(asset1)
	//}
	//
	//fmt.Println("进入市场交易...")
	//state,err:=marketTransaction.StartTransaction(controller.ConsumerApp,user2.UserId,order1.Id)
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}else{
	//	if state==1{
	//		fmt.Println("你所购买的资产达到上限，交易失败")
	//	}else if state==2{
	//		fmt.Println("你所购买的资产中存在上限，交易失败")
	//	}else if state==3{
	//		fmt.Println("你的余额不足，交易失败")
	//	}else{
	//		fmt.Println("交易成功！")
	//	}
	//}
	//msg, err = Org2ServiceSetup.SaveAssetTransactionRecords(transaction1)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println("信息发布成功, 交易编号为: " + msg)
	//}
	//fmt.Println("查询单笔资产交易记录....")
	//var singleTransaction service.SingleTransaction
	//result, err = Org2ServiceSetup.QuerySingleAssetTransactionRecords(transaction.AssetId,transaction.Date)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &singleTransaction)
	//	fmt.Println("根据电力资产Id查询信息成功：")
	//	fmt.Println(singleTransaction)
	//}
	//fmt.Println("查询指定资产所有交易记录....")
	//var assetTransaction []service.SingleTransaction
	//result, err = Org2ServiceSetup.QueryAssetTransactionRecords(asset1.Id)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &assetTransaction)
	//	fmt.Println("根据电力资产Id查询信息成功：")
	//	fmt.Println(assetTransaction)
	//}
	//fmt.Println("查询指定资产交易记录....")
	//result, err = Org2ServiceSetup.QueryAssetTransactionRecords(transaction1.AssetId)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &assetTransaction)
	//	fmt.Println("根据电力资产Id查询信息成功：")
	//	fmt.Println(assetTransaction)
	//}
	//fmt.Println("查询市场全部资产交易记录....")
	//result, err = Org2ServiceSetup.QueryAllAssetsTransactionRecords()
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &assetTransaction)
	//	fmt.Println("查询信息成功：")
	//	fmt.Println(assetTransaction)
	//}
	//fmt.Println("查询总用户...")
	//var users []service.User
	//result, err = Org2ServiceSetup.QueryAllUsers()
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &users)
	//	fmt.Println("查询信息成功：")
	//	fmt.Println(users)
	//}
	//fmt.Println("进行身份溯源...")
	//user1,user2,err=userTracing.Tracing(controller.ConsumerApp,transaction.AssetId,transaction.Date)
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}else{
	//	fmt.Println(user1)
	//	fmt.Println(user2)
	//}
	//fmt.Println("发电用户查询自己余额")
	//result, err = Org2ServiceSetup.QueryUser("431122200103110052")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	json.Unmarshal(result, &producer)
	//	fmt.Println("根据身份证号码查询信息成功：")
	//	fmt.Println(cryptoCode.HomoDecryptData(producer.Surplus))
	//}
	app := controller.Application{
		Setup: Org1ServiceSetup,
	}

	// 初始化数据库
	// if err := model.InitDB(); err != nil {
	// 	fmt.Println("数据库初始化失败")
	// 	return
	// }

	web.WebStart(app)
}
