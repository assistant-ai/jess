package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

func PrintContextIDs(contextIds []string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "+-----------+")
	fmt.Fprintln(w, "| Context ID\t|")

	for _, contextId := range contextIds {
		fmt.Fprintf(w, "| %s\t|\n", contextId)
	}
	fmt.Fprintln(w, "+-----------+")

	w.Flush()
}

func PrintErrorAndExit(err error) {
	// Print the error message to the standard error
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// Provide guidance on how to resolve the issue
	fmt.Fprintln(os.Stderr, "Please check your input and try again.")
	// Exit the program with a non-zero status code to indicate an error
	os.Exit(1)
}

func HandleError(err error) {
	if err != nil {
		cli.Exit(err, 1)
		panic(err)
	}
}
