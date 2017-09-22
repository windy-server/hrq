# hrq
[![build status](https://secure.travis-ci.org/windy-server/hrq.svg?branch=master)](http://travis-ci.org/windy-server/hrq) [![GoDoc](https://godoc.org/github.com/windy-server/hrq?status.png)](http://godoc.org/github.com/windy-server/hrq)  

Http client like requests in Go

```Go
import (
    "fmt"

    "github.com/windy-server/hrq"
)

func main() {
    req, _ := hrq.Get("http://example.com")
    res, _ := req.Send()
    s, _ := res.Text()
    fmt.Print(s)
}
```

## Table of contents

* [Installation](https://github.com/windy-server/hrq#installation)
* [Usage](https://github.com/windy-server/hrq#usage)
  - [Request](https://github.com/windy-server/hrq#request)
  - [Response](https://github.com/windy-server/hrq#response)
  - [Header](https://github.com/windy-server/hrq#header)
  - [Cookie](https://github.com/windy-server/hrq#cookie)
  - [Timeout](https://github.com/windy-server/hrq#timeout)
  - [File](https://github.com/windy-server/hrq#file)
  - [JSON](https://github.com/windy-server/hrq#json)
  - [History](https://github.com/windy-server/hrq#history)
  - [Gzip](https://github.com/windy-server/hrq#gzip)
  - [Session](https://github.com/windy-server/hrq#session)

## Installation

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
params := map[string]string{
    "foo": "123",
    "bar": "456",
}

// http://example.com?foo=123&bar=456
url := hrq.MakeURL("http://example.com", params)
req, _ := hrq.Get(url)
res, _ := req.Send()
s, _ := res.Text()
fmt.Print(s)
```

#### Post

```Go
data := map[string]string{
    "foo": "123",
    "bar": "456",
}
req, _ := hrq.Post("http://example.com", data)
// When Content-Type is "application/x-www-form-urlencoded"(It is default),
// the request data is urlencoded.
// The request data must be a map[string]string instance.
// When Content-Type is "application/json",
// the request data is converted to json string.
// When Content-Type is "multipart/form-data",
// the request data is converted to fields.
res, _ := req.SetApplicationFormUrlencoded().Send()
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
v := res.HeaderValue("foo")
fmt.Print(v)
```

### Cookie

```Go
req, _ := hrq.Get("http://example.com")
req.PutCookie("abc", "efg")
res, _ := req.Send()
v := res.CookieValue("foo")
cm := res.CookiesMap()
```

### Timeout

```Go
req, _ := hrq.Get("http://example.com")
// This sets requset timeout to 30 seconds.
// (Default timeout is 15 seconds.)
req.SetTimeout(30)
res, _ := req.Send()
```

### File

```Go
data := map[string]string{
    "foo": "123",
    "bar": "456",
}
req, _ := hrq.Post("http://example.com", data)
// When Content-Type is "multipart/form-data",
// the request data is converted to fields.
req.SetMultipartFormData()
file, _ := os.Open("foo.gif")
req.AddFile("image/gif", "foo", "foo.gif", file)
res, _ := req.Send()
```

### JSON

```Go
data := map[string][]string{
    "foo": "123",
    "bar": "456",
}
req, _ := hrq.Post("http://example.com", data)
// When Content-Type is "application/json",
// the request data is converted to json string.
req.SetApplicationJSON()
res, _ := req.Send()
var result map[string]string
err := res.JSON(&result)
```

### History

```Go
req, _ := hrq.Get("http://example.com")
res, _ := req.Send()
// The redirect history by http.Request slice
history := req.History
// The recent request
res.Request
```

### Gzip

```Go
data := map[string][]string{
    "foo": "123",
    "bar": "456",
}
req, _ := hrq.Post("http://example.com", data)
// You can send a request compressed by gzip.
req.UseGzip()
res, _ := req.SetApplicationJSON().Send()
```

### Session

```Go
session, _ := NewSession()
data := map[string][]string{
    "foo": "123",
    "bar": "456",
}
req, _ := hrq.Post("http://example.com", data)
res, _ := session.Send(req)
```
