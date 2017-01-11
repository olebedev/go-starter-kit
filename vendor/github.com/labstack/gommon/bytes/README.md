# Bytes

Format bytes to string

## Installation

```go
go get github.com/labstack/gommon/bytes
```

## [Usage](https://github.com/labstack/gommon/blob/master/bytes/bytes_test.go)

```sh
import github.com/labstack/gommon/bytes
```

### Decimal prefix

```go
fmt.Println(bytes.Format(1323))
```

`1.32 KB`

### Binary prefix

```go
bytes.SetBinaryPrefix(true)
fmt.Println(bytes.Format(1323))
```

`1.29 KiB`

### New instance

```go
g := New()
fmt.Println(g.Format(13231323))
```
