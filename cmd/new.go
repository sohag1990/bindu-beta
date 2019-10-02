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

	helper "github.com/bindu-bindu/bindu/Helper"
	new "github.com/bindu-bindu/bindu/New"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new project by given projectName",
	Long:  `Generate a new project using given projectName`,
	Run: func(cmd *cobra.Command, args []string) {
		var cli helper.CommandChain
		var flags = []helper.Flag{
			{Key: "app", Values: []string{fmt.Sprintf("%v", cmd.Flag("app").Value)}},
			{Key: "db", Values: []string{fmt.Sprintf("%v", cmd.Flag("db").Value)}},
			{Key: "port", Values: []string{fmt.Sprintf("%v", cmd.Flag("port").Value)}},
		}

		cli = helper.CLI{Args: args, Flags: flags}
		new.New(cmd, cli)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newCmd.Flags().StringP("port", "p", "8080", "Run on specific port number default 8080. eg --p 9999")
	newCmd.Flags().String("app", "Blank", "Prebuilt app name or slug or url\nAvailable Prebuilt app(Blank, Basic Web, Basic Api, Blog, E-Commerce, GRPC Server, GRPC Client)\neg. --app Blank")
	newCmd.Flags().String("db", "", "db info eg. --db AdapterName:HostName:Port:DbName:DbUserName:DbPass")
}
