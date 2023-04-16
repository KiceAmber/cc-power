package controller

import (
	"encoding/json"
	"fmt"
	"powerTrading/service"
	"powerTrading/web/cryptoCode"
	"powerTrading/web/transactionSign"
	"strconv"
	"time"

	"github.com/developerblockchain/zkrp/bulletproofs"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
)

//交易的修改
func StartTransaction(app Application, userId string, orderId string) (int, error) {
	var records, assets []service.Electricity
	var user, producer service.User
	var asset, assetTemp, order service.Electricity
	var users []service.User
	var proof bulletproofs.ProofBPRP
	state := 0
	b, err := app.Setup.QueryOrderMatchRecords(userId, orderId)
	for err != nil {
		b, err = app.Setup.QueryOrderMatchRecords(userId, orderId)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	err = json.Unmarshal(b, &records)
	if err != nil {
		return 0, err
	}
	if records == nil || len(records) == 0 {
		return 0, nil
	}
	b, err = app.Setup.QueryUser(userId)
	for err != nil {
		b, err = app.Setup.QueryUser(userId)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return 0, err
	}

	b, err = app.Setup.QueryAllUsers()
	for err != nil {
		b, err = app.Setup.QueryAllUsers()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	err = json.Unmarshal(b, &users)
	if err != nil {
		return 0, err
	}
	for i, _ := range records {
		asset.Id = records[i].Id
		if asset.MatchRecords == nil && len(asset.MatchRecords) == 0 {
			assetTemp.Id = orderId
			asset.MatchRecords = append(asset.MatchRecords, assetTemp)
		}
		assets = append(assets, asset)
	}

	for _, order := range user.Purchases {
		if order.Id == orderId {
			err = json.Unmarshal(order.Proof, &proof)
			judge, err := proof.Verify()
			if err != nil {
				return 0, err
			}
			if !judge {
				_, err = app.Setup.CancelPurchase(order)
				for err != nil {
					_, err = app.Setup.CancelPurchase(order)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
				for i, _ := range assets {
					_, err = app.Setup.DeleteOrderMatchRecords(assets[i])
					for err != nil {
						_, err = app.Setup.DeleteOrderMatchRecords(assets[i])
						if err != nil {
							fmt.Println(err.Error())
						}
					}
				}
				return 1, nil
			}
		}
	}

	for i, _ := range records {
		for _, user := range users {
			for _, asset := range user.Sells {
				if asset.Id == records[i].Id {
					records[i].CurrentOwnerId = asset.CurrentOwnerId
					err = json.Unmarshal(asset.Proof, &proof)
					judge, err := proof.Verify()
					if err != nil {
						return 0, err
					}
					if !judge {
						_, err = app.Setup.CancelSell(asset)
						for err != nil {
							_, err = app.Setup.CancelSell(asset)
							if err != nil {
								fmt.Println(err.Error())
							}
						}
						var order, orderTemp service.Electricity
						for j, _ := range asset.MatchRecords {
							order.Id = asset.MatchRecords[j].Id
							if len(order.MatchRecords) == 0 {
								orderTemp.Id = asset.Id
								order.MatchRecords = append(order.MatchRecords, orderTemp)
							}
							_, err = app.Setup.DeleteAssetMatchRecords(order)
							for err != nil {
								_, err = app.Setup.DeleteAssetMatchRecords(order)
								if err != nil {
									fmt.Println(err.Error())
								}
							}
						}
						state = 2
					}
					break
				}
			}
		}
	}
	if state != 0 {
		return state, nil
	}
	totalCashByte := cryptoCode.HomoEncryptData("0")
	for _, asset := range records {
		totalCashByte = paillier.AddCipher(&cryptoCode.PrivKey.PublicKey, totalCashByte, asset.CashProof)
	}
	totalCashStr := cryptoCode.HomoDecryptData(totalCashByte)
	totalCash, err := strconv.ParseInt(totalCashStr, 10, 64)
	if err != nil {
		return 0, err
	}
	surplusStr := cryptoCode.HomoDecryptData(user.Surplus)
	surplus, err := strconv.ParseInt(surplusStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if surplus < totalCash {
		return 3, nil
	} else {
		surplus = surplus - totalCash
	}

	order = service.Electricity{
		Id:             orderId,
		CurrentOwnerId: userId,
	}
	for _, asset := range records {
		transaction := service.SingleTransaction{
			AssetId: asset.Id,
			Date:    time.Now().Format("2006-01-02 15:04:05"),
			Scope:   asset.Scope,
			Cash:    asset.CashProof,
		}
		transaction.SellerSignature, transaction.BuyerSignature, err = transactionSign.GenerateGroupSignature(users, orderId, transaction)
		if err != nil {
			return 0, err
		}
		_, err = app.Setup.SaveAssetTransactionRecords(transaction)
		for err != nil {
			_, err = app.Setup.SaveAssetTransactionRecords(transaction)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		b, err = app.Setup.QueryUser(asset.CurrentOwnerId)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(b, &producer)
		if err != nil {
			return 0, err
		}
		incomeStr := cryptoCode.HomoDecryptData(producer.Surplus)
		income, err := strconv.ParseInt(incomeStr, 10, 64)
		if err != nil {
			return 0, err
		}
		income = income + totalCash
		producer.Surplus = cryptoCode.HomoEncryptData(strconv.FormatInt(income, 10))
		_, err = app.Setup.UserTopUp(producer)
		for err != nil {
			_, err = app.Setup.UserTopUp(producer)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		orderTemp := service.Electricity{
			Date: transaction.Date,
		}
		order.MatchRecords = append(order.MatchRecords, orderTemp)
	}
	user.Surplus = cryptoCode.HomoEncryptData(strconv.FormatInt(surplus, 10))
	_, err = app.Setup.UserTopUp(user)
	for err != nil {
		_, err = app.Setup.UserTopUp(user)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	_, err = app.Setup.UpdateOrderState(order)
	for err != nil {
		_, err = app.Setup.UpdateOrderState(order)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return 4, nil
}
