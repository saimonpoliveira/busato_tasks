package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/config"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/controllers"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/repositories"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/routes"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/services"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	gin.SetMode(cfg.GinMode)

	db, err := repositories.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpirationHours)

	userRepo := repositories.NewUserRepository(db.DB)
	projectRepo := repositories.NewProjectRepository(db.DB)
	ticketRepo := repositories.NewTicketRepository(db.DB)
	taskRepo := repositories.NewTaskRepository(db.DB)
	commentRepo := repositories.NewCommentRepository(db.DB)
	attachmentRepo := repositories.NewAttachmentRepository(db.DB)

	authService := services.NewAuthService(userRepo, jwtManager)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)
	ticketService := services.NewTicketService(ticketRepo, projectRepo)
	taskService := services.NewTaskService(taskRepo, ticketRepo)
	commentService := services.NewCommentService(commentRepo, ticketRepo, taskRepo)
	attachmentService := services.NewAttachmentService(attachmentRepo, ticketRepo, taskRepo, cfg.UploadDir, cfg.MaxUploadSizeMB)

	handlers := &routes.Handlers{
		Auth:       controllers.NewAuthController(authService),
		User:       controllers.NewUserController(userService),
		Project:    controllers.NewProjectController(projectService),
		Ticket:     controllers.NewTicketController(ticketService),
		Task:       controllers.NewTaskController(taskService),
		Comment:    controllers.NewCommentController(commentService),
		Attachment: controllers.NewAttachmentController(attachmentService),
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	routes.Setup(router, handlers, jwtManager, cfg.CORSOrigins)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
