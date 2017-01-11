# http2curl
:triangular_ruler: Convert Golang's http.Request to CURL command line

[![Build Status](https://travis-ci.org/moul/http2curl.svg?branch=master)](https://travis-ci.org/moul/http2curl)
[![GoDoc](https://godoc.org/github.com/moul/http2curl?status.svg)](https://godoc.org/github.com/moul/http2curl)
[![Coverage Status](https://coveralls.io/repos/moul/http2curl/badge.svg)](https://coveralls.io/github/moul/http2curl)

## Example

```go
import "github.com/moul/http2curl"

req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString(`{"hello":"world","answer":42}`))
req.Header.Set("Content-Type", "application/json")

command, _ := GetCurlCommand(req)
fmt.Println(command)
// Output: curl -X PUT -d "{\"hello\":\"world\",\"answer\":42}" -H "Content-Type: application/json" http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu
```

## Install

```php
$ go get github.com/moul/http2curl
```

## License

MIT
