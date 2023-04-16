package controller

import (
	"encoding/json"
	"powerTrading/service"
	"powerTrading/web/cryptoCode"
	"powerTrading/web/zeroProof"
	"time"
)

func GetTimeInterval(datePast string, dateNow string) float64 {
	TimePast, _ := time.ParseInLocation("2006-01-02 15:04:05", datePast, time.Local)
	TimeNow, _ := time.ParseInLocation("2006-01-02 15:04:05", dateNow, time.Local)
	left := TimeNow.Sub(TimePast)
	return left.Minutes()
}

func LaunchAsset(app Application, userId string, assetId string, scopeStr string, priceStr string) (bool, error) {
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return false, err
	}
	dateNow := time.Now().Format("2006-01-02 15:04:05")
	if user.Sells != nil && len(user.Sells) != 0 {
		datePast := user.Sells[len(user.Sells)-1].Date
		timeInterval := GetTimeInterval(datePast, dateNow)
		if timeInterval < 60 {
			return false, nil
		}
	}
	asset := service.Electricity{
		Id:             assetId,
		Type:           "光伏发电",
		State:          true,
		CurrentOwnerId: userId,
		Date:           dateNow,
	}
	asset.Scope = cryptoCode.HomoEncryptData(scopeStr)
	asset.Price = cryptoCode.HomoEncryptData(priceStr)
	asset.Proof, err = zeroProof.GenerateProof(scopeStr, user.Limit)
	if err != nil {
		return false, err
	}
	_, err = app.Setup.AssetSell(asset)
	if err != nil {
		return false, err
	}
	return true, nil
}

func LaunchOrder(app Application, userId string, orderId string, scopeStr string, priceStr string) (bool, error) {
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return false, err
	}
	dateNow := time.Now().Format("2006-01-02 15:04:05")
	if user.Purchases != nil && len(user.Purchases) != 0 {
		datePast := user.Purchases[len(user.Purchases)-1].Date
		timeInterval := GetTimeInterval(datePast, dateNow)
		if timeInterval < 60 {
			return false, nil
		}
	}
	asset := service.Electricity{
		Id:             orderId,
		State:          true,
		CurrentOwnerId: userId,
		Date:           dateNow,
	}
	asset.Scope = cryptoCode.HomoEncryptData(scopeStr)
	asset.Price = cryptoCode.HomoEncryptData(priceStr)
	asset.Proof, err = zeroProof.GenerateProof(scopeStr, user.Limit)
	if err != nil {
		return false, err
	}
	_, err = app.Setup.AssetPurchase(asset)
	if err != nil {
		return false, err
	}
	return true, nil
}

func InitLaunchAsset(app Application, userId string, assetId string, scopeStr string, priceStr string, date string) (bool, error) {
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return false, err
	}
	dateNow := time.Now().Format("2006-01-02 15:04:05")
	if user.Sells != nil {
		datePast := user.Sells[len(user.Sells)-1].Date
		timeInterval := GetTimeInterval(datePast, dateNow)
		if timeInterval < 60 {
			return false, nil
		}
	}
	asset := service.Electricity{
		Id:             assetId,
		Type:           "光伏发电",
		State:          true,
		CurrentOwnerId: userId,
		Date:           date,
	}
	asset.Scope = cryptoCode.HomoEncryptData(scopeStr)
	asset.Price = cryptoCode.HomoEncryptData(priceStr)
	asset.Proof, err = zeroProof.GenerateProof(scopeStr, user.Limit)
	if err != nil {
		return false, err
	}
	_, err = app.Setup.AssetSell(asset)
	if err != nil {
		return false, err
	}
	return true, nil
}

func InitLaunchOrder(app Application, userId string, orderId string, scopeStr string, priceStr string, date string) (bool, error) {
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return false, err
	}
	dateNow := time.Now().Format("2006-01-02 15:04:05")
	if user.Purchases != nil {
		datePast := user.Purchases[len(user.Purchases)-1].Date
		timeInterval := GetTimeInterval(datePast, dateNow)
		if timeInterval < 60 {
			return false, nil
		}
	}
	asset := service.Electricity{
		Id:             orderId,
		State:          true,
		CurrentOwnerId: userId,
		Date:           date,
	}
	asset.Scope = cryptoCode.HomoEncryptData(scopeStr)
	asset.Price = cryptoCode.HomoEncryptData(priceStr)
	asset.Proof, err = zeroProof.GenerateProof(scopeStr, user.Limit)
	if err != nil {
		return false, err
	}
	_, err = app.Setup.AssetPurchase(asset)
	if err != nil {
		return false, err
	}
	return true, nil
}
