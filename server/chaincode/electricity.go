package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// 电力交易系统智能合约
type User struct { //定义用户结构
	Id        string        `json:"Id"`
	Purchases []Electricity //用户所要购买的电力资源
	Sells     []Electricity //用户所发布的电力资源
	Surplus   []byte //用户余额
	An float64               `json:"an"`
	Bn float64               `json:"bn"`
	Limit int64               `json:"limit"`
}
type TotalUser struct {
	Users []User
}
var userKeyId = "02eb3850b7a8c944f22fa3b24c77121ce2c1b025f27e41e7b51a0749a0d4f1b5" //采用公共哈希值ID
type Electricity struct {                                                          //定义电力资产结构
	Id             string          `json:"id"`
	Type           string          `json:"type"`
	Scope          []byte          `json:"scope"`
	State          bool            `json:"state"`
	Price          []byte         `json:"price"`
	CurrentOwnerId string          `json:"currentOwnerId"`
	Date           string          `json:"date"`
	Records        []HistoryRecord //电力资产交易记录
	MatchRecords   []Electricity //市场交易匹配记录
}
type HistoryRecord struct { //定义电力资产历史记录数据结构
	TxId        string //区块连中存储电力资产的Id
	Electricity Electricity
}
type ElectricityChanCode struct {
}

func (t *ElectricityChanCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")
	return shim.Success(nil)
}
func (t *ElectricityChanCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()
	switch fun {
	case "userRegister":
		return t.userRegister(stub, args)
	case "assetSell":
		return t.assetsSell(stub, args)
	case "assetPurchase":
		return t.assetPurchase(stub, args)
	case "queryUser":
		return t.queryUser(stub, args)
	case "queryAsset":
		return t.queryAsset(stub, args)
	case "queryAllAssets":
		return t.queryAllAssets(stub, args)
	case "queryAllOrders":
		return t.queryAllOrders(stub, args)
	case "userTopUp":
		return t.userTopUp(stub,args)
	case "cancelSell":
		return t.cancelSell(stub,args)
	case "cancelPurchase":
		return t.cancelPurchase(stub,args)
	case "queryAssetMatchRecords":
		return t.queryAssetMatchRecords(stub,args)
	case "queryOrderMatchRecords":
		return t.queryOrderMatchRecords(stub,args)
	case "saveAssetMatchRecords":
		return t.saveAssetMatchRecords(stub,args)
	case "saveOrderMatchRecords":
		return t.saveOrderMatchRecords(stub,args)
	default:
		return shim.Error("指定的函数名称错误")
	}
}
func putAsset(stub shim.ChaincodeStubInterface, asset Electricity) ([]byte, bool) { //存储电力资产
	b, err := json.Marshal(asset)
	if err != nil {
		return nil, false
	}
	key := asset.CurrentOwnerId + "#" + asset.Date
	err = stub.PutState(key, b)
	if err != nil {
		return nil, false
	}
	return b, true
}
func putUser(stub shim.ChaincodeStubInterface, user User) ([]byte, bool) { //存储注册用户
	b, err := json.Marshal(user) //序列化
	if err != nil {
		return nil, false
	}
	// 保存edu状态
	err = stub.PutState(user.Id, b)
	if err != nil {
		return nil, false
	}
	return b, true
}
func putAllUser(stub shim.ChaincodeStubInterface, allUser TotalUser) ([]byte, bool) {
	b, err := json.Marshal(allUser) //序列化
	if err != nil {
		return nil, false
	}
	// 保存edu状态
	err = stub.PutState(userKeyId, b)
	if err != nil {
		return nil, false
	}
	return b, true
}
func getAsset(stub shim.ChaincodeStubInterface, userId string, date string) (Electricity, bool) { //获取状态数据库中的电力资产
	var electricity Electricity
	// 根据身份证号码查询信息状态
	key := userId + "#" + date
	b, err := stub.GetState(key)
	if err != nil {
		return electricity, false
	}
	if b == nil {
		return electricity, false
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &electricity)
	if err != nil {
		return electricity, false
	}
	// 返回结果
	return electricity, true
}
func getUser(stub shim.ChaincodeStubInterface, userId string) (User, bool) { //获取状态数据库中的用户
	var user User
	// 根据身份证号码查询信息状态
	b, err := stub.GetState(userId)
	if err != nil {
		return user, false
	}
	if b == nil {
		return user, false
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &user)
	if err != nil {
		return user, false
	}
	// 返回结果
	return user, true
}
func getAllUser(stub shim.ChaincodeStubInterface) (TotalUser, bool) {
	var allUser TotalUser
	// 根据身份证号码查询信息状态
	b, err := stub.GetState(userKeyId)
	if err != nil {
		return allUser, false
	}
	if b == nil {
		return allUser, false
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &allUser)
	if err != nil {
		return allUser, false
	}
	// 返回结果
	return allUser, true
}
func getAllAssets(stub shim.ChaincodeStubInterface) ([]Electricity, bool) {
	var allAssets []Electricity
	allUser, exist := getAllUser(stub)
	if !exist {
		return allAssets, false
	}
	for _, user := range allUser.Users {
		if user.Sells!=nil{
			for _, asset := range user.Sells {
				if asset.State {
					allAssets = append(allAssets, asset)
				}
			}
		}
	}
	return allAssets, true
}
func getAllOrders(stub shim.ChaincodeStubInterface) ([]Electricity, bool) {
	var allOrders []Electricity
	allUser, exist := getAllUser(stub)
	if !exist {
		return allOrders, false
	}
	for _, user := range allUser.Users {
		if user.Purchases!=nil{
			for _, asset := range user.Purchases {
				if asset.State {
					allOrders = append(allOrders, asset)
				}
			}
		}
	}
	return allOrders, true
}
func (t *ElectricityChanCode) userRegister(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	var user User
	allUser,_:=getAllUser(stub)
	err := json.Unmarshal([]byte(args[0]), &user)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}
	_, exist := getUser(stub, user.Id)
	if exist {
		return shim.Error("要添加的用户已存在")
	}
	_, bl := putUser(stub, user)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}
	allUser.Users = append(allUser.Users, user)
	_, bl = putAllUser(stub, allUser)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("用户信息添加成功"))
}
func (t *ElectricityChanCode)userTopUp(stub shim.ChaincodeStubInterface, args []string)peer.Response{ //用户充值
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}
	var userTemp User
	err := json.Unmarshal([]byte(args[0]), &userTemp)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}
	user,exist:=getUser(stub,userTemp.Id)
	if !exist{
		return shim.Error("找不到指定用户！")
	}
	user.Surplus=userTemp.Surplus
	_, judge := putUser(stub, user)
	if !judge {
		return shim.Error("保存信息时发生错误")
	}
	allUser,exist:=getAllUser(stub);if !exist{
		return shim.Error("找不到总用户")
	}
	for i,_:=range allUser.Users{
		if allUser.Users[i].Id==user.Id{
			allUser.Users[i]=user
			break
		}
	}
	_,judge=putAllUser(stub,allUser)
	if !judge{
		return shim.Error("总用户信息存储失败存储失败！")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("用户充值成功！"))
}
func (t *ElectricityChanCode) assetsSell(stub shim.ChaincodeStubInterface, args []string) peer.Response { //电力资源发布
	if len(args) != 2 {
		return shim.Error("参数个数不符合要求")
	}
	var asset Electricity
	err := json.Unmarshal([]byte(args[0]), &asset)
	if err != nil {
		return shim.Error("反序列化信息发生错误")
	}
	var exist, judge bool
	_, exist = getAsset(stub, asset.Id, asset.Date)
	if exist {
		return shim.Error("该电力资产编号已经存在")
	}
	_, judge = putAsset(stub, asset)
	if !judge {
		return shim.Error("电力数据存储失败！")
	}
	user, exist := getUser(stub, asset.CurrentOwnerId)
	if !exist {
		return shim.Error("找不到指定用户")
	}
	user.Sells = append(user.Sells, asset)
	_, judge = putUser(stub, user)
	if !judge {
		return shim.Error("电力数据存储失败！")
	}
	allUser, exist := getAllUser(stub)
	if !exist {
		return shim.Error("找不到总用户")
	}
	for i, _ := range allUser.Users {
		if allUser.Users[i].Id == user.Id {
			allUser.Users[i] = user
			break
		}
	}
	_, judge = putAllUser(stub, allUser)
	if !judge {
		return shim.Error("电力数据存储失败！")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("电力资产发布成功"))
}
func (t *ElectricityChanCode) assetPurchase(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("参数个数不符合要求")
	}
	var asset Electricity
	err := json.Unmarshal([]byte(args[0]), &asset)
	if err != nil {
		return shim.Error("反序列化信息发生错误")
	}
	var exist, judge bool
	user, exist := getUser(stub, asset.CurrentOwnerId)
	if !exist {
		return shim.Error("找不到指定用户")
	}
	user.Purchases = append(user.Purchases, asset)
	_, judge = putUser(stub, user)
	if !judge {
		return shim.Error("用户信息存储失败！")
	}
	allUser, exist := getAllUser(stub)
	if !exist {
		return shim.Error("找不到总用户")
	}
	for i, _ := range allUser.Users {
		if allUser.Users[i].Id == user.Id {
			allUser.Users[i] = user
			break
		}
	}
	_, judge = putAllUser(stub, allUser)
	if !judge {
		return shim.Error("电力数据存储失败！")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("电力资产求购订单发布成功"))
}
func (t *ElectricityChanCode) queryUser(stub shim.ChaincodeStubInterface, args []string) peer.Response { //查询用户相关信息
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	// 根据用户账号进行查询
	user, exist := getUser(stub, args[0])
	if !exist {
		return shim.Error("找不到指定用户")
	}
	result, err := json.Marshal(user)
	if err != nil {
		return shim.Error("序列化edu信息时发生错误")
	}
	return shim.Success(result)
	return shim.Success([]byte("用户信息查询成功"))
}
func (t *ElectricityChanCode) queryAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response { //查询电力资源相关信息(可溯源)
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	asset, exist := getAsset(stub, args[0], args[1])
	if !exist {
		return shim.Error("电力资产信息获取失败！")
	}
	// 获取历史变更数据
	key := asset.CurrentOwnerId + "#" + asset.Date
	iterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error("根据指定的身份证号码查询对应的历史变更数据失败")
	}
	defer iterator.Close()
	// 迭代处理
	var hisElectricity Electricity
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取edu的历史变更数据失败")
		}
		var historyRecord HistoryRecord
		historyRecord.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisElectricity)
		historyRecord.Electricity = hisElectricity
		asset.Records = append(asset.Records, historyRecord)
	}
	// 返回
	result, err := json.Marshal(asset)
	if err != nil {
		return shim.Error("序列化asset信息时发生错误")
	}
	return shim.Success(result)
}
func (t *ElectricityChanCode) queryAssetMatchRecords(stub shim.ChaincodeStubInterface, args []string) peer.Response { //查询电力资源相关信息(可溯源)
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	asset, exist := getAsset(stub, args[0], args[1])
	if !exist {
		return shim.Error("电力资产信息获取失败！")
	}
	var records []Electricity
	for i:=0;i<len(asset.MatchRecords);i++{
		records=append(records,asset.MatchRecords[i])
	}
	result, err := json.Marshal(records)
	if err != nil {
		return shim.Error("序列化asset信息时发生错误")
	}
	return shim.Success(result)
}
func (t *ElectricityChanCode) queryAllAssets(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 0 {
		return shim.Error("给定的参数个数不符合要求")
	}
	allAssets, exist := getAllAssets(stub)
	if !exist {
		return shim.Error("当前暂时还没有人发布资产或者资产已全部售出！")
	}
	result, err := json.Marshal(allAssets)
	if err != nil {
		return shim.Error("序列化allAssets信息时发生错误")
	}
	return shim.Success(result)
}
func (t *ElectricityChanCode) queryAllOrders(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 0 {
		return shim.Error("给定的参数个数不符合要求")
	}
	allOrders, exist := getAllOrders(stub)
	if !exist {
		return shim.Error("当前暂时还没有人发布资产或者资产已全部售出！")
	}
	result, err := json.Marshal(allOrders)
	if err != nil {
		return shim.Error("序列化allAssets信息时发生错误")
	}
	return shim.Success(result)
}
func (t *ElectricityChanCode) queryOrderMatchRecords(stub shim.ChaincodeStubInterface, args []string) peer.Response { //查询电力资源相关信息(可溯源)
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	user, exist := getUser(stub, args[0])
	if !exist {
		return shim.Error("电力资产信息获取失败！")
	}
	var records []Electricity
	for i:=0;i<len(user.Purchases);i++{
		if user.Purchases[i].Id==args[1]{
			for j:=0;j<len(user.Purchases[i].MatchRecords);j++{
				records=append(records,user.Purchases[i].MatchRecords[j])
			}
			break
		}
	}
	result, err := json.Marshal(records)
	if err != nil {
		return shim.Error("序列化asset信息时发生错误")
	}
	return shim.Success(result)
}
func (t *ElectricityChanCode)cancelSell(stub shim.ChaincodeStubInterface, args []string)peer.Response{
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	var asset Electricity
	err:=json.Unmarshal([]byte(args[0]),&asset)
	if err!=nil{
		return shim.Error("反序列化信息发生错误")
	}
	user,exist:=getUser(stub,asset.CurrentOwnerId);if !exist{
		return shim.Error("找不到总用户")
	}
	for i,_:=range user.Sells{
		if user.Sells[i].Id==asset.Id{
			if user.Sells[i].MatchRecords==nil{
				temp:=user.Sells[0:i]
				temp=append(temp,user.Sells[i+1:]...)
				user.Sells=temp
			}else{
				return shim.Error("当前资产已被人预定！")
			}
		}
	}
	_,judge:=putUser(stub,user)
	if !judge{
		return shim.Error("用户信息存储失败！")
	}
	allUser,exist:=getAllUser(stub);if !exist{
		return shim.Error("找不到总用户")
	}
	for i,_:=range allUser.Users{
		if allUser.Users[i].Id==user.Id{
			allUser.Users[i]=user
			break
		}
	}
	_,judge=putAllUser(stub,allUser)
	if !judge{
		return shim.Error("电力数据存储失败！")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("电力资产撤销成功"))
}
func (t *ElectricityChanCode)cancelPurchase(stub shim.ChaincodeStubInterface, args []string)peer.Response{
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	var asset Electricity
	err:=json.Unmarshal([]byte(args[0]),&asset)
	if err!=nil{
		return shim.Error("反序列化信息发生错误")
	}
	user,exist:=getUser(stub,asset.CurrentOwnerId);if !exist{
		return shim.Error("找不到总用户")
	}
	for i,_:=range user.Purchases{
		if user.Purchases[i].Id==asset.Id{
			temp:=user.Purchases[0:i]
			temp=append(temp,user.Purchases[i+1:]...)
			user.Purchases=temp
		}
	}
	_,judge:=putUser(stub,user)
	if !judge{
		return shim.Error("用户信息存储失败！")
	}
	allUser,exist:=getAllUser(stub);if !exist{
		return shim.Error("找不到总用户")
	}
	for i,_:=range allUser.Users{
		if allUser.Users[i].Id==user.Id{
			allUser.Users[i]=user
			break
		}
	}
	_,judge=putAllUser(stub,allUser)
	if !judge{
		return shim.Error("电力数据存储失败！")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("求购订单撤销成功"))
}
func (t *ElectricityChanCode)saveAssetMatchRecords(stub shim.ChaincodeStubInterface, args []string)peer.Response{
	if len(args) != 2 {
		return shim.Error("参数个数不符合要求")
	}
	var assetTemp Electricity
	err := json.Unmarshal([]byte(args[0]), &assetTemp)
	if err != nil {
		return shim.Error("反序列化信息发生错误")
	}
	var judge bool
	user, exist := getUser(stub, assetTemp.CurrentOwnerId)
	if !exist {
		return shim.Error("找不到指定用户")
	}
	asset,exist:=getAsset(stub,assetTemp.CurrentOwnerId,assetTemp.Date)
	if !exist {
		return shim.Error("找不到指定电力资产")
	}
	asset.State=assetTemp.State
	for i,_:=range assetTemp.MatchRecords{
		asset.MatchRecords=append(asset.MatchRecords,assetTemp.MatchRecords[i])
	}
	for i,_:=range user.Sells{
		if user.Sells[i].Id==asset.Id{
			user.Sells[i]=asset
		}
	}
	_, judge = putUser(stub, user)
	if !judge {
		return shim.Error("用户信息存储失败！")
	}
	_, judge = putAsset(stub, asset)
	if !judge {
		return shim.Error("电力资产信息存储失败！")
	}
	allUser, exist := getAllUser(stub)
	if !exist {
		return shim.Error("找不到总用户")
	}
	for i, _ := range allUser.Users {
		if allUser.Users[i].Id == user.Id {
			allUser.Users[i] = user
			break
		}
	}
	_, judge = putAllUser(stub, allUser)
	if !judge {
		return shim.Error("总的用户信息存储失败存储失败！")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("电力资产求购订单发布成功"))
}
func (t *ElectricityChanCode)saveOrderMatchRecords(stub shim.ChaincodeStubInterface, args []string)peer.Response{
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	var orderTemp Electricity
	err:=json.Unmarshal([]byte(args[0]),&orderTemp)
	if err!=nil{
		return shim.Error("反序列化信息发生错误")
	}
	user,exist:=getUser(stub,orderTemp.CurrentOwnerId);if !exist{
		return shim.Error("找不到指定用户")
	}
	for j,_:=range user.Purchases{
		if user.Purchases[j].Id==orderTemp.Id{
			user.Purchases[j].State=orderTemp.State
			user.Purchases[j].MatchRecords=append(user.Purchases[j].MatchRecords,orderTemp.MatchRecords[0])
			break
		}
	}
	_,judge:=putUser(stub,user)
	if !judge{
		return shim.Error("用户信息存储失败！")
	}
	allUser,exist:=getAllUser(stub)
	if !exist{
		return shim.Error("找不到总用户")
	}
	for j,_:=range allUser.Users{
		if allUser.Users[j].Id==user.Id{
			allUser.Users[j]=user
			break
		}
	}
	_,judge=putAllUser(stub,allUser)
	if !judge{
		return shim.Error("总的用户信息存储失败！")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("订单匹配记录保存成功！"))
}
func main() {
	err := shim.Start(new(ElectricityChanCode))
	if err != nil {
		fmt.Printf("启动electricityChaincode时发生错误: %s", err)
	}
}


