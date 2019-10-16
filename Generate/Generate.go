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
func Generator(cmd *cobra.Command, c helper.CommandChain) {
	//--update flag check.
	u, _ = strconv.ParseBool(fmt.Sprintf("%v", cmd.Flag("update").Value))
	// User command args
	args := c.GetArgs()

	// flags := c.GetFlags()
	genItems := []string{"Model", "Controller", "Scaffold", "Routes", "View"}
	// if no agruments Prompt to select action name eg. model, controller, scaffold, view
	if len(args) == 0 {
		_, actionName := helper.AskSelect("Select To Generate", genItems)
		args = append(args, actionName)
		c.SetCliArgs(args)
	}

	// check the first items is matched with predefined arg.. Controller, Model,,,etc
	i, found := helper.ArrayFind(genItems, args[0])
	if found {
		switch i {
		case 0:
			// fmt.Printf("%v\n", "Action to generate new model")
			// Model Generator
			ModelGenerator(cmd, c)

		case 1:
			// fmt.Printf("%v\n", "Controller action")
			ControllerGenerator(args)
			RoutesGenerator(args)
		case 2:
			// fmt.Printf("%v\n", "Scaffold action")
			ModelGenerator(cmd, c)
			ControllerGenerator(args)
			RoutesGenerator(args)
		case 3:
			// fmt.Printf("%v\n", "Routes action")
			RoutesGenerator(args)
			// case 4:
			// 	fmt.Printf("%v\n", "View action")
			// 	fmt.Printf("%v\n", args)
			// default:
			// 	fmt.Printf("%v\n", "No action taken")
			// 	fmt.Printf("%v\n", args)
		}

	}
}

// ModelGenerator to generate the model in project using the user inputs
func ModelGenerator(cmd *cobra.Command, c helper.CommandChain) {

	// get all others args and flags
	args := c.GetArgs()
	flags := c.GetFlags()

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
			return
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
	modelIfNotExistCreate(modelName, newLines)

	for _, flg := range flags {
		key := flg.Key
		// SubcommandChain find the subcommand args
		subcommandChain := helper.SubCommandChain(flg.Values)
		// fmt.Println(values)
		if key == "hasOne" && len(subcommandChain) > 0 {
			// hasOne flags arguments commandChain format then find the args
			hasOneModel := subcommandChain[0]

			// append new lines to the main model
			helper.AppendLinesInFile(fp, lineAfter, []string{"\t" + hasOneModel + " " + hasOneModel})
			// create related hasOne model
			modelIfNotExistCreate(hasOneModel, []string{"\t" + args[1] + "ID uint64"})
		}
		if key == "belongsTo" && len(subcommandChain) > 0 {
			// belongsTo flags arguments commandChain format then find the args
			belongsToModel := subcommandChain[0]

			// belongs to properties for main model
			belongsToProps := []string{
				"\t" + belongsToModel + " " + belongsToModel,
				"\t" + belongsToModel + "ID   uint64",
			}
			// append line to main model
			helper.AppendLinesInFile(fp, lineAfter, belongsToProps)
			// belongs to model
			modelIfNotExistCreate(belongsToModel, nil)
		}
		if key == "hasMany" && len(subcommandChain) > 0 {
			// hasManyModel flags arguments commandChain format then find the args
			hasManyModel := subcommandChain[0]
			// hasMany properties for main model
			hasManyProps := []string{
				"\t" + plural.Plural(hasManyModel) + " []" + hasManyModel,
			}
			// append line to main model
			helper.AppendLinesInFile(fp, lineAfter, hasManyProps)
			// belongs to model
			modelIfNotExistCreate(hasManyModel, []string{"\t" + modelName + "ID   uint64"})
		}
		if key == "manyToMany" && len(subcommandChain) > 0 {
			// manyToManyModel flags arguments commandChain format then find the args
			manyToManyModel := subcommandChain[0]
			// manyToMany properties for main model
			manyToManyProps := []string{
				"\t" + plural.Plural(manyToManyModel) + " []" + manyToManyModel + " `gorm:\"many2many:" + strings.ToLower(args[1]) + "_" + plural.Plural(strings.ToLower(manyToManyModel)) + ";association_foreignkey:id;foreignkey:id\"`",
			}
			// append line to main model
			helper.AppendLinesInFile(fp, lineAfter, manyToManyProps)
			// Crate manyToMany model if not exist
			// manyToMany Relationship
			modelIfNotExistCreate(manyToManyModel, []string{"\t" + plural.Plural(args[1]) + " []" + args[1] + " `gorm:\"many2many:" + strings.ToLower(args[1]) + "_" + plural.Plural(strings.ToLower(manyToManyModel)) + ";association_foreignkey:id;foreignkey:id\"`\n"})

		}
	}

}

// modelIfNotExistCreate create model if not exist input model and props as string line
func modelIfNotExistCreate(model string, props []string) {

	path := "./app/models/" + model + ".go"
	if helper.FileExists(path) {
		fmt.Println("Failed to generate model!!! Model already exist.")
		return
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
}

// ControllerGenerator to generate the controller in project using the user inputs
func ControllerGenerator(args []string) {
	// fmt.Println(args)

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
		return
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

}

// RoutesGenerator Routes generator
func RoutesGenerator(args []string) {
	// routes not working... it should do first, find the fucname and append lines before return
	helper.WriteRoutes("routes/API.go", args[1])
}
