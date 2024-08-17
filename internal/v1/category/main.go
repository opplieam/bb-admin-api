package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/store"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type Storer interface {
	GetAllCategory() ([]store.AllCategoryResult, error)
	GetUnmatchedCategory(filter utils.Filter) ([]store.UnmatchedCategoryResult, utils.MetaData, error)
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

func (h *Handler) GetUnmatchedCategory(c *gin.Context) {
	var filterParam utils.Filter
	if err := c.ShouldBindQuery(&filterParam); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	result, metaData, err := h.Store.GetUnmatchedCategory(filterParam)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "metadata": metaData})
}
