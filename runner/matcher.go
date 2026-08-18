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

func compareMetaState(gmeta []TableMeta, dbmeta []TableMeta) ([]func(drv DbInterface) string, []func(drv DbInterface) string) {
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
			left := &gmeta[v.l]
			right := &dbmeta[v.r]
			fmt.Printf("type of left array %T\n", left)
			columnsCatalog := getCatalog(left.Columns, right.Columns)
			for n, col := range columnsCatalog {
				switch col.delta {
				case 1:
					fmt.Printf("Column presented on left side only %s (%v)\n", n, left.Columns[col.l])
					upFuncList = append(upFuncList, left.Columns[col.l].CreateColumn)
					downFuncList = append(downFuncList, left.Columns[col.l].DeleteColumn)
				case -1:
					fmt.Printf("Column presented on right side only %s (%v)\n", n, right.Columns[col.r])
					upFuncList = append(upFuncList, right.Columns[col.r].DeleteColumn)
					downFuncList = append(downFuncList, right.Columns[col.r].CreateColumn)
				case 0:
					fmt.Printf("Column presented on both sides %s (%v)\n", n, right.Columns[col.r])
					if state := isColumnSchemaChanged(&left.Columns[col.l], &right.Columns[col.r]); state {
						upFuncList = append(upFuncList, left.Columns[col.l].AlterColumn)
						downFuncList = append(downFuncList, right.Columns[col.r].AlterColumn)
					}

				}
			}
			// same shit for indices
			indicesCatalog := getCatalog(left.Indeces, right.Indeces)
			for n, ind := range indicesCatalog {
				switch ind.delta {
				case 1:
					fmt.Printf("Index presented on left side only %s (%v)\n", n, left.Indeces[ind.l])
					upFuncList = append(upFuncList, left.Indeces[ind.l].CreateIndex)
					downFuncList = append(downFuncList, left.Indeces[ind.l].DropIndex)
				case -1:
					fmt.Printf("Index presented on right side only %s (%v)\n", n, right.Indeces[ind.r])
					upFuncList = append(upFuncList, right.Indeces[ind.r].DropIndex)
					downFuncList = append(downFuncList, right.Indeces[ind.r].CreateIndex)
				case 0:
					fmt.Printf("Index presented on both sides %s (%v)\n", n, right.Columns[ind.r])
					if state := isIndexSchemaChanged(&left.Indeces[ind.l], &right.Indeces[ind.r]); state {
						upFuncList = append(upFuncList, right.Indeces[ind.r].DropIndex, left.Indeces[ind.l].CreateIndex)
						downFuncList = append(downFuncList, left.Indeces[ind.l].DropIndex, right.Indeces[ind.r].CreateIndex)
					}

				}
			}
			//same shit for constraints
			refCatalog := getCatalog(left.References, right.References)
			for n, ref := range refCatalog {
				switch ref.delta {
				case 1:
					fmt.Printf("Reference presented on left side only %s (%v)\n", n, left.References[ref.l])
					upFuncList = append(upFuncList, left.References[ref.l].CreateConstraint)
					downFuncList = append(downFuncList, left.References[ref.l].DeleteConstraint)
				case -1:
					fmt.Printf("Reference presented on right side only %s (%v)\n", n, right.References[ref.r])
					upFuncList = append(upFuncList, right.References[ref.r].DeleteConstraint)
					downFuncList = append(downFuncList, right.References[ref.r].CreateConstraint)
				case 0:
					fmt.Printf("Reference presented on both sides %s (%v)\n", n, right.Columns[ref.r])
					if state := isReferenceSchemaChanged(&left.References[ref.l], &right.References[ref.r]); state {
						upFuncList = append(upFuncList, right.References[ref.r].DeleteConstraint, left.References[ref.l].CreateConstraint)
						downFuncList = append(downFuncList, left.References[ref.l].DeleteConstraint, right.References[ref.r].CreateConstraint)
					}

				}
			}

		}
	}
	return upFuncList, downFuncList
}

func isColumnSchemaChanged(left *Column, right *Column) bool {
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

func isReferenceSchemaChanged(l *ReferenceMeta, r *ReferenceMeta) bool {
	if l.RefColumn != r.RefColumn || l.RefTable != r.RefTable || l.RefOptions != r.RefOptions {
		return true
	}
	return false
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
