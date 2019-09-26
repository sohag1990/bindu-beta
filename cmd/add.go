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
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bindu-bindu/bindu/cmd/helper"
	"github.com/spf13/cobra"
)

// bindu add swagger --skip // skipp to download dipendency lib
// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		helper.IsInProjectDir()

		// sanitize user args array
		args = helper.SanitizeUserInput(args)
		switch {
		case args[0] == "Swagger":
			s, _ := strconv.ParseBool(fmt.Sprintf("%v", cmd.Flag("skip").Value))
			// fmt.Println(s)
			if s == false {
				fmt.Println(exec.Command("go", "get", "-u", "github.com/swaggo/swag/cmd/swag").Run())
				fmt.Println(exec.Command("go", "get", "-u", "github.com/swaggo/gin-swagger").Run())
				fmt.Println(exec.Command("go", "get", "-u", "github.com/swaggo/files").Run())
			}
			// initialize swag if not exist
			fmt.Println(exec.Command("swag", "init").Run())

			// Edit main.go then read file and find find if exist swagger or initialize swagger
			theImportPath := helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH")
			// fmt.Println(theImportPath)
			lines, err := helper.ScanLines("./main.go")
			apiLines, err := helper.ScanLines("./routes/API.go")

			helper.ErrorCheck(err)
			// fmt.Println(lines)
			read, err := ioutil.ReadFile("./main.go")
			helper.ErrorCheck(err)
			// fmt.Println(lines)

			// fmt.Println(line)
			if !(strings.Contains(string(read), "/docs")) {

				for i, line := range lines {
					if strings.Contains(line, "import (") || strings.Contains(line, "import(") {
						var newLines []string
						newLines = append(newLines, lines[:i+1]...)
						newLines = append(newLines, "_ \""+theImportPath+"/docs\"")

						newLines = append(newLines, lines[i+1:]...)
						lines = newLines
						// fmt.Println(lines)
					}

				}
			}
			if !(strings.Contains(string(read), "swaggo/files")) {
				for i, line := range lines {
					if strings.Contains(line, "import (") || strings.Contains(line, "import(") {
						var newLines []string
						newLines = append(newLines, lines[:i+1]...)
						newLines = append(newLines, "swaggerFiles \"github.com/swaggo/files\"")

						newLines = append(newLines, lines[i+1:]...)
						lines = newLines
						// fmt.Println(lines)
					}

				}
			}
			if !(strings.Contains(string(read), "swaggo/gin-swagger")) {

				for i, line := range lines {
					if strings.Contains(line, "import (") || strings.Contains(line, "import(") {
						var newLines []string
						newLines = append(newLines, lines[:i+1]...)
						newLines = append(newLines, "ginSwagger \"github.com/swaggo/gin-swagger\"")

						newLines = append(newLines, lines[i+1:]...)
						lines = newLines
						// fmt.Println(lines)
					}

				}
			}
			if !(strings.Contains(string(read), "@title")) {

				for i, line := range lines {
					if strings.Contains(line, "main()") || strings.Contains(line, "func main()") {
						var newLines []string
						newLines = append(newLines, lines[:i+1]...)
						newLines = append(newLines, "\t// eg. http://localhost:"+helper.GetEnvValueByKey(".env", "APP_PORT")+"/swagger/index.html")
						newLines = append(newLines, "\t// @title "+helper.GetEnvValueByKey(".env", "APP_NAME")+" Swagger API")
						newLines = append(newLines, "\t// @version 1.0")
						newLines = append(newLines, "\t// @description bindu auto-generated swagger api documentation for "+helper.GetEnvValueByKey(".env", "APP_NAME")+" server.")
						newLines = append(newLines, "\t// @termsOfService http://swagger.io/terms/")
						newLines = append(newLines, "\t// @contact.name "+helper.GetEnvValueByKey(".env", "APP_NAME")+" API Support")
						newLines = append(newLines, "\t// @contact.url http://www.swagger.io/support")
						newLines = append(newLines, "\t// @contact.email support@swagger.io")
						newLines = append(newLines, "\t// @license.name Apache 2.0")
						newLines = append(newLines, "\t// @license.url http://www.apache.org/licenses/LICENSE-2.0.html")
						newLines = append(newLines, "\t// @host localhost:"+helper.GetEnvValueByKey(".env", "APP_PORT"))
						newLines = append(newLines, "\t// @BasePath /")

						newLines = append(newLines, lines[i+1:]...)
						lines = newLines
					}

				}
			}
			if !(strings.Contains(string(read), "ginSwagger.URL")) {

				for i, line := range lines {
					if strings.Contains(line, "gin.Default()") || strings.Contains(line, "gin.New()") {
						var newLines []string
						newLines = append(newLines, lines[:i+1]...)
						newLines = append(newLines, "\turl := ginSwagger.URL(\"http://localhost:"+helper.GetEnvValueByKey(".env", "APP_PORT")+"/swagger/doc.json\") // The url pointing to API definition")
						newLines = append(newLines, "\tr.GET(\"/swagger/*any\", ginSwagger.WrapHandler(swaggerFiles.Handler, url))")
						newLines = append(newLines, lines[i+1:]...)
						lines = newLines
						// fmt.Println(lines)
					}

				}
			}

			// fmt.Println(lines)
			ioutil.WriteFile("./main.go", []byte(strings.Join(lines, "\n")), 0644)

			// add autogen documentations for controllers according to the routes actions
			for _, apiLine := range apiLines {
				// filter only routes
				if strings.Contains(apiLine, "r.") {
					apiRouteSplit := strings.Split(apiLine, ".")
					routeMethodSplit := strings.Split(apiRouteSplit[1], "(")
					routeMethodName := routeMethodSplit[0]
					fmt.Println(routeMethodName)
					routeNameSplit := strings.Split(routeMethodSplit[1], ",")
					routeName := helper.TrimQuote(routeNameSplit[0])
					fmt.Println(routeName)
					apiControllerSplit := strings.Split(apiRouteSplit[len(apiRouteSplit)-1], ")")
					controllerName := apiControllerSplit[0]
					fmt.Println(controllerName)

					helper.ScanDirFindController("./app/controllers/", controllerName+"(c *gin.Context)", routeName, routeMethodName)

				}
			}

			// fmt.Println("called for update")
			// update swag according to the updated data
			fmt.Println(exec.Command("swag", "init").Run())

		default:
			fmt.Println("command not found!")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().BoolP("skip", "s", false, "Skip go get command")
}
