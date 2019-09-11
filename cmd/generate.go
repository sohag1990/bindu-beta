/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bindu-bindu/bindu/cmd/helper"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		helper.IsInProjectDir()
		// sanitize user args array
		args = helper.SanitizeUserInputReArry(args)
		// fmt.Printf("%v\n", args)

		genItems := []string{"Model", "Controller", "Scaffold", "View"}
		// fmt.Printf("%v", args)
		if len(args) == 0 {
			prompt := promptui.Select{
				Label: "Select To Generate",
				Items: genItems,
			}
			selectedIndex, _, err := prompt.Run()
			helper.ErrorCheck(err)
			args = append(args, genItems[selectedIndex])
			// fmt.Println("Item Selected: ", genItems[selectedIndex])
		}

		// check the first items is matched with predefined arg.. Controller, Model,,,etc
		i, found := helper.ArrayFind(genItems, args[0])
		if found {
			switch i {
			case 0:
				fmt.Printf("%v\n", "Model action")
				ModelGenerator(args)

			case 1:
				fmt.Printf("%v\n", "Controller action")
				fmt.Printf("%v\n", args)
			case 2:
				fmt.Printf("%v\n", "Scaffold action")
				fmt.Printf("%v\n", args)
			case 3:
				fmt.Printf("%v\n", "View action")
				fmt.Printf("%v\n", args)
			default:
				fmt.Printf("%v\n", "No action taken")
				fmt.Printf("%v\n", args)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("controller", "c", false, "Help message for Generate controller")

	// generateCmd.Flags().StringP("controller", "c", "", "Enter new controller name")
}

// ModelGenerator to generate the model in project using the user inputs
func ModelGenerator(args []string) {
	fmt.Println(args)
	if len(args) == 2 {
		fmt.Println("Initializing empty Model " + args[1])
		newpath := filepath.Join(".", "app/models")
		os.MkdirAll(newpath, os.ModePerm)
		f, err := os.Create("./app/models/" + args[1] + ".go")
		helper.ErrorCheck(err)
		f.WriteString("package model\n")
		defer f.Close()
		absPath, _ := filepath.Abs(newpath)
		fmt.Println(absPath + "/" + args[1] + ".go")
	}
}
