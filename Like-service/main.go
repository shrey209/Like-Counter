package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shrey209/Like-Service/controller"
	"github.com/shrey209/Like-Service/repository"
	"github.com/shrey209/Like-Service/service"
)

func main() {
	session := NewSession("127.0.0.1", "like_service")
	repo := repository.NewLikeRepository(session)
	svc := service.NewLikeService(repo)
	ctrl := controller.NewLikeController(svc)

	router := gin.Default()
	ctrl.RegisterRoutes(router)

	router.Run(":8080")
}
