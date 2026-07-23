// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import (
	"sync"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// ClaseMemoria implementa ClaseRepository en memoria.
// Se usa en los tests de reglas de negocio como fake del repositorio real.
type ClaseMemoria struct {
	mu     sync.Mutex
	datos  map[uint]models.Clase
	nextID uint
}

func NuevaClaseMemoria() *ClaseMemoria {
	return &ClaseMemoria{datos: make(map[uint]models.Clase), nextID: 1}
}

func (r *ClaseMemoria) Crear(h *models.Clase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	h.ID = r.nextID
	r.nextID++
	r.datos[h.ID] = *h
	return nil
}

func (r *ClaseMemoria) ObtenerPorID(id uint) (models.Clase, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.datos[id]
	return h, ok
}

func (r *ClaseMemoria) Listar() ([]models.Clase, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := make([]models.Clase, 0, len(r.datos))
	for _, h := range r.datos {
		lista = append(lista, h)
	}
	return lista, nil
}

func (r *ClaseMemoria) Actualizar(h *models.Clase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.datos[h.ID]; !ok {
		return ErrRegistroNoExiste
	}
	r.datos[h.ID] = *h
	return nil
}
