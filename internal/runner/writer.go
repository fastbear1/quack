package runner

import (
	"bufio"
	"fmt"
	"os"
	"time"

	utils "github.com/fastbear1/quack/internal"
)

func writeToFile(conf *utils.ConfigYaml, fileName string, sqlUp []string, sqlDown []string) {
	path := fmt.Sprintf("./%s/%d_%s.sql", conf.Migrations.Path, time.Now().Unix(), fileName)
	f, err := os.Create(path)
	utils.CheckErrLite(err)
	w := bufio.NewWriter(f)
	_, err = w.WriteString("-- +goose Up\n")
	utils.CheckErrLite(err)
	for i := 0; i < len(sqlUp); i++ {
		_, err = w.WriteString(sqlUp[i] + "\n")
		utils.CheckErrLite(err)
	}
	_, err = w.WriteString("\n")
	_, err = w.WriteString("-- +goose Down\n")
	utils.CheckErrLite(err)
	for i := 0; i < len(sqlDown); i++ {
		_, err = w.WriteString(sqlDown[i] + "\n")
		utils.CheckErrLite(err)
	}
	w.Flush()
}
