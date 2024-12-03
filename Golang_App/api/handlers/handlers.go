package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"realtime_chat/api/database"
	"realtime_chat/models"
	"realtime_chat/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	DBHelperProvider database.DBHelperProvider
}

func NewHandlerProvider(dbHelper database.DBHelperProvider) *ApiHandler {

	return &ApiHandler{DBHelperProvider: dbHelper}
}

func (apihandle *ApiHandler) RenderHome(c *gin.Context) {

	response := ConstructResponse(http.StatusOK, models.APIWelcomeMessage, nil, nil)
	c.JSON(response.StatusCode, response)
}

func (apihandle *ApiHandler) UserLogin(c *gin.Context) {

	var user models.User

	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Error decoding the request body", err)
		fmt.Println("Error to decoding the request body", err)
		return
	}

	if user.UserName == "" {
		fmt.Println("Username not provided")
	}
	if user.Password == "" {
		RespondError(c, http.StatusBadRequest, "Password can't be empty", errors.New("password not provided"))
		fmt.Println("password not provided")
		return
	}

	if user.Email == "" {
		RespondError(c, http.StatusBadRequest, "Email can't be empty", errors.New("email not provided"))
		fmt.Println("email not provided")
		return
	}

	userDetails, err := apihandle.DBHelperProvider.GetUserbyEmailId(user.Email)
	if err != nil {
		RespondError(c, http.StatusNotFound, "User is not Registered", errors.New("user is not find in the database"))
		return
	}

	// Now we have to compare the password.
	isPasswordCheck := utils.ComparePasswords(user.Password, userDetails.Password)
	if isPasswordCheck != nil {
		RespondError(c, http.StatusNotFound, models.LoginPasswordIsInCorrect, isPasswordCheck)
		return
	}

	// Now change the status of the user.
	err = apihandle.DBHelperProvider.UpdateUserStatus(user.Email, "Y")
	if err != nil {
		fmt.Println("Failed to update the status of the user", err)
		return
	}

	response := ConstructResponse(http.StatusOK, models.UserLoginCompleted, userDetails, nil)
	c.JSON(response.StatusCode, response)
}

func (apihandle *ApiHandler) Registration(c *gin.Context) {

	var user models.User

	fmt.Println("c.re", c.Request)
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Error decoding the request body", err)
		fmt.Println("Error to decoding the request body", err)
		return
	}

	fmt.Println("User - ", user)

	if user.UserName == "" {
		RespondError(c, http.StatusBadRequest, "UserName can't be empty", errors.New("userName not provided"))
		fmt.Println("Username not provided")
		return
	}
	if user.Password == "" {
		RespondError(c, http.StatusBadRequest, "Password can't be empty", errors.New("password not provided"))
		fmt.Println("password not provided")
		return
	}

	if user.Email == "" {
		RespondError(c, http.StatusBadRequest, "Email can't be empty", errors.New("email not provided"))
		fmt.Println("email not provided")
		return
	}

	err = apihandle.DBHelperProvider.CreateUser(user)
	if err != nil {
		fmt.Println("Erro getting inserting the user", err)
		RespondError(c, http.StatusInternalServerError, "error creating the user", err)
		return
	}

	response := ConstructResponse(http.StatusCreated, models.UserCreated, user, nil)

	c.JSON(response.StatusCode, response)
}

func (apihandle *ApiHandler) IsUserAvailable(c *gin.Context) {

	var response *models.Response
	emailId := c.Param("username")

	if emailId == "" {
		RespondError(c, http.StatusBadRequest, "emailId can't be empty", errors.New("emailId not provided"))
		fmt.Println("emailId not provided")
		return
	}

	isAvailable, err := apihandle.DBHelperProvider.CheckUserExists(emailId)
	if err != nil {
		fmt.Println("Error checking user availability", err)
		RespondError(c, http.StatusInternalServerError, "Error checking user availability", err)
		return
	}

	if isAvailable {
		response = ConstructResponse(http.StatusOK, "User is available, Please login 🙏", isAvailable, nil)
	} else {
		response = ConstructResponse(http.StatusOK, "User is not available", isAvailable, nil)

	}
	c.JSON(response.StatusCode, response)
}

func (apihandle *ApiHandler) UserSessionCheck(c *gin.Context) {

	var status bool
	emailId := c.Param("emailId")
	if emailId == "" {
		RespondError(c, http.StatusBadRequest, "emailId can't be empty", errors.New("emailId not provided"))
		fmt.Println("emailId not provided")
		return
	}

	user, err := apihandle.DBHelperProvider.GetUserStatus(emailId)
	if err != nil {
		fmt.Println("Error getting user session status", err)
		RespondError(c, http.StatusInternalServerError, "Error getting user session status", err)
		return
	}

	if user.Online == "Y" {
		status = true
	} else {
		status = false
	}

	response := ConstructResponse(http.StatusOK, "User session status fetched successfully", status, nil)
	c.JSON(response.StatusCode, response)
}

func (apihandle *ApiHandler) GetMessagesHandler(c *gin.Context) {

	toUserId := c.Param("toUserId")
	fromUserId := c.Param("fromUserId")

	conversation := apihandle.DBHelperProvider.GetConversationBetweenTwoUsers(toUserId, fromUserId)
	response := ConstructResponse(http.StatusOK, "User is Available", conversation, nil)

	c.JSON(response.StatusCode, response)

}
