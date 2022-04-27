# crab 🦀
![License](https://img.shields.io/github/license/Bananenpro/crab)
![Lines of code](https://img.shields.io/tokei/lines/github/Bananenpro/crab)

An interpreted dynamically typed toy programming language.

## [Documentation](https://github.com/Bananenpro/crab/blob/main/DOCUMENTATION.md)

## Installation

### Prerequisites

- [Go](https://go.dev/) 1.18+

```sh
go install github.com/Bananenpro/crab@latest
```

You might need to add $GOBIN to your path.

## Hello World

```go
func main() {
    println("Hello World!");
}
```

```sh
crab helloworld.cb
```

## Features

- dynamic typing
- helpful error messages
- scopes and variable shadowing
- lists
- control flow statements
- ternary conditional
- functions
- multiple return values
- functions as values / closures
- exceptions
- useful builtin functions

## License

MIT License

Copyright (c) 2022 Julian Hofmann

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
