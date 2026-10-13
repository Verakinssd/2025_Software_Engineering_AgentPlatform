package router

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()

	loginRequired := []string{"user", "admin"}
	adminRequired := []string{"admin"}

	api := Router.Group("api")
	{
		api.POST("/register", controller.UserRegister)
		api.POST("/login", controller.UserLogin)
		api.GET("/image/:filename", controller.GetImage)
		api.POST("/image/upload", middleware.Auth(loginRequired), controller.UploadImage)

		user := api.Group("user", middleware.Auth(loginRequired))
		{
			user.GET("/profile", controller.UserProfile)
			user.POST("/logout", controller.UserLogout)
			user.POST("/update_profile", controller.UserUpdateProfile)
			user.POST("/upload_avatar", controller.UserUpdateAvatar)
		}

		admin := api.Group("admin", middleware.Auth(adminRequired))
		{
			admin.GET("/penalty/:user_id", controller.GetUserPenaltyFor)
			admin.GET("/all_user", controller.GetAllUser)
			admin.POST("/penalty", controller.HandleUserPenalty)
			admin.POST("/update_penalty", controller.UpdateUserPenalty)
		}

		notification := api.Group("notification", middleware.Auth(loginRequired))
		{
			notification.GET("/list", controller.GetNotifications)
			notification.POST("/send", controller.SendNotification)
			notification.POST("/broadcast", middleware.Auth(adminRequired), controller.BroadcastNotification)
			notification.POST("/read/:notification_id", controller.ReadNotification)
			notification.POST("/read_all", controller.ReadAllNotifications)
		}

		workflowPublic := api.Group("workflow")
		{
			workflowPublic.GET("/public", controller.GetAllPublicWorkflows)
			workflowPublic.GET("/popular", controller.GetPopularWorkflows)

			// comments related
			workflowPublic.GET("/:workflow_id/comment", controller.GetAllComments)
		}

		workflow := api.Group("workflow", middleware.Auth(loginRequired))
		{
			workflow.POST("/create", middleware.OperationAuth(middleware.PublishWorkflow), controller.CreateWorkflow)
			workflow.GET("/list", controller.GetWorkflowList)
			workflow.GET("/:workflow_id", controller.GetWorkflowById)
			workflow.GET("/:workflow_id/io", controller.GetWorkflowIOStructure)
			workflow.POST("/:workflow_id/update", middleware.OperationAuth(middleware.PublishWorkflow), controller.UpdateWorkflow)
			workflow.POST("/:workflow_id/run", controller.RunWorkflow)
			workflow.POST("/:workflow_id/run_with_inputs", controller.RunWorkflowWithInputs)
			workflow.POST("/:workflow_id/content", controller.UpdateWorkflowContent)
			workflow.POST("/:workflow_id/fork", controller.ForkWorkflow)
			workflow.DELETE("/:workflow_id", controller.DeleteWorkflow)
			workflow.GET("/available", controller.GetAvailableWorkflows)
			workflow.GET("/public2", controller.GetPublicWorkflows)
			workflow.GET("/personal", controller.GetPersonalWorkflows)

			// comments related
			workflow.POST("/comment/:comment_id/like", controller.LikeComment)
			workflow.POST("/comment/:comment_id/delete", controller.DeleteComment)
			workflow.POST("/comment", middleware.OperationAuth(middleware.Comment), controller.CreateComment)

			// favorites related
			workflow.GET("/favorites", controller.GetFavoriteWorkflowList)
			workflow.POST("/:workflow_id/favorite", controller.FavoriteWorkflow)

			// like workflow
			workflow.POST("/:workflow_id/like", controller.LikeWorkflow)
		}
		
		// 知识库相关路由
		knowledge := api.Group("knowledge", middleware.Auth(loginRequired))
		{
			knowledge.POST("/upload", controller.UploadKnowledgeFile)
			knowledge.GET("/files", controller.GetUserKnowledgeFiles)
			knowledge.GET("/file/:file_id", controller.GetKnowledgeFile)
			knowledge.GET("/download/:file_id", controller.DownloadKnowledgeFile)
			knowledge.DELETE("/file/:file_id", controller.DeleteKnowledgeFile)
		}

		//node := api.Group("node", middleware.Auth(loginRequired))
		//{
		//	//node.POST("/:node_id/run", controller.RunNode)
		//}
	}
}
