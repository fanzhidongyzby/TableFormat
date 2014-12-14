# TableFormat

## Introduction

TableFormat is a formatting tool for string in Go programing language, which can convert a string to table style with UTF-8 characters.

## Usage

* Install by `go get github.com/fanzhidongyzby/TableFormat`
* Import package `github.com/fanzhidongyzby/table`
* Use `table.Format` to format your string

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
* `func Format (data string) string` : to format a string to table style<br>
* `func FormatMap (m map[string]string) string` : to format a map to table style<br>

## Options

Follow Options are provided:<br>
* `RowSeparator string = "\n"			//Separate rows`
* `ColumnSeparator string = ""			//Separate columns, empty string means all the space characters`
* `Placeholder string = "_"				//Represent an empty table field`
* `BlankFilling string = ""				//What to be filled in blank table field when row's too short`
* `BlankFillingForHeader string = ""	//What to be filled in blank header field when row's too long`
* `ColOverflow bool = true				//Discard more columns or not when row's too long`
* `UseBoard bool = true					//Use utf8 character to print board`
* `SpaceAlt byte = ' '					//What to replace \n \b \t ...`
* `OverFlowSeparator string = " "		//What to join overflow columns`
* `CenterFilling byte = " "				//What to charactor into field in order to centralize`


