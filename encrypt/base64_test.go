package encrypt

import (
	"log"
	"testing"
)

func TestBase64(t *testing.T) {
	name := "Lucifer"
	newStr := Base64Encode(name)
	log.Println(newStr)
	oldStr := Base64Decode(newStr)
	log.Println(oldStr)
}
