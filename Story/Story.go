package story

import (
	"fmt"
	"io/ioutil"
	"strings"

	helper "github.com/bindu-bindu/bindu/Helper"
	"github.com/spf13/cobra"
)

// Story of commands that used in entire project development
func Story(cmd *cobra.Command, c helper.CommandChain) {
	// User command args
	// args := c.GetArgs()
	fmt.Println("Bindu telling you a nice story about project\n")
	path := "./bindu/story.sh"
	read, _ := ioutil.ReadFile(path)
	if len(read) == 0 {
		fmt.Println("Nothing found ^)^")
	}
	fmt.Println(string(read))
}

// WriteStory Write Command Story
func WriteStory(cName string, args []string, flags []helper.Flag) {
	path := "./bindu/story.sh"
	helper.AppendLastLine(path, "bindu "+cName+" "+strings.Join(args, " "))
}
