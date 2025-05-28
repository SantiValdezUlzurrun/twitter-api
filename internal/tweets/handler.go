package tweets

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"twitter-api/internal/tweets/dto"
)

type Handler struct {
	service Service
}

func NewHandler(tweetsService Service) *Handler {
	server := &Handler{
		service: tweetsService,
	}
	return server
}

func (h *Handler) Create(c *gin.Context) {
	var tweet dto.CreateTweet

	if err := c.ShouldBindJSON(&tweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTweet, err := h.service.Create(c.Request.Context(), &tweet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTweet)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := stouint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Tweet %d deleted", id)})
}

func stouint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}
