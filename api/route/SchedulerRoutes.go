package route

import (
	"Echo/api/controller"
	"Echo/mongo/models"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SchedulerRoutes(router *gin.Engine, repo *controller.SchedulerRepository, userRepo *controller.UserRepository) {
	schedulerRoute := router.Group("/scheduler")
	schedulerRoute.POST("/", func(ctx *gin.Context) {
		var task models.Task
		if err := ctx.ShouldBindJSON(&task); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		repo.AddTask(context.Background(), &task, userRepo)
	})
}
