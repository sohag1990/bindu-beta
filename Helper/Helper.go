package helper

import (
	"bufio"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

// ENV_APP env variable app data
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

// ENV_DB env variable app DB connection data
type ENV_DB struct {
	DB_ADAPTER  string
	DB_HOST     string
	DB_PORT     string
	DB_DATABASE string
	DB_USERNAME string
	DB_PASSWORD string
}

// ENV_LOG env variable app log data
type ENV_LOG struct {
	LOG_CHANNEL string
}

// CLI is the struct of user commands
type CLI struct {
	Args  []string
	Flags []Flag
}

// Flag is the command flag key and value eg. --port 8080 --host localhost
type Flag struct {
	Key    string
	Values []string
}

// CommandChain is the command chain of user input.
type CommandChain interface {
	SetCli(args []string, flags []Flag)
	SetCliArgs(args []string)
	SetCliFlags(flags []Flag)
	GetArgs() []string
	GetFlags() []Flag
}

// SetCli to set cli
func (c *CLI) SetCli(args []string, flags []Flag) {
	c.Args = args
	c.Flags = flags
}

// SetCliArgs to set cli
func (c *CLI) SetCliArgs(args []string) {
	c.Args = args
}

// SetCliFlags to set cli
func (c *CLI) SetCliFlags(flags []Flag) {
	c.Flags = flags
}

// InitialCli Initialize CLI
func InitialCli() *CLI {
	return &CLI{}
}

// GetArgs to find sanitized array of commands args
func (c *CLI) GetArgs() []string {
	var newArgs []string
	for _, arg := range c.Args {

		splitArg := strings.Split(arg, ",")
		if len(splitArg) > 1 {
			for _, cItem := range splitArg {
				newArgs = append(newArgs, cItem)

			}
		} else {
			newArgs = append(newArgs, arg)

		}
	}
	c.Args = newArgs
	return CleanEmptyArray(c.Args)
}

// GetFlags to find sanitized array of commands args
func (c CLI) GetFlags() []Flag {
	var newFlags []Flag
	for _, flag := range c.Flags {
		var f Flag
		f.Key = flag.Key
		f.Values = strings.Split(flag.Values[0], ",")
		newFlags = append(newFlags, f)
	}
	return newFlags

}

//SubCommandChain Subcommand/Flag value args clichain
func SubCommandChain(args []string) []string {
	var cli CommandChain
	cli = InitialCli()
	cli.SetCliArgs(args)
	return cli.GetArgs()
}

// ArrayFind return true or false with int. first check the bool then use index, Because if not found item also return 0
func ArrayFind(arr interface{}, keyword interface{}) (int, bool) {
	arrV := reflect.ValueOf(arr)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if strings.ToLower(fmt.Sprintf("%v", arrV.Index(i).Interface())) == strings.ToLower(fmt.Sprintf("%v", keyword)) {
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
				reArray = append(reArray, cItem)
				j++
			}
		} else {
			reArray = append(reArray, arg)
			j++
		}

	}
	// if missing arg name: suppose to missing Model/Controller/Scafforl/View Name just promt user to input the name
	if (len(reArray) > 2 && StringsContains(reArray[1], ":")) || len(reArray) == 1 {
		fmt.Printf("%v\n", "Forgot to enter")
		promptArgName := promptui.Prompt{
			Label:    reArray[0] + " Name",
			Validate: nil,
		}
		argName, err := promptArgName.Run()
		ErrorCheck(err)
		var newArray []string
		newArray = append(newArray, reArray[0])
		newArray = append(newArray, argName)
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
				reArray = append(reArray, cItem)
				j++
			}
		} else {
			reArray = append(reArray, arg)
			j++
		}

	}
	// // if missing arg name: suppose to missing Model/Controller/Scafforl/View Name just promt user to input the name
	// if (len(reArray) > 2 && StringsContains(reArray[1], ":")) || len(reArray) == 1 {
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
		str = strings.TrimSpace(str)
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
// Gorm attributes(index, varchar(255) unique_index, null, not null, default 0) json:"-"
func PropertyFormatter(prop string) string {
	p := strings.Split(prop, ":")
	if len(p) >= 2 {

		pVal := strings.ToLower(p[1])
		pValLimit := strings.Split(pVal, "_")

		if len(pValLimit) > 1 {
			if pValLimit[0] == "text" {
				pVal = "string" + "\t\t" + "`gorm:\"type:text;\"`"
			} else if pValLimit[0] == "string" {
				pVal = "string" + "\t\t" + "`gorm:\"type:varchar(" + pValLimit[1] + ");\"`"
			}
		} else {
			if pVal == "text" {
				pVal = "string" + "\t\t" + "`gorm:\"type:text;\"`"
			} else if pVal == "time" || pVal == "date" || pVal == "date" || pVal == "datetime" || pVal == "time.Time" {
				pVal = "time.Time"

			}
		}

		return "\t" + strings.Title(p[0]) + "\t\t" + pVal
	}

	return ""
}

// ErrorCheck if error the panic
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println(e)
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
	if len(appName) == 0 {
		return ""
	}
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

// FindAssociationRecursivly Read models by key to find associations
func FindAssociationRecursivly(associatedModel string) []string {
	dir := "./app/models/"
	files, err := ioutil.ReadDir(dir)
	var relModels []string
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		path := dir + f.Name()
		relModels = append(relModels, FindCustomDataType(FindModelProperties(path, associatedModel))...)
	}
	return relModels
}

//CreatePreloadByModels Preload string creator, take models as string array, output as preload string
func CreatePreloadByModels(relModels []string) string {
	var preloadString string
	for _, model := range relModels {
		if StringsContains(model, "[]") {
			// if array then put s with modelname
			// remove the [] from the modelname
			model = strings.Trim(model, "[")
			model = strings.Trim(model, "]")
			model = strings.Trim(model, "*")
			preloadString = preloadString + ".Preload(\"" + model + "s" + "\")"
		} else {
			preloadString = preloadString + ".Preload(\"" + model + "\")"
		}
	}
	return preloadString
}

// FindCustomDataType find custom model from model properties
func FindCustomDataType(prop []string) []string {
	goTypes := []string{"string", "int", "uint", "bool", "*string", "*int", "*bool"}
	var customTypes []string
	for _, p := range prop {
		pSplit := strings.Split(p, " ")
		pSplit = CleanEmptyArray(pSplit)
		if len(pSplit) >= 2 {

			// if is not go define type then its a custom type.
			// so its a relational model for sure.
			if !ItemExistInArray(goTypes, pSplit[1]) {
				customTypes = append(customTypes, pSplit[1])
			}
		}
	}
	return customTypes
}

// ItemExistInArray check if item exist in array
func ItemExistInArray(arr []string, item string) bool {
	for _, dType := range arr {
		if StringsContains(item, dType) {
			return true
		}
	}
	return false
}

// FindModelProperties as array of string. Collect only one matching model properties
func FindModelProperties(path string, model string) []string {

	if !FileExists(path) {
		return nil
	}
	lines, err := ScanLines(path)
	lineAfter := "type " + model + " struct"

	if err != nil {
		fmt.Println("File not found!")
		panic(err)
	}
	var newLines []string
	for i, line := range lines {
		if StringsContains(line, lineAfter) {
			copyLines := lines[i+1:]
			for _, l := range copyLines {
				newLines = append(newLines, l)
				if StringsContains(l, "}") {
					// Remove the last curly } bracket
					// Clear Empty array
					newLines = CleanEmptyArray(newLines[:len(newLines)-1])

					// if struct end loop end for this struct
					break

				}
			}

		}
	}
	return newLines
}

// GetEnvValueByKey Read Env file and find value by key
func GetEnvValueByKey(path string, envKey string) string {
	lines, err := ScanLines(path)
	if err != nil {
		fmt.Println("File not found!")
		panic(err)
	}
	for _, line := range lines {
		if StringsContains(line, envKey) {
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

// ScanLines read file and return a string array. Scan file for lines
func ScanLines(path string) ([]string, error) {

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
			if StringsContains(path, ".go") {
				lines, err := ScanLines(path)
				ErrorCheck(err)
				for _, line := range lines {

					if StringsContains(line, keywordToFind) {

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
			if StringsContains(path, ".go") {
				lines, err := ScanLines(path)
				ErrorCheck(err)
				for _, line := range lines {
					// fmt.Println(line)
					if StringsContains(line, oldImportPath) {
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

// ScanDirFindController find controller to generate auto swagger documentation
func ScanDirFindController(dirToScan string, keywordToFind string, routeName string, routeMethodName string) {
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

			if StringsContains(path, ".go") {
				lines, err := ScanLines(path)
				ErrorCheck(err)
				for i, line := range lines {

					if StringsContains(line, keywordToFind) {

						// lineSplit := strings.Split(line, " ")
						// s := lineSplit[1]
						// // if the method public, to detect public methor check the name first letter is capital
						// if s != "DefaultProperties" {
						// 	f.WriteString("\t" + "db.DB.AutoMigrate(models." + s + "{})\n")
						// }

						var newLines []string
						newLines = append(newLines, lines[:i]...)

						found := false // find if already exist the documentation
						for _, nl := range newLines {
							if StringsContains(nl, "@Router "+routeName+" ") { // if the swagger doc already exist then skip
								found = true
								// fmt.Println(nl)
								// fmt.Println(path)
								break
							}
						}

						newLines = append(newLines, "\t//")
						newLines = append(newLines, "\t// @Summary "+routeName+" a api url")
						newLines = append(newLines, "\t// @Description "+routeName+" "+routeMethodName)
						newLines = append(newLines, "\t// @Accept  json")
						newLines = append(newLines, "\t// @Produce  json")
						newLines = append(newLines, "\t// @Success 200 {string} string	\"ok\"")
						newLines = append(newLines, "\t// @Router "+routeName+" ["+routeMethodName+"]")

						if !found {
							newLines = append(newLines, lines[i:]...)
							lines = newLines
							// fmt.Println(lines)
						}

						// fmt.Println(line)
					}
				}
				ioutil.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func TrimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

// NetCheck checking internet connection by ping to github.com and return bool
func NetCheck() bool {
	fmt.Println("checking internet connection....")
	_, err := net.Dial("tcp", "github.com:443")
	if err != nil {
		fmt.Println("No internet connection")
		return false
	}
	return true
}

// AskString Ask user input in Terminal
func AskString(l string, d string) string {
	prompt := promptui.Prompt{
		Label:    l,
		Validate: nil,
		Default:  d,
	}
	str, err := prompt.Run()
	ErrorCheck(err)
	return str
}

// AskSelect Ask to select from suggestion returns index and itemName
func AskSelect(l string, s []string) (int, string) {
	p := promptui.Select{
		Label: l,
		Items: s,
	}
	i, selected, err := p.Run()

	ErrorCheck(err)
	return i, selected

}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// AppendLinesInFileNext  append lines in file params: filepath, lineAfter: multiple Selector, lines to append
func AppendLinesInFileNext(file string, lineAfter []string, offset int, newlines []string) {
	filelines, err := ScanLines(file)
	read, err := ioutil.ReadFile(file)
	ErrorCheck(err)
	// loop the file lines
	exit := false
	// First loop the kewywords to match line by line
	// This way is perfect to controll the priority of which keyword is best priority for matching.
	// suppose to api routes i am looking to append line after a version line / cors line
	for _, la := range lineAfter {

		for i, fl := range filelines {
			// find the line of the page where append the new lines

			if StringsContains(string(fl), la) {

				var newFile []string
				newFile = append(newFile, filelines[:i+1+offset]...)

				for _, ln := range newlines {
					// split the text to remove any spaces and tabs
					// join the text to back as a word
					// then find contains
					if !StringsContains(strings.Join(strings.Fields(string(read)), ""), strings.Join(strings.Fields(ln), "")) {
						newFile = append(newFile, ln)
						// fmt.Println(ln)
					}

				}

				newFile = append(newFile, filelines[i+1+offset:]...)
				filelines = newFile
				// fmt.Println(newFile)
				// fmt.Println(filelines)
				exit = true
				ioutil.WriteFile(file, []byte(strings.Join(filelines, "\n")), 0644)
				break
			}
		}
		if exit {
			break
		}
	}

}

// AppendLinesInFile  append lines in file params: filepath, lineAfter: multipleCheck, lines to append
func AppendLinesInFile(file string, lineAfter []string, newlines []string) {
	filelines, err := ScanLines(file)
	read, err := ioutil.ReadFile(file)
	ErrorCheck(err)
	// loop the file lines
	exit := false
	// First loop the kewywords to match line by line
	// This way is perfect to controll the priority of which keyword is best priority for matching.
	// suppose to api routes i am looking to append line after a version line / cors line
	for _, la := range lineAfter {

		for i, fl := range filelines {
			// find the line of the page where append the new lines

			if StringsContains(string(fl), la) {

				var newFile []string
				newFile = append(newFile, filelines[:i+1]...)

				for _, ln := range newlines {
					// split the text to remove any spaces and tabs
					// join the text to back as a word
					// then find contains
					if !StringsContains(strings.Join(strings.Fields(string(read)), ""), strings.Join(strings.Fields(ln), "")) {
						newFile = append(newFile, ln)
						// fmt.Println(ln)
					}

				}

				newFile = append(newFile, filelines[i+1:]...)
				filelines = newFile
				// fmt.Println(newFile)
				// fmt.Println(filelines)
				exit = true
				ioutil.WriteFile(file, []byte(strings.Join(filelines, "\n")), 0644)
				break
			}
		}
		if exit {
			break
		}
	}

}

// func AppendLinesInFile(file string, lineAfter []string, newlines []string) {
// 	filelines, err := ScanLines(file)
// 	read, err := ioutil.ReadFile(file)
// 	ErrorCheck(err)
// 	// loop the file lines
// 	exit := false
// 	for i, fl := range filelines {
// 		// find the line of the page where append the new lines
// 		for _, la := range lineAfter {

// 			if StringsContains(string(fl), la) {
// 				fmt.Printf(la)
// 				var newFile []string
// 				newFile = append(newFile, filelines[:i+1]...)

// 				for _, ln := range newlines {
// 					// split the text to remove any spaces and tabs
// 					// join the text to back as a word
// 					// then find contains
// 					if !StringsContains(strings.Join(strings.Fields(string(read)), ""), strings.Join(strings.Fields(ln), "")) {
// 						newFile = append(newFile, ln)
// 						// fmt.Println(ln)
// 					}

// 				}

// 				newFile = append(newFile, filelines[i+1:]...)
// 				filelines = newFile
// 				// fmt.Println(newFile)
// 				// fmt.Println(filelines)
// 				exit = true
// 				ioutil.WriteFile(file, []byte(strings.Join(filelines, "\n")), 0644)
// 				break
// 			}
// 		}
// 		if exit {
// 			break
// 		}
// 	}

// }

// AppendLastLine write logs
func AppendLastLine(path string, line string) {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		log.Println(err)
	}
}

// StringsContains to check if string contains in strings
func StringsContains(str string, m string) bool {
	s := strings.Split(str, " ")
	s = CleanEmptyArray(s)
	str = strings.ToLower(strings.Join(s, ""))

	mm := strings.Split(m, " ")
	mm = CleanEmptyArray(mm)
	m = strings.ToLower(strings.Join(mm, ""))

	r := strings.Contains(str, m)

	return r
}
