package context

type GenerateContext struct {
	DataHash  string  `json:"hash"`
	Threshold *string `json:"threshold"`
}

type ValidateContext struct {
	DataHash  string  `json:"hash"`
	Threshold *string `json:"threshold"`
	Work      string  `json:"work"`
}

type CancelContext struct {
	DataHash string `json:"hash"`
}
