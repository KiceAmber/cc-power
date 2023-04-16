package controller

import (
	"encoding/json"
	"fmt"
	"powerTrading/service"
	"powerTrading/web/cryptoCode"
	"strconv"
)

type Record struct {
	Date  string `json:"date"`
	Scope string `json:"scope"`
	Cash  string `json:"cash"`
}

//查询购电用户所有交易记录
func QueryOrderTransactionRecords(app Application, userId string) (error, []Record) {
	var records []Record
	var record Record
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, nil
	}
	if user.Purchases == nil || len(user.Purchases) == 0 {
		return nil, nil
	}
	for _, order := range user.Purchases {
		if !order.State {
			for _, transaction := range order.MatchRecords {
				record.Scope = cryptoCode.HomoDecryptData(transaction.Scope)
				record.Cash = cryptoCode.HomoDecryptData(transaction.CashProof)
				record.Date = transaction.Date
				records = append(records, record)
			}
		}
	}
	return nil, records
}

//按照时间查询查询购电用户所有交易记录
func QueryOrderTransactionRecordsByTime(app Application, userId string, time string) (error, []Record) {
	var records []Record
	var record Record
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, nil
	}
	if user.Purchases == nil || len(user.Purchases) == 0 {
		return nil, nil
	}
	for _, order := range user.Purchases {
		if !order.State {
			for _, transaction := range order.MatchRecords {
				date := transaction.Date[:7]
				if date == time {
					record.Scope = cryptoCode.HomoDecryptData(transaction.Scope)
					record.Cash = cryptoCode.HomoDecryptData(transaction.CashProof)
					record.Date = transaction.Date
					records = append(records, record)
				}
			}
		}
	}
	return nil, records
}

//查询购电商所有交易记录
func QueryAssetTransactionRecords(app Application, userId string) (error, []Record) {
	var records []Record
	var record Record
	var transactions []service.SingleTransaction
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, nil
	}
	if user.Sells == nil || len(user.Sells) == 0 {
		return nil, nil
	}
	b, err = app.Setup.QueryAllAssetsTransactionRecords()
	for err != nil {
		_, err := app.Setup.QueryAllAssetsTransactionRecords()
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &transactions)
	if err != nil {
		return err, nil
	}
	if transactions == nil && len(transactions) == 0 {
		return nil, nil
	}
	for _, asset := range user.Sells {
		for _, transaction := range transactions {
			if transaction.AssetId == asset.Id {
				record.Scope = cryptoCode.HomoDecryptData(transaction.Scope)
				record.Cash = cryptoCode.HomoDecryptData(transaction.Cash)
				record.Date = transaction.Date
				records = append(records, record)
			}
		}
	}
	return nil, records
}

//按照时间查询购电用户所有交易记录
func QueryAssetTransactionRecordsByTime(app Application, userId string, time string) (error, []Record) {
	var records []Record
	var record Record
	var transactions []service.SingleTransaction
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, nil
	}
	if user.Sells == nil || len(user.Sells) == 0 {
		return nil, nil
	}
	b, err = app.Setup.QueryAllAssetsTransactionRecords()
	for err != nil {
		_, err := app.Setup.QueryAllAssetsTransactionRecords()
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &transactions)
	if err != nil {
		return err, nil
	}
	if transactions == nil || len(transactions) == 0 {
		return nil, nil
	}
	for _, asset := range user.Sells {
		for _, transaction := range transactions {
			if transaction.AssetId == asset.Id {
				date := transaction.Date[:7]
				if date == time {
					record.Scope = cryptoCode.HomoDecryptData(transaction.Scope)
					record.Cash = cryptoCode.HomoDecryptData(transaction.Cash)
					record.Date = transaction.Date
					records = append(records, record)
				}
			}
		}
	}
	return nil, records
}

//查询购电用户每月的金额支出
func QueryOrderTransactionCashMonth(app Application, userId string) (error, map[int64]int64) {
	bills := make(map[int64]int64, 12)
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, bills
	}
	if user.Purchases == nil || len(user.Purchases) == 0 {
		return err, bills
	}
	for _, order := range user.Purchases {
		if !order.State {
			for _, transaction := range order.MatchRecords {
				month := transaction.Date[5:7]
				monthInt, err := strconv.ParseInt(month, 10, 64)
				if err != nil {
					return err, bills
				}
				cashStr := cryptoCode.HomoDecryptData(transaction.CashProof)
				cash, err := strconv.ParseInt(cashStr, 10, 64)
				if err != nil {
					return err, bills
				}
				bills[monthInt] = bills[monthInt] + cash
			}
		}
	}
	return nil, bills
}

//查询供电商每月的金额收入
func QueryAssetTransactionCashByMonth(app Application, userId string) (error, map[int64]int64) {
	bills := make(map[int64]int64, 12)
	var transactions []service.SingleTransaction
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, bills
	}
	if user.Sells == nil || len(user.Sells) == 0 {
		return nil, bills
	}
	b, err = app.Setup.QueryAllAssetsTransactionRecords()
	for err != nil {
		_, err := app.Setup.QueryAllAssetsTransactionRecords()
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &transactions)
	if err != nil {
		return err, bills
	}
	if transactions == nil || len(transactions) == 0 {
		return nil, bills
	}
	for _, asset := range user.Sells {
		for _, transaction := range transactions {
			if transaction.AssetId == asset.Id {
				month := transaction.Date[5:7]
				monthInt, err := strconv.ParseInt(month, 10, 64)
				if err != nil {
					return err, bills
				}
				cashStr := cryptoCode.HomoDecryptData(transaction.Cash)
				cash, err := strconv.ParseInt(cashStr, 10, 64)
				if err != nil {
					return err, bills
				}
				bills[monthInt] = bills[monthInt] + cash
			}
		}
	}
	return nil, bills
}

// 查询购电用户每月的交易笔数
func QueryOrderTransactionTimesByMonth(app Application, userId string) (error, map[int64]int64) {
	bills := make(map[int64]int64, 12)
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, bills
	}
	if user.Purchases == nil || len(user.Purchases) == 0 {
		return err, bills
	}
	for _, order := range user.Purchases {
		if !order.State {
			for _, transaction := range order.MatchRecords {
				month := transaction.Date[5:7]
				monthInt, err := strconv.ParseInt(month, 10, 64)
				if err != nil {
					return err, bills
				}
				bills[monthInt] = bills[monthInt] + 1
			}
		}
	}
	return nil, bills
}

//查询供电商每月的交易笔数
func QueryAssetTransactionTimesByMonth(app Application, userId string) (error, map[int64]int64) {
	bills := make(map[int64]int64, 12)
	var transactions []service.SingleTransaction
	var user service.User
	b, err := app.Setup.QueryUser(userId)
	for err != nil {
		_, err := app.Setup.QueryUser(userId)
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err, bills
	}
	if user.Sells == nil || len(user.Sells) == 0 {
		return nil, bills
	}
	b, err = app.Setup.QueryAllAssetsTransactionRecords()
	for err != nil {
		_, err := app.Setup.QueryAllAssetsTransactionRecords()
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, &transactions)
	if err != nil {
		return err, bills
	}
	if transactions == nil || len(transactions) == 0 {
		return nil, bills
	}
	for _, asset := range user.Sells {
		for _, transaction := range transactions {
			if transaction.AssetId == asset.Id {
				month := transaction.Date[5:7]
				monthInt, err := strconv.ParseInt(month, 10, 64)
				if err != nil {
					return err, bills
				}
				bills[monthInt] = bills[monthInt] + 1
			}
		}
	}
	return nil, bills
}
