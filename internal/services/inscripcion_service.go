package services

import (
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// TAREA (CP2): Implemente InscripcionService con las 5 reglas de negocio.
//
// Las reglas están A LA VISTA en las pantallas (carpeta pantallas/) y los
// tests de acceptance/reglas_test.go las verifican una por una. Devuelva los
// errores de dominio de errores.go: los tests los comprueban con errors.Is.
//
// Reglas:
//   - NO cambie el nombre del tipo, del constructor ni las firmas de los métodos.
//   - Observe que el service recibe TRES repositories: necesita consultar
//     Clase y Cliente para validar, y actualizar Clase al retirar.
type InscripcionService struct {
	inscripciones   storage.InscripcionRepository
	clases storage.ClaseRepository
	clientes     storage.ClienteRepository
}

func NuevaInscripcionService(
	inscripciones storage.InscripcionRepository,
	clases storage.ClaseRepository,
	clientes storage.ClienteRepository,
) *InscripcionService {
	return &InscripcionService{
		inscripciones:   inscripciones,
		clases: clases,
		clientes:     clientes,
	}
}

// Crear registra un nuevo inscripcion aplicando R1, R2 y R3.
// TODO (R1): la clase debe existir y estar activa; el cliente debe existir.
// TODO (R2): la cantidad no puede superar el stock disponible de la clase.
// TODO (R3): calcule el total (observe en las pantallas cuándo aplica descuento).
// TODO: al crear, el stock de la clase se descuenta (mire la pantalla 01
// antes y después de crear una inscripcion; R5 es la operación inversa).
func (s *InscripcionService) Crear(a *models.Inscripcion) error {
	// TODO: implementar.
	return ErrNoImplementado
}

func (s *InscripcionService) ObtenerPorID(id uint) (models.Inscripcion, error) {
	// TODO: implementar.
	return models.Inscripcion{}, ErrNoImplementado
}

func (s *InscripcionService) Listar() ([]models.Inscripcion, error) {
	// TODO: implementar.
	return nil, ErrNoImplementado
}

// Retirar cancela una inscripcion aplicando R4 y R5.
// TODO (R4): solo se puede retirar una inscripcion en estado PENDIENTE.
// TODO (R5): al retirar, la cantidad se repone al stock de la clase.
func (s *InscripcionService) Retirar(id uint) error {
	// TODO: implementar.
	return ErrNoImplementado
}
