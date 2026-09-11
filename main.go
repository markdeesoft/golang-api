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
	"github.com/markdeesoft/golang-api/handler"
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

	db := database.InitDB()
	defer database.DB.Close() //สำหรับเทสรูปแบบ sql ปกติ

	//ประกาศตัวแปรเพื่อประกอบร่าง (Dependency Injection)
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)
	productCategoryRepo := repository.NewProductCategoryRepository(db)
	productCategoryHandler := handler.NewProductCategoryHandler(productCategoryRepo)
	productHashtagRepo := repository.NewProductHashtagRepository(db)
	productHashtagHandler := handler.NewProductHashtagHandler(productHashtagRepo)
	productRepo := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)

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
	app.Get("/user", userHandler.View)
	app.Put("/user/:id", userHandler.Update)
	app.Post("/user/uploadphoto", handler.UploadPhotoUser)
	app.Delete("/user/:id", userHandler.Delete)
	// productcategory handler
	app.Get("/productcategory", productCategoryHandler.List)
	app.Get("/productcategory/:id", productCategoryHandler.View)
	app.Post("/productcategory", productCategoryHandler.Store)
	app.Delete("/productcategory/:id", productCategoryHandler.Delete)
	// product handler
	app.Get("/producthash", productHashtagHandler.List)
	app.Post("/producthash", productHashtagHandler.Store)
	app.Put("/producthash/:id", productHashtagHandler.Update)
	// product handler
	app.Get("/product", productHandler.List)
	app.Get("/product/:id", productHandler.View)
	app.Post("/product", productHandler.Store)
	app.Get("/product/:id/edit", productHandler.View)
	app.Put("/product/:id", productHandler.Update)
	app.Delete("/product/:id", productHandler.Delete)

	// Group routes under /admin
	adminGroup := app.Group("/admin")

	// Apply the isAdmin middleware only to the /admin routes
	adminGroup.Use(isAdmin)

	// admin/user handler
	app.Get("/admin/user", userHandler.List)
	app.Get("/admin/user/:id", userHandler.View)
	app.Post("/admin/user", userHandler.Store)
	app.Put("/admin/user/:id", userHandler.Update)
	app.Delete("/admin/user/:id", userHandler.Delete)
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
	app.Get("/admin/product", productHandler.List)
	app.Get("/admin/product/:id", productHandler.View)
	app.Post("/admin/product", productHandler.Store)
	app.Get("/admin/product/:id/edit", productHandler.View)
	app.Put("/admin/product/:id", productHandler.Update)
	app.Delete("/admin/product/:id", productHandler.Delete)

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
