package types

type OnboardUnit struct {
	OID  string  `json:"oid"`
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}
