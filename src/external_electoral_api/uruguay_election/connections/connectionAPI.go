package connections

import (
	"bufio"
	"external_electoral_api/uruguay_election/connections/configs"
	"external_electoral_api/uruguay_election/connections/routes"
	"external_electoral_api/uruguay_election/controllers"
	"external_electoral_api/uruguay_election/repositories"
	"fmt"
	"os"

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

	fmt.Print("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
