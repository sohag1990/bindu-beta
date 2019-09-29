package new

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	helper "github.com/bindu-bindu/bindu/Helpers"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// New New project create proccess
func New(cmd *cobra.Command, args []string) {
	// fmt.Println(len(args))
	appName := ""
	if len(args) > 0 {
		appName = args[0]
	}
	createNewApp(helper.MakeAppName(appName))
}
func createNewApp(appName string) {
	// ENV VARIABLES START
	var envApp ENV_APP
	var envLog ENV_LOG
	var envDb ENV_DB
	// ENV STRUCTS END
	var UserInputs = make(map[string]string)
	fmt.Println("checking internet connection....")
	_, err := net.Dial("tcp", "github.com:443")
	helper.ErrorCheck(err)
	UserInputs["NetCheck"] = "OK"
	// project name
	if len(appName) == 0 {
		promptAppName := promptui.Prompt{
			Label:    "App Name",
			Validate: nil,
			Default:  "NewProject",
		}
		aName, err := promptAppName.Run()
		helper.ErrorCheck(err)
		envApp.APP_NAME = helper.MakeAppName(aName)
		envApp.APP_IMPORT_PATH = helper.AppImportPath(envApp.APP_NAME)

	} else {

		envApp.APP_NAME = appName
		envApp.APP_IMPORT_PATH = helper.AppImportPath(envApp.APP_NAME)
	}
	type preBuiltApp struct {
		Label    string
		Name     string
		UrlSSH   string
		UrlHTTPS string
	}

	preBuiltApps := []preBuiltApp{
		{Label: "Blank", Name: "bindu-blank", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-blank.git"},
		{Label: "Basic Web", Name: "bindu-basic-web", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-basic-web.git"},
		{Label: "Basic Api", Name: "bindu-basic-api", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-basic-api.git"},
		{Label: "Blog", Name: "bindu-blog", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-blog.git"},
		{Label: "E-Commerce", Name: "bindu-e-commerce", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-e-commerce.git"},
		{Label: "GRPC Server", Name: "bindu-grpc-server", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-grpc-server.git"},
		{Label: "GRPC Client", Name: "bindu-grpc-client", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-grpc-client.git"},
		{Label: "Download Third Party Project", Name: "download", UrlSSH: "", UrlHTTPS: ""},
	}

	var preBuiltAppLabels []string
	for _, app := range preBuiltApps {
		preBuiltAppLabels = append(preBuiltAppLabels, app.Label)
	}

	//prebuilt app selection
	promptPreBuiltApp := promptui.Select{
		Label: "Select pre-built project",
		Items: preBuiltAppLabels,
	}

	preSelectedIndex, result, err := promptPreBuiltApp.Run()

	envApp.APP_PREBUILT = helper.MakeAppName(result)
	helper.ErrorCheck(err)
	// the last item
	if preSelectedIndex == len(preBuiltApps)-1 {
		promptDownLoadLink := promptui.Prompt{
			Label:    preBuiltApps[7].Label,
			Validate: nil,
		}
		downLoadLink, err := promptDownLoadLink.Run()
		helper.ErrorCheck(err)
		UserInputs["APP_SELECTED"] = strconv.Itoa(preSelectedIndex)
		preBuiltApps[7].UrlHTTPS = downLoadLink

	} else {
		//Select Database
		promtDbAdapter := promptui.Select{
			Label: "Select Database",
			Items: []string{"Sqlite", "Mysql", "PGSql", "MongoDB", "Oracle", "None"},
		}
		_, dbAdapterName, err := promtDbAdapter.Run()

		envDb.DB_ADAPTER = strings.ToLower(dbAdapterName)
		if dbAdapterName == "None" {
			fmt.Println("I will set database manually")
		} else {
			// DB Host
			promptDbHost := promptui.Prompt{
				Label:    "Database Host",
				Validate: nil,
				Default:  "localhost",
			}
			dbHost, err := promptDbHost.Run()
			helper.ErrorCheck(err)
			envDb.DB_HOST = dbHost

			// DB Port
			promptDbPort := promptui.Prompt{
				Label:    "Database PORT",
				Validate: nil,
				Default:  "3306",
			}
			dbPort, err := promptDbPort.Run()
			helper.ErrorCheck(err)
			envDb.DB_PORT = dbPort
			// DB Name
			promptDbName := promptui.Prompt{
				Label:    "Database Name",
				Validate: nil,
				Default:  appName,
			}
			dbName, err := promptDbName.Run()
			helper.ErrorCheck(err)
			envDb.DB_DATABASE = dbName
			// DB User Name
			promptDbUserName := promptui.Prompt{
				Label:    "Database User Name",
				Validate: nil,
				Default:  "root",
			}
			dbUserName, err := promptDbUserName.Run()
			helper.ErrorCheck(err)
			envDb.DB_USERNAME = dbUserName
			// DB Password
			promptDbPass := promptui.Prompt{
				Label:    "Database Password",
				Validate: nil,
			}
			dbPass, err := promptDbPass.Run()
			helper.ErrorCheck(err)
			envDb.DB_PASSWORD = dbPass
		}

		// Run Project On Port
		promptProjectPort := promptui.Prompt{
			Label:    "Run project on port",
			Validate: nil,
			Default:  "8080",
		}
		projectPort, err := promptProjectPort.Run()
		helper.ErrorCheck(err)
		pp, _ := strconv.Atoi(projectPort)
		envApp.APP_PORT = pp
	}

	fmt.Println("Creating new project hold tight.")
	appSelectedIndex, _ := strconv.Atoi(UserInputs["APP_SELECTED"])
	appSelected := preBuiltApps[appSelectedIndex]
	// rename if core project exist
	os.Rename(appSelected.Name, envApp.APP_NAME+time.Now().String())
	// rename if old project in same name
	os.Rename(envApp.APP_NAME, envApp.APP_NAME+time.Now().String())
	cmd := exec.Command("git", "clone", appSelected.UrlHTTPS)
	errD := cmd.Run()
	// fmt.Println(errD)
	helper.ErrorCheck(errD)

	// get the old import path from .env file
	oldImportPath := helper.GetEnvValueByKey(helper.PWD()+"/"+appSelected.Name+"/.env", "APP_IMPORT_PATH")
	// rename after download the core project
	errRename := os.Rename(appSelected.Name, envApp.APP_NAME)
	helper.ErrorCheck(errRename)

	//cd to project file
	helper.CD(envApp.APP_NAME)
	// fmt.Println(oldImportPath)
	// Migrate import path. must after CD inside the project directory
	helper.FixImportPath(oldImportPath, envApp.APP_IMPORT_PATH)

	// Create .env file accordingly user data
	f, err := os.Create(".env")
	helper.ErrorCheck(err)

	f.WriteString("APP_NAME    		=" + envApp.APP_NAME + "\n")
	f.WriteString("APP_IMPORT_PATH	=" + envApp.APP_IMPORT_PATH + "\n")
	f.WriteString("APP_ENV     		=" + fmt.Sprintf("%v", helper.IfThenElse(envApp.APP_DEBUG, "Dev", "Prod")) + "\n")
	f.WriteString("APP_KEY     		=" + envApp.APP_KEY + "\n")
	f.WriteString("APP_DEBUG   		=" + envApp.APP_ENV + "\n")
	f.WriteString("APP_URL     		=" + envApp.APP_URL + "\n\n\n")
	f.WriteString("APP_PORT    		=" + strconv.Itoa(envApp.APP_PORT) + "\n")
	f.WriteString("APP_PREBUILT		=" + envApp.APP_PREBUILT + "\n")

	f.WriteString("LOG_CHANNEL=" + envLog.LOG_CHANNEL + "\n\n\n")

	f.WriteString("DB_ADAPTER   =" + envDb.DB_ADAPTER + "\n")
	f.WriteString("DB_HOST      =" + envDb.DB_HOST + "\n")
	f.WriteString("DB_PORT      =" + envDb.DB_PORT + "\n")
	f.WriteString("DB_DATABASE  =" + envDb.DB_DATABASE + "\n")
	f.WriteString("DB_USERNAME  =" + envDb.DB_USERNAME + "\n")
	f.WriteString("DB_PASSWORD  =" + envDb.DB_PASSWORD + "\n\n\n")

	defer f.Close()
	fmt.Println("Done!")
	// i := 1
	// for k, v := range UserInputs {
	// 	fmt.Println(i, ") ", k, ":", v)
	// 	i++
	// }

}

type ENV_APP struct {
	APP_NAME        string
	APP_IMPORT_PATH string
	APP_ENV         string
	APP_KEY         string
	APP_DEBUG       bool
	APP_URL         string
	APP_PORT        int
	APP_PREBUILT    string
}

type ENV_LOG struct {
	LOG_CHANNEL string
}

type ENV_DB struct {
	DB_ADAPTER  string
	DB_HOST     string
	DB_PORT     string
	DB_DATABASE string
	DB_USERNAME string
	DB_PASSWORD string
}
