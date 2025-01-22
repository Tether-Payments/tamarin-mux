# Tamarin Mux

![tamarin.png](tamarin.png)

---

Written in Go ("Golang" for search engines) with zero external dependencies, this package implements a clean, 
non-bloated, HTTP request multiplexer.

---

## Installation & Usage
1. Once confirming you have [Go](https://go.dev/doc/install) installed, the command below will add
   `tamarin` as a dependency to your Go program.
```go
go get -u github.com/tether-payments/tamarin-mux
```
2. Import the package into your code
```go
package main

import (
"github.com/Tether-Payments/tamarin-mux"
)
```
3. [Examples](examples)