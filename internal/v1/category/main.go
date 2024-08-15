package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/store"
)

type Storer interface {
	GetAllCategory() ([]store.AllCategoryResult, error)
}

type Handler struct {
	Store Storer
}

func NewHandler(store Storer) *Handler {
	return &Handler{
		Store: store,
	}
}

func (h *Handler) GetAllCategory(c *gin.Context) {
	result, err := h.Store.GetAllCategory()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
