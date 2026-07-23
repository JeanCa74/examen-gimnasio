// ARCHIVO BLOQUEADO — NO MODIFICAR
package main

import (
	"log"
	"net/http"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/joancema/examen-gimnasio/internal/handlers"
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

func main() {
	db, err := gorm.Open(sqlite.Open("gimnasio.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Clase{},
		&models.Cliente{},
		&models.Inscripcion{},
	); err != nil {
		log.Fatalf("error en la migración: %v", err)
	}

	sembrarClases(db)

	// Repositories (GORM)
	claseRepo := storage.NuevaClaseGORM(db)
	clienteRepo := storage.NuevoClienteGORM(db)
	inscripcionRepo := storage.NuevaInscripcionGORM(db)

	// Services
	claseSvc := services.NuevaClaseService(claseRepo)
	clienteSvc := services.NuevoClienteService(clienteRepo)
	inscripcionSvc := services.NuevaInscripcionService(inscripcionRepo, claseRepo, clienteRepo)

	// Handlers + Router
	router := handlers.NuevoRouter(
		handlers.NuevoClaseHandler(claseSvc),
		handlers.NuevoClienteHandler(clienteSvc),
		handlers.NuevaInscripcionHandler(inscripcionSvc),
	)

	log.Println("API del gimnasio escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

// sembrarClases carga el catálogo inicial solo si la tabla está vacía.
// Los clientes y inscripciones se crean vía API.
func sembrarClases(db *gorm.DB) {
	var total int64
	db.Model(&models.Clase{}).Count(&total)
	if total > 0 {
		return
	}
	iniciales := []models.Clase{
		{Nombre: "Clase de spinning", PrecioUnitario: 8.50, Stock: 10, Activo: true},
		{Nombre: "Clase de yoga", PrecioUnitario: 6.00, Stock: 4, Activo: true},
		{Nombre: "Clase de crossfit", PrecioUnitario: 5.00, Stock: 2, Activo: true},
		{Nombre: "Clase de boxeo", PrecioUnitario: 15.00, Stock: 3, Activo: false},
	}
	for i := range iniciales {
		db.Create(&iniciales[i])
	}
}
