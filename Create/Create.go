package create

import (
	"fmt"
	"log"
	"os"
	"strings"

	helper "github.com/bindu-bindu/bindu/Helper"
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

	db, err := gorm.Open(os.Getenv("DB_ADAPTER"), os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_DATABASE")+"?charset=utf8&parseTime=True&loc=Local")
	helper.ErrorCheck(err)
	defer db.Close()
	status := true
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
		switch {
		case p[0] == "Username" || p[0] == "UserName":
			user.Username = p[1]
		case p[0] == "Password" || p[0] == "PassWord":
			user.Password = p[1]
		case p[0] == "Role" || p[0] == "Rule":
			user.Role = p[1]

		}
	}

	db.Create(&user)
	fmt.Println("User Created Successfully!")
	return status
}

// PolicyCreate to create the Policy CasbinRule
func PolicyCreate(cmd *cobra.Command, cli helper.CommandChain) bool {
	//db connection
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(os.Getenv("DB_ADAPTER"), os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_DATABASE")+"?charset=utf8&parseTime=True&loc=Local")
	helper.ErrorCheck(err)
	defer db.Close()
	status := true
	// get all others args and flags
	args := cli.GetArgs()
	// CasbinRule model

	var rule CasbinRule
	// plural := pluralize.NewClient()
	for _, prop := range args[1:] {
		p := strings.Split(prop, ":")
		switch {
		case p[0] == "Alice" || p[0] == "Role" || p[0] == "Policy" || p[0] == "Rule" || p[0] == "PType" || p[0] == "Ptype":
			rule.PType = strings.ToLower(p[1])
		case p[0] == "Sub" || p[0] == "v0" || p[0] == "V0":
			rule.V0 = strings.ToLower(p[1])
		case p[0] == "Obj" || p[0] == "v1" || p[0] == "V1":
			rule.V1 = strings.ToLower(p[1])
		case p[0] == "Act" || p[0] == "v2" || p[0] == "V2":
			// act like GET, POST so do not lower case
			rule.V2 = p[1]

		}
	}
	fmt.Println(rule)
	db.Create(&rule)
	fmt.Println("Policy Created Successfully!")
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
