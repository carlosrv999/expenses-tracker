package paymentmethod

type createRequest struct {
	MethodName string  `json:"method_name"`
	MethodType string  `json:"method_type"`
	Icon       *string `json:"icon"`
}

type updateRequest struct {
	MethodName string  `json:"method_name"`
	MethodType string  `json:"method_type"`
	Icon       *string `json:"icon"`
}
