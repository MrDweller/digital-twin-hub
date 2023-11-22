package digitaltwinmodels

import physicaltwinmodels "github.com/MrDweller/digital-twin-hub/physical-twin-models"

type DigitalTwinModel struct {
	SensedProperties []SensedPropertyModel  `json:"sensedProperties"`
	ControlCommands  []ControllCommandModel `json:"controlCommands"`

	PhysicalTwinConnectionModel physicaltwinmodels.PhysicalTwinConnectionModel `json:"physicalTwinConnection"`
}
