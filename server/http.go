package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/icodeologist/url_shortner/database"
	"gorm.io/gorm"
)

// when user enteres long url save in db
func ShortURLrouting(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "This should render actual long url")
}

// redirection url logic
func RedirectShortToLongURL(db *gorm.DB, shortURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "GET" {
			http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
			return
		}
		//get the ID only from shortURL
		shortID := strings.Trim(shortURL, "usly/")

		fmt.Println("ShortID ", shortID)
		if shortID == "" {
			fmt.Fprintf(w, "Welcome to URL shortener")
			return
		}

		// fetch long url
		var user database.User
		result := db.Where("shortID = ?", shortID).First(&user)
		if result.Error != nil {
			http.NotFound(w, req)
			return
		}
		http.Redirect(w, req, user.LongUrl, http.StatusMovedPermanently)

	}
}
