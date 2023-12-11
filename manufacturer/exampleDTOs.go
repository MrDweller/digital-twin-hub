package manufacturer

type SensedPropertiesDTO struct {
	ServiceDefinition string `json:"serviceDefinition" default:"temperature"`
	ServiceUri        string `json:"serviceUri" default:"/temperature"`
}

type ControllPropertiesDTO struct {
	ServiceDefinition string `json:"serviceDefinition" default:"lamp"`
	ServiceUri        string `json:"serviceUri" default:"/lamp"`
}

type ConnectionDTO struct {
	ConnectionType  string `json:"connectionType" default:"simple-CoAP"`
	ConnectionModel ConnectionModelDTO
}

type ConnectionModelDTO struct {
	Address string `json:"address" default:"localhost"`
	Port    int    `json:"port" default:"5000"`
}

type DigitalTwinModelDTO struct {
	SensedProperties       []SensedPropertiesDTO
	ControlCommands        []ControllPropertiesDTO
	PhysicalTwinConnection ConnectionDTO
}

type SystemDefinitionDTO struct {
	Address            string `json:"address" default:"localhost"`
	Port               int    `json:"port" default:"5000"`
	SystemName         string `json:"systemName" default:"my-digital-twin"`
	AuthenticationInfo string `json:"authenticationInfo"`
}
