package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/icodeologist/url_shortner/urlconverter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint `gorm:"unique;primarykey;autoincrement"`
	LongUrl  string
	ShortID  string `gorm:"column:shortid;type:varchar(255)"`
	CreateAt time.Time
	Clicks   uint `gorm:"column:shortid;default:0"`
}

func ConnectToPSQL() (*gorm.DB, error) {
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("NAME")

	dsn := fmt.Sprintf("user=%v password=%v dbname=%v host=localhost port=5432 sslmode=disable", user, password, dbname)
	//connect to postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{})

	return db, nil
}

// NOTE newuser is a pointer
func CreateNewUser(db *gorm.DB, newuser User) error {
	result := db.Create(&newuser)
	db.Save(&newuser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByID(db *gorm.DB, userID int) (*User, error) {
	var url User
	// user is empty User object
	result := db.First(&url, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &url, nil
}

func SetUpDb() (*gorm.DB, error) {
	db, err := ConnectToPSQL()
	if err != nil {
		return nil, err
	} else {

		er := db.AutoMigrate(&User{})
		if er != nil {
			log.Fatal(er)
		} else {
			fmt.Println("Migrated successfully")
		}
	}
	return db, nil
}

func MakeShortUrl(userID int) string {
	getUniqueID := urlconverter.Base62Encoding(userID)
	return getUniqueID
}

func UpdateShortID(db *gorm.DB, userID int) error {
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println(MakeShortUrl(userID))
	newShortID := string(MakeShortUrl(userID))
	fmt.Println("New short ID", newShortID)
	res := db.Model(&user).Update("ShortID", newShortID)
	if res.Error != nil {
		return res.Error
	}
	db.Save(&user)
	return nil
}

func DeleteEntitybyID(db *gorm.DB, userID int) {
	var user User
	db.Delete(&user, userID)

}

func GetShortID(db *gorm.DB, userID int) (string, error) {
	usr, err := GetUserByID(db, userID)
	if err != nil {
		return "", err
	}
	return usr.ShortID, nil
}

// this gets the longurl of given shortid
func GetLongURL(db *gorm.DB, shortID string) (string, error) {
	var user User
	res := db.Where("shortID = ?", shortID).First(&user)
	if res.Error != nil {
		return "", res.Error
	}
	return user.LongUrl, nil
}
