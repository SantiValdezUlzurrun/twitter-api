package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"twitter-api/internal/users/dto"
)

type Handler struct {
	service Service
}

func NewHandler(userService Service) *Handler {
	server := &Handler{
		service: userService,
	}
	return server
}

func (h *Handler) Create(c *gin.Context) {
	var user dto.CreateUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.service.Create(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) Follow(c *gin.Context) {
	userID, err := stouint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	followerID, err := stouint(c.Param("followerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Follow(c.Request.Context(), userID, followerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %d follows %d", userID, followerID)})
}

func (h *Handler) Unfollow(c *gin.Context) {
	userID, err := stouint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	followerID, err := stouint(c.Param("followerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Unfollow(c.Request.Context(), userID, followerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %d no longer follows %d", userID, followerID)})
}

func stouint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}
