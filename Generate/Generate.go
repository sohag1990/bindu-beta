package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	helper "github.com/bindu-bindu/bindu/Helper"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
)

// --update or -u flag variable
var u bool

// Generator is a func to generate model controller and scaffold
func Generator(cmd *cobra.Command, cli helper.CommandChain) bool {
	// To return a status true or false, whether code executed or not
	// Deafult true
	status := true
	//--update flag check.
	u, _ = strconv.ParseBool(fmt.Sprintf("%v", cmd.Flag("update").Value))
	// User command args
	args := cli.GetArgs()

	// flags := c.GetFlags()
	genItems := []string{"Model", "Controller", "Scaffold", "Routes", "View"}
	// if no agruments Prompt to select action name eg. model, controller, scaffold, view
	if len(args) == 0 {
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
			status = ModelGenerator(cmd, cli)

		case 1:
			// fmt.Printf("%v\n", "Controller action")
			status = ControllerGenerator(args)
			status = RoutesGenerator(args)
		case 2:
			// fmt.Printf("%v\n", "Scaffold action")
			status = ModelGenerator(cmd, cli)
			status = ControllerGenerator(args)
			status = RoutesGenerator(args)
		case 3:
			// fmt.Printf("%v\n", "Routes action")
			status = RoutesGenerator(args)
			// case 4:
			// 	fmt.Printf("%v\n", "View action")
			// 	fmt.Printf("%v\n", args)
			// default:
			// 	fmt.Printf("%v\n", "No action taken")
			// 	fmt.Printf("%v\n", args)
		}

	}
	return status
}

// ModelGenerator to generate the model in project using the user inputs
func ModelGenerator(cmd *cobra.Command, cli helper.CommandChain) bool {
	status := true
	// get all others args and flags
	args := cli.GetArgs()
	flags := cli.GetFlags()

	plural := pluralize.NewClient()
	modelName := args[1]

	// if directory not present create models dir
	newpath := filepath.Join(".", "app/models")
	os.MkdirAll(newpath, os.ModePerm)
	//file path of the model
	fp := "./app/models/" + modelName + ".go"
	// check if the file exist. if exist then suggest to change the command update instead generate
	if helper.FileExists(fp) {
		// check if genarator should update or modify files
		if !u {
			fmt.Println("Model already exist. If you want to update model,\nuse the --update or -u flag to add new or modify")
			return false
		}
	}
	// Create model if not exist, ofcouse not present at this moment, create one
	// append model properties in the model
	lineAfter := []string{"type " + modelName + " struct {", "type " + modelName}
	var newLines []string
	for _, prop := range args[2:] {
		newLines = append(newLines, helper.PropertyFormatter(prop))
	}

	// if --update flag fired update the model
	if u {
		helper.AppendLinesInFile(fp, lineAfter, newLines)
	}
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
				status = modelIfNotExistCreate(relModel, []string{"\t" + args[1] + "ID uint64"})
			}
			if keyFlag == "belongsTo" {
				// belongs to properties for main model
				belongsToProps := []string{
					"\t" + relModel + " " + relModel,
					"\t" + relModel + "ID   uint64",
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
				status = modelIfNotExistCreate(relModel, []string{"\t" + modelName + "ID   uint64"})
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
		return false
	}
	fmt.Println("Model creating..." + path)
	f, err := os.Create(path)
	helper.ErrorCheck(err)
	f.WriteString("package models\n\n")
	f.WriteString("// " + model + " public model generated by bindu\n")
	f.WriteString("type " + model + " struct {\n")
	// Initialize Primary, CreatedAt, UpdatedAt, DeletedAt property

	for _, p := range props {
		f.WriteString(p + "\n")
	}
	f.WriteString("\tDefaultProperties\n")
	f.WriteString("}")
	defer f.Close()
	fmt.Println(".......Success!")
	return true
}

// ControllerGenerator to generate the controller in project using the user inputs
func ControllerGenerator(args []string) bool {
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
	f.WriteString("\tdb.Find(&" + plural.Plural(modelVar) + ").Offset(page).Limit(limit)\n\n")

	f.WriteString("\tc.JSON(200, " + plural.Plural(modelVar) + ")\n")
	f.WriteString("}\n\n")

	// Get Single Data for Show
	f.WriteString("// Show" + args[1] + " to get single data\n")
	f.WriteString("func Show" + args[1] + "(c *gin.Context) {\n")
	// Controller functionality here
	f.WriteString("\tid := c.Params.ByName(\"id\")\n")
	f.WriteString("\tvar db = db.DB\n")
	f.WriteString("\tvar " + modelVar + " models." + args[1] + "\n")
	f.WriteString("\tdb.Where(\"id=?\",id).Find(&" + modelVar + ")\n\n")

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
	return status
}

// RoutesGenerator Routes generator
func RoutesGenerator(args []string) bool {

	// routes not working... it should do first, find the fucname and append lines before return
	return helper.WriteRoutes("routes/API.go", args[1])

}
