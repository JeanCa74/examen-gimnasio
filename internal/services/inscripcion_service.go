package services

import (
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// Reglas del cálculo de total (pantalla 02: "Desde 5 unidades: 10% de descuento").
const (
	umbralDescuento = 5    // a partir de esta cantidad de cupos aplica el descuento
	factorDescuento = 0.90 // 10% de descuento = pagar el 90%
)

// InscripcionService contiene las 5 reglas de negocio de la inscripción.
//
// Recibe TRES repositorios porque una inscripción cruza tres entidades:
// consulta Clase y Cliente para validar las referencias (R1), lee el stock
// de la Clase (R2) y lo modifica al crear y al retirar (R5). Depende de las
// INTERFACES, no de las implementaciones GORM: por eso reglas_test.go le
// inyecta los repositorios en memoria.
type InscripcionService struct {
	inscripciones storage.InscripcionRepository
	clases        storage.ClaseRepository
	clientes      storage.ClienteRepository
}

func NuevaInscripcionService(
	inscripciones storage.InscripcionRepository,
	clases storage.ClaseRepository,
	clientes storage.ClienteRepository,
) *InscripcionService {
	return &InscripcionService{
		inscripciones: inscripciones,
		clases:        clases,
		clientes:      clientes,
	}
}

// Crear registra una inscripción nueva aplicando R1, R2, R3 y el lado "crear"
// de R5. El orden importa: primero se valida, y solo si todo pasa se toca el
// stock y se persiste.
func (s *InscripcionService) Crear(a *models.Inscripcion) error {
	// R1 — la clase debe existir Y estar activa.
	clase, ok := s.clases.ObtenerPorID(a.ClaseID)
	if !ok || !clase.Activo {
		return ErrReferenciaInvalida
	}

	// R1 — el cliente debe existir.
	if _, ok := s.clientes.ObtenerPorID(a.ClienteID); !ok {
		return ErrReferenciaInvalida
	}

	// R2 — no se puede inscribir más cupos de los que hay en stock.
	if a.Cantidad > clase.Stock {
		return ErrStockInsuficiente
	}

	// R3 — total = cantidad x precio, con descuento desde 5 unidades.
	a.Total = calcularTotal(a.Cantidad, clase.PrecioUnitario)

	// Una inscripción recién creada siempre nace PENDIENTE (pantalla 03).
	a.Estado = models.EstadoPendiente

	// R5 (lado crear) — reservar los cupos: el stock baja por la cantidad.
	clase.Stock -= a.Cantidad
	if err := s.clases.Actualizar(&clase); err != nil {
		return err
	}

	// Persistir la inscripción ya validada y con su total y estado.
	return s.inscripciones.Crear(a)
}

func (s *InscripcionService) ObtenerPorID(id uint) (models.Inscripcion, error) {
	a, ok := s.inscripciones.ObtenerPorID(id)
	if !ok {
		return models.Inscripcion{}, ErrNoEncontrado
	}
	return a, nil
}

func (s *InscripcionService) Listar() ([]models.Inscripcion, error) {
	return s.inscripciones.Listar()
}

// Retirar cancela una inscripción aplicando R4 y el lado "retirar" de R5.
func (s *InscripcionService) Retirar(id uint) error {
	// La inscripción debe existir (esto es lo que el handler mapea a 404).
	a, ok := s.inscripciones.ObtenerPorID(id)
	if !ok {
		return ErrNoEncontrado
	}

	// R4 — solo se retira lo que está PENDIENTE. Una ASISTIDA o una ya
	// RETIRADA no se puede volver a retirar (pantalla 03: el botón Retirar
	// solo está activo en la fila PENDIENTE).
	if a.Estado != models.EstadoPendiente {
		return ErrEstadoInvalido
	}

	// R5 (lado retirar) — devolver los cupos reservados al stock de la clase.
	if clase, ok := s.clases.ObtenerPorID(a.ClaseID); ok {
		clase.Stock += a.Cantidad
		if err := s.clases.Actualizar(&clase); err != nil {
			return err
		}
	}

	// La inscripción queda RETIRADA.
	a.Estado = models.EstadoRetirada
	return s.inscripciones.Actualizar(&a)
}

// calcularTotal aplica R3: precio total con 10% de descuento desde 5 unidades.
func calcularTotal(cantidad uint, precioUnitario float64) float64 {
	total := float64(cantidad) * precioUnitario
	if cantidad >= umbralDescuento {
		total *= factorDescuento
	}
	return total
}
