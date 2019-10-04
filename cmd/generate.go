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

	generate "github.com/bindu-bindu/bindu/Generate"
	helper "github.com/bindu-bindu/bindu/Helper"
	story "github.com/bindu-bindu/bindu/Story"
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

		// to sanitize commands
		var cli helper.CommandChain
		cli = helper.InitialCli()
		var flags = []helper.Flag{
			{Key: "hasOne", Values: []string{fmt.Sprintf("%v", cmd.Flag("hasOne").Value)}},
			{Key: "hasMany", Values: []string{fmt.Sprintf("%v", cmd.Flag("hasMany").Value)}},
			{Key: "manyToMany", Values: []string{fmt.Sprintf("%v", cmd.Flag("manyToMany").Value)}},
			{Key: "hasOneThrough", Values: []string{fmt.Sprintf("%v", cmd.Flag("hasOneThrough").Value)}},
			{Key: "hasManyThrough", Values: []string{fmt.Sprintf("%v", cmd.Flag("hasManyThrough").Value)}},
			{Key: "belongsTo", Values: []string{fmt.Sprintf("%v", cmd.Flag("belongsTo").Value)}},
		}
		//Story writter
		story.WriteStory("generate", args, flags)
		// setter of cli to get from other page
		cli.SetCli(args, flags)
		generate.Generator(cmd, cli)

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

	generateCmd.Flags().String("hasOne", "", "HasOne Relationship")
	generateCmd.Flags().String("hasMany", "", "HasMany Relationship")
	generateCmd.Flags().String("manyToMany", "", "HasOne Relationship")
	generateCmd.Flags().String("hasOneThrough", "", "Has One Through Relationship")
	generateCmd.Flags().String("hasManyThrough", "", "Has Many Through Relationship")
	generateCmd.Flags().String("belongsTo", "", "Belongs to other model")
	generateCmd.Flags().BoolP("update", "u", false, "Skip go get command")

}
