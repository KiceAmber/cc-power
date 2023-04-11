package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"powerTrading/tools"
	"powerTrading/web/cryptoCode"
	"strconv"
	"time"

	"powerTrading/service"
	"powerTrading/web/global"
	"powerTrading/web/middleware"
	"powerTrading/web/model"
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
	Status bool   `json:"status"`
}

// UserRegister 用户注册
func (app *Application) UserRegister(w http.ResponseWriter, r *http.Request) {

	middleware.Cors(&w)
	if r.Method == "Option" {
		return
	}

	// 绑定参数
	user := new(model.User)
	user.Role = r.FormValue("role")
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.UserId = tools.GenUUID()

	row := global.InsertUser(user)
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

	fmt.Println("currMoney = ", currMoney)

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

	electricity := service.Electricity{
		Id:             tools.GenUUID(),
		Type:           "光伏发电",
		State:          true,
		CurrentOwnerId: currUser.UserId,
		Date:           time.Now().Format("2006-01-02 15:04:05"),
	}

	// 绑定参数，用户需要的电量，能接受的价格
	electricity.Scope = cryptoCode.HomoEncryptData(r.FormValue("scope"))
	if err != nil {
		fmt.Println("绑定参数 scope 错误")
		return
	}

	electricity.Price = cryptoCode.HomoEncryptData(r.FormValue("price"))
	if err != nil {
		fmt.Println("绑定参数 price 错误")
		return
	}

	_, err = app.Setup.AssetPurchase(electricity)
	if err != nil {
		fmt.Println("发布订单失败,Error:", err)
		return
	}

	retData := new(RetData)
	retData.Code = 200
	retData.Flag = true
	retData.Data = nil
	retDataJson, _ := json.Marshal(retData)
	io.WriteString(w, string(retDataJson))

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
			Id:     electricity.Id,
			Scope:  cryptoCode.HomoDecryptData(electricity.Scope),
			Price:  cryptoCode.HomoDecryptData(electricity.Price),
			Date:   electricity.Date,
			Status: true,
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

	fmt.Println("查询用户信息如下:")
	for _, user := range global.UserList {
		fmt.Printf("用户role为: %s, 用户id为: %s, 用户名为: %s, 用户密码为: %s\n",
			user.Role, user.UserId, user.Username, user.Password)
	}
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
			Id:     electricity.Id,
			Scope:  cryptoCode.HomoDecryptData(electricity.Scope),
			Price:  cryptoCode.HomoDecryptData(electricity.Price),
			Date:   electricity.Date,
			Status: true,
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

	var err error
	// 绑定参数
	electricity := service.Electricity{
		Id:             tools.GenUUID(),
		Type:           "光伏发电",
		State:          true,
		CurrentOwnerId: currUser.UserId,
		Date:           time.Now().Format("2006-01-02 15:04:05"),
	}
	electricity.Scope = cryptoCode.HomoEncryptData(r.FormValue("scope"))
	if err != nil {
		fmt.Println("绑定参数 scope 错误")
		return
	}

	electricity.Price = cryptoCode.HomoEncryptData(r.FormValue("price"))
	if err != nil {
		fmt.Println("绑定参数 price 错误")
		return
	}

	_, err = app.Setup.AssetSell(electricity)

	// 下面代码用来测试结果
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("用户发布的电力资源信息为：%#v\n", electricity)
}

// 查询用户自身发布的电量数据
func (app *Application) QueryUserSelfOrder(w http.ResponseWriter, r *http.Request) {

}

// 查询所有发布的电力数据(供电商)
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

	fmt.Println("=============================")
	fmt.Println("所有的订单信息为:")
	for _, electricity := range electricityList {
		fmt.Printf("订单ID:%v, 电力类型:%v, 电量:%v, 订单价格:%v, 订单发布者id:%v, 发布日期:%v\n",
			electricity.Id,
			electricity.Type,
			electricity.Scope,
			electricity.Price,
			electricity.CurrentOwnerId,
			electricity.Date)
	}

	var retDataElectricity []ElectricityRespData
	for _, electricity := range electricityList {
		var temp = ElectricityRespData{
			Scope: cryptoCode.HomoDecryptData(electricity.Scope),
			Price: cryptoCode.HomoDecryptData(electricity.Price),
			Date:  electricity.Date,
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
