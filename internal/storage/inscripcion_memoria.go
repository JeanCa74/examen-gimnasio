// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import (
	"sync"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// InscripcionMemoria implementa InscripcionRepository en memoria.
// Se usa en los tests de reglas de negocio como fake del repositorio real.
type InscripcionMemoria struct {
	mu     sync.Mutex
	datos  map[uint]models.Inscripcion
	nextID uint
}

func NuevaInscripcionMemoria() *InscripcionMemoria {
	return &InscripcionMemoria{datos: make(map[uint]models.Inscripcion), nextID: 1}
}

func (r *InscripcionMemoria) Crear(a *models.Inscripcion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a.ID = r.nextID
	r.nextID++
	r.datos[a.ID] = *a
	return nil
}

func (r *InscripcionMemoria) ObtenerPorID(id uint) (models.Inscripcion, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.datos[id]
	return a, ok
}

func (r *InscripcionMemoria) Listar() ([]models.Inscripcion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := make([]models.Inscripcion, 0, len(r.datos))
	for _, a := range r.datos {
		lista = append(lista, a)
	}
	return lista, nil
}

func (r *InscripcionMemoria) Actualizar(a *models.Inscripcion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.datos[a.ID]; !ok {
		return ErrRegistroNoExiste
	}
	r.datos[a.ID] = *a
	return nil
}
