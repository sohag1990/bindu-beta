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
	"net"
	"os/exec"

	helper "github.com/bindu-bindu/bindu/Helper"
	story "github.com/bindu-bindu/bindu/Story"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		helper.IsInProjectDir()
		// to sanitize commands
		var cli helper.CommandChain
		cli = helper.InitialCli()
		cli.SetCliArgs(args)

		fmt.Println("checking internet connection....")
		_, err := net.Dial("tcp", "github.com:443")
		helper.ErrorCheck(err)
		fmt.Println("Bindu installing app, Wait.....")
		out, err := exec.Command("go", "get").Output()
		fmt.Println("Output: ", out, err)

		// Story writter
		// if the command execute return true,
		// so the story can know that command was success or failed
		story.WriteStory("install", cli)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	// installCmd.Flags().StringP("migration", "m", "", "Set your name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
