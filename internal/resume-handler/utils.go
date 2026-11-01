package resumebuilderparser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// Extracts pdf text based on a path.
func ExtractPDFText(file *os.File) (string, error){

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error when reading bytes from file: %w", err)
	}

	reader, err := pdf.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return "", err
	}

	var text strings.Builder

	for pageNumber := 1; pageNumber <= reader.NumPage(); pageNumber++ {

		page := reader.Page(pageNumber)
		if page.V.IsNull(){
			continue
		}

		pageText, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}

		text.WriteString("\n" + pageText)
	}

	return text.String(), nil

}

