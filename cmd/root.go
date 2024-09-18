// Description: This file contains the root command for the CLI application.
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/navnitms/csvtohtml/converter"
	"github.com/spf13/cobra"
)

var title, delimiter, quotechar, height string
var displayLength int
var overwrite, serve, pagination, noHeader, export bool

var rootCmd = &cobra.Command{
	Use:   "csvtohtml",
	Short: "CLI Application to convert CSV file to HTML table",
	Long: `csvtohtml is a CLI tool designed to convert CSV files into HTML tables with ease. 
It provides various customization options for generating the HTML table, including control 
over the table's appearance and features such as pagination, virtual scrolling, and export functionality.

You can use this application by specifying the input CSV file and optionally the output HTML file.

Examples:
  Convert a CSV file to an HTML table:
    csvtohtml input.csv output.html

  Convert and display the HTML table directly in a browser:
    csvtohtml input.csv --serve
    
    Flags for customization:
      -t, --title:        	Set a title for the HTML table.
      -d, --delimiter:      Specify the delimiter used in the CSV file. Default is ','.
      -q, --quotechar:      Specify the character used for quoting fields with special characters. Default is '"'.
      -dl, --display-length: Set the number of rows to display by default. Default is -1 (show all rows).
      -o, --overwrite:      Overwrite the output file if it already exists.
      -s, --serve:          Open the generated HTML in the browser instead of saving to a file.
      -h, --height:         Set the table height (in px or %) for the generated HTML.
      -p, --pagination:     Enable or disable pagination for the table. Default is enabled.
      -nh, --no-header:     Disable treating the first row as headers.
      -e, --export:         Enable export options for filtered rows. Default is enabled.`,

	Run: func(cmd *cobra.Command, args []string) {
			inputFile := args[0]
			outputFile := ""
			if len(args) > 1 {
				outputFile = args[1]
			}

			options := map[string]interface{}{
				"title":          title,
				"delimiter":      delimiter,
				"quotechar":      quotechar,
				"display_length": displayLength,
				"height":         height,
				"pagination":     pagination,
				"no_header":      noHeader,
				"export":         export,
			}

			// Convert CSV to HTML content
			content, err := converter.Convert(inputFile, options)
			if err != nil {
				log.Fatalf("Failed to convert: %v", err)
			}

			if serve {
				// Serve the HTML in the browser
				converter.Serve(content)
			} else if outputFile != "" {
				// Check if file should be overwritten
				if !overwrite && !promptOverwrite(outputFile) {
					log.Fatal("Operation aborted.")
				}
				// Save the output
				err := converter.Save(outputFile, content)
				if err != nil {
					log.Fatalf("Failed to save file: %v", err)
				}
				fmt.Printf("File converted successfully: %s\n", outputFile)
			} else {
				log.Fatal("Missing argument \"output_file\".")
			}
		},
}


func init() {
rootCmd.Flags().StringVarP(&title, "title", "t", "", "Table title")
rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "CSV delimiter")
rootCmd.Flags().StringVarP(&quotechar, "quotechar", "q", `"`, "String used to quote fields containing special characters")
rootCmd.Flags().IntVarP(&displayLength, "length", "l", -1, "Number of rows to show by default. Defaults to -1 (show all rows)")
rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite the output file if exists")
rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "Open output HTML in browser instead of writing to file")
rootCmd.Flags().StringVarP(&height, "height", "H", "", "Table height in px or %")  // Changed shorthand to "H"
rootCmd.Flags().BoolVarP(&pagination, "pagination", "p", true, "Enable/disable table pagination")
rootCmd.Flags().BoolVar(&noHeader, "no-header", false, "Disable displaying first row as headers")
rootCmd.Flags().BoolVarP(&export, "export", "e", true, "Enable filtered rows export options")
}


func promptOverwrite(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("File (%s) already exists. Do you want to overwrite? (y/n): ", fileName)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "y" && input != "Y" {
		return false
	}

	return true
}

func Execute() error {
	return rootCmd.Execute()
}