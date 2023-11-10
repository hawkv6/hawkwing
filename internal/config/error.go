package config

type validationError struct {
	Namespace   string `json:"namespace"`
	Field       string `json:"field"`
	ActualValue string `json:"actual_value"`
	Message     string `json:"message"`
}
