package constants

type Estado string

const (
	EstadoExitoso Estado = "SUCCESS"
	EstadoGateway Estado = "GATEWAY_ERROR"
	EstadoFallido Estado = "FAILED"
)
