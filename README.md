# TableFormat

## Introduction

TableFormat is a formatting tool for type in Go programing language, which can convert any type to table style with UTF-8 characters.

## Usage

* Install by `go get github.com/fanzhidongyzby/TableFormat`
* Import package `github.com/fanzhidongyzby/table`
* Use `table.Format` to format your type

## Demo

For the code below:
```go
package main

import (
	"fmt"
	"github.com/fanzhidongyzby/table"
)

func main() {
	str := ` ID _ Num Digit
		1 2 3你好
		4 _ 5 
		7 8 9 10 11`
	fmtStr := table.Format(str)
	fmt.Print(fmtStr)
}
```

Its output in console is:<br>
![](https://github.com/fanzhidongyzby/TableFormat/blob/master/image/output.jpg)<br>
In default, table rows are separated by '\n' and columns are separated by space character, including ' ', '\t', '\v', '\b', '\f' and son on.<br>
If you need to define your own separators, some options are provided by table package. See the options below for details.<br>

## APIs

Following APIs are provided:<br>
* `func Format (obj interface{}) string` : to format anything to table style<br>

## Options

Follow Options are provided:<br>
* `RowSeparator string = "\n"           //Separate rows`
* `ColumnSeparator string = ""          //Separate columns, empty string means all the space characters`
* `Placeholder string = "_"             //Represent an empty table field`
* `BlankFilling string = ""             //What to be filled in blank table field when row is too short`
* `BlankFillingForHeader string = ""    //What to be filled in blank header field`
* `ColOverflow bool = true              //Do not discard more columns or not when row is too long`
* `UseBoard bool = true                 //Use utf8 character to print board`
* `SpaceAlt byte = ' '                  //What to replace \n \b \t ...`
* `OverFlowSeparator string = " "       //What to join overflow columns`
* `CenterFilling byte = " "             //What to be filled into field in order to centralize`
* `IgnoreEmptyHeader bool = true		//Whether ignore empty header when all header fields are placeholder`
<br>
Use `defer table.Reset()` to confirm all the options set to default after your last configuration.
