package base62

import (
	"fmt"
	"testing"
)

func TestBase62(t *testing.T) {
	fmt.Println(GetBase62(6347))
}

func TestStrToint64(t *testing.T) {
	fmt.Println(StrToint64("fuck"))
}
