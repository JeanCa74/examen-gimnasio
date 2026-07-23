// ARCHIVO BLOQUEADO — NO MODIFICAR
//
// Las 5 reglas de negocio se verifican aquí usando los repositorios EN MEMORIA
// (ya implementados en el repo base) como fakes. Así, estos tests miden solo
// la lógica de su InscripcionService, sin depender de su implementación GORM.
package acceptance

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

type entornoReglas struct {
	svc          *services.InscripcionService
	clases *storage.ClaseMemoria
	clientes     *storage.ClienteMemoria
	inscripciones   *storage.InscripcionMemoria
	principal      models.Clase
	ana          models.Cliente
}

func nuevoEntornoReglas(t *testing.T) entornoReglas {
	t.Helper()
	hm := storage.NuevaClaseMemoria()
	cm := storage.NuevoClienteMemoria()
	am := storage.NuevaInscripcionMemoria()

	principal := models.Clase{Nombre: "Clase de spinning", PrecioUnitario: 8.5, Stock: 10, Activo: true}
	require.NoError(t, hm.Crear(&principal))
	ana := models.Cliente{Nombre: "Ana Zambrano", Cedula: "1310000001", Telefono: "0990000001"}
	require.NoError(t, cm.Crear(&ana))

	return entornoReglas{
		svc:          services.NuevaInscripcionService(am, hm, cm),
		clases: hm,
		clientes:     cm,
		inscripciones:   am,
		principal:      principal,
		ana:          ana,
	}
}

// R1: no se crea una inscripcion si la clase no existe o está inactiva,
// o si el cliente no existe.
func TestCP2_R1_ReferenciasValidas(t *testing.T) {
	e := nuevoEntornoReglas(t)

	a := models.Inscripcion{ClaseID: 99999, ClienteID: e.ana.ID, Cantidad: 1}
	require.ErrorIs(t, e.svc.Crear(&a), services.ErrReferenciaInvalida,
		"crear con una clase inexistente debe devolver ErrReferenciaInvalida")

	extra := models.Clase{Nombre: "Clase de boxeo", PrecioUnitario: 15, Stock: 3, Activo: false}
	require.NoError(t, e.clases.Crear(&extra))
	a = models.Inscripcion{ClaseID: extra.ID, ClienteID: e.ana.ID, Cantidad: 1}
	require.ErrorIs(t, e.svc.Crear(&a), services.ErrReferenciaInvalida,
		"crear con una clase INACTIVA debe devolver ErrReferenciaInvalida")

	a = models.Inscripcion{ClaseID: e.principal.ID, ClienteID: 99999, Cantidad: 1}
	require.ErrorIs(t, e.svc.Crear(&a), services.ErrReferenciaInvalida,
		"crear con un cliente inexistente debe devolver ErrReferenciaInvalida")
}

// R2: la cantidad no puede superar el stock disponible.
func TestCP2_R2_StockInsuficiente(t *testing.T) {
	e := nuevoEntornoReglas(t)

	a := models.Inscripcion{ClaseID: e.principal.ID, ClienteID: e.ana.ID, Cantidad: 11}
	require.ErrorIs(t, e.svc.Crear(&a), services.ErrStockInsuficiente,
		"pedir 11 unidades con stock 10 debe devolver ErrStockInsuficiente")
}

// R3: Total = Cantidad x PrecioUnitario, con 10% de descuento desde 5 unidades.
func TestCP2_R3_CalculoTotal(t *testing.T) {
	e := nuevoEntornoReglas(t)

	sinDescuento := models.Inscripcion{ClaseID: e.principal.ID, ClienteID: e.ana.ID, Cantidad: 3}
	require.NoError(t, e.svc.Crear(&sinDescuento),
		"crear una inscripcion válida no debe devolver error")
	require.InDelta(t, 25.50, sinDescuento.Total, 0.001,
		"3 x 8.50 = 25.50 (sin descuento)")
	require.Equal(t, models.EstadoPendiente, sinDescuento.Estado,
		"una inscripcion recién creada debe quedar en estado PENDIENTE")

	conDescuento := models.Inscripcion{ClaseID: e.principal.ID, ClienteID: e.ana.ID, Cantidad: 5}
	require.NoError(t, e.svc.Crear(&conDescuento))
	require.InDelta(t, 38.25, conDescuento.Total, 0.001,
		"5 x 8.50 = 42.50, con 10% de descuento = 38.25")
}

// R4: solo se puede retirar una inscripcion en estado PENDIENTE.
func TestCP2_R4_RetirarSoloPendiente(t *testing.T) {
	e := nuevoEntornoReglas(t)

	asistida := models.Inscripcion{
		ClaseID: e.principal.ID,
		ClienteID:     e.ana.ID,
		Cantidad:      1,
		Estado:        models.EstadoAsistida,
		Total:         8.5,
	}
	require.NoError(t, e.inscripciones.Crear(&asistida))
	require.ErrorIs(t, e.svc.Retirar(asistida.ID), services.ErrEstadoInvalido,
		"retirar una inscripcion ASISTIDA debe devolver ErrEstadoInvalido")

	require.ErrorIs(t, e.svc.Retirar(99999), services.ErrNoEncontrado,
		"retirar una inscripcion inexistente debe devolver ErrNoEncontrado")
}

// R5: al crear se descuenta el stock; al retirar, se repone.
func TestCP2_R5_ReposicionStock(t *testing.T) {
	e := nuevoEntornoReglas(t)

	a := models.Inscripcion{ClaseID: e.principal.ID, ClienteID: e.ana.ID, Cantidad: 3}
	require.NoError(t, e.svc.Crear(&a))

	h, ok := e.clases.ObtenerPorID(e.principal.ID)
	require.True(t, ok)
	require.Equal(t, uint(7), h.Stock,
		"al crear una inscripcion de 3 unidades, el stock debe bajar de 10 a 7")

	require.NoError(t, e.svc.Retirar(a.ID), "retirar una inscripcion PENDIENTE debe funcionar")

	retirada, ok := e.inscripciones.ObtenerPorID(a.ID)
	require.True(t, ok)
	require.Equal(t, models.EstadoRetirada, retirada.Estado,
		"tras retirar, la inscripcion debe quedar en estado RETIRADA")

	h, ok = e.clases.ObtenerPorID(e.principal.ID)
	require.True(t, ok)
	require.Equal(t, uint(10), h.Stock,
		"al retirar, las 3 unidades deben reponerse al stock (7 -> 10)")
}
