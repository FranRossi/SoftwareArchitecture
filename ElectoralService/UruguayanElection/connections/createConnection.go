package connections

import (
	"ElectoralService/UruguayanElection/connections/configs"
	"ElectoralService/UruguayanElection/connections/routes"
	"ElectoralService/UruguayanElection/controllers"
	"ElectoralService/UruguayanElection/repositories"
	"github.com/gofiber/fiber/v2"
)

func Connection() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	//// Middlewares.
	//// Aqui defino middlewares que quiera para mi app, solo uso un logger por ahora
	//middlewares.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Repositories
	// Defino mi instancia del repositorio para inyectarla a los controladores a utilizar
	// Tengo que sacar la direccion de memoria a mano para pasarlo a las funciones
	// que utilizan mi instancia de repo, al ser metodo y no funcion se pierde el syntax sugar
	repo := &repositories.ElectionRepo{}

	// Creo una instancia de mis controladores con mi instancia de repo
	controller := controllers.NewElectionController(repo)

	// Routes.
	// Aqui defino cuales van a ser las rutas accesibles
	routes.PublicRoutes(app, controller) // Register a public routes for app.

	// Aqui inicializamos el servidor en el puerto 8080
	app.Listen(":8080")
}
