package drivers

import (
	"fmt"
	"regexp"
	"strings"
)

type IndexOption struct {
	Field      string
	Expression string
	Sort       string
	Collate    string
	Priority   int
}

type Index struct {
	Name   string
	Unique bool
	Type   string
	Where  string
	Option string
	Fields []IndexOption
}

func ParseDatabaseIndices(indexdef string) interface{} {

	var test = []string{
		"CREATE INDEX idx_auth_users_id ON public.auth_users USING btree (id) NULLS DISTINCT",
		"CREATE INDEX idx_auth_users_password ON public.auth_users USING btree (password) WITH (fillfactor = 70)",
		"CREATE INDEX idx_auth_users_password_partial_id ON public.auth_users USING btree (password) WHERE (NOT ((created_at)::date > '2026-01-01'::date))",
		"CREATE UNIQUE INDEX idx_auth_users_password_user_id ON public.auth_users USING btree (password DESC NULLS FIRST) INCLUDE (user_id)",
		`CREATE INDEX idx_auth_users_password_collation ON public.auth_users USING btree (password COLLATE "de-DE-x-icu") WITH (fillfactor = 70)`,
		"CREATE UNIQUE INDEX idx_auth_users_password_user_id_2 ON public.auth_users USING btree (password, user_id) NULLS NOT DISTINCT",
		"CREATE INDEX idx_auth_users_password_lower ON public.auth_users USING btree (upper(password))",
	}

	var IndexMeta = []Index{}

	main_part := `CREATE(?<Unique> UNIQUE)? INDEX (?<Name>\w+) ON (?:.+) USING (?<Type>\w+)`
	field_part := ` \((?<Exp>\w+\(\w+\))?(?<Fieldlist>[a-z_\s,]+)?(?<Collate>COLLATE [\w\d"-]+)?(?<Sort>[A-Z\s]+)?\)`
	partial_part := `(?<Where> WHERE \([\w\d\s<>=:'"\-\(\)_]+\))?(?<Include> INCLUDE \(\w+\))?(?<With> WITH \([\w\d\s=]+\))?(?<Distinct> NULLS DISTINCT| NULLS NOT DISTINCT)?`

	r := regexp.MustCompile(main_part + field_part + partial_part)

	type paramsMap map[string]string
	var idxraw = []paramsMap{}

	for _, tt := range test {
		pm := make(paramsMap)
		match := r.FindStringSubmatch(tt)
		for i, name := range r.SubexpNames() {
			if i > 0 && i <= len(match) {
				pm[name] = match[i]
			}
		}
		idxraw = append(idxraw, pm)
	}

	for _, idx := range idxraw {
		b := true
		if idx["Unique"] == "" {
			b = false
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

		IndexMeta = append(IndexMeta, Index{
			Name:   idx["Name"],
			Unique: b,
			Type:   idx["Type"],
			Where:  idx["Where"],
			Option: idx["Include"] + idx["With"] + idx["Distinct"],
			Fields: idxopts,
		})
	}
	for _, im := range IndexMeta {
		fmt.Println(im)
	}
	return IndexMeta
}
