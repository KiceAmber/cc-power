package cryptoCode

import (
	"crypto/rand"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"math/big"
	"strconv"
)

var privKey *paillier.PrivateKey //定义服务端全局私钥

// 同态加密
func HomoEncryptData(dataStr string) []byte {
	if privKey == nil {
		privKey, _ = paillier.GenerateKey(rand.Reader, 128)
	}
	data, _ := strconv.ParseInt(dataStr, 10, 64)
	cData, _ := paillier.Encrypt(&privKey.PublicKey, new(big.Int).SetInt64(data).Bytes())
	return cData
}

//同态解密
func HomoDecryptData(cData []byte) string {
	if privKey == nil {
		privKey, _ = paillier.GenerateKey(rand.Reader, 128)
	}
	dData, _ := paillier.Decrypt(privKey, cData)
	dataStr := new(big.Int).SetBytes(dData).String()
	return dataStr
}

//func homoDecryptDataList(cDataList []service.Electricity)([]Electricity) {
//	if privKey==nil{
//		privKey, _ = paillier.GenerateKey(rand.Reader, 128)
//	}
//	var dataList []Electricity
//	for i,_:=range cDataList{
//		dataList[i].id=cDataList[i].Id
//		dataList[i].price,_=strconv.ParseFloat(homoDecryptData(cDataList[i].Price),64)
//		dataList[i].scope,_=strconv.ParseFloat(homoDecryptData(cDataList[i].Scope),64)
//	}
//	return dataList
//}
