package urlconverter

import (
	"fmt"
)

const Base62Characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// ID of ulr model will we encoded with b B62 algo

func Base62Encoding(urlID int) string {
	if urlID == 0 {
		return "0"

	}
	//stores ID in storeID
	var storeID []byte
	for urlID > 0 {
		remainder := urlID % 62
		storeID = append(storeID, Base62Characters[remainder])
		urlID = urlID / 62
	}

	//reversing the encoded string to get correct base62
	for i, j := 0, len(storeID)-1; i < j; i, j = i+1, j-1 {
		storeID[i], storeID[j] = storeID[j], storeID[i]
	}
	fmt.Println(storeID)
	res := string(storeID)
	fmt.Println("res", res)
	return res
}
