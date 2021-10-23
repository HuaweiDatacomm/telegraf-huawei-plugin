package main

import (
	"encoding/base64"
	"fmt"
)

func main() {

	var password string
	fmt.Println("please enter password:")
	fmt.Scanln(&password)
	enc_password := base64.URLEncoding.EncodeToString([]byte(password))
	fmt.Println("Encrypted password:" + enc_password)

}
