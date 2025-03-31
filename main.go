package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/icodeologist/url_shortner/database"
	"github.com/icodeologist/url_shortner/server"
	"github.com/icodeologist/url_shortner/urlconverter"
	"gorm.io/gorm"
)

// repetitive actions written with simple funcitons
func printDBData(db *gorm.DB) {
	var users []database.User

	res := db.Find(&users)

	if res.Error != nil {
		log.Fatal(res.Error)
	}

	for _, u := range users {
		fmt.Printf("ID %v Longurl %v ShortID %v clicks %v\n", u.ID, u.LongUrl, u.ShortID, u.Clicks)
	}
}

func returnsSHORTID(db *gorm.DB, userID int) (string, error) {
	shortID, err := database.GetShortID(db, userID)
	if err != nil {
		return "", err
	}
	if shortID == "" {
		newShortID, err := database.UpdateShortID(db, userID)
		if err != nil {
			return "", err
		}
		return newShortID, nil
	}
	return shortID, nil
}

func main() {
	db, err := database.SetUpDb()
	if err != nil {
		log.Fatal(err)
	}
	printDBData(db)
	shortID, err := returnsSHORTID(db, 22)
	if err != nil {
		log.Fatal(err)
	}
	shortURL, fullShortURL, err := urlconverter.GenerateShortURL(shortID)
	fmt.Println(shortURL)
	fmt.Println(fullShortURL)

	http.Handle("/usly/", server.RedirectShortToLongURL(db, shortURL))
	fmt.Println("Server at 6969 is running")

	http.ListenAndServe(":6969", nil)

}

//auto create users with html forms
//save that to db
//generate shorrID then short URL
//redirect
