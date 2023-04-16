package controller

import (
	"encoding/json"
	"powerTrading/service"
)

func CancelAsset(app Application, userId string, assetId string) (bool, error) {
	var records []service.Electricity
	b, err := app.Setup.QueryAssetMatchRecords(userId, assetId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &records)
	if err != nil {
		return false, err
	}
	if records != nil && len(records) != 0 {
		return false, nil
	}
	asset := service.Electricity{
		Id:             assetId,
		CurrentOwnerId: userId,
	}
	_, err = app.Setup.CancelSell(asset)
	for records != nil {
		_, err = app.Setup.CancelSell(asset)
	}
	return true, nil
}

func CancelOrder(app Application, userId string, orderId string) (bool, error) {
	var records, assets []service.Electricity
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return false, err
	}
	for _, order := range user.Purchases {
		if order.Id == orderId {
			if !order.State {
				return false, nil
			}
		}
	}
	b, err = app.Setup.QueryOrderMatchRecords(userId, userId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &records)
	if err != nil {
		return false, err
	}
	asset := service.Electricity{
		Id:             orderId,
		CurrentOwnerId: userId,
	}
	_, err = app.Setup.CancelPurchase(asset)
	for err != nil {
		_, err = app.Setup.CancelPurchase(asset)
	}
	if records != nil && len(records) != 0 {
		var asset, assetTemp service.Electricity
		for i, _ := range records {
			asset.Id = records[i].Id
			if asset.MatchRecords == nil {
				assetTemp.Id = orderId
				asset.MatchRecords = append(asset.MatchRecords, assetTemp)
			}
			assets = append(assets, asset)
		}
		for i, _ := range assets {
			_, err = app.Setup.DeleteOrderMatchRecords(assets[i])
			for err != nil {
				_, err = app.Setup.DeleteOrderMatchRecords(assets[i])
			}
		}

	}
	return true, nil
}
