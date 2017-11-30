package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	no             []int
	RelatedRows    []int
	RelatedCols    []int
	RelatedSquares []int
	options        []int
	v              int
}

// assumed that the NOs have been calculated
func (data *Data) Options() []int {
	ret := []int{}

	if data.v > 0 {
		return []int{data.v}
	}

	for i := 1; i < 10; i++ {
		found := false
		for _, v := range data.no {
			if v == i {
				found = true
			}
		}
		if !found {
			ret = append(ret, i)
		}
	}
	return ret
}

type Row struct {
	RowIndex int
	ColIndex int
	Val      [9]*Data
}

func (row Row) Needs() []int {
	ret := []int{}
	found := map[int]bool{}
	for _, val := range row.Val {
		if val.v > 0 {
			found[val.v] = true
		}
	}
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		if _, ok := found[i]; !ok {
			ret = append(ret, i)
		}
	}
	return ret
}

func (row Row) Check() bool {
	m := map[int]int{}
	for _, d := range row.Val {
		if d.v < 1 {
			continue
		}
		if _, ok := m[d.v]; ok {
			return false
		} else {
			m[d.v] = 1
		}
	}
	return true
}

type Board struct {
	Src     string
	Rows    [9]Row
	Cols    [9]Row
	Squares [9]Row
	Entries []*Data
}

func (board *Board) Check() int {
	incomplete := false
	incorrect := false
	for _, d := range board.Entries {
		if d.v < 1 {
			incomplete = true
			break
		}
	}

	for _, row := range board.Rows {
		if !row.Check() {
			incorrect = true
			break
		}
	}
	for _, row := range board.Cols {
		if !row.Check() {
			incorrect = true
			break
		}
	}
	for _, row := range board.Squares {
		if !row.Check() {
			incorrect = true
			break
		}
	}

	if incorrect {
		return -1
	}

	if incomplete {
		return 0
	}

	return 1
}

func (board *Board) Print() {
	rowIdx := 0
	for _, row := range board.Rows {
		rowStr := "|"
		idx := 0
		for _, d := range row.Val {
			if d.v < 1 {
				rowStr = rowStr + "."
			} else {
				rowStr = rowStr + strconv.Itoa(d.v)
			}
			idx += 1
			rowStr = rowStr + "|"
		}
		rowIdx += 1
		fmt.Println(rowStr)
	}
}

func (board *Board) PopulateNos() bool {
	changed := false
	for _, entry := range board.Entries {
		seen := []int{}
		seenMap := map[int]bool{}
		if entry.v > 0 {
			no := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
			entry.no = append(no[:(entry.v-1)], no[entry.v:]...)
			continue
		}
		for _, idx := range entry.RelatedRows {
			for _, d := range board.Rows[idx].Val {
				if d.v > 0 {
					if _, ok := seenMap[d.v]; !ok {
						seen = append(seen, d.v)
						seenMap[d.v] = true
					}
				}
			}
		}

		for _, idx := range entry.RelatedCols {
			for _, d := range board.Cols[idx].Val {
				if d.v > 0 {
					if _, ok := seenMap[d.v]; !ok {
						seen = append(seen, d.v)
						seenMap[d.v] = true
					}
				}
			}
		}

		for _, idx := range entry.RelatedSquares {
			for _, d := range board.Squares[idx].Val {
				if d.v > 0 {
					if _, ok := seenMap[d.v]; !ok {
						seen = append(seen, d.v)
						seenMap[d.v] = true
					}
				}
			}
		}

		if len(seen) == 8 {
			m := map[int]bool{}
			for _, i := range seen {
				m[i] = true
			}

			for i := 1; i < 10; i++ {
				if _, ok := m[i]; !ok {
					entry.v = i
					changed = true
					break
				}
			}
		}
		entry.no = seen
		entry.options = entry.Options()
	}
	return changed
}

func (board *Board) Copy() *Board {
	rowStr := ""
	rowIdx := 0
	for _, row := range board.Rows {
		idx := 0
		for _, d := range row.Val {
			if d.v < 1 {
				rowStr = rowStr + "."
			} else {
				rowStr = rowStr + strconv.Itoa(d.v)
			}
			idx += 1
			if !(rowIdx == 8 && idx == 9) {
				rowStr = rowStr + ","
			}
		}
		rowIdx += 1
	}

	ret := &Board{}
	ret.Init(rowStr)
	return ret
}

func (board *Board) Init(src string) {
	l := strings.Split(src, ",")
	board.Entries = make([]*Data, 0)

	for idx, data := range l {
		entry := &Data{}
		entry.no = make([]int, 0)
		entry.RelatedCols = make([]int, 0)
		entry.RelatedRows = make([]int, 0)
		entry.RelatedSquares = make([]int, 0)
		board.Entries = append(board.Entries, entry)
		d := strings.TrimSpace(data)
		if d == "." {
			entry.v = -1
		} else {
			v, err := strconv.Atoi(d)
			if err != nil {
				panic("Failing to read file:" + err.Error())
			}
			entry.v = v
		}
		col := idx % 9
		row := idx / 9
		sqrRow := row / 3
		sqrCol := col / 3
		board.Rows[row].Val[col] = entry
		entry.RelatedRows = append(entry.RelatedRows, row)
		board.Cols[col].Val[row] = entry
		entry.RelatedCols = append(entry.RelatedCols, col)

		sqrRowIndex := row - (sqrRow * 3)
		sqrColIndex := col - (sqrCol * 3)

		board.Squares[(sqrRow*3)+(sqrCol)].Val[sqrRowIndex*3+(sqrColIndex)] = entry
		entry.RelatedSquares = append(entry.RelatedSquares, (sqrRow*3)+(sqrCol))
	}
}

func (board *Board) SolveOnes() bool {
	found := true
	changed := false
	for found {
		found = false
		for _, row := range board.Rows {
			if val := row.Needs(); len(val) == 1 {
				for _, r := range row.Val {
					if r.v == -1 {
						r.v = val[0]
						found = true
						changed = true
						break
					}
				}
			}
		}
		for _, row := range board.Cols {
			if val := row.Needs(); len(val) == 1 {
				for _, r := range row.Val {
					if r.v == -1 {
						r.v = val[0]
						found = true
						changed = true
						break
					}
				}
			}
		}
		for _, row := range board.Squares {
			if val := row.Needs(); len(val) == 1 {
				for _, r := range row.Val {
					if r.v == -1 {
						r.v = val[0]
						found = true
						changed = true
						break
					}
				}
			}
		}
	}
	return changed
}

func (row Row) SolveNos() bool {
	changed := false
	noCount := map[int]int{}
	for _, r := range row.Val {
		for _, no := range r.no {
			if v, ok := noCount[no]; ok {
				noCount[no] = v + 1
			} else {
				noCount[no] = 1
			}
		}
	}
	for k, v := range noCount {
		if v == 8 {
			for _, r := range row.Val {
				if r.v < 1 {
					found := false
					for _, no := range r.no {
						if no == k {
							found = true
							break
						}
					}
					if !found {
						r.v = k
						changed = true
						break
					}
				}
			}
		}
	}

	return changed
}

func (board *Board) CountNos() bool {
	changed := false
	for _, row := range board.Rows {
		if row.SolveNos() {
			changed = true
		}
	}
	for _, row := range board.Cols {
		if row.SolveNos() {
			changed = true
		}
	}
	for _, row := range board.Squares {
		if row.SolveNos() {
			changed = true
		}
	}
	return changed

}

func (board *Board) Brute(done chan bool) {
	// find the Data with minimum options
	var firstEntry *Data
	chosenIdx := 0
	for i, entry := range board.Entries {
		// find the first unfilled entry
		if entry.v < 1 {
			if firstEntry == nil {
				firstEntry = entry
				chosenIdx = 0
			} else {
				if len(firstEntry.options) > len(entry.options) {
					firstEntry = entry
					chosenIdx = i
				}
			}
		}
	}

	donedanadone := []chan bool{}
	for idx, v := range firstEntry.options {
		donedanadone = append(donedanadone, make(chan bool))
		newbrd := board.Copy()
		newbrd.Entries[chosenIdx].v = v
		fmt.Println("DESCENDING", idx)
		go newbrd.SolveWrapper(donedanadone[idx])
	}

	ret := false

	for idx, v := range donedanadone {
		result := <-v
		if result {
			ret = true
		}
		fmt.Println("ASCENDING", idx)
	}

	done <- ret
	return
}

func (board *Board) SolveWrapper(done chan bool) {
	done <- board.Solve()
}

func (board *Board) Solve() bool {
	// check if any of them have only one missing, if so add it.
	i := 0
	for true {
		i += 1
		changed := false
		if board.SolveOnes() {
			changed = true
		}
		if board.PopulateNos() {
			changed = true
		}
		if board.CountNos() {
			changed = true
		}
		check := board.Check()
		if !changed {
			if check == 1 {
				fmt.Println("SUCCESS SOLUTION" + strconv.Itoa(i))
				board.Print()
				return true
			}

			if check == 0 {
				// now need to resort to brute force
				done := make(chan bool)
				go board.Brute(done)
				fmt.Println("WAITING")
				dval := <-done
				fmt.Println("DONE")
				return dval
			}

			if check < 0 {
				return false
			}

		} else {

			if check == 1 {
				fmt.Println("SUCCESS SOLUTION" + strconv.Itoa(i))
				board.Print()
				return true
			}

			if check < 0 {
				return false
			}

		}
	}
	// can never reach here.
	return false
}

func loadBoard(fname string) (*Board, error) {
	board := &Board{}
	src, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	board.Init(string(src))
	return board, nil
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		panic("Invalid number of arguments: Usage 'sudoku <initial board file>'")
	}

	board, err := loadBoard(args[0])
	if err != nil {
		panic("Failed to load board:" + err.Error())
	}

	board.Solve()
}
