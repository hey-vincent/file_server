package link

import (
	"encoding/base64"
	"fmt"
)

func LinkAuthorize(auth string) {
	encoding := base64.StdEncoding
	cipher := []byte("1234")
	fmt.Println(base64.StdEncoding.StroEncodeToString(cipher))
	data := make([]byte , base64.StdEncoding.DecodedLen(len(cipher)))
	encoding.Decode(data , cipher)

	fmt.Printf("%s" ,data)


}


