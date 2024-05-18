package todo

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// create struct
type Todo struct {
	Title string `json:"text"`
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

	// Logging
	if todo.Title == "sleep" {
		// show header for mark problem possition by transaction in backend
		transactionID := c.Request.Header.Get("TransactionID")
		aud, _ := c.Get("aud")
		log.Println(transactionID, aud, "not allowed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not allowed",
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

func (t *TodoHandler) List(c *gin.Context) {
	var todos []Todo
	r := t.db.Find(&todos)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) Remove(c *gin.Context){
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r := t.db.Delete(&Todo{}, id)
	if err = r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}