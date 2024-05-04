package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// create struct
type Todo struct {
	Title string "json: title"
	gorm.Model
}

// create gorm type for use db
type TodoHandler struct {
	db *gorm.DB
}

// create func handler for frontend send
func NewHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{
		db: db,
	}
}

// create func newTask for get data form path
func (t *TodoHandler) NewTask(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// database
	r := t.db.Create(&todo)
	// return
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ID": todo.Model.ID,
	})
}
