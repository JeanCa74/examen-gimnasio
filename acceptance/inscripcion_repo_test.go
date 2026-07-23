// ARCHIVO BLOQUEADO — NO MODIFICAR
package acceptance

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// TestCP2_InscripcionRepositorioGORM verifica la implementación GORM de Inscripcion
// contra la interfaz bloqueada InscripcionRepository.
func TestCP2_InscripcionRepositorioGORM(t *testing.T) {
	db := nuevaDB(t)

	// Registros padre para que las claves foráneas tengan sentido.
	principal := models.Clase{Nombre: "Clase de spinning", PrecioUnitario: 8.5, Stock: 10, Activo: true}
	require.NoError(t, db.Create(&principal).Error)
	ana := models.Cliente{Nombre: "Ana Zambrano", Cedula: "1310000001", Telefono: "0990000001"}
	require.NoError(t, db.Create(&ana).Error)

	// La asignación fuerza el contrato: si las firmas no coinciden, no compila.
	var repo storage.InscripcionRepository = storage.NuevaInscripcionGORM(db)

	a := models.Inscripcion{
		ClaseID: principal.ID,
		ClienteID:     ana.ID,
		Cantidad:      2,
		Estado:        models.EstadoPendiente,
		Total:         17.0,
	}
	require.NoError(t, repo.Crear(&a), "Crear debe persistir la inscripcion sin error")
	require.NotZero(t, a.ID, "tras Crear, la inscripcion debe tener ID asignado")

	obtenido, ok := repo.ObtenerPorID(a.ID)
	require.True(t, ok, "ObtenerPorID debe encontrar la inscripcion recién creado")
	require.Equal(t, models.EstadoPendiente, obtenido.Estado)
	require.Equal(t, uint(2), obtenido.Cantidad)

	_, ok = repo.ObtenerPorID(99999)
	require.False(t, ok, "ObtenerPorID de un ID inexistente debe devolver ok=false")

	obtenido.Estado = models.EstadoAsistida
	require.NoError(t, repo.Actualizar(&obtenido), "Actualizar debe guardar los cambios")
	releido, ok := repo.ObtenerPorID(a.ID)
	require.True(t, ok)
	require.Equal(t, models.EstadoAsistida, releido.Estado,
		"tras Actualizar, el estado debe quedar persistido")

	lista, err := repo.Listar()
	require.NoError(t, err)
	require.Len(t, lista, 1, "Listar debe devolver el único inscripcion creado")
}
