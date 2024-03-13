package models

import (
	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
)

type DigitalTwinModel struct {
	CertificateId string `json:"certificateId"`

	SensedProperties        []SensedPropertyModel                      `json:"sensedProperties"`
	ControlCommands         []ControllCommandModel                     `json:"controlCommands"`
	AdditionalServiceModels []additionalservice.AdditionalServiceModel `json:"additionalServiceModel"`

	PhysicalTwinConnectionModel PhysicalTwinConnectionModel `json:"physicalTwinConnection"`
}
