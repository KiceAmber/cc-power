package groupSignature

func GenerateMemberPrivateKey() *Cert {
	if GPK == nil {
		GPK = GenerateGroup(160, 512)
	}
	membersk := GeneratePrivateKeyofMemberInGroup(GPK)
	return membersk
}
