/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bindu-bindu/bindu/cmd/helper"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// exmpl
		// bindu db create con1
		// bindu db migrate con1
		// bindu db migrate con1:tableName
		// bindu db rollBack con1
		// bindu db rollback con1:tablename
		// bindu db drop con1
		// bindu db drop con1:tablename
		helper.IsInProjectDir()
		args = helper.SanitizeUserInput(args)
		argsLen := len(args)
		switch {
		case args[0] == "Create":
			if argsLen > 1 {
				if strings.Contains(args[1], ":") {
					fmt.Println("Wrong argument provided!")
				} else {
					fmt.Println("Create command not ready yet...")
					// conName := args[1]
					// fmt.Println(conName)
					// Create database connection file accordingly user data
					// f, err := os.Create("./db/" + conName + ".go")
					// helper.ErrorCheck(err)
					// defer f.Close()
				}
			}
		case strings.Contains(args[0], "Migrate"):
			if strings.Contains(args[0], ":") {
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
				helper.ScanDirFindLine("./app/", "type", f)
				f.WriteString("\tdefer db.DB.Close()\n")
				f.WriteString("}")
				defer f.Close()
				fmt.Println("migration file created: " + fullPathFileName)
				// Lets run the migration
				serverRunCMD := exec.Command("go", "run", "./"+fullPathFileName)
				err := serverRunCMD.Run()
				helper.ErrorCheck(err)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// dbCmd.PersistentFlags().StringP("connection", "c", "", "Chose a connection name")
}
