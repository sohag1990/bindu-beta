package create

import (
	"fmt"
	"log"
	"os"
	"strings"

	helper "github.com/bindu-bindu/bindu/Helper"
	story "github.com/bindu-bindu/bindu/Story"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// --update or -u flag variable
var u bool

// Create is a func to create User, policy etc
func Create(cmd *cobra.Command, cli helper.CommandChain) bool {
	// To return a status true or false, whether code executed or not
	// Deafult true
	status := true
	//--update flag check.
	// u, _ = strconv.ParseBool(fmt.Sprintf("%v", cmd.Flag("update").Value))
	// User command args
	args := cli.GetArgs()

	// flags := c.GetFlags()
	items := []string{"User", "Policy"}
	// if no agruments Prompt to select action name eg. model, controller, scaffold, view
	if len(args) == 0 {
		_, actionName := helper.AskSelect("Select To Create", items)
		args = append(args, actionName)
		cli.SetCliArgs(args)
	}

	// check the first items is matched with predefined arg.. Controller, Model,,,etc
	i, found := helper.ArrayFind(items, args[0])
	if found {
		switch i {
		case 0:
			// fmt.Printf("%v\n", "Action to generate new model")
			// Model Generator
			status = UserCreate(cmd, cli)

		case 1:
			// fmt.Printf("%v\n", "Controller action")
			status = PolicyCreate(cmd, cli)

		}

	}
	return status
}

// UserCreate to create the user
func UserCreate(cmd *cobra.Command, cli helper.CommandChain) bool {
	//db connection
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	status := true
	db, err := gorm.Open(os.Getenv("DB_ADAPTER"), os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_DATABASE")+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		helper.ErrorCheck(err)
		status = false
	}

	defer db.Close()

	// get all others args and flags
	args := cli.GetArgs()
	// flags := cli.GetFlags()
	type User struct {
		Username string
		Password string
		Role     string
	}
	var user User
	// plural := pluralize.NewClient()
	for _, prop := range args[1:] {
		p := strings.Split(prop, ":")

		key := strings.ToLower(p[0])
		switch {
		case key == "username":
			user.Username = strings.ToLower(p[1]) //username always lowercase
		case key == "password":
			user.Password = p[1]
		case key == "role":
			user.Role = strings.ToLower(p[1]) // role always lowercase
		}
	}

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("Error Happend! called ", err)
		// story.UpdateThisStoryStatus("false: " + err.Error())
	} else {
		fmt.Println("User Created Successfully!")
		// story.UpdateThisStoryStatus("true: User Created")
	}
	return status
}

// PolicyCreate to create the Policy CasbinRule
func PolicyCreate(cmd *cobra.Command, cli helper.CommandChain) bool {
	//db connection
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		story.UpdateThisStoryStatus("false: .env file can't load, " + err.Error())
	}

	db, err := gorm.Open(os.Getenv("DB_ADAPTER"), os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_DATABASE")+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		helper.ErrorCheck(err)
		story.UpdateThisStoryStatus("false: .env file can't load, " + err.Error())
	}
	defer db.Close()
	status := true
	// get all others args and flags
	args := cli.GetArgs()
	// CasbinRule model

	var rule CasbinRule
	// plural := pluralize.NewClient()
	for _, prop := range args[1:] {

		p := strings.Split(strings.ToLower(prop), ":")

		switch {
		case p[0] == "alice" || p[0] == "role" || p[0] == "policy" || p[0] == "rule" || p[0] == "ptype":
			rule.PType = strings.ToLower(p[1])
		case p[0] == "sub" || p[0] == "v0":
			rule.V0 = strings.ToLower(p[1])
		case p[0] == "obj" || p[0] == "v1":
			rule.V1 = strings.ToLower(p[1])
		case p[0] == "act" || p[0] == "v2":
			// act like GET, POST so do not lower case
			rule.V2 = p[1]
		}
	}
	var ruleExistCount int
	db.Model(&CasbinRule{}).Where("p_type=? and v0=? and v1=? and v2=?", rule.PType, rule.V0, rule.V1, rule.V2).Count(&ruleExistCount)
	if ruleExistCount > 0 {
		fmt.Println("Error Plicy already exist")
		story.UpdateThisStoryStatus("false: Policy Already Exist")
	} else {
		if err := db.Create(&rule).Error; err != nil {
			story.UpdateThisStoryStatus("false: Err " + err.Error())
			fmt.Println("Failed to create policy")
		} else {
			story.UpdateThisStoryStatus("true: Policy Created")
			fmt.Println("Policy Created Successfully!")
		}
	}
	return status
}

// CasbinRule model for policy set
type CasbinRule struct {
	PType string `gorm:"column:p_type"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
}

// TableName Set TableName name to be `casbin_rule`
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
