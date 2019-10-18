package add

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	helper "github.com/bindu-bindu/bindu/Helper"
	"github.com/spf13/cobra"
)

// AddSwagger to add swagger documentation, returns bool
func AddSwagger(cmd *cobra.Command, cli helper.CommandChain) bool {

	// Get sanitize user args array
	args := cli.GetArgs()
	switch {
	case args[0] == "Swagger":
		// to skip download lib from web
		s, _ := strconv.ParseBool(fmt.Sprintf("%v", cmd.Flag("skip").Value))
		// fmt.Println(s)
		if s == false {
			fmt.Println("Downloading.... swag")
			fmt.Println(exec.Command("go", "get", "-u", "github.com/swaggo/swag/cmd/swag").Run())
			fmt.Println("github.com/swaggo/swag/cmd/swag  ---done")

			fmt.Println("Downloading.... gin-swagger")
			fmt.Println(exec.Command("go", "get", "-u", "github.com/swaggo/gin-swagger").Run())
			fmt.Println("github.com/swaggo/gin-swagger  ---done")

			fmt.Println("Downloading.... Docs")
			fmt.Println(exec.Command("go", "get", "-u", "github.com/swaggo/files").Run())
			fmt.Println("github.com/swaggo/files  ---done")
		}
		// initialize swag if not exist
		fmt.Println("Initializing Swag... ---done")
		fmt.Println(exec.Command("swag", "init").Run())

		// Edit main.go then read file and find find if exist swagger or initialize swagger
		theImportPath := helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH")
		// fmt.Println(theImportPath)
		lines, err := helper.ScanLines("./main.go")

		fmt.Println("Reading Routes API... ---done")
		apiLines, err := helper.ScanLines("./routes/API.go")

		helper.ErrorCheck(err)

		fmt.Println("Implementing swagger.....")
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
		fmt.Println("Congretulation your application successfully documented by swagger!!!")
	default:
		fmt.Println("command not found!")
		return false
	}
	return true
}
