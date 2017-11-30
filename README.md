## A Sudoku Solver in Go

This solver uses the relationships between rows, columns and squares to fill up as many numbers as it can. Once that is exhausted, it calculates remaining options for the empty squares and parallely tries them out, automatically stopping further recursion when a collision is seen.

There are 2 test files in the repository `test1.dat` and `test2.dat`.

### Format of Input
Comma seperated values. White spaces are ignored. Last value should not have `,` after it.

`.` means an empty square while any other number is the value of that file.

For example
```
.,2,.,.,.,.,.,3,.,
1,.,5,.,.,.,9,.,8,
9,3,.,5,.,2,.,7,6,
7,.,.,4,.,5,.,.,2,
.,.,.,.,.,.,.,.,.,
2,.,.,8,.,6,.,.,9,
5,6,.,7,.,1,.,8,4,
8,.,1,.,.,.,6,.,7,
.,4,.,.,.,.,.,9,.
```

### Building

`go build sudoku.go`

### Running

`sudoku <filename>`

### Output
```
WAITING
DESCENDING 0
DESCENDING 1
ASCENDING 0
SUCCESS SOLUTION10
|1|6|9|5|7|8|2|4|3|
|5|4|7|2|3|1|6|8|9|
|2|8|3|4|6|9|1|5|7|
|3|9|8|1|4|2|5|7|6|
|4|1|6|3|5|7|8|9|2|
|7|5|2|9|8|6|4|3|1|
|8|7|5|6|1|3|9|2|4|
|9|3|1|8|2|4|7|6|5|
|6|2|4|7|9|5|3|1|8|
ASCENDING 1
DONE
```

If there are multiple solutions, they will also be shown.