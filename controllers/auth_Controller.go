package controllers

import (
	"fmt"
	"golang/dao"
	"golang/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthController struct {
	userDao dao.UserDao
}

func NewAuthController(dao dao.UserDao) *AuthController {
	return &AuthController{userDao: dao}
}

func (ac *AuthController) Login(c *gin.Context) {
	var requestBody models.LoginRequest

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user *models.User
	var err error
	// initializers.DB.First(&user, "email = ?", requestBody.Email)
	user, err = ac.userDao.FindByEmail(requestBody.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
		"role": models.Role.String(user.Role),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	log.Println("token string : ", tokenString)
	fmt.Println(tokenString, err)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error generating token",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
