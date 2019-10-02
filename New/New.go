package new

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	helper "github.com/bindu-bindu/bindu/Helper"
	"github.com/spf13/cobra"
)

// New project create proccess
func New(cmd *cobra.Command, c helper.CommandChain) {
	if !helper.NetCheck() {
		return
	}

	// ENV VARIABLES to catch user input data
	var envApp helper.ENV_APP
	var envLog helper.ENV_LOG
	var envDb helper.ENV_DB
	var selectedIndex = 0
	// Set value from CLI after read flags command
	for _, flg := range c.GetFlags() {
		if flg.Key == "app" && len(flg.Values[0]) > 0 {
			i, f := FindPrebuitAppIndex(flg.Values[0])
			if f {
				// if download link
				if i == len(PreBuiltApps)-1 {
					envApp.APP_PREBUILT = PreBuiltApps[i].Label
					envApp.APP_URL = flg.Values[0]
					selectedIndex = i
				} else {
					envApp.APP_PREBUILT = PreBuiltApps[i].Label
					envApp.APP_URL = PreBuiltApps[i].UrlHTTPS
					selectedIndex = i
				}
			}
		}
		if flg.Key == "db" && len(flg.Values[0]) > 0 {
			// fmt.Println("has db value")
			dbStr := strings.Split(flg.Values[0], ":")
			for i, d := range dbStr {
				switch i {
				case 0:
					envDb.DB_ADAPTER = d
				case 1:
					envDb.DB_HOST = d
				case 2:
					envDb.DB_PORT = d
				case 3:
					envDb.DB_DATABASE = d
				case 4:
					envDb.DB_USERNAME = d
				case 5:
					envDb.DB_PASSWORD = d
				}

			}

		}
		if flg.Key == "port" && len(flg.Values[0]) > 0 {
			fmt.Println("has port value")
			pp, _ := strconv.Atoi(flg.Values[0])
			envApp.APP_PORT = pp
		}
	}

	// Promt to cli to get input from user show options that available
	args := c.GetArgs()
	if len(args) > 0 {
		envApp.APP_NAME = args[0]
		envApp.APP_IMPORT_PATH = helper.AppImportPath(envApp.APP_NAME)
	}

	// var UserInputs = make(map[string]string)

	// if appName empty then as in terminal for app Name
	if len(envApp.APP_NAME) == 0 {
		askAppName := helper.AskString("App Name", "NewProject")
		envApp.APP_NAME = helper.MakeAppName(askAppName)
		// Make Import path for this project according to the GO IMPORT rules
		envApp.APP_IMPORT_PATH = helper.AppImportPath(envApp.APP_NAME)
	}
	var preBuiltAppLabels []string
	for _, app := range PreBuiltApps {
		preBuiltAppLabels = append(preBuiltAppLabels, app.Label)
	}
	//Ask user to select a prebuilt app
	// If db adapter is wrong then ask everything again
	_, f := helper.ArrayFind(DbAdapters, envDb.DB_ADAPTER)
	if !f {
		fmt.Println(envDb.DB_ADAPTER + " is not a DB adapter, Select again... ")
	}
	if len(envDb.DB_ADAPTER) == 0 || !f {
		selectedIndex, result := helper.AskSelect("Select a app", preBuiltAppLabels)
		envApp.APP_PREBUILT = helper.MakeAppName(result)
		// if user want to download app from server, the last item (download from remote)
		// ask for remote source url
		// Notes: prebuilt app has built in env data
		if selectedIndex == len(PreBuiltApps)-1 {
			PreBuiltApps[selectedIndex].UrlHTTPS = helper.AskString(PreBuiltApps[selectedIndex].Label, "eg. https://...")

		} else {
			//Get User Input
			//Ask to Select Database Adapter

			_, dbAdapterName := helper.AskSelect("Select Database", DbAdapters)
			envDb.DB_ADAPTER = strings.ToLower(dbAdapterName)
			if dbAdapterName == "None" {
				fmt.Println("I will set database manually")
			} else {
				// Ask DB Host
				envDb.DB_HOST = helper.AskString("Database Host", "localhost")
				// Ask DB Port
				envDb.DB_PORT = helper.AskString("Database PORT", "3306")
				// Ask DB Name
				envDb.DB_DATABASE = helper.AskString("Database Name", envApp.APP_NAME)
				// Ask DB User Name
				envDb.DB_USERNAME = helper.AskString("Database Username", "root")
				// Ask DB Password
				envDb.DB_PASSWORD = helper.AskString("Database Password", "")
				// Ask Run Project On Port
				pp, _ := strconv.Atoi(helper.AskString("Run Project On PORT", strconv.Itoa(envApp.APP_PORT)))
				envApp.APP_PORT = pp
			}

		}
	}
	fmt.Println("Creating new project hold tight.......")

	appSelected := PreBuiltApps[selectedIndex]
	// rename if core project exist
	if os.Rename(appSelected.Name, envApp.APP_NAME+time.Now().String()) == nil {
		fmt.Println(appSelected.Name+" name already exsits, renamed to ", envApp.APP_NAME+time.Now().String())
	}

	// rename if old project in same name
	if os.Rename(envApp.APP_NAME, envApp.APP_NAME+time.Now().String()) == nil {
		fmt.Println(appSelected.Name+" name already exsits, renamed to ", envApp.APP_NAME+time.Now().String())
	}
	// Download project from remote
	cmdDownload := exec.Command("git", "clone", appSelected.UrlHTTPS)
	dError := cmdDownload.Run()
	helper.ErrorCheck(dError)
	fmt.Println("Project downloaded from ", appSelected.UrlHTTPS)

	// To fix import path, get the old import path from .env file
	oldImportPath := helper.GetEnvValueByKey(helper.PWD()+"/"+appSelected.Name+"/.env", "APP_IMPORT_PATH")
	// fmt.Println("Old Import Path ", oldImportPath)
	// Rename after download the core project
	errRename := os.Rename(appSelected.Name, envApp.APP_NAME)
	helper.ErrorCheck(errRename)
	// fmt.Println(appSelected.Name + " Renamed to... " + envApp.APP_NAME)
	//cd to project file
	helper.CD(envApp.APP_NAME)
	// fmt.Println("Enter into the new project " + envApp.APP_NAME)
	// Migrate import path. must after CD inside the project directory
	helper.FixImportPath(oldImportPath, envApp.APP_IMPORT_PATH)
	// Createing env file according to the user data
	writtingEnvFileForNewProject(envApp, envDb, envLog)
	fmt.Println("Env file created")
	fmt.Println("Congretulations!!! your new project successfully created.")

}

// Writting env file for new project
func writtingEnvFileForNewProject(envApp helper.ENV_APP, envDb helper.ENV_DB, envLog helper.ENV_LOG) {
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
	f.WriteString("DB_DATABASE  =" + strings.ToLower(envDb.DB_DATABASE) + "\n")
	f.WriteString("DB_USERNAME  =" + envDb.DB_USERNAME + "\n")
	f.WriteString("DB_PASSWORD  =" + envDb.DB_PASSWORD + "\n\n\n")

	defer f.Close()
}
