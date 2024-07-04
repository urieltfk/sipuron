package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	bookMap = make(map[uint]string)
	mu      sync.Mutex
	cond    = sync.NewCond(&mu)
)

func main() {
	r := gin.Default()

	r.GET("/book/:id", getBook)
	r.POST("/book/:id", postBook)

	r.Run(":8080")
}

func getBook(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	for bookMap[id] == "" {
		cond.Wait()
	}
	value := bookMap[id]
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"book": value})
}

func postBook(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var requestBody struct {
		Text string `json:"text"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	mu.Lock()
	bookMap[id] = requestBody.Text
	cond.Broadcast()
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func parseID(idParam string) (uint, error) {
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

/*
func TestGet() {
	time.Sleep(3 * time.Second)

	resp, err := http.Get("http://localhost:8080/book/1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response Body:", string(body))

	time.Sleep(3 * time.Second)
}

func TestPost() {
	time.Sleep(3 * time.Second)
	data := []byte(`{"Text":"This is our new story, hello from the new story"}`)

	resp, err := http.Post("http://localhost:8080/book/1", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	time.Sleep(3 * time.Second)
}
*/
