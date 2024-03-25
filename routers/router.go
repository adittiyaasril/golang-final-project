package routers

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	users := r.Group("/users")
	{
		users.POST("/register", controllers.RegisterUser)
		users.POST("/login", controllers.Login)
		users.PUT("/:userId", middlewares.Authenticate(), controllers.UpdateUser)
		users.DELETE("/", middlewares.Authenticate(), controllers.DeleteUser)

	}

	photos := r.Group("/photos")
	{
		photos.POST("/", middlewares.Authenticate(), controllers.CreatePhoto)
		photos.GET("/", middlewares.Authenticate(), controllers.GetPhotos)
		photos.GET("/:photoId", middlewares.Authenticate(), controllers.GetPhotoByID)
		photos.PUT("/:photoId", middlewares.Authenticate(), middlewares.AuthorizePhoto(), controllers.UpdatePhoto)
		photos.DELETE("/:photoId", middlewares.Authenticate(), middlewares.AuthorizePhoto(), controllers.DeletePhoto)
	}

	comments := r.Group("/comments")
	{
		comments.POST("/", middlewares.Authenticate(), controllers.CreateComment)
		comments.GET("/", middlewares.Authenticate(), controllers.GetComments)
		comments.GET("/:commentId", middlewares.Authenticate(), controllers.GetCommentByID)
		comments.PUT("/:commentId", middlewares.Authenticate(), middlewares.AuthorizeComment(), controllers.UpdateComment)
		comments.DELETE("/:commentId", middlewares.Authenticate(), middlewares.AuthorizeComment(), controllers.DeleteComment)
	}

	socialMedia := r.Group("/socialmedias")
	{
		socialMedia.POST("/", middlewares.Authenticate(), controllers.CreateSocialMedia)
		socialMedia.GET("/", middlewares.Authenticate(), controllers.GetSocialMedias)
		socialMedia.GET("/:socialMediaId", middlewares.Authenticate(), controllers.GetSocialMediaByID)
		socialMedia.PUT("/:socialMediaId", middlewares.Authenticate(), middlewares.AuthorizeSocialMedia(), controllers.UpdateSocialMedia)
		socialMedia.DELETE("/:socialMediaId", middlewares.Authenticate(), middlewares.AuthorizeSocialMedia(), controllers.DeleteSocialMedia)
	}

	return r
}
