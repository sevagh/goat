[![Build Status](https://api.travis-ci.org/sevagh/stdcap.svg?branch=master)](https://travis-ci.org/sevagh/stdcap)
[![Coverage Status](https://coveralls.io/repos/github/sevagh/stdcap/badge.svg?branch=master)](https://coveralls.io/github/sevagh/stdcap?branch=master)
[![ReportCard][ReportCard-Image]][ReportCard-Url]

[ReportCard-Url]: http://goreportcard.com/report/sevagh/stdcap
[ReportCard-Image]: http://goreportcard.com/badge/sevagh/stdcap

# stdcap
Package to capture stdout, stderr in Go

### Install

    $ go install github.com/sevagh/stdcap

### Usage

    import (
        "testing"
        "github.com/sevagh/stdcap"
    )

    func PrintsHello() string{
        fmt.Println("Hello world!")
        return "value"
    }

    ...

    var retval string
    sc := stdcap.StdoutCapture()

    out := sc.Capture(func() {
        retval = PrintsHello()
    })

    if out != "Hello world!\n" {
        t.Errorf("printed wrong thing to stdout")
    }

    if retval != "value" {
        t.Errorf("Returned incorrect value")
    }

