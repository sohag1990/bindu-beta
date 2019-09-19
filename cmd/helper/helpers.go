package helper

import (
	"bufio"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/manifoldco/promptui"
)

// ArrayFind return true or false with int. first check the bool then use index, Because if not found item also return 0
func ArrayFind(s interface{}, elem interface{}) (int, bool) {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return i, true
			}
		}
	}

	return 0, false
}

// SanitizeUserInputReArry to sanitize user input for generate command: sometime user will input with comma sometime without comma or vice varsa.
// args passed make them formated upppercase
func SanitizeUserInputReArry(args []string) []string {
	var reArray []string

	var j int
	for _, arg := range args {

		commaSepToArray := strings.Split(arg, ",")
		// fmt.Println(len(commaSepToArray))
		if len(commaSepToArray) > 1 {
			for _, cItem := range commaSepToArray {
				reArray = append(reArray, strings.Title(cItem))
				j++
			}
		} else {
			reArray = append(reArray, strings.Title(arg))
			j++
		}

	}
	// if missing arg name: suppose to missing Model/Controller/Scafforl/View Name just promt user to input the name
	if (len(reArray) > 2 && strings.Contains(reArray[1], ":")) || len(reArray) == 1 {
		fmt.Printf("%v\n", "Forgot to enter")
		promptArgName := promptui.Prompt{
			Label:    reArray[0] + " Name",
			Validate: nil,
		}
		argName, err := promptArgName.Run()
		ErrorCheck(err)
		var newArray []string
		newArray = append(newArray, reArray[0])
		newArray = append(newArray, strings.Title(argName))
		newArray = append(newArray, reArray[1:]...)
		reArray = newArray
	}

	return CleanEmptyArray(reArray)
}

// SanitizeUserInput to sanitize user input for generate command:
//  sometime user will input with comma sometime without comma or vice varsa.
// args passed make them formated upppercase
func SanitizeUserInput(args []string) []string {
	var reArray []string

	var j int
	for _, arg := range args {

		commaSepToArray := strings.Split(arg, ",")
		// fmt.Println(len(commaSepToArray))
		if len(commaSepToArray) > 1 {
			for _, cItem := range commaSepToArray {
				reArray = append(reArray, strings.Title(cItem))
				j++
			}
		} else {
			reArray = append(reArray, strings.Title(arg))
			j++
		}

	}
	// // if missing arg name: suppose to missing Model/Controller/Scafforl/View Name just promt user to input the name
	// if (len(reArray) > 2 && strings.Contains(reArray[1], ":")) || len(reArray) == 1 {
	// 	fmt.Printf("%v\n", "Forgot to enter")
	// 	promptArgName := promptui.Prompt{
	// 		Label:    reArray[0] + " Name",
	// 		Validate: nil,
	// 	}
	// 	argName, err := promptArgName.Run()
	// 	ErrorCheck(err)
	// 	var newArray []string
	// 	newArray = append(newArray, reArray[0])
	// 	newArray = append(newArray, strings.Title(argName))
	// 	newArray = append(newArray, reArray[1:]...)
	// 	reArray = newArray
	// }

	return CleanEmptyArray(reArray)
}

// CleanEmptyArray Remove empty item from array
func CleanEmptyArray(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

type Property struct {
	Key   string
	Value string
}

// PropertyFormatter formate the correct property and type
func PropertyFormatter(prop string) string {
	p := strings.Split(prop, ":")
	if len(p) >= 2 {
		return p[0] + "\t\t" + strings.ToLower(p[1])
	}

	return ""
}

// ErrorCheck if error the panic
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

// IsInProjectDir if run command outside the project directory
func IsInProjectDir() {
	if _, err := os.Stat("./.env"); os.IsNotExist(err) {
		fmt.Println("You are not in project directory")
		ErrorCheck(err)
	}
}

// letter fix if someone inside project directory and want to create a new project
//
//

// IfThenElse evaluates a condition, if true returns the first parameter otherwise the second
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

// MakeAppName to sanitize appName
func MakeAppName(arg string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	ErrorCheck(err)
	appName := reg.ReplaceAllString(arg, "")
	return strings.Trim(appName, "")
}

// AppImportPath Project import path generate
func AppImportPath(appName string) string {
	projectPath := strings.Replace(PWD(), build.Default.GOPATH+"/src/", "", -1)
	return projectPath + "/" + appName
}

// PWD get present working directory
func PWD() string {
	pwd, err := os.Getwd()
	ErrorCheck(err)
	return pwd
}

// CD change directory
func CD(appName string) {
	errCD := os.Chdir(filepath.Join(PWD(), appName))
	ErrorCheck(errCD)
}

// GetEnvValueByKey Read Env file and find value by key
func GetEnvValueByKey(path string, envKey string) string {
	lines, err := scanLines(path)
	if err != nil {
		fmt.Println("File not found!")
		panic(err)
	}
	for _, line := range lines {
		if strings.Contains(line, envKey) {
			// found the line
			// get the value only
			splitLine := strings.Split(line, "=")
			value := splitLine[1]
			// fmt.Println(line)
			return value
		}
	}
	return ""
}

func WriteRoutes(path string, routeName string) {
	lines, err := scanLines(path)
	ErrorCheck(err)
	var lineNumber int
	var newLines []string
	for i, line := range lines {

		if strings.Contains(line, "func API") {
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
					newLines = append(newLines, "\tr."+r.key+"(\"/"+strings.ToLower(plural.Plural(routeName))+"/:page\", controllers."+r.value+routeName+")")
				case "Create":
					newLines = append(newLines, "\tr."+r.key+"(\"/"+strings.ToLower(routeName)+"\", controllers."+r.value+routeName+")")

				case "Show":
					newLines = append(newLines, "\tr."+r.key+"(\"/"+strings.ToLower(routeName)+"/:id\", controllers."+r.value+routeName+")")

				case "Update":
					newLines = append(newLines, "\tr."+r.key+"(\"/"+strings.ToLower(routeName)+"/:id\", controllers."+r.value+routeName+")")
				case "Destroy":
					newLines = append(newLines, "\tr."+r.key+"(\"/"+strings.ToLower(routeName)+"/:id\", controllers."+r.value+routeName+")")

				}
			}

		}
	}

	if lineNumber > 0 {
		newLines = append(newLines, lines[lineNumber+1:]...)
		ioutil.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
	}

}

func scanLines(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

// ScanDirFindLine fix import path when app migrate to new name or new host
func ScanDirFindLine(dirToScan string, keywordToFind string, f *os.File) {
	err := filepath.Walk(dirToScan,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}
			// fmt.Println(path, info.Size())
			// if extenstion .go then read file and find old import path and replace new import path
			if strings.Contains(path, ".go") {
				lines, err := scanLines(path)
				ErrorCheck(err)
				for _, line := range lines {

					if strings.Contains(line, keywordToFind) {

						lineSplit := strings.Split(line, " ")
						s := lineSplit[1]
						// if the method public, to detect public methor check the name first letter is capital
						if s != "DefaultProperties" {
							f.WriteString("\t" + "db.DB.AutoMigrate(models." + s + "{})\n")
						}

					}
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

// FixImportPath fix import path when app migrate to new name or new host
func FixImportPath(oldImportPath string, newImportPath string) {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}
			// fmt.Println(path, info.Size())
			// if extenstion .go then read file and find old import path and replace new import path
			if strings.Contains(path, ".go") {
				lines, err := scanLines(path)
				ErrorCheck(err)
				for _, line := range lines {
					// fmt.Println(line)
					if strings.Contains(line, oldImportPath) {
						// found the line
						read, err := ioutil.ReadFile(path)
						ErrorCheck(err)
						newContents := strings.Replace(string(read), oldImportPath, newImportPath, -1)
						// fmt.Println(newContents)
						err = ioutil.WriteFile(path, []byte(newContents), 0)
						ErrorCheck(err)
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
