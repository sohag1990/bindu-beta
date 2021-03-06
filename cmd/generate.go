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
		helper.IsInProjectDir()
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
			{Key: "method", Values: []string{fmt.Sprintf("%v", cmd.Flag("method").Value)}},
			{Key: "methods", Values: []string{fmt.Sprintf("%v", cmd.Flag("methods").Value)}},
			{Key: "middleware", Values: []string{fmt.Sprintf("%v", cmd.Flag("middleware").Value)}},
			{Key: "group", Values: []string{fmt.Sprintf("%v", cmd.Flag("group").Value)}},
		}
		// setter of cli to get from other page

		cli.SetCli(args, flags)
		// Story writter
		// if the command execute return true,
		// so the story can know that command was success or failed
		story.WriteStory("generate", cli)
		generate.Generator(cmd, cli)

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().String("hasOne", "", "HasOne Relationship")
	generateCmd.Flags().String("hasMany", "", "HasMany Relationship")
	generateCmd.Flags().String("manyToMany", "", "HasOne Relationship")
	generateCmd.Flags().String("hasOneThrough", "", "Has One Through Relationship")
	generateCmd.Flags().String("hasManyThrough", "", "Has Many Through Relationship")
	generateCmd.Flags().String("belongsTo", "", "Belongs to other model")
	generateCmd.Flags().String("method", "", "Get-Post-Put-Delete methods available, All is shortcut for all methods")
	generateCmd.Flags().String("methods", "", "Get-Post-Put-Delete methods available, All is shortcut for all methods")
	generateCmd.Flags().BoolP("update", "u", false, "Skip go get command. -u=true")
	generateCmd.Flags().StringP("middleware", "m", "", "Define middleware which. -a=false")
	generateCmd.Flags().StringP("group", "g", "", "API Group: Define a api group. eg. v1 or v2 etc, -g=v2")

}
