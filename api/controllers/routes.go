package controllers

func (server *Server) InitializeRouter() {
	v1 := server.Router.Group("/users")
	{
		v1.POST("/create-user", server.CreateUser)
		v1.GET("/all-users", server.GetUsers)
	}
}
