package digitaltwinregistry

import serviceModels "github.com/MrDweller/service-registry-connection/models"

type DigitalTwinRegistry struct {
	Address         string
	Port            int
	CertificateInfo serviceModels.CertificateInfo
}
