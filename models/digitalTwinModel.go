package models

import (
	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/google/uuid"
)

type DigitalTwinModel struct {
	CertificateId uuid.UUID `json:"certificateId"`

	SensedProperties        []SensedPropertyModel                      `json:"sensedProperties"`
	ControlCommands         []ControllCommandModel                     `json:"controlCommands"`
	AdditionalServiceModels []additionalservice.AdditionalServiceModel `json:"additionalServiceModel"`

	PhysicalTwinConnectionModel PhysicalTwinConnectionModel `json:"physicalTwinConnection"`
}
