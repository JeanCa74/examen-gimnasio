// ARCHIVO BLOQUEADO — NO MODIFICAR
//
// Carpeta acceptance/: estos tests SON la rúbrica del examen.
// Ejecútelos cuantas veces quiera con:
//
//	go test ./acceptance/... -v
//
// Para ver solo un checkpoint:
//
//	go test ./acceptance/... -v -run TestCP1
//	go test ./acceptance/... -v -run TestCP2
//	go test ./acceptance/... -v -run TestCP3
//
// NOTA: hasta que usted complete los campos de los modelos (CP1), este
// paquete no compila. Los errores de compilación le indican exactamente
// qué campos faltan.
package acceptance

import (
	"net/http"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joancema/examen-gimnasio/internal/handlers"
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// nuevaDB abre una base SQLite en memoria y migra los tres modelos.
func nuevaDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("no se pudo abrir SQLite en memoria: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Clase{},
		&models.Cliente{},
		&models.Inscripcion{},
	); err != nil {
		t.Fatalf("la migración de los modelos falló: %v", err)
	}
	return db
}

// nuevoRouterCompleto arma el cableado completo de la API sobre la base dada:
// repositories GORM -> services -> handlers -> router.
func nuevoRouterCompleto(t *testing.T, db *gorm.DB) http.Handler {
	t.Helper()
	claseRepo := storage.NuevaClaseGORM(db)
	clienteRepo := storage.NuevoClienteGORM(db)
	inscripcionRepo := storage.NuevaInscripcionGORM(db)
	return handlers.NuevoRouter(
		handlers.NuevoClaseHandler(services.NuevaClaseService(claseRepo)),
		handlers.NuevoClienteHandler(services.NuevoClienteService(clienteRepo)),
		handlers.NuevaInscripcionHandler(services.NuevaInscripcionService(inscripcionRepo, claseRepo, clienteRepo)),
	)
}
