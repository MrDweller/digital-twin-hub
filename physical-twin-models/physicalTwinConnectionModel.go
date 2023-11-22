package physicaltwinmodels

type PhysicalTwinConnectionModel struct {
	ConnectionType  PhysicalTwinConnectionType `json:"connectionType"`
	ConnectionModel map[string]any             `json:"connectionModel"`
}
