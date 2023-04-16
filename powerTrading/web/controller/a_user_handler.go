package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"powerTrading/web/cryptoCode"
	"powerTrading/web/groupSignature"
	"strconv"
	"time"

	"powerTrading/service"
	"powerTrading/web/global"
	"powerTrading/web/middleware"
	"powerTrading/web/model"
	"powerTrading/web/tools"
)

var err error

var currUser *model.User

var currMoney int64

type RetData struct {
	Code int32       `json:"code"`
	Flag bool        `json:"flag"`
	Data interface{} `json:"data"`
}

type ElectricityRespData struct {
	Id     string `json:"id"`
	Scope  string `json:"scope"`
	Price  string `json:"price"`
	Date   string `json:"date"`
	Status string `json:"status"`
}

type UserRespData struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// 购电方交易记录详情响应
type ConsumerRespData struct {
	Scope string `json:"scope"`
	Price string `json:"price"`
	Date  string `json:"date"`
}

// 供电商交易记录详情响应
type ProducerRespData struct {
	Scope string `json:"scope"`
	Price string `json:"price"`
	Date  string `json:"date"`
}

// 显示资产溯源交易响应结构体
type AssetTraceBackData struct {
	AssetId string `json:"asset_id"`
	Scope   string `json:"scope"`
	Price   string `json:"price"`
	Date    string `json:"date"`
}

// 溯源结果响应结构体
type TraceBackDataResp struct {
	ProducerId        string `json:"producer_id"`
	ProducerName      string `json:"producer_name"`
	ProducerCreatedAt string `json:"producer_created_at"`
	ConsumerId        string `json:"consumer_id"`
	ConsumerName      string `json:"consumer_name"`
	ConsumerCreatedAt string `json:"consumer_created_at"`
}

type RecordDataResp struct {
	Id    string `json:"id"`
	Scope string `json:"scope"`
	Cash  string `json:"cash"`
	Date  string `json:"date"`
}

// UserRegister 用户注册
func (app *Application) UserRegister(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	// 初始化用户
	for _, user := range global.UserList {
		user.GroupPrivateKey = groupSignature.GenerateMemberPrivateKey()
		global.UserMap[user.UserId] = user
	}

	// 绑定参数
	user := new(model.User)
	user.Role = r.FormValue("role")
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.UserId = tools.GenUUID()
	user.GroupPrivateKey = groupSignature.GenerateMemberPrivateKey()
	user.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	row := global.InsertUser(user)
	// dao.InsertUser(user)

	if row == 0 {
		fmt.Println("注册用户失败")
		return
	}

	if user.Role == "供电商" {
		_, err = app.Setup.ProducerRegister(service.User{
			Id:        user.UserId,
			Purchases: nil,
			Sells:     make([]service.Electricity, 0),
			Surplus:   cryptoCode.HomoEncryptData("0"),
		})
	} else {
		_, err = app.Setup.ConsumerRegister(service.User{
			Id:        user.UserId,
			Purchases: make([]service.Electricity, 0),
			Sells:     nil,
			Surplus:   cryptoCode.HomoEncryptData("0"),
		})
	}

	if err != nil {
		fmt.Println("区块用户注册失败")
		return
	}

	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = nil

	retDataJson, _ := json.Marshal(retData)

	io.WriteString(w, string(retDataJson))
}

// UserLogin 用户登陆
func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	// 绑定参数
	param := new(model.User)
	param.Username = r.FormValue("username")
	param.Password = r.FormValue("password")
	user := global.QueryUserByName(param.Username, param.Password)
	if user == nil {
		fmt.Println("用户不存在")
		return
	}

	currUser = user

	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = currUser.UserId
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// QueryUserBaseInfo 查询用户基本信息
func (app *Application) QueryUserBaseInfo(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	// 首先从登陆用户的ID(currUser.userID) 查询到用户名和身份信息
	var user = global.QueryUser(currUser.UserId)
	if user == nil {
		fmt.Println("通过用户userId查询失败,Error")
		return
	}

	userByte, _ := app.Setup.QueryUser(user.UserId)
	var srvUser service.User
	_ = json.Unmarshal(userByte, &srvUser)

	// 再查询用户的余额
	balanceStr := cryptoCode.HomoDecryptData(srvUser.Surplus)
	currMoney, _ = strconv.ParseInt(balanceStr, 10, 64)

	// 使用结构体存储
	var respUser = struct {
		Role     string `json:"role"`
		Username string `json:"username"`
		Balance  int64  `json:"balance"`
	}{
		Role:     user.Role,
		Username: user.Username,
		Balance:  currMoney,
	}

	// 放在 Response 中返回
	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = respUser

	// fmt.Printf("用户信息为:%#v\n", user)
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// UserTopUp 用户充值
func (app *Application) UserTopUp(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	var user service.User
	userByte, _ := app.Setup.QueryUser(currUser.UserId)

	_ = json.Unmarshal(userByte, &user)

	temp := cryptoCode.HomoDecryptData(user.Surplus)
	currBalance, _ := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		return
	}

	// 绑定参数
	balance, err := strconv.ParseInt(r.FormValue("balance"), 10, 64)
	if err != nil {
		fmt.Println("解析参数转换失败：", err)
		return
	}
	user.Surplus = cryptoCode.HomoEncryptData(fmt.Sprintf("%d", balance+currBalance))
	_, err = app.Setup.UserTopUp(user)
	if err != nil {
		fmt.Println("充值失败")
		return
	}

	currMoney = currBalance + balance

	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = nil
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// PublishBuyOrder 发布求购订单
func (app *Application) PublishBuyOrder(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	// 绑定参数，用户需要的电量，能接受的价格
	scope := r.FormValue("scope")
	if err != nil {
		fmt.Println("绑定参数 scope 错误")
		return
	}

	price := r.FormValue("price")
	if err != nil {
		fmt.Println("绑定参数 price 错误")
		return
	}

	// zeroProof
	success, err := LaunchOrder(ConsumerApp, currUser.UserId, tools.GenUUID(), scope, price)

	if err != nil {
		fmt.Println("发布订单失败,Error:", err)
		return
	}

	if !success {
		fmt.Println("zeroProof.LaunchOrder 发布失败, Error: ", "一小时内只能发布一次求购订单...")
		retData := new(RetData)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "发布失败,一小时内只能发布一次求购订单"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	} else {
		fmt.Println("发布成功，进入匹配队列...")
		// 匹配操作
		StartMarketMatch(ProducerApp, ConsumerApp)
		retData := new(RetData)
		retData.Code = 200
		retData.Flag = true
		retData.Data = nil
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
	}
}

// 查询自身发布的求购订单
func (app *Application) QuerySelfBuyOrder(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	var user service.User
	userByte, _ := app.Setup.QueryUser(currUser.UserId)
	err = json.Unmarshal(userByte, &user)
	if err != nil {
		fmt.Println("buyOrderList 反序列化失败, Error:", err)
		return
	}

	var retDataElectricity []ElectricityRespData
	for _, electricity := range user.Purchases {
		var temp = ElectricityRespData{
			Id:    electricity.Id,
			Scope: cryptoCode.HomoDecryptData(electricity.Scope),
			Price: cryptoCode.HomoDecryptData(electricity.Price),
			Date:  electricity.Date,
		}

		if electricity.State {
			if len(electricity.MatchRecords) == 0 {
				temp.Status = "匹配中"
			} else {
				temp.Status = "待支付"
			}
		} else {
			temp.Status = "已完成"
		}

		retDataElectricity = append(retDataElectricity, temp)
	}

	// 将数据传至前端
	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = retDataElectricity
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// QueryAllUser 查询所有区块的用户
func (app *Application) QueryAllUser(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	userList := global.QueryAllUser()
	var userRespData []UserRespData
	for _, user := range userList {
		temp := UserRespData{
			Id:       user.UserId,
			Username: user.Username,
			Role:     user.Role,
		}

		userRespData = append(userRespData, temp)
	}

	// 将数据传至前端
	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = userRespData
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// QueryUserSelfElectricity 查询用户发布的电力资产
func (app *Application) QueryUserSelfElectricity(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	userByte, err := app.Setup.QueryUser(currUser.UserId)
	if err != nil {
		fmt.Printf("查询用户信息失败：%v\n", err)
		return
	}

	var user = new(service.User)
	err = json.Unmarshal(userByte, user)
	if err != nil {
		fmt.Printf("反序列化用户信息失败. Error:%v\n", err)
		return
	}

	var retDataElectricity []ElectricityRespData
	for _, electricity := range user.Sells {
		var temp = ElectricityRespData{
			Id:    electricity.Id,
			Scope: cryptoCode.HomoDecryptData(electricity.Scope),
			Price: cryptoCode.HomoDecryptData(electricity.Price),
			Date:  electricity.Date,
		}

		if electricity.State == true {
			temp.Status = "匹配中"
		} else {
			temp.Status = "已出售"
		}
		retDataElectricity = append(retDataElectricity, temp)
	}

	// 将数据传至前端
	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = retDataElectricity
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// SellElectricity 发布电量数据
func (app *Application) SellElectricity(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	// 绑定参数
	scope := r.FormValue("scope")
	if err != nil {
		fmt.Println("绑定参数 scope 错误")
		retData := new(RetData)
		retData.Code = 200
		retData.Flag = true
		retData.Data = "服务器繁忙"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	price := r.FormValue("price")
	if err != nil {
		fmt.Println("绑定参数 price 错误")
		retData := new(RetData)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "参数错误"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	success, err := LaunchAsset(ProducerApp, currUser.UserId, tools.GenUUID(), scope, price)
	if err != nil {
		fmt.Println("zeroProof.LaunchAsset() 错误，Error: ", err)
		retData := new(RetData)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "参数错误"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	if !success {
		fmt.Println("zeroProof.LaunchAsset 发布失败, Error: ", "一小时内只能发布一次求购订单...")
		// 将数据传至前端
		retData := new(RetData)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "一小时内只能发布一次求购订单"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	} else {
		fmt.Println("111")
		fmt.Println("发布成功，进入匹配队列")
		StartMarketMatch(ProducerApp, ConsumerApp)
		// 将数据传至前端
		retData := new(RetData)
		retData.Code = 200
		retData.Flag = true
		retData.Data = nil
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
	}
}

// 查询所有发布的电力数据(电力中心)
func (app *Application) QueryAllOrder(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	var electricityList = []service.Electricity{}
	listByte, err := app.Setup.QueryAllAssets()
	if err != nil {
		fmt.Println("查询所有订单失败,Error: ", err)
		return
	}

	err = json.Unmarshal(listByte, &electricityList)
	if err != nil {
		fmt.Println("electricityList 反序列化失败, Error:", err)
		return
	}

	var retDataElectricity []ElectricityRespData
	for _, electricity := range electricityList {
		var temp = ElectricityRespData{
			Scope: cryptoCode.HomoDecryptData(electricity.Scope),
			Price: cryptoCode.HomoDecryptData(electricity.Price),
			Date:  electricity.Date,
		}
		scopeNum, _ := strconv.ParseInt(temp.Scope, 10, 64)
		for _, elec := range electricity.MatchRecords {
			tempScope, _ := strconv.ParseInt(cryptoCode.HomoDecryptData(elec.Scope), 10, 64)
			scopeNum -= tempScope
		}

		retDataElectricity = append(retDataElectricity, temp)
	}

	// 将数据传至前端
	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = retDataElectricity
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// 查询用户某个时间段发布的电量订单
func (app *Application) QueryUserOrder(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}
	userId := r.FormValue("user_id")
	date := r.FormValue("date")

	electricityByte, err := app.Setup.QueryAsset(userId, date)
	if err != nil {
		fmt.Println("查询用户发布的订单失败, Error:", err)
		return
	}

	var electricity = new(service.Electricity)

	err = json.Unmarshal(electricityByte, electricity)
	if err != nil {
		fmt.Println("反序列化失败, Error: ", err)
		return
	}
	fmt.Printf("查询用户Id为 %v 发布的订单信息为:\n", userId)
	fmt.Printf("订单ID:%v, 电力类型:%v, 电量:%v, 订单价格:%v, 订单发布者id:%v, 发布日期:%v\n",
		electricity.Id,
		electricity.Type,
		electricity.Scope,
		electricity.Price,
		electricity.CurrentOwnerId,
		electricity.Date)
}

// 购电用户支付订单
func (app *Application) PayOrder(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	// 绑定参数
	orderId := r.FormValue("order_id")

	status, err := StartTransaction(ConsumerApp, currUser.UserId, orderId)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	retData := new(RetData)
	retData.Code = 500
	retData.Flag = false
	switch status {
	case 0:
		retData.Data = "还未匹配相应对资产"
	case 1:
		retData.Data = "购买电量达到上限"
	case 2:
		retData.Data = "供电商售卖电量达到上限"
	case 3:
		retData.Data = "余额不足，请充值"
	case 4:
		retData.Code = 200
		retData.Flag = true
		retData.Data = "支付成功"
	}
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// 购电用户撤销交易
func (app *Application) BuyUserCancelOrder(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	order_id := r.FormValue("order_id")

	retData := new(RetData)
	success, err := CancelOrder(ProducerApp, currUser.UserId, order_id)
	if err != nil {
		fmt.Println("BuyUserCancelOrder 错误, Error: ", err)

		retData.Code = 500
		retData.Flag = false
		retData.Data = "BuyUserCancelOrder 错误"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	if !success {
		fmt.Println("您已支付，无法撤销订单")
		retData.Code = 500
		retData.Flag = false
		retData.Data = "您已支付，无法撤销订单"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}
	fmt.Println("订单撤销成功...")
	// TODO:
	retData.Code = 200
	retData.Flag = true
	retData.Data = nil
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// 供电商撤销交易
func (app *Application) SellUserCancelOrder(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	assetId := r.FormValue("asset_id")

	retData := new(RetData)
	success, err := CancelAsset(ProducerApp, currUser.UserId, assetId)
	if err != nil {
		fmt.Println("SellUserCancelOrder 错误, Error: ", err)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "SellUserCancelOrder 错误"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	if !success {
		fmt.Println("该电力资产已被购电用户预购")
		retData.Code = 500
		retData.Flag = false
		retData.Data = "该电力资产已被购电用户预购"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}
	fmt.Println("撤销成功...")
	retData.Code = 200
	retData.Flag = true
	retData.Data = nil
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// 购电用户交易记录
func (app *Application) FindOrderRecord(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	orderId := r.FormValue("order_id")

	var records []service.Electricity

	retData := new(RetData)
	listByte, err := app.Setup.QueryOrderMatchRecords(currUser.UserId, orderId)
	if err != nil {
		fmt.Println("app.Setup.QueryOrderMatchRecords Error: ", err)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "app.Setup.QueryOrderMatchRecords Error"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	err = json.Unmarshal(listByte, &records)
	if err != nil {
		fmt.Println("order json.Unmarshal(listByte, &records) Error: ", err)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "order json.Unmarshal(listByte, &records) Error"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	if len(records) == 0 {
		fmt.Println("暂无交易记录，该订单还在匹配中...")
		retData.Code = 500
		retData.Flag = false
		retData.Data = "暂无交易记录，该订单还在匹配中"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}
	fmt.Println("订单匹配成功...")

	var consumerRespData = []ConsumerRespData{}

	for _, value := range records {
		var temp = ConsumerRespData{
			Scope: cryptoCode.HomoDecryptData(value.Scope),
			Price: cryptoCode.HomoDecryptData(value.CashProof),
			Date:  value.Date,
		}

		consumerRespData = append(consumerRespData, temp)
	}

	retData.Code = 200
	retData.Flag = true
	retData.Data = consumerRespData
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// 供电商查看交易记录
func (app *Application) FindAssetRecord(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}
	assetId := r.FormValue("asset_id")

	var records []service.SingleTransaction

	retData := new(RetData)
	listByte, err := app.Setup.QueryAssetTransactionRecords(assetId)
	if err != nil {
		fmt.Println("app.Setup.QueryAssetTransactionRecords Error: ", err)
		retData.Code = 200
		retData.Flag = true
		retData.Data = records
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	err = json.Unmarshal(listByte, &records)
	if err != nil {
		fmt.Println("asset json.Unmarshal(listByte, &recorDisplayAllOrderWithTraceBackds) Error: ", err)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "asset json.Unmarshal(listByte, &recorDisplayAllOrderWithTraceBackds) Error:"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	if len(records) == 0 {
		fmt.Println("暂无交易记录，该资产还在匹配中...")
		retData.Code = 500
		retData.Flag = false
		retData.Data = "暂无交易记录，该资产还在匹配中"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}
	fmt.Println("资产匹配成功...")

	respData := []ProducerRespData{}

	for _, value := range records {
		var temp = ProducerRespData{
			Scope: cryptoCode.HomoDecryptData(value.Scope),
			Price: cryptoCode.HomoDecryptData(value.Cash),
			Date:  value.Date,
		}
		respData = append(respData, temp)
	}

	retData.Code = 200
	retData.Flag = true
	retData.Data = respData
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// 管理员显示所有交易(溯源)
func (app *Application) DisplayAllOrderWithTraceBack(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	var records []service.SingleTransaction

	retData := new(RetData)
	// 操作
	listByte, err := app.Setup.QueryAllAssetsTransactionRecords()
	if err != nil {
		fmt.Println("app.Setup.QueryAllAssetsTransactionRecords Error: ", err)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "app.Setup.QueryAllAssetsTransactionRecords Error: "
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}
	err = json.Unmarshal(listByte, &records)
	if err != nil {
		fmt.Println("DisplayAllOrderWithTraceBack json.Unmarshal(listByte, &records) Error: ", err)
		retData.Code = 500
		retData.Flag = false
		retData.Data = "DisplayAllOrderWithTraceBack json.Unmarshal(listByte, &records) Error: "
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	if len(records) == 0 {
		fmt.Println("暂无交易记录")
		retData.Code = 200
		retData.Flag = false
		retData.Data = "暂无交易记录"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	fmt.Println("DisplayAllOrderWithTraceBack显示成功")

	assetTraceBackData := []AssetTraceBackData{}
	for _, value := range records {
		var temp = AssetTraceBackData{
			AssetId: value.AssetId,
			Scope:   cryptoCode.HomoDecryptData(value.Scope),
			Price:   cryptoCode.HomoDecryptData(value.Cash),
			Date:    value.Date,
		}
		assetTraceBackData = append(assetTraceBackData, temp)
	}
	retData.Code = 200
	retData.Flag = true
	retData.Data = assetTraceBackData
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

// 交易溯源操作
func (app *Application) TraceBackOrder(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	assetId := r.FormValue("asset_id")
	date := r.FormValue("date")

	retData := new(RetData)
	buyer, seller, err := Tracing(ProducerApp, assetId, date)
	if err != nil {
		fmt.Println("溯源操作失败")
		retData.Code = 500
		retData.Flag = false
		retData.Data = "溯源操作失败"
		retDataJson, _ := json.Marshal(retData)
		io.WriteString(w, string(retDataJson))
		return
	}

	fmt.Println("该笔交易双方分别为：", buyer, seller)

	traceBackDataResp := [1]TraceBackDataResp{}

	traceBackDataResp[0].ProducerId = seller.UserId
	traceBackDataResp[0].ProducerName = seller.Username
	traceBackDataResp[0].ProducerCreatedAt = seller.CreatedAt

	traceBackDataResp[0].ConsumerId = buyer.UserId
	traceBackDataResp[0].ConsumerName = buyer.Username
	traceBackDataResp[0].ConsumerCreatedAt = buyer.CreatedAt

	retData.Code = 200
	retData.Flag = true
	retData.Data = traceBackDataResp
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}

func (app *Application) QueryAllRecords(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	records := []Record{}
	retData := new(RetData)
	if currUser.Role == "供电商" {
		err, records = QueryAssetTransactionRecords(ProducerApp, currUser.UserId)
		if err != nil {
			fmt.Println("供电商查询所有历史交易记录失败", err)
			return
		}
	} else {
		err, records = QueryOrderTransactionRecords(ProducerApp, currUser.UserId)
		if err != nil {
			fmt.Println("购电用户查询所有历史交易记录失败", err)
			return
		}
	}

	retData.Code = 200
	retData.Flag = true
	retData.Data = records
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))
}
