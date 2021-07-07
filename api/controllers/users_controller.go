package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saranrajv123/levelupblog/api/auth"
	"github.com/saranrajv123/levelupblog/api/models"
	"github.com/saranrajv123/levelupblog/api/security"
	"github.com/saranrajv123/levelupblog/api/utils"
	"golang.org/x/crypto/bcrypt"
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

func (server *Server) GetUserById(context *gin.Context) {
	errList := map[string]string{}

	userId := context.Param("id")
	uId, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		errList["Invalid Request"] = "Invalid request"
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user := models.User{}
	foundUser, err := user.FindUserById(server.DB, uint32(uId))
	if err != nil {
		errList["No_user"] = "No User Found"
		context.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": foundUser,
	})

}

func (server *Server) GetUser(context *gin.Context) {
	errList := map[string]string{}

	param := context.Param("id")

	userId, err := strconv.ParseUint(param, 10, 20)
	if err != nil {
		errList["Invalid Request"] = "Invalid Request"
		context.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	user := models.User{}
	userGotten, err := user.FindUserById(server.DB, uint32(userId))
	if err != nil {
		errList["no user found"] = "No user found"
		context.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userGotten,
	})
}

func (server *Server) Updateuser(context *gin.Context) {
	errList := map[string]string{}
	userId := context.Param("id")
	fmt.Println("user id", userId)

	uid, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		errList["Invalid Request"] = "Invalid Request"
		context.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	tokenId, err := auth.ExtractTokenID(context.Request)
	fmt.Println("toeknId, err", tokenId, err)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		context.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	fmt.Println("token id and user id ", tokenId, uid)

	if tokenId != 0 && tokenId != uint32(uid) {
		errList["Unauthorized"] = "Unauthorized"
		context.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	previousUserModel := models.User{}
	newUser := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&previousUserModel).Error
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	fmt.Println("current password", requestBody["current_password"])

	if requestBody["current_password"] == "" && requestBody["new_password"] != "" ||
		requestBody["current_password"] != "" && requestBody["new_password"] == "" {
		errList["Empty_new"] = "Please Provide new password"
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return

	}

	if requestBody["current_password"] != "" && requestBody["new_password"] != "" {
		if len(requestBody["new_password"]) < 6 {
			errList["Invalid_password"] = "Password should be atleast 6 characters"
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}

		err = security.VerifyPassword(previousUserModel.Password, requestBody["new_password"])
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			errList["Password_mismatch"] = "The password not correct"
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}

		newUser.FullName = previousUserModel.FullName
		newUser.Password = requestBody["new_password"]
		newUser.Email = requestBody["email"]
	}

	newUser.FullName = previousUserModel.FullName
	newUser.Email = requestBody["email"]

	newUser.Prepare()
	errorMessages := newUser.Validate("update")

	if len(errorMessages) > 0 {
		errList = errorMessages
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	updatedUser, err := newUser.UpdateUser(server.DB, uint32(uid))
	if err != nil {
		errList := utils.FormatError(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedUser,
	})

}

func (server *Server) DeleteUser(context *gin.Context) {
	errList := map[string]string{}

	userId := context.Param("id")
	uid, err := strconv.ParseUint(userId, 10, 32)

	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		context.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	tokenId, err := auth.ExtractTokenID(context.Request)
	fmt.Println("tokenID, err", tokenId, err)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		context.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	if tokenId != 0 && tokenId != uint32(uid) {
		errList["Unauthorized"] = "Unauthorized"
		context.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	user := models.User{}
	_, err = user.DeleteUser(server.DB, uint32(uid))

	if err != nil {
		errList["Other_error"] = "Please try again later"
		context.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "User Deleted Successfully",
	})

}
