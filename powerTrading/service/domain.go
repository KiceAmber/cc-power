
package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"powerTrading/sdkInit"
	"time"
)

type User struct { //定义用户结构
	Id        string        `json:"id"` //用户ID
	Purchases []Electricity //用户所要购买的电力资源
	Sells     []Electricity //用户所发布的电力资源
	Surplus   []byte //用户余额
	An float64               `json:"an"`
	Bn float64               `json:"bn"`
	Limit int64               `json:"limit"`
}
type TotalUser struct {
	Users []User //区块连总的用户
}
type Electricity struct { //定义电力资产结构
	Id             string        `json:"id"`             //资产ID或者订单编号
	Type           string        `json:"type"`           //电力类型
	Scope          []byte        `json:"scope"`          //初始电量
	State          bool          `json:"state"`          //当前状态
	Price          []byte        `json:"price"`          //资产价格或者订单最低的可接受价格
	CurrentOwnerId string          `json:"currentOwnerId"` //该资产发布者或者该订单的发布者
	Date           string          `json:"date"`           //发布日期
	HisRecords     []HistoryRecord //电力资产交易记录
	MatchRecords   []Electricity //市场交易匹配记录
	Proof  []byte //零知识证明材料
	CashProof  []byte //金额相等性材料
}
type HistoryRecord struct { //定义电力资产历史记录数据结构
	TxId        string //区块连中存储电力资产的Id
	Electricity Electricity
}
type SingleTransaction struct{
	AssetId string `json:"asset_id"`
	BuyerSignature   []byte `json:"buyer_signature"`
	SellerSignature  []byte `json:"seller_signature"`
	Date     string     `json:"date"`
	Scope  []byte    `json:"scope"`
	Cash   []byte     `json:"cash"`
}
type AssetTransaction struct{
	AssetId string `json:"asset_id"`
	Records []SingleTransaction  //电力资产交易记录
}
type TotalTransactions struct {
	Records []SingleTransaction  //所有电力资产交易记录
}

var TransactionKeyId="41777342161ab8615eed1e3d199d4609094d631bdb878fb5062fc3d38a7f938f"
type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}
type Parameters struct {
	An float64               `json:"an"`
	Bn float64               `json:"bn"`
	Limit int64               `json:"limit"`
}
var consumerParameters []Parameters
var producerParameters []Parameters
var cIndex int
var pIndex int
//func getParameter(peer string)(float64,float64){
//	switch peer {
//	case "peer0.org1.example.com:7051":return 0.067,64
//	case "peer1.org1.example.com:7061":return 0.047,79
//	case "peer2.org1.example.com:7071":return 0.067,71
//	case "peer3.org1.example.com:7081":return 0.053,62
//	case "peer4.org1.example.com:7091":return 0.082,65
//	case "peer5.org1.example.com:7050":return 0.052,83
//	case "peer6.org1.example.com:7050":return 0.087,63
//	case "peer7.org1.example.com:7050":return 0.057,81
//	case "peer8.org1.example.com:7050":return 0.064,79
//	case "peer9.org1.example.com:7050":return 0.057,60
//	case "peer10.org1.example.com:7050":return 0.071,62
//	case "peer11.org1.example.com:7050":return 0.064,79
//	case "peer12.org1.example.com:7050":return 0.057,60
//	case "peer13.org1.example.com:7050":return 0.082,80
//	case "peer14.org1.example.com:7050":return 0.069,78
//	case "peer15.org1.example.com:7050":return 0.069,70
//	case "peer16.org1.example.com:7050":return 0.086,62
//	case "peer17.org1.example.com:7050":return 0.054,70
//	case "peer18.org1.example.com:7050":return 0.078,66
//	case "peer19.org1.example.com:7050":return 0.081,70
//	case "peer20.org1.example.com:7050":return 0.059,71
//	case "peer21.org1.example.com:7050":return 0.089,80
//	case "peer22.org1.example.com:7050":return 0.067,63
//	case "peer23.org1.example.com:7050":return 0.055,75
//	case "peer24.org1.example.com:7050":return 0.082,66
//	case "peer25.org1.example.com:7050":return 0.076,81
//	case "peer27.org1.example.com:7050":return 0.084,61
//	case "peer28.org1.example.com:7050":return 0.077,76
//	case "peer29.org1.example.com:7050":return 0.051,79
//	case "peer30.org1.example.com:7050":return 0.087,69
//	case "peer31.org1.example.com:7050":return 0.062,62
//	case "peer32.org1.example.com:7050":return 0.059,71
//	case "peer0.org2.example.com:8051":return 0.077,17
//	case "peer1.org2.example.com:8061":return 0.065,33
//	case "peer2.org2.example.com:8071":return 0.082,29
//	case "peer3.org2.example.com:8081":return 0.043,18
//	case "peer4.org2.example.com:7091":return 0.051,20
//	case "peer5.org2.example.com:7050":return 0.063,31
//	case "peer6.org2.example.com:7050":return 0.08,23
//	case "peer7.org2.example.com:7050":return 0.059,18
//	case "peer8.org2.example.com:7050":return 0.071,21
//	case "peer9.org2.example.com:7050":return 0.075,37
//	case "peer10.org2.example.com:7050":return 0.085,25
//	case "peer11.org2.example.com:7050":return 0.06,17
//	case "peer12.org2.example.com:7050":return 0.057,38
//	case "peer13.org2.example.com:7050":return 0.079,28
//	case "peer14.org2.example.com:7050":return 0.054,36
//	case "peer15.org2.example.com:7050":return 0.089,38
//	case "peer16.org2.example.com:7050":return 0.047,19
//	}
//	return 0,0
//}
func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error{
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent.SourceURL)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client
	return handler, nil
}