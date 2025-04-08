package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/icodeologist/url_shortner/database"
	"github.com/icodeologist/url_shortner/urlconverter"
	"gorm.io/gorm"
)

var mu sync.Mutex

func getDB() (*gorm.DB, error) {
	db, err := database.ConnectToPSQL()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func HandleCreate(c *gin.Context) {
	var user database.URL
	//bind the json file for now to user
	longUrl := c.PostForm("longurl")
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}
	//update the long url
	db.Model(&user).Update("LongUrl", longUrl)
	res := db.Create(&user)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	c.HTML(http.StatusOK, "res.html", gin.H{
		"longurl": user.LongUrl,
		"ID":      user.ID,
	})
}

func HandleGetEverything(c *gin.Context) {
	var users []database.User
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}
	res := db.Find(&users)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func atoiForID(id string) uint {
	uId, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return uint(uId)

}

func HandleGetUrl(c *gin.Context) {
	// get the user Id
	id := c.Param("id")
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}
	uId := atoiForID(id)

	user, err := getUserbyId(db, uint(uId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"err":     err,
			"message": "Given ID was not found in the Database.",
		})
	}
	c.IndentedJSON(http.StatusOK, user)
}

func getUserbyId(db *gorm.DB, id uint) (*database.URL, error) {
	var user database.URL
	res := db.First(&user, id)
	if res.Error != nil {
		return nil, res.Error

	}
	return &user, nil
}

// update the short Id
func UpdateShortID(c *gin.Context) {
	id := c.Param("id")
	uId := atoiForID(id)
	db, err := getDB()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error code": http.StatusInternalServerError,
			"message":    "Could not connect to Database.",
		})
	}
	user, err := getUserbyId(db, uId)
	// manually update users short ID
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   http.StatusNotFound,
			"message": "Given ID was not found in the Database.",
		})
	}
	shortId := urlconverter.Base62Encoding(uId)
	user.ShortID = shortId
	db.Save(user)
	c.IndentedJSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	uId := atoiForID(id)
	db, err := getDB()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error code": http.StatusInternalServerError,
			"message":    "Could not connect to Database.",
		})
	}
	user, err := getUserbyId(db, uId)
	// manually update users short ID
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   http.StatusNotFound,
			"message": "Given ID was not found in the Database.",
		})
		return
	}
	db.Delete(&user)
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%v got deleted.", user.ID),
	})
}

func displayDB(c *gin.Context) {
	var users []database.User
	db, err := getDB()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error code": http.StatusInternalServerError,
			"message":    "Could not connect to Database.",
		})
	}
	db.Find(&users)
	fmt.Println(users)
	c.HTML(http.StatusOK, "home.html", gin.H{
		"users": users,
	})
}

func main() {
	router := gin.Default()
	//POST
	router.LoadHTMLFiles("templates/form.html", "templates/res.html", "templates/home.html")
	//form in home
	router.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form.html", nil)
	})

	router.POST("/add", HandleCreate)
	//GET
	router.GET("/urls", HandleGetEverything)
	router.GET("/", displayDB)
	router.GET("urls/:id", HandleGetUrl)
	//UPDATE
	router.PUT("/urls/:id", UpdateShortID)
	//DELETE
	router.DELETE("/urls/:id", DeleteUser)
	router.Run("localhost:2000")
}
