// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	helper "github.com/bindu-bindu/bindu/Helper"
	modify "github.com/bindu-bindu/bindu/Modify"
	story "github.com/bindu-bindu/bindu/Story"
	"github.com/spf13/cobra"
)

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var cli helper.CommandChain
		cli = helper.InitialCli()
		var flags = []helper.Flag{
			{Key: "addGorm", Values: []string{fmt.Sprintf("%v", cmd.Flag("addGorm").Value)}},
			{Key: "removeGorm", Values: []string{fmt.Sprintf("%v", cmd.Flag("removeGorm").Value)}},
		}
		cli.SetCli(args, flags)
		// Story writter
		// if the command execute return true,
		// so the story can know that command was success or failed
		story.WriteStory("modify", cli)
		modify.Modifier(cmd, cli)
	},
}

func init() {
	rootCmd.AddCommand(modifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	modifyCmd.Flags().String("addGorm", "addGorm", "Add new gorm option/s")
	modifyCmd.Flags().String("removeGorm", "addGorm", "Remove any gorm option/s")
}
