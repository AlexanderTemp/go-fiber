package services

import (
	"context"
	"log"
	"restapi/dto"
	"restapi/repository"
)

type BitacoraService struct {
	bitacoraRepo *repository.BitacoraRepository
}

func NewBitacoraService(repo *repository.BitacoraRepository) *BitacoraService {
	return &BitacoraService{
		bitacoraRepo: repo,
	}
}

func (s *BitacoraService) ObtenerTendenciaConsumo(ctx context.Context, filtros dto.FiltroTendenciaConsumoDto) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic recuperado en ObtenerTendenciaConsumo:", r)
		}
	}()

	res, err := s.bitacoraRepo.ObtenerDatosTendencia(ctx, filtros)
	if err != nil {
		log.Println("Error al obtener tendencia de consumo:", err)
		return map[string]interface{}{
			"fechas":    []string{},
			"entidades": []interface{}{},
		}, err
	}

	return res, nil
}

func (s *BitacoraService) ObtenerDatosExtra(ctx context.Context, filtros dto.FiltroExtraDto) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic recuperado en ObtenerDatosExtra:", r)
		}
	}()

	res, err := s.bitacoraRepo.ObtenerDatosExtra(ctx, filtros)
	if err != nil {
		log.Println("Error al obtener datos extra:", err)
		return []map[string]interface{}{}, err
	}

	return res, nil
}
