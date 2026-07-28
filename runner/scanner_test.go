package runner

import (
	"os"
	"path"
	"testing"

	utils "github.com/fastbear1/quack/internal"
	"github.com/stretchr/testify/assert"
)

func TestScanFunction(t *testing.T) {
	var conf utils.ConfigYaml
	conf.ReadConfig()

	wd, err := os.Getwd()
	assert.Nil(t, err)

	testDir := path.Join(wd, conf.Models.Path.String())
	err = os.Mkdir(testDir, 0777)
	assert.Nil(t, err)
	defer utils.CleanUpDir(testDir)

	src := `package test` + "\n" +
		`import (` + "\n" +
		`	"time"` + "\n" +
		`	uuid "github.com/satori/go.uuid"` + "\n" +
		`)` + "\n\n" +

		`type TestBase struct {` + "\n" +
		`	ID uuid.UUID` + "`" + `gorm:"index;type:uuid;primary_key;default:gen_random_uuid()"` + "`\n" +
		// Test comment
		`	CreatedAt time.Time` + "`" + `gorm:"type:timestamp;not null;default:now();<-:create"` + "`\n" +
		`	UpdatedAt time.Time` + "`" + `gorm:"type:timestamp;not null;default:now()"` + "`\n" + `}`

	file, err := os.Create(path.Join(testDir, "test_model.go"))
	assert.Nil(t, err)
	defer file.Close()

	_, err = file.WriteString(src)
	assert.Nil(t, err)

	res, err := Scan(&conf)

	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
}

func TestScanFunctionResults(t *testing.T) {
	// TODO: too big test case
	var conf utils.ConfigYaml
	conf.ReadConfig()

	wd, err := os.Getwd()
	assert.Nil(t, err)

	testDir := path.Join(wd, conf.Models.Path.String())
	err = os.Mkdir(testDir, 0777)
	assert.Nil(t, err)
	defer utils.CleanUpDir(testDir)

	src1 := `package test` + "\n" +
		`import (` + "\n" +
		`	"time"` + "\n" +
		`	uuid "github.com/satori/go.uuid"` + "\n" +
		`)` + "\n\n" +

		`type Base struct {` + "\n" +
		`	ID uuid.UUID` + "`" + `gorm:"index;type:uuid;primary_key;default:gen_random_uuid()"` + "`\n" +
		// Test comment
		`	CreatedAt time.Time` + "`" + `gorm:"type:timestamp;not null;default:now();<-:create"` + "`\n" +
		`	UpdatedAt time.Time` + "`" + `gorm:"type:timestamp;not null;default:now()"` + "`\n" + `}`

	src2 := `package test` + "\n" +
		`import (` + "\n" +
		`	"time"` + "\n" +
		`	uuid "github.com/satori/go.uuid"` + "\n" +
		`)` + "\n\n" +
		`type TestUsers struct {` + "\n" +
		`	Base` + "\n" +
		`	UserID   uuid.UUID` + "`" + `gorm:"type:uuid;not null"` + "`\n" +
		`	Password string` + "`" + `gorm:"type:text;not null"` + "`\n" +
		`	// Constraints` + "\n" +
		`	Users Users` + "`" + `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:SET NULL;"` + "`\n" + `}`

	src3 := `type Users struct {` + "\n" +
		`	Base` + "\n" +
		`	Name   string` + "`" + `gorm:"not null"` + "`" + `// full name field` + "\n" +
		`	Email  string` + "`" + `gorm:"not null"` + "`" + `// User email field` + "\n" +
		`	Status string` + "`" + `gorm:"default:active"` + "`" + `// user statuses` + "\n"

	file, err := os.Create(path.Join(testDir, "test_model.go"))
	assert.Nil(t, err)
	defer file.Close()

	fileEmbed, err := os.Create(path.Join(testDir, "test_model_embed.go"))
	assert.Nil(t, err)
	defer fileEmbed.Close()

	_, err = file.WriteString(src2)
	assert.Nil(t, err)

	_, err = fileEmbed.WriteString(src1)
	assert.Nil(t, err)

	res, err := Scan(&conf)

	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)

	model := res[0]
	assert.Equal(t, model.Name, "TestUsers")
	assert.Equal(t, len(model.Fields), 6)
	assert.Equal(t, model.Fields[0].FieldName, "ID")
	assert.Equal(t, model.Fields[3].FieldName, "UserID")
	assert.Equal(t, model.EmbedFields, []string{"Base"})
	// Reference struct not initalized, so there is no constraints
	assert.Equal(t, len(model.ReferenceFields), 0)

	// add definition of reference struct
	_, err = file.WriteString("\n" + src3)
	assert.Nil(t, err)

	res, err = Scan(&conf)
	assert.Nil(t, err)

	model = res[0]
	assert.Equal(t, model.Name, "TestUsers")
	assert.Equal(t, len(model.Fields), 5)
	assert.Equal(t, model.Fields[0].FieldName, "ID")
	assert.Equal(t, model.Fields[3].FieldName, "UserID")
	assert.Equal(t, model.EmbedFields, []string{"Base"})
	assert.Equal(t, len(model.ReferenceFields), 1)
	assert.Equal(t, model.ReferenceFields[0].FieldName, "Users")

}
