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
func WriteStory(cName string, cli helper.CommandChain, status bool) {

	args := cli.GetArgs()
	flags := cli.GetFlags()

	line := strings.Join(args, " ") + " "
	fl := ""
	for _, f := range flags {

		for _, fv := range f.Values {
			// if the flag has value then add in story line
			if len(fv) > 0 {
				fl = fl + "--" + f.Key + " " + fv + " "
			}
		}
	}
	line = line + fl + "#" + fmt.Sprintf("%v", status)
	path := "./bindu/story.sh"
	helper.AppendLastLine(path, "bindu "+cName+" "+line)
}
