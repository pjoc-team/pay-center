package center

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	response := Request("http://127.0.0.1:18080")
	fmt.Println(response)
}
