package runner

import (
	. "github.com/fastbear1/quack/schema"
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

func compareMetaState(dbmeta []TableMeta, gmeta []TableMeta) ([]func(drv DbInterface) string, []func(drv DbInterface) string) {
	var (
		upFuncList   []func(drv DbInterface) string
		downFuncList []func(drv DbInterface) string
		dbmap        = map[string]TableMeta{}
		gmap         = map[string]TableMeta{}
	)

	if len(dbmeta) == 0 {
		// return create table for all objects in gmeta
		for _, str := range gmeta {
			upFuncList = append(upFuncList, str.CreateTable)
			downFuncList = append(downFuncList, str.DeleteTable)
			for _, i := range str.Indeces {
				upFuncList = append(upFuncList, i.CreateIndex)
				downFuncList = append(downFuncList, i.DropIndex)
			}
		}
		// job done, not enough data for comparing
		return upFuncList, downFuncList
	}
	var (
		left  = make([]string, 0)
		right = make([]string, 0)
	)
	for _, l := range gmeta {
		left = append(left, l.Name)
		gmap[l.Name] = l
	}
	for _, r := range dbmeta {
		right = append(right, r.Name)
		dbmap[r.Name] = r
	}

	toDelete, toCreate := getCatalogData(left, right)

	if len(toDelete) > 0 || len(toCreate) > 0 {
		//TODO: inefficient
		for _, cr := range toCreate {
			for _, l := range gmeta {
				if cr == l.Name {
					upFuncList = append(upFuncList, l.CreateTable)
					downFuncList = append(downFuncList, l.DeleteTable)
					for _, idx := range l.Indeces {
						upFuncList = append(upFuncList, idx.CreateIndex)
						downFuncList = append(downFuncList, idx.DropIndex)
					}
				}
			}
		}

		for _, dt := range toDelete {
			for _, r := range dbmeta {
				if dt == r.Name {
					upFuncList = append(upFuncList, r.DeleteTable)
					downFuncList = append(downFuncList, r.CreateTable)
				}
			}
		}
	}

	// Check columns
	// check all column are exists
	for name, gtable := range gmap {
		if dbtable, ok := dbmap[name]; !ok {
			// Skipping tables that are not exists for now
			continue
		} else {
			toCreateCol, toDeleteCol, toAlterCol := StateDifference(gtable.Columns, dbtable.Columns)
			for _, c := range toCreateCol {
				upFuncList = append(upFuncList, c.CreateColumn)
				downFuncList = append(downFuncList, c.DeleteColumn)
			}
			for _, c := range toDeleteCol {
				upFuncList = append(upFuncList, c.DeleteColumn)
				downFuncList = append(downFuncList, c.CreateColumn)
			}
			// TODO: toAlterCol should return both table objects, first for altering Up, second for down migration
			for _, c := range toAlterCol {
				upFuncList = append(upFuncList, c.AlterColumn)
			}
		}
	}
	// now check that columns has same parameters

	// Same shit for references
	for name, gtable := range gmap {
		if dbtable, ok := dbmap[name]; !ok {
			// Skipping tables that are not exists for now
			continue
		} else {
			toCreateRef, toDeleteRef, toAlterRef := referenceStateChanged(gtable.References, dbtable.References)
			for _, c := range toCreateRef {
				upFuncList = append(upFuncList, c.CreateConstraint)
				downFuncList = append(downFuncList, c.DeleteConstraint)
			}
			for _, c := range toDeleteRef {
				upFuncList = append(upFuncList, c.DeleteConstraint)
				downFuncList = append(downFuncList, c.CreateConstraint)
			}
			for _, c := range toAlterRef {
				upFuncList = append(upFuncList, c[1].DeleteConstraint, c[0].CreateConstraint)
				downFuncList = append(downFuncList, c[0].DeleteConstraint, c[1].CreateConstraint)
			}
		}
	}

	// Same shit for indices
	for name, gtable := range gmap {
		if dbtable, ok := dbmap[name]; !ok {
			// Skipping tables that are not exists for now
			continue
		} else {
			toCreateIdx, toDeleteIdx, toAlterIdx := indicesStateChanged(gtable.Indeces, dbtable.Indeces)
			for _, c := range toCreateIdx {
				upFuncList = append(upFuncList, c.CreateIndex)
				downFuncList = append(downFuncList, c.DropIndex)
			}
			for _, c := range toDeleteIdx {
				upFuncList = append(upFuncList, c.DropIndex)
				downFuncList = append(downFuncList, c.CreateIndex)
			}
			for _, c := range toAlterIdx {
				upFuncList = append(upFuncList, c[1].DropIndex, c[0].CreateIndex)
				downFuncList = append(downFuncList, c[0].DropIndex, c[1].CreateIndex)
			}
		}
	}

	// Not implemented
	return upFuncList, downFuncList
}

// TODO: expensive code
func StateDifference[T Meta](leftArray []T, rightArray []T) ([]T, []T, []T) {
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
			lcol := any(lv).(Column)
			rcol := any(rv).(Column)
			if state := columnSchemaChanged(&lcol, &rcol); state {
				alterColumns = append(alterColumns, any(lcol).(T))
			}
		}
	}

	return retLeft, retRight, alterColumns
}

func columnSchemaChanged(left *Column, right *Column) bool {
	var changed bool = false
	if left.DataType != right.DataType {
		left.AlterState.Type = 0
		left.AlterState.DataType = right.DataType
		changed = true
	}
	if left.IsNullable != right.IsNullable {
		left.AlterState.Type = 1
		left.AlterState.IsNullable = right.IsNullable
		changed = true
	}
	if left.ColumnDefault != right.ColumnDefault {
		left.AlterState.Type = 2
		left.AlterState.ColumnDefault = right.ColumnDefault
		changed = true
	}
	return changed
}

func referenceStateChanged(leftArray []ReferenceMeta, rightArray []ReferenceMeta) ([]ReferenceMeta, []ReferenceMeta, [][]ReferenceMeta) {
	var retLeft, retRight []ReferenceMeta
	var leftMap, rightMap = map[string]ReferenceMeta{}, map[string]ReferenceMeta{}
	var leftNames, rightNames []string

	for _, i := range leftArray {
		leftMap[i.GetName()] = i
		leftNames = append(leftNames, i.GetName())
	}
	for _, j := range rightArray {
		rightMap[j.GetName()] = j
		rightNames = append(rightNames, j.GetName())
	}

	missedRight, missedLeft := getCatalogData(leftNames, rightNames)

	for _, lname := range missedRight {
		retRight = append(retRight, rightMap[lname])
	}
	for _, rname := range missedLeft {
		retLeft = append(retLeft, leftMap[rname])
	}

	var alterColumns [][]ReferenceMeta
	// compare column parameters
	for k, lv := range leftMap {
		if rv, ok := rightMap[k]; ok {
			if state := isReferenceSchemaChanged(&lv, &rv); state {
				alterColumns = append(alterColumns, []ReferenceMeta{lv, rv})
			}
		}
	}

	return retLeft, retRight, alterColumns

}

func isReferenceSchemaChanged(l *ReferenceMeta, r *ReferenceMeta) bool {
	if l.RefColumn != r.RefColumn || l.RefTable != r.RefTable || l.RefOptions != r.RefOptions {
		return true
	}
	return false
}

func indicesStateChanged(leftArray []IndexMeta, rightArray []IndexMeta) ([]IndexMeta, []IndexMeta, [][]IndexMeta) {
	var retLeft, retRight []IndexMeta
	var leftMap, rightMap = map[string]IndexMeta{}, map[string]IndexMeta{}
	var leftNames, rightNames []string

	for _, i := range leftArray {
		leftMap[i.GetName()] = i
		leftNames = append(leftNames, i.GetName())
	}
	for _, j := range rightArray {
		rightMap[j.GetName()] = j
		rightNames = append(rightNames, j.GetName())
	}

	missedRight, missedLeft := getCatalogData(leftNames, rightNames)

	for _, lname := range missedRight {
		retRight = append(retRight, rightMap[lname])
	}
	for _, rname := range missedLeft {
		retLeft = append(retLeft, leftMap[rname])
	}

	var alterColumns [][]IndexMeta
	// compare column parameters
	for k, lv := range leftMap {
		if rv, ok := rightMap[k]; ok {
			if state := isIndexSchemaChanged(&lv, &rv); state {
				alterColumns = append(alterColumns, []IndexMeta{lv, rv})
			}
		}
	}

	return retLeft, retRight, alterColumns

}

func isIndexSchemaChanged(l *IndexMeta, r *IndexMeta) bool {
	if l.Unique != r.Unique || l.Type != r.Type || len(l.Columns) != len(r.Columns) {
		return true
	}
	// additionaly check column in expression
	if l.Columns[0].Expression != r.Columns[0].Expression {
		return true
	}

	lfields := map[string]uint16{}
	rfields := map[string]uint16{}

	for _, lv := range l.Columns {
		lfields[lv.Field] = uint16(lv.Priority)
	}
	for _, rv := range r.Columns {
		rfields[rv.Field] = uint16(rv.Priority)
	}

	for lk, lv := range lfields {
		rv, ok := rfields[lk]
		if !ok || lv != rv {
			return true
		}
	}
	return false
}
