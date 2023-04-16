package controller

import (
	"encoding/json"
	"math"
	"powerTrading/service"
	"powerTrading/web/cryptoCode"
	"sort"
	"strconv"
)

//市场匹配修改
var Max float64 = 999999999999999999999999

type Electricity struct {
	Id             string
	Scope          float64
	State          bool
	Price          float64
	CurrentOwnerId string
	Date           string
	An             float64
	Bn             float64
}
type ElectricitySlice []Electricity

func (slice ElectricitySlice) Len() int {
	return len(slice)
}
func (slice ElectricitySlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
func (slice ElectricitySlice) Less(i, j int) bool {
	return slice[j].Price < slice[i].Price
}
func CollectiveDecrypt(assets []service.Electricity, orders []service.Electricity, producerApp Application, consumerApp Application) ([]Electricity, []Electricity) {
	dAssets := make([]Electricity, len(assets))
	dOrders := make([]Electricity, len(orders))
	var user service.User
	for i := 0; i < len(assets); i++ {
		result, err := producerApp.Setup.QueryUser(assets[i].CurrentOwnerId)
		for err != nil {
			result, err = producerApp.Setup.QueryUser(assets[i].CurrentOwnerId)
		}
		json.Unmarshal(result, &user)
		dAssets[i].Id = assets[i].Id
		dAssets[i].Price, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(assets[i].Price), 64)
		dAssets[i].Scope, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(assets[i].Scope), 64)
		dAssets[i].An = user.An
		dAssets[i].Bn = user.Bn
		dAssets[i].CurrentOwnerId = assets[i].CurrentOwnerId
		dAssets[i].Date = assets[i].Date
		if assets[i].MatchRecords != nil && len(assets[i].MatchRecords) != 0 {
			for j, _ := range assets[i].MatchRecords {
				scope, _ := strconv.ParseFloat(cryptoCode.HomoDecryptData(assets[i].MatchRecords[j].Scope), 64)
				dAssets[i].Scope -= scope
			}
		}
	}
	for i := 0; i < len(orders); i++ {
		result, err := consumerApp.Setup.QueryUser(orders[i].CurrentOwnerId)
		for err != nil {
			result, err = producerApp.Setup.QueryUser(orders[i].CurrentOwnerId)
		}
		json.Unmarshal(result, &user)
		dOrders[i].Id = orders[i].Id
		dOrders[i].Price, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(orders[i].Price), 64)
		dOrders[i].Scope, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(orders[i].Scope), 64)
		dOrders[i].An = user.An
		dOrders[i].Bn = user.Bn
		dOrders[i].CurrentOwnerId = orders[i].CurrentOwnerId
		if orders[i].MatchRecords != nil && len(orders[i].MatchRecords) != 0 {
			for j, _ := range orders[i].MatchRecords {
				scope, _ := strconv.ParseFloat(cryptoCode.HomoDecryptData(orders[i].MatchRecords[j].Scope), 64)
				dOrders[i].Scope -= scope
			}
		}
	}
	return dAssets, dOrders
}
func ProducerCostComputation(n float64, a float64, b float64) float64 {
	return (0.5*a*n*n + b*n)
}
func ConsumerCostComputation(n float64, a float64, b float64) float64 {
	return (0.5*a*n*n - b*n)
}
func PHPTransactionCost(asset Electricity, order Electricity) float64 {
	if order.Scope > asset.Scope {
		cost := ProducerCostComputation(asset.Scope, asset.An, asset.Bn) + ConsumerCostComputation(asset.Scope, order.An, order.Bn)
		return cost
	}
	cost := ProducerCostComputation(order.Scope, asset.An, asset.Bn) + ConsumerCostComputation(order.Scope, order.An, order.Bn)
	return cost

}
func getCostMinimumIndex(costs []float64) int {
	var min float64
	min = Max
	var index int
	for i := 0; i < len(costs); i++ {
		if min > costs[i] {
			min = costs[i]
			index = i
		}
	}
	return index
}
func saveMatchRecords(i int, j int, assets []Electricity, orders []Electricity) (service.Electricity, service.Electricity) {
	var asset, order, cAsset, cOrder service.Electricity
	price := (assets[j].Price + orders[i].Price) / 2
	priceInt := int64(math.Ceil(price - 0.4))
	if orders[i].Scope > assets[j].Scope {
		cash := int64(math.Ceil(assets[j].Scope*price - 0.4))
		asset.Id = assets[j].Id
		asset.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Scope, 'g', -1, 64))
		asset.Price = cryptoCode.HomoEncryptData(strconv.FormatInt(priceInt, 10))
		asset.CashProof = cryptoCode.HomoEncryptData(strconv.FormatInt(cash, 10))
		order.Id = orders[i].Id
		order.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Scope, 'g', -1, 64))
		order.Price = cryptoCode.HomoEncryptData(strconv.FormatInt(priceInt, 10))
		order.CashProof = cryptoCode.HomoEncryptData(strconv.FormatInt(cash, 10))
		cOrder.Id = orders[i].Id
		cOrder.CurrentOwnerId = orders[i].CurrentOwnerId
		cOrder.MatchRecords = append(cOrder.MatchRecords, asset)
		cAsset.CurrentOwnerId = assets[j].CurrentOwnerId
		cAsset.Date = assets[j].Date
		cAsset.MatchRecords = append(cAsset.MatchRecords, order)
		orders[i].Scope = orders[i].Scope - assets[j].Scope
		assets[j].Scope = 0
		return cAsset, cOrder
	}
	cash := int64(math.Ceil(orders[i].Scope*price - 0.4))
	asset.Id = assets[j].Id
	asset.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(orders[i].Scope, 'g', -1, 64))
	asset.Price = cryptoCode.HomoEncryptData(strconv.FormatInt(priceInt, 10))
	asset.CashProof = cryptoCode.HomoEncryptData(strconv.FormatInt(cash, 10))
	order.Id = orders[i].Id
	order.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(orders[i].Scope, 'g', -1, 64))
	order.Price = cryptoCode.HomoEncryptData(strconv.FormatInt(priceInt, 10))
	order.CashProof = cryptoCode.HomoEncryptData(strconv.FormatInt(cash, 10))
	cOrder.Id = orders[i].Id
	cOrder.CurrentOwnerId = orders[i].CurrentOwnerId
	cOrder.MatchRecords = append(cOrder.MatchRecords, asset)
	cAsset.CurrentOwnerId = assets[j].CurrentOwnerId
	cAsset.Date = assets[j].Date
	cAsset.MatchRecords = append(cAsset.MatchRecords, order)
	if orders[i].Scope < assets[j].Scope {
		cAsset.State = true
		assets[j].Scope = assets[j].Scope - orders[i].Scope
		orders[i].Scope = 0
	} else {
		orders[i].Scope = 0
		assets[j].Scope = 0
	}
	return cAsset, cOrder
}
func checkCondition(order []Electricity, assets []Electricity) bool {
	var supply, need float64
	supply, need = 0, 0
	for i, _ := range assets {
		supply = supply + assets[i].Scope
	}
	for i, _ := range order {
		need = need + order[i].Scope
	}
	if supply < need {
		return false
	}
	return true
}
func StartMarketMatch(producerApp Application, consumerApp Application) {
	var cAssets, cOrders, assetsRecords, ordersRecords []service.Electricity
	bAssets, err := producerApp.Setup.QueryAllAssets()
	for err != nil {
		bAssets, err = producerApp.Setup.QueryAllAssets()
	}
	json.Unmarshal(bAssets, &cAssets)
	bOrders, err := consumerApp.Setup.QueryAllOrders()
	for err != nil {
		bOrders, err = consumerApp.Setup.QueryAllOrders()
	}
	json.Unmarshal(bOrders, &cOrders)
	if len(cAssets) == 0 || len(cOrders) == 0 {
		return
	}
	assets, orders := CollectiveDecrypt(cAssets, cOrders, producerApp, consumerApp)
	if assets == nil || orders == nil {
		return
	}
	sort.Sort(sort.Reverse(ElectricitySlice(assets)))
	sort.Sort(ElectricitySlice(orders))
	if checkCondition(orders, assets) {
		for i := 0; i < len(orders); i++ {
			for {
				var costs []float64
				for j := 0; j < len(assets); j++ {
					cost := PHPTransactionCost(assets[j], orders[i])
					costs = append(costs, cost)
				}
				j := getCostMinimumIndex(costs)
				cAsset, cOrder := saveMatchRecords(i, j, assets, orders)
				if assets[j].Scope == 0 {
					temp := assets[0:j]
					temp = append(temp, assets[j+1:]...)
					assets = temp
				}
				//fmt.Println(assets)
				//fmt.Println(orders)
				assetsRecords = append(assetsRecords, cAsset)
				ordersRecords = append(ordersRecords, cOrder)
				if orders[i].Scope == 0 {
					break
				}
			}
		}
	}
	if len(assetsRecords) != 0 && len(ordersRecords) != 0 {
		for i, _ := range assetsRecords {
			_, err = producerApp.Setup.SaveAssetMatchRecords(assetsRecords[i])
			for err != nil {
				_, err = producerApp.Setup.SaveAssetMatchRecords(assetsRecords[i])

			}
		}
		for j, _ := range ordersRecords {
			_, err = consumerApp.Setup.SaveOrderMatchRecords(ordersRecords[j])
			for err != nil {
				_, err = consumerApp.Setup.SaveOrderMatchRecords(ordersRecords[j])
			}
		}
	}
	return
}
