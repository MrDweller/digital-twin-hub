package sensoranomalyhandler

type Work struct {
	WorkId string `json:"workId"`

	ProductId string `json:"productId"`
}

type WorkDTO struct {
	WorkId string `json:"workId"`

	ProductId string `json:"productId"`
	EventType string `json:"eventType"`
}
