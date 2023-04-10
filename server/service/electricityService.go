package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func InitConsumerParameters(consumerParameters []Parameters, index int) []Parameters {
	if consumerParameters != nil {
		return consumerParameters
	}
	consumerParameters = []Parameters{
		{0.067, 64, 146},
		{0.047, 79, 483},
		{0.047, 71, 750},
		{0.053, 62, 350},
		{0.082, 65, -782},
		{0.052, 83, 10},
		{0.087, 63, 13},
		{0.057, 81, 480},
		{0.05, 73, 493},
		{0.052, 69, 237},
		{0.071, 62, 1020},
		{0.064, 79, 411},
		{0.057, 60, 371},
		{0.082, 80, 462},
		{0.069, 78, 336},
		{0.069, 70, 208},
		{0.086, 62, 421},
		{0.054, 70, 309},
		{0.078, 66, 425},
		{0.081, 70, 14},
		{0.059, 71, 1656},
		{0.089, 80, 251},
		{0.067, 63, 86},
		{0.055, 75, 510},
		{0.082, 66, 459},
		{0.088, 70, 357},
		{0.076, 81, 155},
		{0.084, 61, 751},
		{0.077, 76, 209},
		{0.051, 79, 17},
		{0.087, 69, 321},
		{0.063, 62, 158},
		{0.059, 71, 57},
	}
	index = 0
	return consumerParameters
}

func InitProducerParameters(producerParameters []Parameters, index int) []Parameters {
	if producerParameters != nil {
		return producerParameters
	}
	producerParameters = []Parameters{
		{0.077, 17, 1040},
		{0.065, 33, 646},
		{0.082, 29, 725},
		{0.043, 18, 682},
		{0.051, 20, 508},
		{0.063, 31, 687},
		{0.08, 23, 580},
		{0.059, 18, 564},
		{0.071, 21, 865},
		{0.075, 37, 1100},
		{0.085, 25, 792},
		{0.06, 17, 681},
		{0.057, 38, 846},
		{0.079, 28, 955},
		{0.054, 36, 582},
		{0.089, 38, 658},
		{0.047, 19, 1005},
	}
	index = 0
	return producerParameters
}

func (t *ServiceSetup) ConsumerRegister(user User) (string, error) {
	//初始化用电用户参数
	consumerParameters = InitConsumerParameters(consumerParameters, cIndex)
	user.An = consumerParameters[cIndex].An
	user.Bn = consumerParameters[cIndex].Bn
	user.Limit = consumerParameters[cIndex].Limit
	cIndex++
	eventID := "eventConsumerRegister"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将user对象序列化成为字节数组
	b, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "userRegister", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *ServiceSetup) ProducerRegister(user User) (string, error) {
	//初始化发电用户参数
	producerParameters = InitProducerParameters(producerParameters, pIndex)
	user.An = producerParameters[pIndex].An
	user.Bn = producerParameters[pIndex].Bn
	user.Limit = producerParameters[pIndex].Limit
	pIndex++
	eventID := "eventProducerRegister"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将user对象序列化成为字节数组
	b, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "userRegister", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

// UserTopUp 用户充值
func (t *ServiceSetup) UserTopUp(user User) (string, error) {
	eventID := "eventUserTopUp"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将user对象序列化成为字节数组
	b, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "userTopUp", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *ServiceSetup) QueryUser(userId string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryUser", Args: [][]byte{[]byte(userId)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) QueryAllAssets() ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryAllAssets", Args: [][]byte{}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) QueryAllOrders() ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryAllOrders", Args: [][]byte{}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

func (t *ServiceSetup) AssetSell(asset Electricity) (string, error) {
	eventID := "eventAssetSell"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将edu对象序列化成为字节数组
	bAsset, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("指定的asset对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "assetSell", Args: [][]byte{bAsset, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *ServiceSetup) AssetPurchase(asset Electricity) (string, error) {
	eventID := "eventAssetPurchase"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将edu对象序列化成为字节数组
	bAsset, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("指定的asset对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "assetPurchase", Args: [][]byte{bAsset, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *ServiceSetup) CancelSell(asset Electricity) (string, error) {
	eventID := "eventCancelSell"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("指定的asset对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "cancelSell", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *ServiceSetup) CancelPurchase(asset Electricity) (string, error) {
	eventID := "eventCancelPurchase"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("指定的asset对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "cancelPurchase", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

// 查询用户自己发布的电力资产
func (t *ServiceSetup) QueryAsset(userId string, date string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryAsset", Args: [][]byte{[]byte(userId), []byte(date)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) QueryAssetMatchRecords(userId string, date string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryAssetMatchRecords", Args: [][]byte{[]byte(userId), []byte(date)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) QueryOrderMatchRecords(userId string, orderId string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryOrderMatchRecords", Args: [][]byte{[]byte(userId), []byte(orderId)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) SaveAssetMatchRecords(asset Electricity) (string, error) {
	eventID := "eventSaveAssetMatchRecords"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("指定的asset对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "saveAssetMatchRecords", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *ServiceSetup) SaveOrderMatchRecords(order Electricity) (string, error) {
	eventID := "eventSaveOrderMatchRecords"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(order)
	if err != nil {
		return "", fmt.Errorf("指定的asset对象序列化时发生错误")
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "saveOrderMatchRecords", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}
