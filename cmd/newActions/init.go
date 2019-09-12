package newActions

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/bindu-bindu/bindu/cmd/helper"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Init New project create proccess
func Init(cmd *cobra.Command, args []string) {
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
	var envDbMysql ENV_DB_MYSQL
	var envDbSql ENV_DB_SQL
	var envDbMongo ENV_DB_MONGO
	var envDbPgsql ENV_DB_PGSQL
	var envDbOracle ENV_DB_ORACLE
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
		UserInputs["APP_NAME"] = helper.MakeAppName(aName)
		envApp.APP_NAME = helper.MakeAppName(aName)
		envApp.APP_IMPORT_PATH = helper.AppImportPath(envApp.APP_NAME)
	} else {
		UserInputs["APP_NAME"] = appName
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
	UserInputs["PRE_BUILT_APP"] = helper.MakeAppName(result)
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
		UserInputs["DATABASE_ADAPTER"] = helper.MakeAppName(dbAdapterName)
		envDbMysql.MYSQL_DB_CONNECTION = true
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
			UserInputs["DATABASE_HOSTS"] = dbHost
			envDbMysql.MYSQL_DB_HOST = dbHost

			// DB Port
			promptDbPort := promptui.Prompt{
				Label:    "Database PORT",
				Validate: nil,
				Default:  "3306",
			}
			dbPort, err := promptDbPort.Run()
			helper.ErrorCheck(err)
			UserInputs["DATABASE_PORT"] = dbPort
			P, _ := strconv.Atoi(dbPort)
			envDbMysql.MYSQL_DB_PORT = P
			// DB Name
			promptDbName := promptui.Prompt{
				Label:    "Database Name",
				Validate: nil,
				Default:  UserInputs["APP_NAME"],
			}
			dbName, err := promptDbName.Run()
			helper.ErrorCheck(err)
			UserInputs["DATABASE_NAME"] = dbName
			envDbMysql.MYSQL_DB_DATABASE = dbName
			// DB User Name
			promptDbUserName := promptui.Prompt{
				Label:    "Database User Name",
				Validate: nil,
				Default:  "root",
			}
			dbUserName, err := promptDbUserName.Run()
			helper.ErrorCheck(err)
			UserInputs["DATABASE_USERNAME"] = dbUserName
			envDbMysql.MYSQL_DB_USERNAME = dbUserName
			// DB Password
			promptDbPass := promptui.Prompt{
				Label:    "Database Password",
				Validate: nil,
			}
			dbPass, err := promptDbPass.Run()
			helper.ErrorCheck(err)
			UserInputs["DATABASE_PASSWORD"] = dbPass
			envDbMysql.MYSQL_DB_PASSWORD = dbPass
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
		UserInputs["APP_PORT"] = projectPort
		envApp.APP_PORT = pp
	}

	fmt.Println("Creating new project hold tight.")
	appSelectedIndex, _ := strconv.Atoi(UserInputs["APP_SELECTED"])
	appSelected := preBuiltApps[appSelectedIndex]
	// rename if core project exist
	os.Rename(appSelected.Name, UserInputs["APP_NAME"]+time.Now().String())
	// rename if old project in same name
	os.Rename(UserInputs["APP_NAME"], UserInputs["APP_NAME"]+time.Now().String())
	cmd := exec.Command("git", "clone", appSelected.UrlHTTPS)
	errD := cmd.Run()
	// fmt.Println(errD)
	helper.ErrorCheck(errD)

	// get the old import path from .env file
	oldImportPath := helper.GetEnvValueByKey(helper.PWD()+"/"+appSelected.Name+"/.env", "APP_IMPORT_PATH")
	// rename after download the core project
	errRename := os.Rename(appSelected.Name, UserInputs["APP_NAME"])
	helper.ErrorCheck(errRename)

	//cd to project file
	helper.CD(UserInputs["APP_NAME"])
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

	f.WriteString("MYSQL_DB_CONNECTION=" + strconv.FormatBool(envDbMysql.MYSQL_DB_CONNECTION) + "\n")
	f.WriteString("MYSQL_DB_HOST      =" + envDbMysql.MYSQL_DB_DATABASE + "\n")
	f.WriteString("MYSQL_DB_PORT      =" + envDbMysql.MYSQL_DB_HOST + "\n")
	f.WriteString("MYSQL_DB_DATABASE  =" + envDbMysql.MYSQL_DB_PASSWORD + "\n")
	f.WriteString("MYSQL_DB_USERNAME  =" + strconv.Itoa(envDbMysql.MYSQL_DB_PORT) + "\n")
	f.WriteString("MYSQL_DB_PASSWORD  =" + envDbMysql.MYSQL_DB_USERNAME + "\n\n\n")

	f.WriteString("MONGO_DB_CONNECTION=" + strconv.FormatBool(envDbMysql.MYSQL_DB_CONNECTION) + "\n")
	f.WriteString("MONGO_DB_HOST      =" + envDbMongo.MONGO_DB_DATABASE + "\n")
	f.WriteString("MONGO_DB_PORT      =" + envDbMongo.MONGO_DB_HOST + "\n")
	f.WriteString("MONGO_DB_DATABASE  =" + envDbMongo.MONGO_DB_PASSWORD + "\n")
	f.WriteString("MONGO_DB_USERNAME  =" + strconv.Itoa(envDbMongo.MONGO_DB_PORT) + "\n")
	f.WriteString("MONGO_DB_PASSWORD  =" + envDbMongo.MONGO_DB_USERNAME + "\n\n\n")

	f.WriteString("SQL_DB_CONNECTION=" + strconv.FormatBool(envDbSql.SQL_DB_CONNECTION) + "\n")
	f.WriteString("SQL_DB_HOST      =" + envDbSql.SQL_DB_DATABASE + "\n")
	f.WriteString("SQL_DB_PORT      =" + envDbSql.SQL_DB_HOST + "\n")
	f.WriteString("SQL_DB_DATABASE  =" + envDbSql.SQL_DB_PASSWORD + "\n")
	f.WriteString("SQL_DB_USERNAME  =" + strconv.Itoa(envDbSql.SQL_DB_PORT) + "\n")
	f.WriteString("SQL_DB_PASSWORD  =" + envDbSql.SQL_DB_USERNAME + "\n\n\n")

	f.WriteString("PGSQL_DB_CONNECTION=" + strconv.FormatBool(envDbPgsql.PGSQL_DB_CONNECTION) + "\n")
	f.WriteString("PGSQL_DB_HOST      =" + envDbPgsql.PGSQL_DB_DATABASE + "\n")
	f.WriteString("PGSQL_DB_PORT      =" + envDbPgsql.PGSQL_DB_HOST + "\n")
	f.WriteString("PGSQL_DB_DATABASE  =" + envDbPgsql.PGSQL_DB_PASSWORD + "\n")
	f.WriteString("PGSQL_DB_USERNAME  =" + strconv.Itoa(envDbPgsql.PGSQL_DB_PORT) + "\n")
	f.WriteString("PGSQL_DB_PASSWORD  =" + envDbPgsql.PGSQL_DB_USERNAME + "\n\n\n")

	f.WriteString("ORACLE_DB_CONNECTION=" + strconv.FormatBool(envDbOracle.ORACLE_DB_CONNECTION) + "\n")
	f.WriteString("ORACLE_DB_HOST      =" + envDbOracle.ORACLE_DB_DATABASE + "\n")
	f.WriteString("ORACLE_DB_PORT      =" + envDbOracle.ORACLE_DB_HOST + "\n")
	f.WriteString("ORACLE_DB_DATABASE  =" + envDbOracle.ORACLE_DB_PASSWORD + "\n")
	f.WriteString("ORACLE_DB_USERNAME  =" + strconv.Itoa(envDbOracle.ORACLE_DB_PORT) + "\n")
	f.WriteString("ORACLE_DB_PASSWORD  =" + envDbOracle.ORACLE_DB_USERNAME + "\n\n\n")
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

type ENV_DB_MYSQL struct {
	MYSQL_DB_CONNECTION bool
	MYSQL_DB_HOST       string
	MYSQL_DB_PORT       int
	MYSQL_DB_DATABASE   string
	MYSQL_DB_USERNAME   string
	MYSQL_DB_PASSWORD   string
}

type ENV_DB_SQL struct {
	SQL_DB_CONNECTION bool
	SQL_DB_HOST       string
	SQL_DB_PORT       int
	SQL_DB_DATABASE   string
	SQL_DB_USERNAME   string
	SQL_DB_PASSWORD   string
}

type ENV_DB_MONGO struct {
	MONGO_DB_CONNECTION bool
	MONGO_DB_HOST       string
	MONGO_DB_PORT       int
	MONGO_DB_DATABASE   string
	MONGO_DB_USERNAME   string
	MONGO_DB_PASSWORD   string
}

type ENV_DB_PGSQL struct {
	PGSQL_DB_CONNECTION bool
	PGSQL_DB_HOST       string
	PGSQL_DB_PORT       int
	PGSQL_DB_DATABASE   string
	PGSQL_DB_USERNAME   string
	PGSQL_DB_PASSWORD   string
}

type ENV_DB_ORACLE struct {
	ORACLE_DB_CONNECTION bool
	ORACLE_DB_HOST       string
	ORACLE_DB_PORT       int
	ORACLE_DB_DATABASE   string
	ORACLE_DB_USERNAME   string
	ORACLE_DB_PASSWORD   string
}
