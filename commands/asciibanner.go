package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//DrawASCIIBanner draws a little ASCII banner to make the Goat logs easier to look at
func DrawASCIIBanner(headLine string, debug bool) string {
	if debug {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press enter to continue: ")
		reader.ReadString('\n')
	}

	return fmt.Sprintf("\n%[1]s\n# %[2]s #\n%[1]s\n",
		strings.Repeat("#", len(headLine)+4),
		headLine)
}
