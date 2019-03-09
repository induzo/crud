# A little opinionated CRUD interface

[![Documentation](https://godoc.org/github.com/induzo/crud?status.svg)](http://godoc.org/github.com/induzo/crud) [![Go Report Card](https://goreportcard.com/badge/github.com/induzo/crud)](https://goreportcard.com/report/github.com/induzo/crud) [![Maintainability](https://api.codeclimate.com/v1/badges/937d15e44061eeb32877/maintainability)](https://codeclimate.com/github/induzo/crud/maintainability)[![Coverage Status](https://coveralls.io/repos/github/induzo/crud/badge.svg?branch=master)](https://coveralls.io/github/induzo/crud?branch=master) [![CircleCI](https://circleci.com/gh/induzo/crud.svg?style=svg)](https://circleci.com/gh/induzo/crud)

## Opinions

- Your entity ids should be using github.com/rs/xid
- Your crud errors should be handlable by an httpresponse

## Benchmarks (i7, 16GB)

```bash
    goos: linux
    goarch: amd64
    pkg: github.com/induzo/crud/rest
    BenchmarkPOSTHandler-8            300000              4190 ns/op            1912 B/op         20 allocs/op
    BenchmarkGETListHandler-8         300000              4694 ns/op            2131 B/op         31 allocs/op
    BenchmarkGETHandler-8             500000              2856 ns/op            1138 B/op         14 allocs/op
    BenchmarkDELETEHandler-8         3000000               596 ns/op              80 B/op          2 allocs/op
    BenchmarkPUTHandler-8             300000              4241 ns/op            1842 B/op         21 allocs/op
    BenchmarkPATCHHandler-8           500000              2772 ns/op            1544 B/op         17 allocs/op
```
