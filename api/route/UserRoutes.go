package route

import (
	"Echo/api/controller"
	"Echo/mongo/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func UserRoutes(router *gin.Engine, repo *controller.UserRepository) {
	userGroup := router.Group("/user")

	userGroup.GET("/all", func(context *gin.Context) {
		users, err := repo.GetAllUsers(context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, users)

	})

	userGroup.POST("/", func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		if err := user.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.Timestamp = time.Now()

		if err := repo.CreateUser(context.Background(), user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"user": user})

	})

	userGroup.GET("/:id", func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		user, err := repo.FindUserByID(context.Background(), id)
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		ctx.JSON(http.StatusOK, user)
	})

	userGroup.PUT("/:id", func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// Set the existing ID and the Timestamp (optional)
		user.ID = id

		if err := repo.UpdateUser(context.Background(), user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"user": user})
	})

	userGroup.DELETE("/:id", func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		if err := repo.DeleteUser(context.Background(), id); err != nil {
			if err.Error() == "mongo: no documents in result" {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

	})

}
