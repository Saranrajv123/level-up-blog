package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saranrajv123/levelupblog/api/models"
	"github.com/saranrajv123/levelupblog/api/utils"
)

func (server *Server) CreateUser(context *gin.Context) {
	errList := map[string]string{}
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		errList["Invalid body"] = "Invalid Body"
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user := models.User{}

	if err = json.Unmarshal(body, &user); err != nil {
		errList["UnMarshal Error"] = "Cannot Un Marshal body"
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user.Prepare()
	//errorMessages := user.Validate("")
	//if len(errorMessages) > 0 {
	//	errList = errorMessages
	//	context.JSON(http.StatusUnprocessableEntity, gin.H{
	//		"status": http.StatusUnprocessableEntity,
	//		"error": errList,
	//	})
	//	return
	//
	//}

	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		//errList = formattedError /* todo  change to errList format*/
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  formattedError,
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"response": userCreated,
		"status":   http.StatusCreated,
	})

}

func (server *Server) GetUsers(context *gin.Context) {
	errList := map[string]string{}

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		errList["No user available"] = "No User available"
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": users,
	})

}
