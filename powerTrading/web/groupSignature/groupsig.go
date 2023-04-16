package groupSignature

import (
	"crypto/sha256"
	"fmt"
	"github.com/Nik-U/pbc"
)

var GPK *PrivateKey
// 群公钥
type Group struct {
	g1, g2, h, u, v, w *pbc.Element
	Pairing            *pbc.Pairing
}

// 群主私钥
type PrivateKey struct {
	*Group
	xi1, xi2, gamma *pbc.Element
}

// 群成员私钥, 群成员私钥中有保护群公钥
type Cert struct {
	*Group
	A *pbc.Element `json:"a"`
	X *pbc.Element  `json:"x"`
}

type Signature struct {
	T1      *pbc.Element `json:"t1"`
	T2      *pbc.Element `json:"t2"`
	T3      *pbc.Element `json:"t3"`
	C       *pbc.Element `json:"c"`
	SAplpha *pbc.Element `json:"salpha"`
	SBeta   *pbc.Element `json:"sbeta"`
	Sa      *pbc.Element `json:"sa"`
	SDelta1 *pbc.Element `json:"sdelta1"`
	SDelta2 *pbc.Element `json:"sdelta2"`
}

type SignatureByte struct {
	T1       []byte
	T2       []byte
	T3       []byte
	C        []byte
	SAplpha  []byte
	SBeta    []byte
	Sa       []byte
	SDelta1  []byte
	SDelta2  []byte
}
// 对A进行Linear Encryption 后得到的参数
type LE_Ts struct {
	T1 *pbc.Element
	T2 *pbc.Element
	T3 *pbc.Element
}

type BlindValues struct {
	ralpha  *pbc.Element
	rbeta   *pbc.Element
	rdelta1 *pbc.Element
	rdelta2 *pbc.Element
	rx      *pbc.Element
}

type Rs struct {
	R1 *pbc.Element
	R2 *pbc.Element
	R3 *pbc.Element
	R4 *pbc.Element
	R5 *pbc.Element
}

type Ss struct {
	Salpha  *pbc.Element
	Sbeta   *pbc.Element
	Sx      *pbc.Element
	Sdelta1 *pbc.Element
	Sdelta2 *pbc.Element
}

// 生成群
/*
  推荐默认值rbits = 160, qbits = 512
*/
func GenerateGroup(rbits uint32, qbits uint32) *PrivateKey {
	// 生成群主私钥
	privkey := new(PrivateKey)
	group := new(Group)
	privkey.Group = group

	params := pbc.GenerateA(rbits, qbits)
	Pairing := params.NewPairing()

	privkey.Pairing = Pairing

	// 生成参数g1, g2
	g1 := Pairing.NewG1().Rand()
	g2 := Pairing.NewG2().Rand()
	privkey.g1 = g1
	privkey.g2 = g2

	// 实现各个群参数

	// 参数gamma初始化
	privkey.gamma = privkey.Pairing.NewZr().Rand()
	// 参数xi1初始化
	privkey.xi1 = privkey.Pairing.NewZr().Rand()
	// 参数xi2初始化
	privkey.xi2 = privkey.Pairing.NewZr().Rand()
	// 参数h初始化
	privkey.h = privkey.Pairing.NewG1().Rand()

	// 参数u,v初始化
	xi1_invert := privkey.Pairing.NewZr().Invert(privkey.xi1)
	xi2_invert := privkey.Pairing.NewZr().Invert(privkey.xi2)
	privkey.u = privkey.Pairing.NewG1().PowZn(privkey.h, xi1_invert)
	privkey.v = privkey.Pairing.NewG1().PowZn(privkey.h, xi2_invert)

	// 参数w初始化
	privkey.w = privkey.Pairing.NewG2().PowZn(privkey.g2, privkey.gamma)

	return privkey

}

// 生成群中成员的私钥
func GeneratePrivateKeyofMemberInGroup(gpk *PrivateKey) *Cert {
	cert := new(Cert)

	cert.Group = gpk.Group
	// 生成x_i
	cert.X = gpk.Pairing.NewZr().Rand()

	// 生成A
	gammaPlusX := gpk.Pairing.NewZr().Add(gpk.gamma, cert.X)
	gammaPlusX_invert := gpk.Pairing.NewZr().Invert(gammaPlusX)
	cert.A = gpk.Pairing.NewG1().PowZn(gpk.g1, gammaPlusX_invert)

	return cert
}

func calculateLinearEncryptionArgumentsofA(cert *Cert, alpha, beta *pbc.Element) *LE_Ts {
	// 计算T1,T2,T3
	t1 := cert.Pairing.NewG1().PowZn(cert.u, alpha)
	t2 := cert.Pairing.NewG1().PowZn(cert.v, beta)
	alphaPlusBeta := cert.Pairing.NewZr().Add(alpha, beta)
	exponent := cert.Pairing.NewG1().PowZn(cert.h, alphaPlusBeta)
	t3 := cert.Pairing.NewG1().Mul(cert.A, exponent)

	return &LE_Ts{
		t1,
		t2,
		t3,
	}
}

func calculateHelperValues(cert *Cert, alpha, beta *pbc.Element) (*pbc.Element, *pbc.Element) {

	// 计算两个辅助值(helper values)
	delta1 := cert.Pairing.NewZr().Mul(cert.X, alpha)
	delta2 := cert.Pairing.NewZr().Mul(cert.X, beta)

	return delta1, delta2
}

func pickBlindValues(cert *Cert) *BlindValues {
	ralpha := cert.Pairing.NewZr().Rand()
	rbeta := cert.Pairing.NewZr().Rand()
	rdelta1 := cert.Pairing.NewZr().Rand()
	rdelta2 := cert.Pairing.NewZr().Rand()
	rx := cert.Pairing.NewZr().Rand()

	return &BlindValues{
		ralpha,
		rbeta,
		rdelta1,
		rdelta2,
		rx,
	}
}

func calculateRs(cert *Cert, Ts *LE_Ts, bvs *BlindValues) *Rs {
	// 计算R1
	r1 := cert.Pairing.NewG1().PowZn(cert.u, bvs.ralpha)
	// 计算R2
	r2 := cert.Pairing.NewG1().PowZn(cert.v, bvs.rbeta)
	// 计算R3
	temp1 := cert.Pairing.NewGT().Pair(Ts.T3, cert.g2)
	r3_e1 := cert.Pairing.NewGT().PowZn(temp1, bvs.rx)
	temp2 := cert.Pairing.NewGT().Pair(cert.h, cert.w)
	negRalpha := cert.Pairing.NewZr().Neg(bvs.ralpha)
	negRbeta := cert.Pairing.NewZr().Neg(bvs.rbeta)
	neg_add := cert.Pairing.NewZr().Add(negRalpha, negRbeta)
	r3_e2 := cert.Pairing.NewGT().PowZn(temp2, neg_add)
	temp3 := cert.Pairing.NewGT().Pair(cert.h, cert.g2)
	negRdelta1 := cert.Pairing.NewZr().Neg(bvs.rdelta1)
	negRdelta2 := cert.Pairing.NewZr().Neg(bvs.rdelta2)
	neg_add1 := cert.Pairing.NewZr().Add(negRdelta1, negRdelta2)
	r3_e3 := cert.Pairing.NewGT().PowZn(temp3, neg_add1)
	r3 := cert.Pairing.NewGT().Mul(cert.Pairing.NewGT().Mul(r3_e1, r3_e2), r3_e3)
	r4_left := cert.Pairing.NewG1().PowZn(Ts.T1, bvs.rx)
	r4_right := cert.Pairing.NewG1().PowZn(cert.u, negRdelta1)
	r4 := cert.Pairing.NewG1().Mul(r4_left, r4_right)
	r5_left := cert.Pairing.NewG1().PowZn(Ts.T2, bvs.rx)
	r5_right := cert.Pairing.NewG1().PowZn(cert.v, negRdelta2)
	r5 := cert.Pairing.NewG1().Mul(r5_left, r5_right)

	return &Rs{
		r1,
		r2,
		r3,
		r4,
		r5,
	}
}

func calculateChallenge(cert *Cert, ts *LE_Ts, rs *Rs, m string) *pbc.Element {
	var s string
	s += ts.T1.String()
	s += ts.T2.String()
	s += ts.T3.String()
	s += rs.R1.String()
	s += rs.R2.String()
	s += rs.R3.String()
	s += rs.R4.String()
	s += rs.R5.String()
	s += m
	c := cert.Pairing.NewZr().SetFromStringHash(s, sha256.New())

	return c
}

func calculateSs(cert *Cert, bvs *BlindValues, c, delta1, delta2, alpha, beta *pbc.Element) *Ss {
	ss := new(Ss)

	ss.Salpha = cert.Pairing.NewZr().Add(bvs.ralpha, cert.Pairing.NewZr().Mul(c, alpha))
	ss.Sbeta = cert.Pairing.NewZr().Add(bvs.rbeta, cert.Pairing.NewZr().Mul(c, beta))
	ss.Sx = cert.Pairing.NewZr().Add(bvs.rx, cert.Pairing.NewZr().Mul(c, cert.X))
	ss.Sdelta1 = cert.Pairing.NewZr().Add(bvs.rdelta1, cert.Pairing.NewZr().Mul(c, delta1))
	ss.Sdelta2 = cert.Pairing.NewZr().Add(bvs.rdelta2, cert.Pairing.NewZr().Mul(c, delta2))

	return ss
}

// 签名
/*
gpk : 群私钥
cert : 签名者私钥
M : 待签名信息
*/
func Sign(gpk *PrivateKey, cert *Cert, M string) *Signature {
	sig := new(Signature)
	// alpha, beta <- Z_{p}^{*}
	alpha := cert.Pairing.NewZr().Rand()
	beta := cert.Pairing.NewZr().Rand()

	le_Ts := calculateLinearEncryptionArgumentsofA(cert, alpha, beta)
	sig.T1 = le_Ts.T1
	sig.T2 = le_Ts.T2
	sig.T3 = le_Ts.T3
	delta1, delta2 := calculateHelperValues(cert, alpha, beta)

	blindvalues := pickBlindValues(cert)

	rs := calculateRs(cert, le_Ts, blindvalues)

	c := calculateChallenge(cert, le_Ts, rs, M)
	sig.C = c

	ss := calculateSs(cert, blindvalues, c, delta1, delta2, alpha, beta)
	sig.Sa = ss.Sx
	sig.SBeta = ss.Sbeta
	sig.SAplpha = ss.Salpha
	sig.SDelta1 = ss.Sdelta1
	sig.SDelta2 = ss.Sdelta2

	return sig
}

func Verify(gpk *PrivateKey, M string, sign *Signature) bool {
	rs_ := &Rs{}
	// 先计算R1_,R2_,R4_,R5_，最复杂的R3_放在最后计算
	uHatSalpha := gpk.Pairing.NewG1().PowZn(gpk.u, sign.SAplpha)
	neg_c := gpk.Pairing.NewZr().Neg(sign.C)
	t1HatNeg_c := gpk.Pairing.NewG1().PowZn(sign.T1, neg_c)
	rs_.R1 = gpk.Pairing.NewG1().Mul(uHatSalpha, t1HatNeg_c)

	vHatSbeta := gpk.Pairing.NewG1().PowZn(gpk.v, sign.SBeta)
	t2HatNeg_c := gpk.Pairing.NewG1().PowZn(sign.T2, neg_c)
	rs_.R2 = gpk.Pairing.NewG1().Mul(vHatSbeta, t2HatNeg_c)

	t1HatSx := gpk.Pairing.NewG1().PowZn(sign.T1, sign.Sa)
	neg_Sdelta1 := gpk.Pairing.NewZr().Neg(sign.SDelta1)
	uHatNeg_Sdelta1 := gpk.Pairing.NewG1().PowZn(gpk.u, neg_Sdelta1)
	rs_.R4 = gpk.Pairing.NewG1().Mul(t1HatSx, uHatNeg_Sdelta1)

	t2HatSx := gpk.Pairing.NewG1().PowZn(sign.T2, sign.Sa)
	neg_Sdelta2 := gpk.Pairing.NewZr().Neg(sign.SDelta2)
	vHatNeg_Sdelta2 := gpk.Pairing.NewG1().PowZn(gpk.v, neg_Sdelta2)
	rs_.R5 = gpk.Pairing.NewG1().Mul(t2HatSx, vHatNeg_Sdelta2)

	temp1 := gpk.Pairing.NewGT().Pair(sign.T3, gpk.g2)
	r3_e1 := gpk.Pairing.NewGT().PowZn(temp1, sign.Sa)
	temp2 := gpk.Pairing.NewGT().Pair(gpk.h, gpk.w)
	neg_Salpha := gpk.Pairing.NewZr().Neg(sign.SAplpha)
	neg_SBeta := gpk.Pairing.NewZr().Neg(sign.SBeta)
	sumOfTwoNeg := gpk.Pairing.NewZr().Add(neg_Salpha, neg_SBeta)
	r3_e2 := gpk.Pairing.NewGT().PowZn(temp2, sumOfTwoNeg)

	NegSdelta1PlusNegSdelta2 := gpk.Pairing.NewZr().Add(neg_Sdelta1, neg_Sdelta2)
	temp3 := gpk.Pairing.NewGT().Pair(gpk.h, gpk.g2)
	r3_e3 := gpk.Pairing.NewGT().PowZn(temp3, NegSdelta1PlusNegSdelta2)

	e_T3_w := gpk.Pairing.NewGT().Pair(sign.T3, gpk.w)
	e_g1_g2 := gpk.Pairing.NewGT().Pair(gpk.g1, gpk.g2)
	e_g1_g2_invert := gpk.Pairing.NewGT().Invert(e_g1_g2)
	temp4 := gpk.Pairing.NewGT().Mul(e_T3_w, e_g1_g2_invert)
	r3_e4 := gpk.Pairing.NewGT().PowZn(temp4, sign.C)

	inter_res1 := gpk.Pairing.NewGT().Mul(r3_e1, r3_e2)
	inter_res2 := gpk.Pairing.NewGT().Mul(inter_res1, r3_e3)
	inter_res3 := gpk.Pairing.NewGT().Mul(inter_res2, r3_e4)

	rs_.R3 = inter_res3

	// 开始验证签名
	var s string
	s += sign.T1.String()
	s += sign.T2.String()
	s += sign.T3.String()
	s += rs_.R1.String()
	s += rs_.R2.String()
	s += rs_.R3.String()
	s += rs_.R4.String()
	s += rs_.R5.String()
	s += M

	c_ := gpk.Pairing.NewZr().SetFromStringHash(s, sha256.New())
	if c_.Equals(sign.C) {
		fmt.Println("verify_sign true")
		return true
	} else {
		fmt.Println("verify_sign false")
		return false
	}
}

func Open(gpk *PrivateKey, sign *Signature) *pbc.Element {
	temp1 := gpk.Pairing.NewG1().PowZn(sign.T1, gpk.xi1)
	temp2 := gpk.Pairing.NewG1().PowZn(sign.T2, gpk.xi2)
	temp3 := gpk.Pairing.NewG1().Mul(temp1, temp2)
	A := gpk.Pairing.NewG1().Mul(sign.T3, gpk.Pairing.NewG1().Invert(temp3))

	return A
}
func main() {
	priv := GenerateGroup(160, 512)
	cert := GeneratePrivateKeyofMemberInGroup(priv)
	fmt.Printf("%v\n", cert)
	var str []byte = cert.A.Bytes()
	//A := priv.Pairing.NewG1().SetBytes(str)
	////if !ok {
	////	fmt.Printf("初始化错误\n")
	////}
	//if A.Equals(cert.A) {
	//	fmt.Printf("签名生产呢个成功\n")
	//}
	mes:="你好！"
	sign:=Sign(priv,cert,mes)
	str=sign.T1.Bytes()
	T1:=priv.Pairing.NewG1().SetBytes(str)
	if T1.Equals(sign.T1) {
		fmt.Printf("签名生产呢个成功\n")
	}
	A:=Open(priv,sign)
	if A.Equals(cert.A){
		fmt.Printf("签名生产呢个成功\n")
	}
	//mes:="你好！"
	//singn:=Sign(priv,cert,mes)
	//sign1:=Sign(priv,cert1,mes)
	//A:=Open(priv,singn)
	//A1:=Open(priv,sign1)
	//if A.Equals(cert1.A){
	//	fmt.Println("bad")
	//}
	//if A1.Equals(cert.A){
	//	fmt.Println("bad")
	//}
	//gpk := GenerateGroup(160, 512)
	//membersk := GeneratePrivateKeyofMemberInGroup(gpk)
	//membersk1:=GeneratePrivateKeyofMemberInGroup(gpk)
	//fmt.Println(membersk.A.String())
	//fmt.Println(membersk1)
	//b1,_:=json.Marshal(membersk)
	//b2,_:=json.Marshal(membersk1)
	//fmt.Println(b1)
	//fmt.Println(b2)
	//var s1 *Cert
	//err:=json.Unmarshal(b1,&s1)
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}
	//mes:="HelloWord"
	//sign := Sign(gpk, membersk, mes)
	//fmt.Println(sign)
	//signByte,err :=json.Marshal(*sign)
	//if err != nil {
	//	fmt.Println("err = ",err)
	//	return
	//}
	//var sign1 transactionSign
	//err =json.Unmarshal(signByte,&sign1)
	//
	//if err != nil {
	//	fmt.Println("err = ",err)
	//	return
	//}
	//fmt.Println(sign1)
	//fmt.Println(s1)
	//err=json.Unmarshal(b2,&s2)
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(s1)
	//fmt.Println(s2)
	//byte,_:=json.Marshal(membersk.A)
	//var mes string = "Hello world"
	//var mes1 string = "Hello world!"
	//sign1 := Sign(gpk, membersk1, mes1)
	//Verify(gpk, mes, sign)
	//Verify(gpk, mes, sign1)
	//A:=Open(gpk,sign)
	//byte2,_:=json.Marshal(A)
	//A1:=Open(gpk,sign1)
	//byte3,_:=json.Marshal(A1)
	//fmt.Println(A)
	//fmt.Println(A1)
	//fmt.Println(byte2)
	//fmt.Println(byte3)
	//if string(byte)==string(byte3){
	//	fmt.Println("good")
	//}
	//fmt.Println(string(byte))
	//fmt.Println(string(byte3))

}