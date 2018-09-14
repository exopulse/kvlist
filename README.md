# exopulse kvlist package
Golang package kvlist contains key-value list.

[![CircleCI](https://circleci.com/gh/exopulse/kvlist.svg?style=svg)](https://circleci.com/gh/exopulse/kvlist)
[![Build Status](https://travis-ci.org/exopulse/kvlist.svg?branch=master)](https://travis-ci.org/exopulse/kvlist)
[![GitHub license](https://img.shields.io/github/license/exopulse/kvlist.svg)](https://github.com/exopulse/kvlist/blob/master/LICENSE)

# Overview

This package contains KeyValue and KeyValueList types for easier usage of such lists.

## Features

### KeyValue

Simple wrapper around key and value. Both key and value are strings.
 
### KeyValueList

Container for key-value pairs. Multiple keys are allowed.

# Using kvlist package

## Installing package

Use go get to install the latest version of the library.

    $ go get github.com/exopulse/kvlist
 
Include kvlist in your application.
```go
import "github.com/exopulse/kvlist"
```

## Sample statements
```go
// create new KeyValue
kv := kvlist.KeyValue{"key1", `value 日本語 "quote"`}

// create list and add multiple key-values
l := kvlist.New().Add("key1", "value1").Add("key1", "value2")

// create list and override key-values
l := kvlist.New().Put("key1", "value1").Put("key1", "value2")

// get count
c := l.Count()

// delete first occurrence of specified key
l.DeleteKey("key2")

// delete all key-value pairs for specified key
l.DeleteKeys("key1")

// remote all key-value pairs
l.Clear()

// get all key-value pairs for specified key
kvs := l.GetKeys("key1")

// read into string variable
var v string

if ok := l.ScanKey(&v, "key2"); ok {
	panic("key2 found")
}

// write to io.Writer
l := kvlist.New()
b := new(bytes.Buffer)

if _, err := b.ReadFrom(l); err != nil {
	tpanic(err)
}

// build from string
s := `key1="value \u65e5\u672c\u8a9e \"quote\""  key2="value 2"`
l, err := NewFromString(s)

if err != nil {
	panic(err)
}
```

# About the project

## Contributors

* [exopulse](https://github.com/exopulse)

## License

Kvlist package is released under the MIT license. See
[LICENSE](https://github.com/exopulse/kvlist/blob/master/LICENSE)
