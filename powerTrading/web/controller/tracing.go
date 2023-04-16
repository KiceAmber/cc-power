package controller

import (
	"encoding/json"
	"fmt"
	"powerTrading/service"
	"powerTrading/web/global"
	"powerTrading/web/groupSignature"
	"powerTrading/web/model"
)

func Tracing(app Application, assetId string, date string) (*model.User, *model.User, error) { //根据签名进行身份溯源
	//查询电力资产所对应的交易
	var singleTransaction service.SingleTransaction
	b, err := app.Setup.QuerySingleAssetTransactionRecords(assetId, date)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(b, &singleTransaction)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(singleTransaction)
	//获取买方和卖方的群签名
	var buyerSign, sellerSign groupSignature.SignatureByte
	err = json.Unmarshal(singleTransaction.BuyerSignature, &buyerSign)
	if err != nil {
		return nil, nil, err
	}
	buyerSignature := &groupSignature.Signature{
		T1:      groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.T1),
		T2:      groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.T2),
		T3:      groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.T3),
		C:       groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.C),
		SAplpha: groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.SAplpha),
		SBeta:   groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.SBeta),
		Sa:      groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.Sa),
		SDelta1: groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.SDelta1),
		SDelta2: groupSignature.GPK.Pairing.NewG1().SetBytes(buyerSign.SDelta2),
	}

	err = json.Unmarshal(singleTransaction.SellerSignature, &sellerSign)
	if err != nil {
		return nil, nil, err
	}
	sellerSignature := &groupSignature.Signature{
		T1:      groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.T1),
		T2:      groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.T2),
		T3:      groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.T3),
		C:       groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.C),
		SAplpha: groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.SAplpha),
		SBeta:   groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.SBeta),
		Sa:      groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.Sa),
		SDelta1: groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.SDelta1),
		SDelta2: groupSignature.GPK.Pairing.NewG1().SetBytes(sellerSign.SDelta2),
	}
	buyerA := groupSignature.Open(groupSignature.GPK, buyerSignature)
	sellerA := groupSignature.Open(groupSignature.GPK, sellerSignature)
	fmt.Println(sellerSign)
	fmt.Println(buyerSign)
	//依次遍历寻找对应用户Id
	var buyer, seller *model.User
	allUers := global.QueryAllUser()
	fmt.Println(allUers)
	for _, user := range allUers {
		fmt.Println("==================", user.GroupPrivateKey.A)
		if user.GroupPrivateKey.A != nil {
			if user.GroupPrivateKey.A.Equals(buyerA) {
				buyer = user
			}
		}
		if user.GroupPrivateKey.A != nil {
			if user.GroupPrivateKey.A.Equals(sellerA) {
				seller = user
			}
		}
		if buyer != nil && seller != nil {
			break
		}
	}
	return buyer, seller, nil
}
