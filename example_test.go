package hrq_test

import (
	"fmt"

	"github.com/windy-server/hrq"
)

func Example() {
	req, _ := hrq.Get("http://example.com")
	res, _ := req.Send()
	s, _ := res.Content()
	fmt.Print(s)
}
