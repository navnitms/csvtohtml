package converter

import (
	"os"
	"fmt"
	"time"
	"encoding/csv"
	"path/filepath"
	"regexp"
	"html/template"
	"bytes"
	"strings"
	"github.com/pkg/browser"
)

var (
	packagePath, _ = os.Getwd()
	templatesDir   = filepath.Join(packagePath, "templates")
	jsSrcPattern   = regexp.MustCompile(`<script.*?src=\"(.*?)\".*?<\/script>`)
	jsFilesPath    = filepath.Join(packagePath, "templates")
	templateFile   = "template.txt"
)

func Convert(inputFileName string, options map[string]interface{}) (string, error) {
	delimiter := getStringOption(options, "delimiter", ",")
	// quotechar := getStringOption(options, "quotechar", "\"")
	
	file, err := os.Open(inputFileName)
	if err != nil {
		return "", fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	reader.Comma = rune(delimiter[0])
	reader.LazyQuotes = true
	
	var csvHeaders []string
	var csvRows [][]string
	
	fmt.Println(options["title"])
	if options["no_header"] == false {
		csvHeaders, err = reader.Read()
		if err != nil {
			return "", fmt.Errorf("failed to read CSV header: %v", err)
		}
	}

	csvRows, err = reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read CSV rows: %v", err)
	}
	html, err := renderTemplate(csvHeaders, csvRows, options)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %v", err)
	}
	
	return replaceScript(html) , nil
}


func renderTemplate(headers []string, rows [][]string, options map[string]interface{}) (string, error) {
	tmplPath := filepath.Join(templatesDir, templateFile)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}
	

	data := map[string]interface{}{
		"Title":            getStringOption(options, "title", "Table"),
		"Headers":          headers,
		"Rows":             rows, // Now rows is [][]string
		"Pagination":       options["pagination"],
		"TableHeight":      getStringOption(options, "height", "70vh"),
		"DisplayLength":    getIntOption(options, "display_length", -1),
	}

	fmt.Print(data["Title"])
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %v", err)
	}

	return buf.String(), nil
}




func replaceScript(html string) string {
    matches := jsSrcPattern.FindAllStringSubmatch(html, -1)

    if len(matches) == 0 {
        return html
    }

    for _, match := range matches {
        src := match[1]
        
        if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
            continue
        }

        filePath := filepath.Join(jsFilesPath, src)
        fileContent, err := os.ReadFile(filePath)
        if err != nil {
			fmt.Println(fmt.Errorf("failed to read JS file: %v", err))
        }

        jsContent := fmt.Sprintf("<script type=\"text/javascript\">%s</script>", fileContent)
        html = strings.Replace(html, match[0], jsContent, 1)
    }

    return html
}

// Save saves content to a file
func Save(fileName, content string) error {
	return os.WriteFile(fileName, []byte(content), 0644)
}

// Serve serves the content in a temporary HTML file and opens it in a browser
func Serve(content string) {
	tmpFile, err := os.CreateTemp("", "csvtotable_*.html")
	if err != nil {
		fmt.Printf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		fmt.Printf("Failed to write to temp file: %v", err)
	}

	browser.OpenFile(tmpFile.Name())

	for {
		time.Sleep(time.Second)
	}
}