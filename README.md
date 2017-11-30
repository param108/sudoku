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