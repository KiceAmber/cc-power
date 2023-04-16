package transactionSign

import (
	"encoding/json"
	"powerTrading/service"
	"powerTrading/web/global"
	"powerTrading/web/groupSignature"
)

func GenerateGroupSignature(users []service.User, orderId string, singleTransaction service.SingleTransaction) ([]byte, []byte, error) {
	var sellerSignature, buyerSignature groupSignature.SignatureByte
	assetId := singleTransaction.AssetId
	singleTransactionByte, err := json.Marshal(singleTransaction)
	if err != nil {
		return nil, nil, err
	}
	var sellerId, buyerId string
	for _, user := range users {
		if sellerId == "" {
			for _, asset := range user.Sells {
				if asset.Id == assetId {
					sellerId = user.Id
				}
			}
		}
		if buyerId == "" {
			for _, order := range user.Purchases {
				if order.Id == orderId {
					buyerId = user.Id
				}
			}
		}
		if sellerId != "" && buyerId != "" {
			break
		}
	}
	seller := global.QueryUser(sellerId)
	buyer := global.QueryUser(buyerId)
	sellerSign := groupSignature.Sign(groupSignature.GPK, seller.GroupPrivateKey, string(singleTransactionByte))
	sellerSignature = groupSignature.SignatureByte{
		T1:      sellerSign.T1.Bytes(),
		T2:      sellerSign.T2.Bytes(),
		T3:      sellerSign.T3.Bytes(),
		C:       sellerSign.C.Bytes(),
		SAplpha: sellerSign.SAplpha.Bytes(),
		SBeta:   sellerSign.SBeta.Bytes(),
		Sa:      sellerSign.Sa.Bytes(),
		SDelta1: sellerSign.SDelta1.Bytes(),
		SDelta2: sellerSign.SDelta2.Bytes(),
	}
	sellerSignatureByte, err := json.Marshal(sellerSignature)
	if err != nil {
		return nil, nil, err
	}
	buyerSign := groupSignature.Sign(groupSignature.GPK, buyer.GroupPrivateKey, string(singleTransactionByte))
	buyerSignature = groupSignature.SignatureByte{
		T1:      buyerSign.T1.Bytes(),
		T2:      buyerSign.T2.Bytes(),
		T3:      buyerSign.T3.Bytes(),
		C:       buyerSign.C.Bytes(),
		SAplpha: buyerSign.SAplpha.Bytes(),
		SBeta:   buyerSign.SBeta.Bytes(),
		Sa:      buyerSign.Sa.Bytes(),
		SDelta1: buyerSign.SDelta1.Bytes(),
		SDelta2: buyerSign.SDelta2.Bytes(),
	}
	buyerSignatureByte, err := json.Marshal(buyerSignature)
	if err != nil {
		return nil, nil, err
	}
	return sellerSignatureByte, buyerSignatureByte, nil
}
