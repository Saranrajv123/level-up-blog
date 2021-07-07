package controllers

import "github.com/saranrajv123/levelupblog/api/middleware"

func (server *Server) InitializeRouter() {
	v1 := server.Router.Group("/users")
	{
		v1.POST("/login", server.Login)
		v1.POST("/create-user", server.CreateUser)
		v1.GET("/all-users", server.GetUsers)
		v1.GET("/user/:id", server.GetUserById)
		v1.PUT("/update-user/:id", middleware.TokenAuthMiddleware(), server.Updateuser)
		v1.DELETE("/delete-user/:id", middleware.TokenAuthMiddleware(), server.DeleteUser)
	}
}
