package ai

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yohanc3/resumemaxxer/internal/config"
	resumebuilderparser "github.com/yohanc3/resumemaxxer/internal/resume-builder-parser"
)

func TestResumeParsing(t *testing.T){

	resume_path := filepath.Join("testdata", "resume_test.pdf")

	info, err := os.Stat(resume_path)
	if err != nil {
		t.Fatalf("cannot access resume pdf: %v", err)
	}

	if info.Size() == 0 {
		t.Fatalf("resume pdf is empty")
	}

	resumeText, err := resumebuilderparser.ExtractPDFText(resume_path)
	if err != nil {
		t.Fatalf("error when extracting pdf text: %v", err)
	}

	if err := config.LoadConfig(); err != nil {
		t.Fatal("error when loading env variables")	
	}

	println(resumeText)

	// deepseekclient, err := NewDeepseekClient("https://api.deepseek.com", config.Cfg.DeepSeekApiKey)
	// if err != nil {
	// 	t.Errorf("error when requesting new deepseek client: %w", err)
	// }
	// 

}

