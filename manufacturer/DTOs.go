package manufacturer

type SensedPropertyDTO struct {
	ServiceDefinition  string `json:"serviceDefinition" default:"temperature"`
	ServiceUri         string `json:"serviceUri" default:"/temperature"`
	SensorEndpointMode string `json:"sensorEndpointMode" default:"INTERVAL_RETRIEVAL"`
	IntervalTime       int    `json:"intervalTime" default:"10"`
}

type ControllCommandDTO struct {
	ServiceDefinition string `json:"serviceDefinition" default:"lamp"`
	ServiceUri        string `json:"serviceUri" default:"/lamp"`
}

type ConnectionDTO struct {
	ConnectionType  string         `json:"connectionType" default:"simple-CoAP"`
	ConnectionModel map[string]any `json:"connectionModel" `
}

type AnomalyDTO struct {
	AnomalyType string `json:"anomalyType" default:"STUCK"`
}

type DigitalTwinDTO struct {
	SensedProperties       []SensedPropertyDTO
	ControlCommands        []ControllCommandDTO
	HandleableAnomalies    []AnomalyDTO
	PhysicalTwinConnection ConnectionDTO
}

type SystemDefinitionDTO struct {
	Address            string `json:"address" default:"localhost"`
	Port               int    `json:"port" default:"5000"`
	SystemName         string `json:"systemName" default:"my-digital-twin"`
	AuthenticationInfo string `json:"authenticationInfo"`
}
