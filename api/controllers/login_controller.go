package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saranrajv123/levelupblog/api/auth"
	"github.com/saranrajv123/levelupblog/api/models"
	"github.com/saranrajv123/levelupblog/api/utils"
)

func (server *Server) Login(context *gin.Context) {
	// errList := map[string]string{}

	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Unable to get request",
		})
		return
	}

	user := models.User{}
	if err = json.Unmarshal(body, &user); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Unable to get request",
		})
		return
	}

	user.Prepare()
	errMessage := user.Validate("login")
	if len(errMessage) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errMessage,
		})
		return
	}

	userData, err := server.Signin(user.Email, user.Password)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  formattedError,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userData,
	})

}

func (server *Server) Signin(email, password string) (map[string]interface{}, error) {
	user := models.User{}
	userData := make(map[string]interface{})

	err := server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		fmt.Println("this is the error getting the user: ", err)
		return nil, err
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		fmt.Println("this is the error creating the token: ", err)
		return nil, err
	}

	userData["token"] = token
	userData["id"] = user.ID
	userData["email"] = user.Email
	userData["fullname"] = user.FullName

	return userData, nil

}
