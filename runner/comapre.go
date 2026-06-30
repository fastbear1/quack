package runner

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

func GetCatalogData(left []string, right []string) (strUp []string, strDown []string) {

	cat := Catalog{Names: make(map[string]int8)}
	cat.fillData(left, -1)
	cat.fillData(right, 1)

	var both, onleft, onright []string
	for k, v := range cat.Names {
		if v == 1 {
			both = append(both, k)
		} else if v == 2 {
			onright = append(onright, k)
		} else {
			onleft = append(onleft, k)
		}
	}
	return onleft, onright
}
