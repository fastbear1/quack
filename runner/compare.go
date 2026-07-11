package runner

import (
	d "github.com/fastbear1/quack/drivers"
)

/*
Catalog - map[string]int
default value for each column is 1
left array uses -1 delta, right array uses +1 delta
Final values:

	1 [1 + (-1) + (+1)] denotes that object presented in both arrays
	0 [1 + (-1)] only in left array
	2 [1 + (+1)] only in right array
*/
type Catalog struct {
	Names map[string]int8
}

func (c *Catalog) fillData(arr []string, delta int8) {
	for i := 0; i < len(arr); i++ {
		if _, ok := c.Names[arr[i]]; ok != true {
			c.Names[arr[i]] = 1 + delta
		} else {
			c.Names[arr[i]] += delta
		}
	}

}

func getCatalogData(left []string, right []string) (strUp []string, strDown []string) {

	cat := Catalog{Names: make(map[string]int8)}
	cat.fillData(left, -1)
	cat.fillData(right, 1)

	var both, onleft, onright []string
	for k, v := range cat.Names {
		if v == 1 {
			both = append(both, k)
		} else if v == 0 {
			onright = append(onright, k)
		} else {
			onleft = append(onleft, k)
		}
	}
	return onleft, onright
}

// TODO: expensive code
func StateDifference[T d.Meta](leftArray []T, rightArray []T) ([]T, []T, []T) {
	var retLeft, retRight []T
	var leftMap, rightMap = map[string]T{}, map[string]T{}

	for _, i := range leftArray {
		leftMap[i.GetName()] = i
	}
	for _, j := range rightArray {
		rightMap[j.GetName()] = j
	}

	var leftNames, rightNames []string

	for l, _ := range leftMap {
		leftNames = append(leftNames, l)
	}
	for r, _ := range rightMap {
		rightNames = append(rightNames, r)
	}

	missedRight, missedLeft := getCatalogData(leftNames, rightNames)

	for _, lname := range missedRight {
		retRight = append(retRight, rightMap[lname])
	}
	for _, rname := range missedLeft {
		retLeft = append(retLeft, leftMap[rname])
	}

	var alterColumns []T
	// compare column parameters
	for k, lv := range leftMap {
		if rv, ok := rightMap[k]; ok {
			state := compareColumnState(lv, rv)

			if !state {
				alterColumns = append(alterColumns, lv)
			}
		}
	}

	return retLeft, retRight, alterColumns
}

func compareColumnState(l d.Meta, r d.Meta) bool {
	left := l.(d.Column)
	right := r.(d.Column)
	if left.DataType != right.DataType || left.IsNullable != right.IsNullable || left.IsPrimary != right.IsPrimary || left.PrimaryConstraint != right.PrimaryConstraint {
		return false
	}
	return true
}
