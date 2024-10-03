package controllers

import (
	"golang/models"
	"golang/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (uc *UserController) GetUserById(c *gin.Context) {
	id, e := c.Params.Get("id")
	if e != true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := uc.userService.GetByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// func (uc *UserController) GetAllUsers(c *gin.Context) {
// 	users, err := uc.userService.GetAll()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, users)
// }

// func (uc *UserController) GetAllUsers(c *gin.Context) {
// 	page, _ := strconv.Atoi(c.Query("page"))
// 	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

// 	users, err := uc.userService.GetAll(page, pageSize)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, users)
// }

func (uc *UserController) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	searchQuery := c.Query("search")

	users, err := uc.userService.GetAll(page, pageSize, searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := c.Params.Get("id")
	if err != true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID" + id})
		return
	}

	// userId, er := strconv.ParseUint(id, 10, 64)
	// if er != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID "})
	// 	return
	// }

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.Update(&updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := c.Params.Get("id")
	if err != true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userId, errorr := strconv.ParseUint(id, 10, 64)
	if errorr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := uc.userService.Delete(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
