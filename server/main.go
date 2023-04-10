package main

import (
	"fmt"
	"os"
	"powerTrading/sdkInit"
	"powerTrading/service"
	"powerTrading/web"
	"powerTrading/web/controller"
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
		fmt.Printf(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// // invoke chaincode set status
	// fmt.Println(">> 通过链码外部服务设置链码状态......")
	// producer := service.User{
	// 	Id:       "431122200103110052",
	// 	Purchases:   nil,
	// 	Sells:    nil,
	// }
	// producer.Surplus= cryptoCode.HomoEncryptData(strconv.FormatInt(0,10))
	// consumer := service.User{
	// 	Id:       "431122200203110052",
	// 	Purchases:   nil,
	// 	Sells: nil,
	// }
	// consumer.Surplus= cryptoCode.HomoEncryptData(strconv.FormatInt(0,10))
	// asset1 := service.Electricity{
	// 	Id:              "123",
	// 	Type:            "光能发电",
	// 	State:           true                                                                                                                                                                                                                                                  ,
	// 	CurrentOwnerId:  "431122200103110052",
	// 	Date:            "2023-2-2 19:55",
	// }
	// asset1.Scope= cryptoCode.HomoEncryptData("100")
	// asset1.Price= cryptoCode.HomoEncryptData("2")
	// asset2 := service.Electricity{
	// 	Id:              "456",
	// 	Type:            "光能发电",
	// 	State:           true                                                                                                                                                                                                                                                  ,
	// 	CurrentOwnerId:  "431122200103110052",
	// 	Date:            "2023-2-3 19:55",
	// }
	// asset2.Scope= cryptoCode.HomoEncryptData("200")
	// asset2.Price= cryptoCode.HomoEncryptData("3")
	// order1:=service.Electricity{
	// 	Id: "abc",
	// 	Date:            "2023-2-1 19:55",
	// 	State:           true,
	// 	CurrentOwnerId:  "431122200203110052",
	// }
	// order1.Scope= cryptoCode.HomoEncryptData("350")
	// order1.Price= cryptoCode.HomoEncryptData("1")
	// order2:=service.Electricity{
	// 	Id: "def",
	// 	Date:            "2023-2-3 20:55",
	// 	State:           true,
	// 	CurrentOwnerId:  "431122200203110052",
	// }
	// order2.Scope= cryptoCode.HomoEncryptData("100")
	// order2.Price= cryptoCode.HomoEncryptData("1")
	// Org1ServiceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	// if err != nil {
	// 	fmt.Println()
	// 	os.Exit(-1)
	// }
	Org2ServiceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[1], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	// fmt.Println("发电用户注册...")
	// msg, err := Org2ServiceSetup.ProducerRegister(producer)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("信息发布成功, 交易编号为: " + msg)
	// }
	// fmt.Println("注册成功！")
	// fmt.Println("发电用户查询个人信息...")
	// result, err := Org2ServiceSetup.QueryUser("431122200103110052")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &producer)
	// 	fmt.Println("根据身份证号码查询信息成功：")
	// 	fmt.Println(producer)
	// }
	// fmt.Println("用电用户注册...")
	// msg, err = Org1ServiceSetup.ConsumerRegister(consumer)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("信息发布成功, 交易编号为: " + msg)
	// }
	// fmt.Println("注册成功！")
	// fmt.Println("用电用户查询个人信息...")
	// result, err = Org1ServiceSetup.QueryUser("431122200203110052")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &consumer)
	// 	fmt.Println("根据身份证号码查询信息成功：")
	// 	fmt.Println(consumer)
	// }
	// fmt.Println("发电用户发布第一笔电力资产...")
	// msg, err = Org2ServiceSetup.AssetSell(asset1)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("信息发布成功, 交易编号为: " + msg)
	// }
	// fmt.Println("电力资产发布成功！")
	// fmt.Println("查询电力资产...")
	// result, err = Org2ServiceSetup.QueryAsset(asset1.CurrentOwnerId,asset1.Date)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &asset1)
	// 	fmt.Println("根据发布者的身份证件号和时间组合健名查询信息成功：")
	// 	fmt.Println(asset1)
	// }
	// fmt.Println("发电用户发布第二笔电力资产...")
	// msg, err = Org2ServiceSetup.AssetSell(asset2)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("信息发布成功, 交易编号为: " + msg)
	// }
	// fmt.Println("电力资产发布成功！")
	// fmt.Println("查询电力资产...")
	// result, err = Org2ServiceSetup.QueryAsset(asset2.CurrentOwnerId,asset2.Date)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &asset2)
	// 	fmt.Println("根据发布者的身份证件号和时间组合健名查询信息成功：")
	// 	fmt.Println(asset2)
	// }
	// fmt.Println("发电用户用户查询个人信息...")
	// result, err = Org2ServiceSetup.QueryUser("431122200103110052")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &producer)
	// 	fmt.Println("根据身份证号码查询信息成功：")
	// 	fmt.Println(producer)
	// }
	// fmt.Println("用电用户发布第一笔求购订单...")
	// msg, err = Org1ServiceSetup.AssetPurchase(order1)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("信息发布成功, 交易编号为: " + msg)
	// }
	// fmt.Println("订单发布成功！")
	// fmt.Println("用电用户发布第二笔求购订单...")
	// msg, err = Org1ServiceSetup.AssetPurchase(order2)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("信息发布成功, 交易编号为: " + msg)
	// }
	// fmt.Println("订单发布成功！")
	// fmt.Println("用电用户查询个人信息...")
	// result, err = Org1ServiceSetup.QueryUser("431122200203110052")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &consumer)
	// 	fmt.Println("根据身份证号码查询信息成功：")
	// 	fmt.Println(consumer)
	// }
	// fmt.Println("查询区块链所有资产...")
	// var allAssets []service.Electricity
	// result,err =Org2ServiceSetup.QueryAllAssets()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &allAssets)
	// 	fmt.Println(allAssets)
	// }
	// fmt.Println("查询区块链所有订单...")
	// result,err =Org2ServiceSetup.QueryAllOrders()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &allAssets)
	// 	fmt.Println(allAssets)
	// }
	// fmt.Println("进入市场匹配...")
	// consumerApp := controller.Application{
	// 	Setup: Org1ServiceSetup,
	// }
	producerApp := controller.Application{
		Setup: Org2ServiceSetup,
	}
	// err= marketMatch.StartMarketMatch(producerApp,consumerApp)
	// if err!=nil{
	// 	fmt.Println(err.Error())
	// }else{
	// 	fmt.Println("市场匹配成功！")
	// }
	// fmt.Println("查询第一笔电力资产所匹配订单...")
	// var records []service.Electricity
	// result, err = Org2ServiceSetup.QueryAssetMatchRecords(asset1.CurrentOwnerId,asset1.Date)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &records)
	// 	fmt.Println("根据发布者的Id和时间组合健名查询信息成功：")
	// 	fmt.Println(records)
	// }
	// fmt.Println("查询第二笔电力资产所匹配订单...")
	// result, err = Org2ServiceSetup.QueryAssetMatchRecords(asset2.CurrentOwnerId,asset2.Date)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &records)
	// 	fmt.Println("根据发布者的Id和时间组合健名查询信息成功：")
	// 	fmt.Println(records)
	// }
	// fmt.Println("查询第一笔订单所匹配电力资产...")
	// result, err = Org2ServiceSetup.QueryOrderMatchRecords(order1.CurrentOwnerId,order1.Id)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &records)
	// 	fmt.Println("根据发布者的Id和订单Id组合健名查询信息成功：")
	// 	fmt.Println(records)
	// }
	// fmt.Println("查询第二笔订单所匹配电力资产...")
	// result, err = Org2ServiceSetup.QueryOrderMatchRecords(order2.CurrentOwnerId,order2.Id)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &records)
	// 	fmt.Println("根据发布者的Id和订单Id组合健名查询信息成功：")
	// 	fmt.Println(records)
	// }
	// fmt.Println("查询第一笔电力资产...")
	// result, err = Org2ServiceSetup.QueryAsset(asset1.CurrentOwnerId,asset1.Date)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	json.Unmarshal(result, &asset1)
	// 	fmt.Println("根据发布者的身份证件号和时间组合健名查询信息成功：")
	// 	fmt.Println(asset1)
	// }

	web.WebStart(producerApp)
}
