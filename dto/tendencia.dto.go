package dto

import "time"

type FiltroTendenciaConsumoDto struct {
	EntidadPublicador string     `query:"entidadPublicador"`
	EntidadConsumidor string     `query:"entidadConsumidor"`
	SistemaPublicador string     `query:"sistemaPublicador"`
	Servicio          string     `query:"servicio"`
	EntidadActual     string     `query:"entidadActual"`
	Tipo              string     `query:"tipo"`
	FechaInicio       *time.Time `query:"fechaInicio"`
	FechaFin          *time.Time `query:"fechaFin"`
}
