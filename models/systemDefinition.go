package models

type SystemDefinition struct {
	Address    string `json:"address"`
	Port       int    `json:"port"`
	SystemName string `json:"systemName"`
}
