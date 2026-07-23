package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// TestInscripcionService_CrearConDescuento es un test PROPIO (no de acceptance/)
// de la capa service, usando los repositorios EN MEMORIA como fakes. Reproduce
// la fila "Luis Mero — Clase de yoga — 5 uds — $27.00 (−10%)" de la pantalla 03.
//
// Prueba tres efectos de Crear en un solo escenario:
//   - R3: el total lleva el 10% de descuento (5 × 6.00 = 30 → 27.00).
//   - una inscripción nace en estado PENDIENTE.
//   - R5 (lado crear): el stock de la clase baja por la cantidad reservada.
func TestInscripcionService_CrearConDescuento(t *testing.T) {
	claseRepo := storage.NuevaClaseMemoria()
	clienteRepo := storage.NuevoClienteMemoria()
	inscRepo := storage.NuevaInscripcionMemoria()

	clase := models.Clase{Nombre: "Clase de yoga", PrecioUnitario: 6.00, Stock: 10, Activo: true}
	require.NoError(t, claseRepo.Crear(&clase))
	cliente := models.Cliente{Nombre: "Luis Mero", Cedula: "1310000002"}
	require.NoError(t, clienteRepo.Crear(&cliente))

	svc := services.NuevaInscripcionService(inscRepo, claseRepo, clienteRepo)

	insc := models.Inscripcion{ClaseID: clase.ID, ClienteID: cliente.ID, Cantidad: 5}
	require.NoError(t, svc.Crear(&insc))

	require.InDelta(t, 27.00, insc.Total, 0.001, "5 × 6.00 = 30.00, con 10% de descuento = 27.00")
	require.Equal(t, models.EstadoPendiente, insc.Estado, "una inscripción recién creada nace PENDIENTE")

	actualizada, ok := claseRepo.ObtenerPorID(clase.ID)
	require.True(t, ok)
	require.Equal(t, uint(5), actualizada.Stock, "reservar 5 de 10 deja el stock en 5")
}
