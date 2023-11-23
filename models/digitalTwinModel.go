package models

type DigitalTwinModel struct {
	SensedProperties []SensedPropertyModel  `json:"sensedProperties"`
	ControlCommands  []ControllCommandModel `json:"controlCommands"`

	PhysicalTwinConnectionModel PhysicalTwinConnectionModel `json:"physicalTwinConnection"`
}
