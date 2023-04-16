package zeroProof

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/developerblockchain/zkrp/bulletproofs"
)

func GenerateProof(scopeStr string, limit int64) ([]byte, error) {
	params, err := bulletproofs.SetupGeneric(0, limit)
	if err != nil {
		return nil, err
	}
	scope, _ := strconv.ParseInt(scopeStr, 10, 64)
	scopeSecret := new(big.Int).SetInt64(scope)
	proof, err := bulletproofs.ProveGeneric(scopeSecret, params)
	if err != nil {
		return nil, err
	}
	proofByte, _ := json.Marshal(proof)
	return proofByte, nil
}

// func GetTimeInterval(datePast string, dateNow string) float64 {
// 	TimePast, _ := time.ParseInLocation("2006-01-02 15:04:05", datePast, time.Local)
// 	TimeNow, _ := time.ParseInLocation("2006-01-02 15:04:05", dateNow, time.Local)
// 	left := TimeNow.Sub(TimePast)
// 	return left.Minutes()
// }

// func LaunchAsset(app controller.Application, userId string, assetId string, scopeStr string, priceStr string) (bool, error) {
// 	var user service.User
// 	b, err := app.Setup.QueryUser(userId)
// 	if err != nil {
// 		return false, err
// 	}
// 	err = json.Unmarshal(b, &user)
// 	if err != nil {
// 		return false, err
// 	}
// 	dateNow := time.Now().Format("2006-01-02 15:04:05")
// 	if user.Sells != nil {
// 		datePast := user.Sells[len(user.Sells)-1].Date
// 		timeInterval := GetTimeInterval(datePast, dateNow)
// 		if timeInterval < 60 {
// 			return false, nil
// 		}
// 	}
// 	asset := service.Electricity{
// 		Id:             assetId,
// 		Type:           "光伏发电",
// 		State:          true,
// 		CurrentOwnerId: userId,
// 		Date:           dateNow,
// 	}
// 	asset.Scope = cryptoCode.HomoEncryptData(scopeStr)
// 	asset.Price = cryptoCode.HomoEncryptData(priceStr)
// 	asset.Proof, err = GenerateProof(scopeStr, user.Limit)
// 	if err != nil {
// 		return false, err
// 	}
// 	_, err = app.Setup.AssetSell(asset)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func LaunchOrder(app controller.Application, userId string, orderId string, scopeStr string, priceStr string) (bool, error) {
// 	var user service.User
// 	b, err := app.Setup.QueryUser(userId)
// 	if err != nil {
// 		return false, err
// 	}
// 	err = json.Unmarshal(b, &user)
// 	if err != nil {
// 		return false, err
// 	}
// 	dateNow := time.Now().Format("2006-01-02 15:04:05")
// 	if user.Purchases != nil {
// 		datePast := user.Purchases[len(user.Purchases)-1].Date
// 		timeInterval := GetTimeInterval(datePast, dateNow)
// 		if timeInterval < 60 {
// 			return false, nil
// 		}
// 	}
// 	asset := service.Electricity{
// 		Id:             orderId,
// 		State:          true,
// 		CurrentOwnerId: userId,
// 		Date:           dateNow,
// 	}
// 	asset.Scope = cryptoCode.HomoEncryptData(scopeStr)
// 	asset.Price = cryptoCode.HomoEncryptData(priceStr)
// 	asset.Proof, err = GenerateProof(scopeStr, user.Limit)
// 	if err != nil {
// 		return false, err
// 	}
// 	_, err = app.Setup.AssetPurchase(asset)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

//func main() {
//	// Set up the range, [18, 200) in this case.
//	// We want to prove that we are over 18, and less than 200 years old.
//	// This information is shared between the prover and the verifier.
//	params, _ := bulletproofs.SetupGeneric(18, 200)
//
//	// Our secret age is 40
//	bigSecret := new(big.Int).SetInt64(int64(40))
//
//	// Create the zero-knowledge range proof
//	proof, _ := bulletproofs.ProveGeneric(bigSecret, params)
//
//	// Encode the proof to JSON
//	jsonEncoded, _ := json.Marshal(proof)
//
//	// It this stage, the proof is passed to the verifier, possibly over a network.
//
//	// Decode the proof from JSON
//	var decodedProof bulletproofs.ProofBPRP
//	_ = json.Unmarshal(jsonEncoded, &decodedProof)
//
//	// Verify the proof
//	ok, _ := decodedProof.Verify()
//
//	if ok == true {
//		println("Age verified to be [18, 200)")
//	}
//}
