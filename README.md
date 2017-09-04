# hrq
[![build status](https://secure.travis-ci.org/windy-server/hrq.svg?branch=master)](http://travis-ci.org/windy-server/hrq) [![GoDoc](https://godoc.org/github.com/windy-server/hrq?status.png)](http://godoc.org/github.com/windy-server/hrq)  

Http client like requests in Go

```Go
import (
	"fmt"

	hrq "github.com/windy-server/hrq"
)

func main() {
    req, _ := hrq.Get("http://example.com")
    res, _ := req.Send()
    s, _ := res.Text()
    fmt.Print(s)
}
```

## Table of contents

* [Install](https://github.com/windy-server/hrq#install)
* [Usage](https://github.com/windy-server/hrq#usage)
  - [Request](https://github.com/windy-server/hrq#request)
  - [Request](https://github.com/windy-server/hrq#response)
  - [Header](https://github.com/windy-server/hrq#header)
  - [Cookie](https://github.com/windy-server/hrq#cookie)
  - [Timeout](https://github.com/windy-server/hrq#timeout)

## Install

```
dep ensure -add github.com/windy-server/hrq
```

or 

```
go get -u github.com/windy-server/hrq
```


## Usage

### Request

hrq.Request inherits http.Request.

#### Get

```Go
req, _ := hrq.Get("http://example.com")
res, _ := req.Send()
s, _ := res.Text()
fmt.Print(s)
```

#### Post

```Go
data := map[string][]string{
    "foo": []string{"123"},
    "bar": []string{"456"},
}
req, _ := hrq.Post("http://example.com", data)
// When Content-Type is "application/x-www-form-urlencoded"(It is default),
// the request data is urlencoded.
// When Content-Type is "application/json",
// the request data is converted to json string.
req.SetHeader("Content-Type", "application/json")
res, _ := req.Send()
s, _ := res.Text()
fmt.Print(s)
```

### Response

hrq.Response inherits http.Response.


```Go
req, _ := hrq.Get("http://example.com")
res, _ := req.Send()
// get request body by byte
b, _ := res.Content()
// get request body by string
s, _ := res.Text()
fmt.Print(s)
```

### Header

```Go
req, _ := hrq.Get("http://example.com")
req.SetHeader("abc", "efg")
res, _ := req.Send()
v := res.GetHeader("foo")
fmt.Print(v)
```

### Timeout

```Go
req, _ := hrq.Get("http://example.com")
// This sets requset timeout to 30 seconds.
// (Default timeout is 15 seconds.)
req.SetTimeout(30)
res, _ := req.Send()
cm := res.CookiesMap()
```
