package drivers

import (
	"regexp"
	"strings"
)

func ParseDatabaseIndices(indexdef string) (IndexMeta, error) {
	main_part := `CREATE(?<Unique> UNIQUE)? INDEX (?<Name>\w+) ON (?:.+) USING (?<Type>\w+)`
	field_part := ` \((?<Exp>\w+\(\w+\))?(?<Fieldlist>[a-z_\s,]+)?(?<Collate>COLLATE [\w\d"-]+)?(?<Sort>[A-Z\s]+)?\)`
	partial_part := `(?<Where> WHERE \([\w\d\s<>=:'"\-\(\)_]+\))?(?<Include> INCLUDE \(\w+\))?(?<With> WITH \([\w\d\s=]+\))?(?<Distinct> NULLS DISTINCT| NULLS NOT DISTINCT)?`

	r := regexp.MustCompile(main_part + field_part + partial_part)

	type paramsMap map[string]string

	idx := make(paramsMap, 0)
	match := r.FindStringSubmatch(indexdef)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			idx[name] = match[i]
		}
	}

	uniq := true
	if idx["Unique"] == "" {
		uniq = false
	}

	var fields = []string{}
	if idx["Fieldlist"] != "" {
		for f := range strings.SplitSeq(idx["Fieldlist"], ", ") {
			fields = append(fields, f)
		}
	}
	if idx["Fieldlist"] == "" && idx["Exp"] != "" {
		rexp := regexp.MustCompile(`\w+\((?<Field>\w+)\)`)
		match := rexp.FindStringSubmatch(idx["Exp"])
		if len(match) > 1 {
			fields = append(fields, match[1])
		}
	}

	var idxopts = []IndexOption{}
	for i, fd := range fields {
		var field = IndexOption{
			Field:      fd,
			Expression: idx["Exp"],
			Sort:       idx["Sort"],
			Collate:    idx["Collate"],
			Priority:   i,
		}
		idxopts = append(idxopts, field)
	}

	parsable := true
	if len(fields) == 0 {
		parsable = false
	}

	var IndexMeta = IndexMeta{
		Name:    idx["Name"],
		Unique:  uniq,
		Type:    idx["Type"],
		Where:   idx["Where"],
		Option:  idx["Include"] + idx["With"] + idx["Distinct"],
		Parsed:  parsable,
		Columns: idxopts,
	}
	return IndexMeta, nil
}
