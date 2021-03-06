package modify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	helper "github.com/bindu-bindu/bindu/Helper"
	story "github.com/bindu-bindu/bindu/Story"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
)

// Modifier is a func to Modifier model controller and scaffold
func Modifier(cmd *cobra.Command, cli helper.CommandChain) bool {
	// To return a status true or false, whether code executed or not
	// Deafult true
	status := true

	// User command args
	args := cli.GetArgs()

	// flags := c.GetFlags()
	genItems := []string{"Model"}
	// if no agruments Prompt to select action name eg. model
	if len(args) == 0 {
		fmt.Println("Arguments not enough to generate Model, Controller, Scaffold, Routes, View or Auth")
		story.UpdateThisStoryStatus("false: No enough agruments for generate command")
		return false
		//user simulation fix letter
		_, actionName := helper.AskSelect("Select To Generate", genItems)
		args = append(args, actionName)
		cli.SetCliArgs(args)
	}

	// check the first items is matched with predefined arg.. Controller, Model,,,etc
	i, found := helper.ArrayFind(genItems, args[0])
	if found {
		switch i {
		case 0:
			// fmt.Printf("%v\n", "Action to generate new model")
			// Model Generator
			status = ModelModifier(cmd, cli)

		}

	} else {

		fmt.Println("Error: Incorrect arguments, Do you meen? 'bindu modify model " + args[0] + "...'")
		fmt.Println("Available arguments: ", genItems)
		story.UpdateThisStoryStatus("False: incorrect command, Missing argument")
	}
	fmt.Println(status)
	return true
}

// ModelModifier to generate the model in project using the user inputs
func ModelModifier(cmd *cobra.Command, cli helper.CommandChain) bool {
	status := true
	// get all others args and flags
	args := cli.GetArgs()
	flags := cli.GetFlags()

	plural := pluralize.NewClient()
	modelName := args[1]

	//file path of the model
	fp := "./app/models/" + strings.Title(modelName) + ".go"
	// check if the file exist. if exist then suggest to change the command update instead generate
	if !helper.FileExists(fp) {
		// check if genarator should update or modify files
		story.UpdateThisStoryStatus("false: modelName not found to modify")
		return false
	}
	// append model properties in the model
	lineAfter := []string{"type " + modelName + " struct {", "type " + modelName}
	var newLines []string
	for _, prop := range args[2:] {
		// newLines = append(newLines, helper.PropertyFormatter(prop))
		fmt.Println(prop)
	}
	return false

	// Create the main model
	modelName = strings.Title(modelName)
	status = modelIfNotExistCreate(modelName, newLines)

	for _, flg := range flags {
		keyFlag := flg.Key
		valueFlag := flg.Values
		// SubcommandChain find the subcommand args
		subcommandChain := helper.SubCommandChain(valueFlag)
		// Check if subcommand exist then procced
		if len(subcommandChain) > 0 {
			// flags arguments commandChain first item is relationship model
			relModel := subcommandChain[0]
			if keyFlag == "hasOne" {
				// append new lines to the main model
				helper.AppendLinesInFile(fp, lineAfter, []string{"\t" + relModel + " " + relModel})
				// create related hasOne model
				status = modelIfNotExistCreate(relModel, []string{"\t" + args[1] + "ID uint64 `json:\"-\"`"})
			}
			if keyFlag == "belongsTo" {
				// belongs to properties for main model
				belongsToProps := []string{
					"\t" + relModel + " " + relModel,
					"\t" + relModel + "ID   uint64 `json:\"-\"`",
				}
				// append line to main model
				helper.AppendLinesInFile(fp, lineAfter, belongsToProps)
				// belongs to model
				status = modelIfNotExistCreate(relModel, nil)
			}
			if keyFlag == "hasMany" {
				// hasMany properties for main model
				hasManyProps := []string{
					"\t" + plural.Plural(relModel) + " []" + relModel,
				}
				// append line to main model
				helper.AppendLinesInFile(fp, lineAfter, hasManyProps)
				// belongs to model
				status = modelIfNotExistCreate(relModel, []string{"\t" + modelName + "ID   uint64 `json:\"-\"`"})
			}
			if keyFlag == "manyToMany" {
				// manyToMany properties for main model
				manyToManyProps := []string{
					"\t" + plural.Plural(relModel) + " []" + relModel + " `gorm:\"many2many:" + strings.ToLower(args[1]) + "_" + plural.Plural(strings.ToLower(relModel)) + ";association_foreignkey:id;foreignkey:id\"`",
				}
				// append line to main model
				helper.AppendLinesInFile(fp, lineAfter, manyToManyProps)
				// Crate manyToMany model if not exist
				// manyToMany Relationship
				status = modelIfNotExistCreate(relModel, []string{"\t" + plural.Plural(args[1]) + " []" + args[1] + " `gorm:\"many2many:" + strings.ToLower(args[1]) + "_" + plural.Plural(strings.ToLower(relModel)) + ";association_foreignkey:id;foreignkey:id\"`\n"})

			}
		}
	}
	return status
}

// modelIfNotExistCreate create model if not exist input model and props as string line
func modelIfNotExistCreate(model string, props []string) bool {

	path := "./app/models/" + model + ".go"
	if helper.FileExists(path) {
		fmt.Println("Failed to generate model!!! Model already exist.")
		story.UpdateThisStoryStatus("false: '" + path + "' Model Already exist!")
		return false
	}
	fmt.Println("Model creating..." + path)
	f, err := os.Create(path)
	helper.ErrorCheck(err)
	f.WriteString("package models\n\n")
	f.WriteString("// " + model + " public model generated by bindu\n")
	f.WriteString("type " + model + " struct {\n")

	for _, p := range props {
		f.WriteString(p + "\n")
	}

	// Initialize Primary, CreatedAt, UpdatedAt, DeletedAt property
	f.WriteString("\tDefaultProperties\n")
	f.WriteString("}")
	defer f.Close()
	fmt.Println(".......Success!")
	story.UpdateThisStoryStatus("true: '" + model + "' Created")
	return true
}

// ControllerGenerator to generate the controller in project using the user inputs
func ControllerGenerator(args []string) bool {
	relModels := helper.FindAssociationRecursivly(args[1])
	preloads := helper.CreatePreloadByModels(relModels)
	// fmt.Println(args)
	status := true
	fmt.Println("Initializing Controller " + args[1])
	newpath := filepath.Join(".", "app/controllers")
	os.MkdirAll(newpath, os.ModePerm)
	// Model name

	plural := pluralize.NewClient()
	modelVar := strings.ToLower(args[1])
	modelVar2 := strings.ToLower(args[1]) + "2"
	path := "./app/controllers/" + args[1] + "Controller.go"
	if helper.FileExists(path) {
		fmt.Println("Failed to generate controller!!!\nIf you want to update Controller,\nuse the --update or -u flag to add new or modify")
		story.UpdateThisStoryStatus("false: Already Exist (-u To Update) '" + path + "' ")
		return false
	}
	// If controller not exist then generate
	f, err := os.Create(path)
	helper.ErrorCheck(err)
	f.WriteString("package controllers\n\n")
	f.WriteString("import (\n\n")
	f.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/app/models\"\n")
	f.WriteString("\t\"" + helper.GetEnvValueByKey(".env", "APP_IMPORT_PATH") + "/db\"\n")
	f.WriteString(")\n\n")

	f.WriteString("// " + args[1] + " public controllers generated by bindu\n\n")
	// Generate CRUD
	// Get All Data for Index
	f.WriteString("// Index" + args[1] + " to get all data\n")
	f.WriteString("func Index" + args[1] + "(c *gin.Context) {\n")
	// Controller functionality here
	// model collection initialize
	f.WriteString("\tpage := c.Params.ByName(\"page\")\n")
	f.WriteString("\tlimit := 10\n")
	f.WriteString("\tvar db = db.DB\n")

	f.WriteString("\tvar " + plural.Plural(modelVar) + " []models." + args[1] + "\n")
	f.WriteString("\tdb" + preloads + ".Find(&" + plural.Plural(modelVar) + ").Offset(page).Limit(limit)\n\n")

	f.WriteString("\tc.JSON(200, " + plural.Plural(modelVar) + ")\n")
	f.WriteString("}\n\n")

	// Get Single Data for Show
	f.WriteString("// Show" + args[1] + " to get single data\n")
	f.WriteString("func Show" + args[1] + "(c *gin.Context) {\n")
	// Controller functionality here
	f.WriteString("\tid := c.Params.ByName(\"id\")\n")
	f.WriteString("\tvar db = db.DB\n")
	f.WriteString("\tvar " + modelVar + " models." + args[1] + "\n")
	f.WriteString("\tdb.Where(\"id=?\",id)" + preloads + ".Find(&" + modelVar + ")\n\n")

	f.WriteString("\tc.JSON(200, " + modelVar + ")\n")
	f.WriteString("}\n\n")

	// Post Single Data for New
	f.WriteString("// Create" + args[1] + " to a new data\n")
	f.WriteString("func Create" + args[1] + "(c *gin.Context) {\n")
	// Controller functionality here
	f.WriteString("\tvar db = db.DB\n")
	f.WriteString("\tvar " + modelVar + " models." + args[1] + "\n")
	f.WriteString("\tc.BindJSON(&" + modelVar + ")\n")
	f.WriteString("\tdb.Create(&" + modelVar + ")\n\n")
	// if the model is user then create login
	// if args[1] == "User" {

	// 	f.WriteString("\tvar login models.Login\n")
	// 	f.WriteString("\tc.BindJSON(&login)\n")
	// 	f.WriteString("\tif len(login.Username) == 0 {\n")
	// 	f.WriteString("\t\tlogin.Username=" + modelVar + ".UserName\n\t}\n")

	// 	f.WriteString("\tif len(login.Password) == 0 {\n")
	// 	f.WriteString("\t\tlogin.Password=login.Username+\"123\"\n\t}\n")

	// 	f.WriteString("\tdb.Create(&login)\n\n")
	// }
	f.WriteString("\tc.JSON(200, " + modelVar + ")\n")
	f.WriteString("}\n\n")

	// // Post Single Data for Create
	// f.WriteString("func New" + args[1] + "(c *gin.Context) {\n")
	// // Controller functionality here
	// f.WriteString("\tc.JSON(200, \"" + args[1] + " Create page\")\n")
	// f.WriteString("}\n\n")

	// Create New Data for Update
	f.WriteString("// Update" + args[1] + " to Update data\n")
	f.WriteString("func Update" + args[1] + "(c *gin.Context) {\n")
	// Controller functionality here
	f.WriteString("\tid := c.Params.ByName(\"id\")\n")
	f.WriteString("\tvar db = db.DB\n")
	f.WriteString("\tvar " + modelVar + " models." + args[1] + "\n")
	f.WriteString("\tvar " + modelVar2 + " models." + args[1] + "\n")
	f.WriteString("\tc.BindJSON(&" + modelVar + ")\n")

	f.WriteString("\tif err := db.Where(\"id=?\",id).Find(&" + modelVar2 + ").Error; err != nil {\n")
	f.WriteString("\t\tc.JSON(404, " + modelVar2 + ")\n")
	f.WriteString("\t\treturn\n\t}\n\n")

	f.WriteString("\tdb.Model(&" + modelVar2 + ").Update(&" + modelVar + ")\n\n")

	f.WriteString("\tc.JSON(200, " + modelVar + ")\n")
	f.WriteString("}\n\n")

	// Get Single Data for Destroy
	f.WriteString("// Destroy" + args[1] + " to delete single data\n")
	f.WriteString("func Destroy" + args[1] + "(c *gin.Context) {\n")
	// Controller functionality here
	f.WriteString("\tid := c.Params.ByName(\"id\")\n")
	f.WriteString("\tvar db = db.DB\n")
	f.WriteString("\tvar " + modelVar + " models." + args[1] + "\n")

	f.WriteString("\tif err := db.Where(\"id=?\",id).Find(&" + modelVar + ").Error; err != nil {\n")
	f.WriteString("\t\tc.JSON(404, " + modelVar + ")\n")
	f.WriteString("\t\treturn\n\t}\n\n")

	f.WriteString("\tdb.Delete(&" + modelVar + ")\n\n")

	f.WriteString("\tc.JSON(200, " + modelVar + ")\n")
	f.WriteString("}\n\n")

	defer f.Close()
	absPath, _ := filepath.Abs(newpath)
	fmt.Println(absPath + "/" + args[1] + "Controller.go\n")

	// if code succesfully exicute the return true
	status = true
	story.UpdateThisStoryStatus("true: Controller Created '" + args[1] + "'")
	return status
}

// RoutesGenerator Routes generator
func RoutesGenerator(cmd *cobra.Command, cli helper.CommandChain) bool {

	// routes not working... it should do first, find the fucname and append lines before return
	// return helper.WriteRoutes("routes/API.go", args[1])

	args := cli.GetArgs()
	path := "routes/API.go"
	routeName := args[1]
	var middleware string
	var apiGroup string
	flags := cli.GetFlags()
	for _, flg := range flags {
		if flg.Key == "middleware" {
			middleware = flg.Values[0]
		}
		if flg.Key == "group" {
			apiGroup = flg.Values[0]
		}
	}
	// if api Group flag raised then check or create api group
	if len(apiGroup) > 0 {
		// first create api group
		var newLinesAPIGroup []string
		newLinesAPIGroup = append(newLinesAPIGroup, "\n\t"+apiGroup+" := api.Group(\"/"+apiGroup+"\")\n")
		// append line to main Routes
		lineAfterAPIGroup := []string{
			"r.Group(\"/api\")",
			"Use(middlewares.Cors())",
		}

		// Generate Routes
		helper.AppendLinesInFile(path, lineAfterAPIGroup, newLinesAPIGroup)
	}
	if len(middleware) > 0 {
		rName := helper.IfThenElse(len(apiGroup) > 0, apiGroup, "api")
		var newLinesMiddleware []string
		if middleware == "auth" {
			newLinesMiddleware = append(newLinesMiddleware, "\n\t"+fmt.Sprintf("%v", rName)+".Use(authMiddleware.MiddlewareFunc())")
			newLinesMiddleware = append(newLinesMiddleware, "\t{//"+apiGroup+" Auth route Start \n\n")
			newLinesMiddleware = append(newLinesMiddleware, "\n\t} //"+apiGroup+" Auth Route end\n")
		}

		// append line to main Routes
		lineAfterAPIMiddleWare := []string{
			"api.Group(\"/" + fmt.Sprintf("%v", rName) + "\")",
		}
		// Generate Routes
		helper.AppendLinesInFile(path, lineAfterAPIMiddleWare, newLinesMiddleware)
	}
	lines, err := helper.ScanLines(path)
	helper.ErrorCheck(err)
	var lineNumber int
	var newLines []string
	newLines = append(newLines, "\n\t//"+apiGroup+" Routes generated by bindu for "+routeName+"\n")

	for i, line := range lines {

		if helper.StringsContains(line, "func API") {
			// fmt.Println("Hello")
			lineNumber = i
			newLines = append(newLines, lines[:i+1]...)
			newLines = append(newLines, "\n")
			type route struct {
				key   string
				value string
			}
			routes := []route{
				{key: "GET", value: "Index"},
				{key: "POST", value: "Create"},
				{key: "GET", value: "Show"},
				{key: "PUT", value: "Update"},
				{key: "DELETE", value: "Destroy"},
			}

			for _, r := range routes {
				plural := pluralize.NewClient()
				switch r.value {
				case "Index":
					newLines = append(newLines, "\t"+fmt.Sprintf("%v", helper.IfThenElse(len(apiGroup) > 0, apiGroup, "api"))+"."+r.key+"(\"/"+strings.ToLower(plural.Plural(routeName))+"/:page\", controllers."+r.value+routeName+")")
				case "Create":
					newLines = append(newLines, "\t"+fmt.Sprintf("%v", helper.IfThenElse(len(apiGroup) > 0, apiGroup, "api"))+"."+r.key+"(\"/"+strings.ToLower(routeName)+"\", controllers."+r.value+routeName+")")

				case "Show":
					newLines = append(newLines, "\t"+fmt.Sprintf("%v", helper.IfThenElse(len(apiGroup) > 0, apiGroup, "api"))+"."+r.key+"(\"/"+strings.ToLower(routeName)+"/:id\", controllers."+r.value+routeName+")")

				case "Update":
					newLines = append(newLines, "\t"+fmt.Sprintf("%v", helper.IfThenElse(len(apiGroup) > 0, apiGroup, "api"))+"."+r.key+"(\"/"+strings.ToLower(routeName)+"/:id\", controllers."+r.value+routeName+")")
				case "Destroy":
					newLines = append(newLines, "\t"+fmt.Sprintf("%v", helper.IfThenElse(len(apiGroup) > 0, apiGroup, "api"))+"."+r.key+"(\"/"+strings.ToLower(routeName)+"/:id\", controllers."+r.value+routeName+")")

				}
			}

		}

	}

	if lineNumber > 0 {
		newLines = append(newLines, lines[lineNumber+1:]...)
		// ioutil.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
		// append line to main Routes
		// first check if the flag raised the group or middleware
		// if flag exist apiGroup
		newLineAfter := helper.IfThenElse(len(apiGroup) > 0, "api.Group(\"/"+apiGroup+"\")", "r.Group(\"/api\")")
		// if flag exist middleware
		newLineAfter = helper.IfThenElse(len(middleware) > 0, "api.Use(authMiddleware.MiddlewareFunc())", newLineAfter)
		// if flag exist apigrou and middleware both
		newLineAfter = helper.IfThenElse(len(middleware) > 0 && len(apiGroup) > 0, apiGroup+".Use(authMiddleware.MiddlewareFunc())", newLineAfter)
		lineAfter := []string{
			fmt.Sprintf("%v", newLineAfter),
		}

		fmt.Println("Generating Routes....")
		// Generate Routes
		// inpur path, lineAfter, offset, newlines
		if len(middleware) > 0 {
			// because of middleware { } bracket need to offset 1
			helper.AppendLinesInFileNext(path, lineAfter, 1, newLines)
		} else {
			helper.AppendLinesInFileNext(path, lineAfter, 0, newLines)
		}

		fmt.Println("Done!")
	}
	return true

}
