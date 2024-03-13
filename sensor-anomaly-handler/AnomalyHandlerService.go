package sensoranomalyhandler

type AnomalyHandlerService interface {
	HandleAnomaly(handleableAnomaly HandleableAnomaly, work Work) error
}
