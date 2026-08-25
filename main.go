package main

import (
	"fmt"
	"log"
	"os"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/markdeesoft/golang-api/database"
	handler "github.com/markdeesoft/golang-api/handlers"
	"github.com/markdeesoft/golang-api/repository"
)

var ()

type Config struct {
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {

	app := fiber.New()

	// Apply CORS middleware
	allow_origin := os.Getenv("ALLOW_ORIGIN")
	if allow_origin == "" {
		allow_origin = "*" // Default port if not specified
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{allow_origin}, // Adjust this to be more restrictive if needed
		// AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	database.InitDB()
	defer database.DB.Close()

	//ประกาศตัวแปรเพื่อประกอบร่าง (Dependency Injection)
	userRepo := repository.NewUserRepository()
	userHandler := handler.NewUserHandler(userRepo)
	productCategoryRepo := repository.NewProductCategoryRepository()
	productCategoryHandler := handler.NewProductCategoryHandler(productCategoryRepo)
	productHashtagRepo := repository.NewProductHashtagRepository()
	productHashtagHandler := handler.NewProductHashtagHandler(productHashtagRepo)

	// Login route
	app.Post("/auth/login", userHandler.Login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		Extractor:  extractors.FromAuthHeader("Bearer"),
	}))

	// Restricted Routes
	app.Get("/restricted", restricted)

	// user handler
	app.Get("/user", userHandler.GetUser)
	app.Put("/user/:id", userHandler.UpdateUser)
	app.Post("/user/uploadphoto", handler.UploadPhotoUser)
	// productcategory handler
	app.Get("/productcategory", productCategoryHandler.List)
	app.Get("/productcategory/:id", productCategoryHandler.View)
	// product handler
	app.Get("/producthash", productHashtagHandler.List)
	// product handler
	app.Get("/product", productCategoryHandler.List)
	app.Get("/product/:id", productCategoryHandler.View)
	app.Post("/product", productCategoryHandler.Store)
	app.Get("/product/:id/edit", productCategoryHandler.View)
	app.Put("/product/:id", productCategoryHandler.Update)
	app.Delete("/product/:id", productCategoryHandler.Delete)

	// Group routes under /admin
	adminGroup := app.Group("/admin")

	// Apply the isAdmin middleware only to the /admin routes
	adminGroup.Use(isAdmin)

	// admin/user handler
	app.Get("/admin/user", userHandler.GetUsers)
	app.Get("/admin/user/:id", userHandler.GetUser)
	app.Post("/admin/user", userHandler.StoreUser)
	app.Put("/admin/user/:id", userHandler.UpdateUser)
	app.Delete("/admin/user/:id", userHandler.DeleteUser)
	app.Post("/admin/user/uploadphoto", handler.UploadPhotoUser)
	app.Patch("/admin/user/resetpass/:id", userHandler.ResetPasswordUser)
	// admin/productcategory handler
	app.Get("/admin/productcategory", productCategoryHandler.List)
	app.Get("/admin/productcategory/:id", productCategoryHandler.View)
	app.Post("/admin/productcategory", productCategoryHandler.Store)
	app.Put("/admin/productcategory/:id", productCategoryHandler.Update)
	app.Delete("/admin/productcategory/:id", productCategoryHandler.Delete)
	// admin/product handler
	app.Get("/admin/producthash", productHashtagHandler.List)
	app.Get("/admin/producthash/:id", productHashtagHandler.View)
	app.Post("/admin/producthash", productHashtagHandler.Store)
	app.Put("/admin/producthash/:id", productHashtagHandler.Update)
	// admin/product handler
	app.Get("/admin/product", productCategoryHandler.List)
	app.Get("/admin/product/:id", productCategoryHandler.View)
	app.Post("/admin/product", productCategoryHandler.Store)
	app.Get("/admin/product/:id/edit", productCategoryHandler.View)
	app.Put("/admin/product/:id", productCategoryHandler.Update)
	app.Delete("/admin/product/:id", productCategoryHandler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Fatal(app.Listen(":" + port))
}

func restricted(c fiber.Ctx) error {
	user := jwtware.FromContext(c)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func isAdmin(c fiber.Ctx) error {

	user := jwtware.FromContext(c)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println("claims", claims)
	if claims["role"].(string) != "admin" {
		return fiber.NewError(fiber.StatusUnauthorized, "ไม่มีสิทธิ์เข้าถึง")
	}

	return c.Next()
}
