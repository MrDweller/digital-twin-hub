package sensoranomalyhandler

type AnomalyHandlerService interface {
	HandleAnomaly(handleableAnomaly HandleableAnomaly) error
}
