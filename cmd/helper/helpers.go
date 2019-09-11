package helper

import (
	"fmt"
	"os"
	"reflect"
	"strings"

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
