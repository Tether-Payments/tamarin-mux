<div align="center">

# Tamarin Mux

![tamarin.png](tamarin.png)

</div>

---

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/Tether-Payments/tamarin-mux)](https://goreportcard.com/report/Tether-Payments/tamarin-mux)
[![codecov](https://codecov.io/gh/Tether-Payments/tamarin-mux/graph/badge.svg?token=ZBQI4PN2CA)](https://codecov.io/gh/Tether-Payments/tamarin-mux)
[![Maintainability](https://api.codeclimate.com/v1/badges/0221a6290e3ca9fca370/maintainability)](https://codeclimate.com/github/Tether-Payments/tamarin-mux/maintainability)
[![CodeQL](https://github.com/Tether-Payments/tamarin-mux/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/Tether-Payments/tamarin-mux/actions/workflows/github-code-scanning/codeql)

</div>

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