package db

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	helper "github.com/bindu-bindu/bindu/Helper"
	"github.com/spf13/cobra"
)

// DbMigrate to migrate database
func DbMigrate(cmd *cobra.Command, cli helper.CommandChain) bool {
	args := cli.GetArgs()
	argsLen := len(args)
	switch {
	case args[0] == "Create":
		if argsLen > 1 {
			if helper.StringsContains(args[1], ":") {
				fmt.Println("Wrong argument provided!")
			} else {
				fmt.Println("Create command not ready yet...")
				// conName := args[1]
				// fmt.Println(conName)
				// Create database connection file accordingly user data
				// f, err := os.Create("./db/" + conName + ".go")
				// helper.ErrorCheck(err)
				// defer f.Close()
				return false
			}
		}
	case helper.StringsContains(args[0], "Migrate"):
		if helper.StringsContains(args[0], ":") {
			splitStr := strings.Split(args[0], ":")
			tableName := strings.Title(splitStr[1])
			// Migrate database connection file accordingly user data

			newpath := filepath.Join(".", "db/migrations")
			// create directory if not exists
			os.MkdirAll(newpath, os.ModePerm)
			y, m, d := time.Now().Date()
			nanoSeco := time.Now().Nanosecond()
			fileName := args[0] + strconv.Itoa(y) + m.String() + strconv.Itoa(d) + "-" + strconv.Itoa(nanoSeco) + ".go"
			fullPathFileName := filepath.Join(".", newpath+"/"+fileName)
			f, e := os.Create(fullPathFileName)
			helper.ErrorCheck(e)
			// fmt.Printf("%T", f)
			f.WriteString("package main\n\n")
			f.WriteString("import (\n")
			f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/app/models\"\n")
			f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/bindu\"\n")
			f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/db\"\n")
			f.WriteString(")\n\n")

			f.WriteString("func main() {\n")
			f.WriteString("\tbindu.Init()\n")
			f.WriteString("\tdb.Con()\n")
			f.WriteString("\t" + "db.DB.AutoMigrate(models." + tableName + "{})\n")
			f.WriteString("\tdefer db.DB.Close()\n")
			f.WriteString("}")
			defer f.Close()
			fmt.Println("migration file created: " + fullPathFileName)
			// Lets run the migration
			serverRunCMD := exec.Command("go", "run", "./"+fullPathFileName)
			err := serverRunCMD.Run()
			helper.ErrorCheck(err)
		} else {
			newpath := filepath.Join(".", "db/migrations")
			// create directory if not exists
			os.MkdirAll(newpath, os.ModePerm)
			y, m, d := time.Now().Date()
			nanoSeco := time.Now().Nanosecond()
			fileName := args[0] + strconv.Itoa(y) + m.String() + strconv.Itoa(d) + "-" + strconv.Itoa(nanoSeco) + ".go"
			fullPathFileName := filepath.Join(".", newpath+"/"+fileName)
			f, e := os.Create(fullPathFileName)
			helper.ErrorCheck(e)
			// fmt.Printf("%T", f)
			f.WriteString("package main\n\n")
			f.WriteString("import (\n")
			f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/app/models\"\n")
			f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/bindu\"\n")
			f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/db\"\n")
			f.WriteString(")\n\n")

			f.WriteString("func main() {\n")
			f.WriteString("\tbindu.Init()\n")
			f.WriteString("\tdb.Con()\n")
			helper.ScanDirFindLine("./app/", "struct", f)
			f.WriteString("\tdefer db.DB.Close()\n")
			f.WriteString("}")
			defer f.Close()
			fmt.Println("migration file created: " + fullPathFileName)
			// Lets run the migration
			serverRunCMD := exec.Command("go", "run", "./"+fullPathFileName)
			err := serverRunCMD.Run()
			if err != nil {
				fmt.Println("Database error, First check database connection, if db connection is ok then debug the migration file.")
				helper.ErrorCheck(err)
				return false
			} else {
				fmt.Println("DB Migration succesfull!")
			}

		}
	}
	return true
}
