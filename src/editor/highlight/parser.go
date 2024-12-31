package highlight

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type SyntaxPatterns struct {
	Keywords struct {
		Patterns []string `toml:"patterns"`
		Color    string   `toml:"color"`
	} `toml:"keywords"`
	Types struct {
		Patterns []string `toml:"patterns"`
		Color    string   `toml:"color"`
	} `toml:"types"`
	Comments struct {
		Patterns []string `toml:"patterns"`
		Color    string   `toml:"color"`
	} `toml:"comments"`
	Strings struct {
		Patterns []string `toml:"patterns"`
		Color    string   `toml:"color"`
	} `toml:"strings"`
	Numbers struct {
		Patterns []string `toml:"patterns"`
		Color    string   `toml:"color"`
	} `toml:"numbers"`
	Operators struct {
		Patterns []string `toml:"patterns"`
		Color    string   `toml:"color"`
	} `toml:"operators"`
	Methods struct {
		Patterns []string `toml:"patterns"`
		Color    string   `toml:"color"`
	} `toml:"methods"`
}

func LoadSyntaxPatterns(filePath string) (*SyntaxPatterns, error) {
	var patterns SyntaxPatterns
	if _, err := toml.DecodeFile(filePath, &patterns); err != nil {
		return nil, err
	}
	return &patterns, nil
}

func LoadAllPatterns() map[string]*SyntaxPatterns {
	patternsMap := make(map[string]*SyntaxPatterns)

	files, err := os.ReadDir("src/editor/highlight/langs")
	if err != nil {
		log.Fatalf("Failed to read langs directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			patterns, err := LoadSyntaxPatterns("src/editor/highlight/langs/" + file.Name())
			if err != nil {
				log.Printf("Failed to load patterns from %s: %v", file.Name(), err)
				continue
			}
			patternsMap[file.Name()] = patterns
		}
	}
	return patternsMap
}
