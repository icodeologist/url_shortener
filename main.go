package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/icodeologist/url_shortner/database"
	"github.com/icodeologist/url_shortner/server"
)

func main() {
	db, err := database.SetUpDb()
	if err != nil {
		log.Fatal(err)
	}
	var users []database.User

	res := db.Find(&users)

	if res.Error != nil {
		log.Fatal(res.Error)
	}

	db.AutoMigrate(&database.User{})
	for _, u := range users {
		fmt.Printf("ID %v Longurl %v ShortID %v clicks %v\n", u.ID, u.LongUrl, u.ShortID, u.Clicks)
	}

	shortID, err := database.GetShortID(db, 6)

	if err != nil {
		log.Fatal(err)
	}
	//encode the shortID
	database.UpdateShortID(db, 6)
	// get the long URL this to redirect
	longURL, err := database.GetLongURL(db, shortID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("longurl:shortID", longURL, shortID)

}
func handleServers() {
	//server logic
	fmt.Println("Server at 6868 is running")
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.ShortURLrouting)
	http.ListenAndServe(":5757", mux)
}
