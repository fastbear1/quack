package runner

import (
	"fmt"

	. "github.com/fastbear1/quack/schema"
)

// delta - difference between two arrays. if item present in left array a
type Item struct {
	delta int8  // difference between two arrays
	l     int16 // position on left array
	r     int16 // position in right array
}

const (
	leftShift  = 1
	rightShift = -1
)

type catalog map[string]*Item

func getCatalog[T Meta](leftData []T, rightData []T) catalog {
	ct := make(catalog, 0)
	for p, v := range leftData {
		name := v.GetName()
		if _, ok := ct[name]; !ok {
			ct[name] = &Item{leftShift, int16(p), 0}
		} else {
			ct[name].l = int16(p)
			ct[name].delta += leftShift
		}
	}

	for p, v := range rightData {
		name := v.GetName()
		if _, ok := ct[name]; !ok {
			ct[name] = &Item{rightShift, 0, int16(p)}
		} else {
			ct[name].r = int16(p)
			ct[name].delta += rightShift
		}
	}
	return ct
}

/*
type Item struct {
	delta int8  // diffence between two arrays
	l     int16 // position on left array
	r     int16 // position in right array
}

type catalog map[string]*Item

func fillCatalog(*ct catalog, name string, pos int16, direction ShiftType) {
	if direction == leftShift {
		if _, ok := ct[name]; !ok {
			ct[name] = &Item{leftShift, pos, 0}
		} else {
			ct[name].l = p
			ct[name].delta += leftShift
		}
	} else {
		if _, ok := ct[name]; !ok {
			ct[name] = &Item{rightShift, 0, pos}
		} else {
			ct[name].l = p
			ct[name].delta += rightShift
		}
	}
}
func getCatalog(leftdata, rightData) {
	for leftData:
		fillCatalog(TableCatalog)

	for rightData:
		fillCatalog(TableCatalog)
	return &Catalog
}

switch TableCatalog:
	case 0:
		ColumnsCatalog := getCatalog(leftData[l].Columns, rightData[r].Columns)
		switch ColumnsCatalog:
			case 0:
				checlColumnState()
			case 1:
				createColumnCommand()
			case -1:
				dropColumnCommand()
		IndexCatalog := getCatalog(leftData[l].Indices, rightData[r].Indices)
		switch IndexCatalog
			case 0:
				checklIndexState()
			case 1:
				createIndexCommand()
			case -1:
				dropIndexCommand()
		RefCatalog := getCatalog(leftData[l].References, rightData[r].References)
		switch RefCatalog
			case 0:
				checklRefState()
			case 1:
				createRefCommand()
			case -1:
				dropRefCommand()
	case 1:
		createTableCommand()
	case -1:
		dropTableCommand()


*/
/*
Catalog - map[string]int
default value for each element is 1
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

	tablesCatalog := getCatalog(gmeta, dbmeta)
	for k, v := range tablesCatalog {
		switch v.delta {
		case 1:
			fmt.Printf("Present on left side only %s (%v)\n", k, gmeta[v.l].Name)
			upFuncList = append(upFuncList, gmeta[v.l].CreateTable)
			downFuncList = append(downFuncList, gmeta[v.l].DeleteTable)
			for _, i := range gmeta[v.l].Indeces {
				upFuncList = append(upFuncList, i.CreateIndex)
				downFuncList = append(downFuncList, i.DropIndex)
			}
		case -1:
			fmt.Printf("Present on right side only %s (%v)\n", k, dbmeta[v.r].Name)
			upFuncList = append(upFuncList, dbmeta[v.r].DeleteTable)
			downFuncList = append(downFuncList, dbmeta[v.r].CreateTable)
		case 0:
			fmt.Printf("Check items in both arrays %s (%v - %v)\n", k, gmeta[v.l].Name, dbmeta[v.r].Name)
			// Check columns first
			columnsCatalog := getCatalog(gmeta[v.l].Columns, dbmeta[v.r].Columns)
			for j, b := range columnsCatalog {
				switch b.delta {
				case 1:
					fmt.Printf("Column presented on left side only %s (%v)\n", j, gmeta[v.l].Columns[b.l])
					upFuncList = append(upFuncList, gmeta[v.l].Columns[b.l].CreateColumn)
					downFuncList = append(downFuncList, gmeta[v.l].Columns[b.l].DeleteColumn)
				case -1:
					fmt.Printf("Column presented on right side only %s (%v)\n", j, dbmeta[v.r].Columns[b.r])
					upFuncList = append(upFuncList, dbmeta[v.r].Columns[b.r].DeleteColumn)
					downFuncList = append(downFuncList, dbmeta[v.r].Columns[b.r].CreateColumn)
				case 0:
					fmt.Printf("Column presented on both sides %s (%v)\n", j, dbmeta[v.r].Columns[b.r])
					if state := columnSchemaChanged(&gmeta[v.l].Columns[b.l], &dbmeta[v.r].Columns[b.r]); state {
						upFuncList = append(upFuncList, gmeta[v.l].Columns[b.l].AlterColumn)
						downFuncList = append(downFuncList, dbmeta[v.r].Columns[b.r].AlterColumn)
					}

				}
			}
			// same shit for indices
			indicesCatalog := getCatalog(gmeta[v.l].Indeces, dbmeta[v.r].Indeces)
			for n, ind := range indicesCatalog {
				switch ind.delta {
				case 1:
					fmt.Printf("Index presented on left side only %s (%v)\n", n, gmeta[v.l].Indeces[ind.l])
					upFuncList = append(upFuncList, gmeta[v.l].Indeces[ind.l].CreateIndex)
					downFuncList = append(downFuncList, gmeta[v.l].Indeces[ind.l].DropIndex)
				case -1:
					fmt.Printf("Index presented on right side only %s (%v)\n", n, dbmeta[v.r].Indeces[ind.r])
					upFuncList = append(upFuncList, dbmeta[v.r].Indeces[ind.r].DropIndex)
					downFuncList = append(downFuncList, dbmeta[v.r].Indeces[ind.r].CreateIndex)
				case 0:
					fmt.Printf("Index presented on both sides %s (%v)\n", n, dbmeta[v.r].Columns[ind.r])
					if state := isIndexSchemaChanged(&gmeta[v.l].Indeces[ind.l], &dbmeta[v.r].Indeces[ind.r]); state {
						upFuncList = append(upFuncList, dbmeta[v.r].Indeces[ind.r].DropIndex, gmeta[v.l].Indeces[ind.l].CreateIndex)
						downFuncList = append(downFuncList, gmeta[v.l].Indeces[ind.l].DropIndex, dbmeta[v.r].Indeces[ind.r].CreateIndex)
					}

				}
			}
			//same shit for constraints
			refCatalog := getCatalog(gmeta[v.l].References, dbmeta[v.r].References)
			for n, ref := range refCatalog {
				switch ref.delta {
				case 1:
					fmt.Printf("Reference presented on left side only %s (%v)\n", n, gmeta[v.l].References[ref.l])
					upFuncList = append(upFuncList, gmeta[v.l].References[ref.l].CreateConstraint)
					downFuncList = append(downFuncList, gmeta[v.l].References[ref.l].DeleteConstraint)
				case -1:
					fmt.Printf("Reference presented on right side only %s (%v)\n", n, dbmeta[v.r].References[ref.r])
					upFuncList = append(upFuncList, dbmeta[v.r].References[ref.r].DeleteConstraint)
					downFuncList = append(downFuncList, dbmeta[v.r].References[ref.r].CreateConstraint)
				case 0:
					fmt.Printf("Reference presented on both sides %s (%v)\n", n, dbmeta[v.r].Columns[ref.r])
					if state := isReferenceSchemaChanged(&gmeta[v.l].References[ref.l], &dbmeta[v.r].References[ref.r]); state {
						upFuncList = append(upFuncList, dbmeta[v.r].References[ref.r].DeleteConstraint, gmeta[v.l].References[ref.l].CreateConstraint)
						downFuncList = append(downFuncList, gmeta[v.l].References[ref.l].DeleteConstraint, dbmeta[v.r].References[ref.r].CreateConstraint)
					}

				}
			}

		}
	}

	/*
		var (
			dbmap = map[string]TableMeta{}
			gmap  = map[string]TableMeta{}
		)
	*/

	/*
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
	*/

	/*
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
	*/

	// Check columns
	// check all column are exists
	/*
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
	*/
	// now check that columns has same parameters

	/*
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
	*/
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
