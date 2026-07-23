// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import (
	"sync"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// ClienteMemoria implementa ClienteRepository en memoria.
type ClienteMemoria struct {
	mu     sync.Mutex
	datos  map[uint]models.Cliente
	nextID uint
}

func NuevoClienteMemoria() *ClienteMemoria {
	return &ClienteMemoria{datos: make(map[uint]models.Cliente), nextID: 1}
}

func (r *ClienteMemoria) Crear(c *models.Cliente) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = r.nextID
	r.nextID++
	r.datos[c.ID] = *c
	return nil
}

func (r *ClienteMemoria) ObtenerPorID(id uint) (models.Cliente, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.datos[id]
	return c, ok
}

func (r *ClienteMemoria) Listar() ([]models.Cliente, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := make([]models.Cliente, 0, len(r.datos))
	for _, c := range r.datos {
		lista = append(lista, c)
	}
	return lista, nil
}
