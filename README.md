# go-file

Package file is a Go library to open files with file locking depending on the system.

[![Build Status](https://travis-ci.org/mithrandie/go-file.svg?branch=master)](https://travis-ci.org/mithrandie/go-file)
[![GoDoc](https://godoc.org/github.com/mithrandie/go-file?status.svg)](http://godoc.org/github.com/mithrandie/go-file)
[![License: MIT](https://img.shields.io/badge/License-MIT-lightgrey.svg)](https://opensource.org/licenses/MIT)

## Install

```sql
go get github.com/mithrandie/go-file
```

## Supported Systems

Currently file locking on the following systems are supported.

### darwin dragonfly freebsd linux netbsd openbsd

Advisory Lock

### windows

Mandatory Lock

### android nacl plan9 solaris zos

Not Supported

## Example

```go
package main

import (
	"bufio"
	"fmt"
	 
	"github.com/mithrandie/go-file"
)

func main() {
	fp, err := file.OpenToRead("/path/to/file")
	if err != nil {
		panic(err)
	}
	defer file.Close(fp)

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
```
