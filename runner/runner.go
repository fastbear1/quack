package runner

import (
	"context"

	utils "github.com/fastbear1/quack/internal"
	. "github.com/fastbear1/quack/schema"
)

var clog utils.Logger

func init() {
	clog = utils.GetLogger("default", nil)
}

// quacking migration file pipeline
func Run(ctx context.Context, conf *utils.ConfigYaml, fileName string) int {
	// init logger
	logLevel := utils.INFO
	if conf.Verbose {
		logLevel = utils.DEBUG
	}
	clog = utils.GetLogger("default", logLevel)

	// step 1: check connection to database
	var dbTablesMeta []TableMeta

	drv, err := GetDriver(ctx, conf)
	if err != nil {
		clog.Info("Error occured in proccess of getting database driver: %s", err)
		return 1
	}

	// step 1.1: get avaiable tables from database
	dbTables, err := drv.GetTablesList()
	if err != nil {
		clog.Info("Error occured while getting tables list: %s", err)
		return 1
	}
	for _, tableName := range dbTables {
		dbTablesMeta = append(dbTablesMeta,
			TableMeta{
				Name: tableName,
			},
		)
	}

	// step 1.2: get table columns, indices and constraint information
	for i := 0; i < len(dbTablesMeta); i++ {
		res, err := drv.GetTableColumnsMeta(dbTablesMeta[i].Name)
		if err != nil {
			clog.Info("Error occured while getting tables column list: %s", err)
			return 1
		}
		dbTablesMeta[i].Columns = res
		// step 1.3: get table indices information
		idx, err := drv.GetTableIndices(dbTablesMeta[i].Name)
		if err != nil {
			clog.Info("Error occured while getting indices list: %s", err)
			return 1
		}
		dbTablesMeta[i].Indeces = idx
		// step 1.4: get table references
		ref, err := drv.GetTableReferences(dbTablesMeta[i].Name)
		if err != nil {
			clog.Info("Error occured while getting constraint list: %s", err)
			return 1
		}
		dbTablesMeta[i].References = ref
	}

	// step 2: Scan models directory for gorm struct definitions
	var gormStructMeta []TableMeta

	StructRaw, err := Scan(conf)
	if err != nil {
		clog.Info("Error parsing directory with gorm models: %s", err)
		return 1
	}

	for i := 0; i < len(StructRaw); i++ {
		// step 2.1: get gorm tables metadata (column and indices)
		gsmeta := parseModelStruct(drv, StructRaw[i])
		// step 2.2: parse FK and embed structs
		for _, f := range StructRaw[i].ReferenceFields {
			reference := parseReferenceEmbedStructs(drv, StructRaw[i].Name, f.FieldType, f.FieldTag)
			gsmeta.References = append(gsmeta.References, reference)
		}
		gormStructMeta = append(gormStructMeta, gsmeta)
	}
	// step3: Compare current state of metadata for database tables and gorm structures
	funcList, downFuncList := compareMetaState(gormStructMeta, dbTablesMeta)
	var sqlUp, sqlDown []string

	for _, fup := range funcList {
		sqlUp = append(sqlUp, fup(drv))
	}
	for _, fdown := range downFuncList {
		sqlDown = append(sqlDown, fdown(drv))
	}

	// step 4: Write sql-Up and sql-Down commands to file
	if len(sqlUp) != 0 && len(sqlDown) != 0 {
		writeToFile(conf, fileName, sqlUp, sqlDown)
	} else {
		clog.Info("Gorm struct and DB tables already synchronized")
	}

	return 0
}
