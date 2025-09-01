package dto

type FiltroExtraDto struct {
	EntidadConsumidora *int  `query:"entidadConsumidora"`
	EntidadPublicadora *int  `query:"entidadPublicadora"`
	SistemaConsumidor  *int  `query:"sistemaConsumidor"`
	Servicio           *int  `query:"servicio"`
	Todo               *bool `query:"todo"`
}
