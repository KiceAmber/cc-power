package marketMatch

import (
	"encoding/json"
	"fmt"
	"powerTrading/service"
	"powerTrading/web/controller"
	"powerTrading/web/cryptoCode"
	"sort"
	"strconv"
)

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
func CollectiveDecrypt(assets []service.Electricity, orders []service.Electricity, producerApp controller.Application, consumerApp controller.Application) ([]Electricity, []Electricity) {
	dAssets := make([]Electricity, len(assets))
	dOrders := make([]Electricity, len(orders))
	var user service.User
	for i := 0; i < len(assets); i++ {
		result, err := producerApp.Setup.QueryUser(assets[i].CurrentOwnerId)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			json.Unmarshal(result, &user)
		}
		dAssets[i].Id = assets[i].Id
		dAssets[i].Price, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(assets[i].Price), 64)
		dAssets[i].Scope, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(assets[i].Scope), 64)
		dAssets[i].An = user.An
		dAssets[i].Bn = user.Bn
		dAssets[i].CurrentOwnerId = assets[i].CurrentOwnerId
		dAssets[i].Date = assets[i].Date
		if assets[i].MatchRecords != nil {
			for j, _ := range assets[i].MatchRecords {
				scope, _ := strconv.ParseFloat(cryptoCode.HomoDecryptData(assets[i].MatchRecords[j].Scope), 64)
				dAssets[i].Scope -= scope
			}
		}
	}
	for i := 0; i < len(orders); i++ {
		result, err := consumerApp.Setup.QueryUser(assets[i].CurrentOwnerId)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			json.Unmarshal(result, &user)
		}
		dOrders[i].Id = orders[i].Id
		dOrders[i].Price, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(orders[i].Price), 64)
		dOrders[i].Scope, _ = strconv.ParseFloat(cryptoCode.HomoDecryptData(orders[i].Scope), 64)
		dOrders[i].An = user.An
		dOrders[i].Bn = user.Bn
		dOrders[i].CurrentOwnerId = orders[i].CurrentOwnerId
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
	if orders[i].Scope > assets[j].Scope {
		asset.Id = assets[j].Id
		asset.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Scope, 'g', -1, 64))
		asset.Price = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Price, 'g', -1, 64))
		order.Id = orders[i].Id
		order.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Scope, 'g', -1, 64))
		order.Price = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Price, 'g', -1, 64))
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
	asset.Id = assets[j].Id
	asset.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(orders[j].Scope, 'g', -1, 64))
	asset.Price = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Price, 'g', -1, 64))
	order.Id = orders[i].Id
	order.Scope = cryptoCode.HomoEncryptData(strconv.FormatFloat(orders[j].Scope, 'g', -1, 64))
	order.Price = cryptoCode.HomoEncryptData(strconv.FormatFloat(assets[j].Price, 'g', -1, 64))
	cOrder.Id = orders[i].Id
	cOrder.CurrentOwnerId = orders[i].CurrentOwnerId
	cOrder.MatchRecords = append(cOrder.MatchRecords, asset)
	cAsset.CurrentOwnerId = assets[j].CurrentOwnerId
	cAsset.Date = assets[j].Date
	cAsset.MatchRecords = append(cAsset.MatchRecords, order)
	if orders[i].Scope < assets[j].Scope {
		assets[j].Scope = assets[j].Scope - orders[i].Scope
		orders[i].Scope = 0
	} else {
		orders[i].Scope = 0
		assets[j].Scope = 0
	}
	return cAsset, cOrder
}
func checkCondition(order Electricity, assets []Electricity) bool {
	var supply, need float64
	for i, _ := range assets {
		supply = supply + assets[i].Scope
	}
	need = order.Scope
	if supply < need {
		return false
	}
	return true
}
func StartMarketMatch(producerApp controller.Application, consumerApp controller.Application) error {
	var cAssets, cOrders, assetsRecords, ordersRecords []service.Electricity
	bAssets, err := producerApp.Setup.QueryAllAssets()
	if err != nil {
		return err
	} else {
		json.Unmarshal(bAssets, &cAssets)
	}
	bOrders, err := consumerApp.Setup.QueryAllOrders()
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		json.Unmarshal(bOrders, &cOrders)
	}
	assets, orders := CollectiveDecrypt(cAssets, cOrders, producerApp, consumerApp)
	if assets == nil || orders == nil {
		return nil
	}
	sort.Sort(sort.Reverse(ElectricitySlice(assets)))
	sort.Sort(ElectricitySlice(orders))
	for i := 0; i < len(orders); i++ {
		if checkCondition(orders[i], assets) {
			for {
				var costs []float64
				for j := 0; j < len(assets); j++ {
					cost := PHPTransactionCost(assets[j], orders[i])
					costs = append(costs, cost)
				}
				index := getCostMinimumIndex(costs)
				cAsset, cOrder := saveMatchRecords(i, index, assets, orders)
				if assets[index].Scope == 0 {
					temp := assets[0:index]
					temp = append(temp, assets[index+1:]...)
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
	if assetsRecords != nil && ordersRecords != nil {
		for i, _ := range assetsRecords {
			_, err = producerApp.Setup.SaveAssetMatchRecords(assetsRecords[i])
			if err != nil {
				return err
			}
		}
		for j, _ := range ordersRecords {
			_, err = consumerApp.Setup.SaveOrderMatchRecords(ordersRecords[j])
			if err != nil {
				return err
			}
		}
	}
	return nil
}