package repository

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"restapi/constants"
	"restapi/dto"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type BitacoraRepository struct {
	conn clickhouse.Conn
}

func NewBitacoraRepository(conn clickhouse.Conn) *BitacoraRepository {
	return &BitacoraRepository{conn: conn}
}

type ServicioData struct {
	Fecha         time.Time `json:"fecha"`
	Transacciones uint64    `json:"transacciones"`
}

type ServicioResultado struct {
	Servicio string         `json:"servicio"`
	Data     []ServicioData `json:"data"`
	Total    uint64         `json:"total"`
}

func (r *BitacoraRepository) ObtenerDatosExtra(ctx context.Context, filtros dto.FiltroExtraDto) ([]ServicioResultado, error) {
	condiciones := make([]string, 0, 4)

	if filtros.EntidadConsumidora != nil {
		condiciones = append(condiciones, fmt.Sprintf("entidadConsumidora = %d", *filtros.EntidadConsumidora))
	}
	if filtros.EntidadPublicadora != nil {
		condiciones = append(condiciones, fmt.Sprintf("entidadPublicadora = %d", *filtros.EntidadPublicadora))
	}
	if filtros.SistemaConsumidor != nil {
		condiciones = append(condiciones, fmt.Sprintf("sistemaConsumidor = %d", *filtros.SistemaConsumidor))
	}
	if filtros.Servicio != nil {
		condiciones = append(condiciones, fmt.Sprintf("servicio = %d", *filtros.Servicio))
	}

	whereClause := "WHERE 1 = 1"
	if len(condiciones) > 0 {
		whereClause = "WHERE " + strings.Join(condiciones, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			toString(servicio) AS servicioId,
			toDate(fechayHora) AS fecha,
			countDistinct(idTransaccion) AS transacciones
		FROM bitacora
		%s
		AND fechayHora >= today() - INTERVAL 30 DAY
		GROUP BY servicioId, fecha
		ORDER BY servicioId, fecha
	`, whereClause)

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error en query: %w", err)
	}
	defer rows.Close()

	dataPorServicio := make(map[string][]ServicioData)

	for rows.Next() {
		var servicioId string
		var fecha time.Time
		var transacciones uint64

		if err := rows.Scan(&servicioId, &fecha, &transacciones); err != nil {
			return nil, fmt.Errorf("error en scan: %w", err)
		}

		dataPorServicio[servicioId] = append(dataPorServicio[servicioId], ServicioData{
			Fecha:         fecha,
			Transacciones: transacciones,
		})
	}

	servicios := make([]ServicioResultado, 0, len(dataPorServicio))
	for servicio, data := range dataPorServicio {
		var total uint64
		for _, d := range data {
			total += d.Transacciones
		}
		servicios = append(servicios, ServicioResultado{
			Servicio: servicio,
			Data:     data,
			Total:    total,
		})
	}

	sort.Slice(servicios, func(i, j int) bool {
		return servicios[i].Total > servicios[j].Total
	})

	if len(servicios) > 4 {
		servicios = servicios[:4]
	}

	return servicios, nil
}

func (r *BitacoraRepository) ObtenerDatosTendencia(ctx context.Context, filtros dto.FiltroTendenciaConsumoDto) (map[string]interface{}, error) {
	condiciones := []string{}
	params := map[string]any{}

	var fechaInicio, fechaFin time.Time

	if filtros.FechaInicio != nil && filtros.FechaFin != nil {
		fechaInicio = *filtros.FechaInicio
		fechaFin = *filtros.FechaFin
		condiciones = append(condiciones, "fechayHora >= parseDateTimeBestEffort(@inicio) AND fechayHora <= parseDateTimeBestEffort(@fin)")
		params["inicio"] = fechaInicio
		params["fin"] = fechaFin
	} else {
		fechaFin = time.Now()
		fechaInicio = fechaFin.AddDate(0, 0, -29)
		condiciones = append(condiciones, "fechayHora >= now() - interval 30 day")
	}

	days := int(fechaFin.Sub(fechaInicio).Hours()/24) + 1
	rangoFechas := make([]time.Time, 0, days)
	for i := 0; i < days; i++ {
		rangoFechas = append(rangoFechas, fechaInicio.AddDate(0, 0, i))
	}

	condiciones = append(condiciones, fmt.Sprintf("estado != '%s'", constants.EstadoGateway))

	if filtros.EntidadActual != "" {
		condiciones = append(condiciones, fmt.Sprintf("entidadPublicadoraGobBo = '%s'", escapeString(filtros.EntidadActual)))
	} else {
		if filtros.EntidadPublicador != "" {
			condiciones = append(condiciones, fmt.Sprintf("entidadPublicadoraGobBo = '%s'", escapeString(filtros.EntidadPublicador)))
		}
		if filtros.EntidadConsumidor != "" {
			condiciones = append(condiciones, fmt.Sprintf("entidadConsumidoraGobBo = '%s'", escapeString(filtros.EntidadConsumidor)))
		}
	}

	if filtros.SistemaPublicador != "" {
		condiciones = append(condiciones, fmt.Sprintf("sistemaPublicadorGobBo = '%s'", escapeString(filtros.SistemaPublicador)))
	}
	if filtros.Servicio != "" {
		condiciones = append(condiciones, fmt.Sprintf("servicioId = '%s'", escapeString(filtros.Servicio)))
	}

	whereClause := ""
	if len(condiciones) > 0 {
		whereClause = "WHERE " + strings.Join(condiciones, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT
			idTransaccion,
			any(servicioId) as servicioId_1, 
			any(entidadPublicadoraGobBo) as entidadId,
			any(entidadConsumidoraGobBo) as consumidorId, 
			any(sistemaPublicadorGobBo) as sistemaId, 
			any(entidadPublicadoraSigla) as entidadSigla,
			any(entidadPublicadoraNombre) as entidadNombre,
			any(toDate(fechayHora)) as fecha
		FROM bitacora_gestion
		%s
		GROUP BY idTransaccion
	`, whereClause)

	log.Println(query)

	rows, err := r.conn.Query(ctx, query, params)
	if err != nil {
		log.Println("Error en query:", err)
		return map[string]interface{}{
			"fechas":    formatDates(rangoFechas),
			"entidades": []interface{}{},
		}, err
	}
	defer rows.Close()

	entidadesData := map[string]struct {
		Sigla   string
		Nombre  string
		Consumo map[time.Time]int
	}{}

	for rows.Next() {
		var idTransaccion, servicioId string
		var entidadId, consumidorId, sistemaId uint32
		var entidadSigla, entidadNombre string
		var fecha time.Time

		if err := rows.Scan(&idTransaccion, &servicioId, &entidadId, &consumidorId, &sistemaId, &entidadSigla, &entidadNombre, &fecha); err != nil {
			log.Println("Error al escanear fila:", err)
			continue
		}

		entidadIdStr := fmt.Sprintf("%d", entidadId)
		if _, ok := entidadesData[entidadIdStr]; !ok {
			entidadesData[entidadIdStr] = struct {
				Sigla   string
				Nombre  string
				Consumo map[time.Time]int
			}{
				Sigla:   entidadSigla,
				Nombre:  entidadNombre,
				Consumo: map[time.Time]int{},
			}
		}

		data := entidadesData[entidadIdStr]
		data.Consumo[fecha]++
		entidadesData[entidadIdStr] = data
	}

	resultado := map[string]any{
		"fechas":    formatDates(rangoFechas),
		"entidades": []map[string]any{},
	}

	for entidadId, data := range entidadesData {
		consumoArray := make([]int, 0, len(rangoFechas))
		for _, f := range rangoFechas {
			consumoArray = append(consumoArray, data.Consumo[f])
		}
		resultado["entidades"] = append(resultado["entidades"].([]map[string]any), map[string]any{
			"id":      entidadId,
			"sigla":   data.Sigla,
			"nombre":  data.Nombre,
			"consumo": consumoArray,
		})
	}

	entidades := resultado["entidades"].([]map[string]any)
	sort.Slice(entidades, func(i, j int) bool {
		return sumArray(entidades[i]["consumo"].([]int)) > sumArray(entidades[j]["consumo"].([]int))
	})
	if len(entidades) > 10 {
		resultado["entidades"] = entidades[:10]
	}

	return resultado, nil
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func formatDates(rango []time.Time) []string {
	out := make([]string, 0, len(rango))
	for _, f := range rango {
		out = append(out, f.Format("2006-01-02"))
	}
	return out
}

func sumArray(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}
