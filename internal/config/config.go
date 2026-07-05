package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)


type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}


type Blog struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
	URL   string `yaml:"url"`
	Body  string `yaml:"body"`
}


type Link struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
}

type Content struct {
	Name        string    `yaml:"name"`
	Tagline     string    `yaml:"tagline"`
	About       string    `yaml:"about"`
	LastUpdated string    `yaml:"last_updated"`
	Projects    []Project `yaml:"projects"`
	Blogs       []Blog    `yaml:"blogs"`
	Links       []Link    `yaml:"links"`


	ASCIIArts []string `yaml:"-"`
}

 
func Load(path string) (*Content, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read content file: %w", err)
	}

	var c Content
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse content file: %w", err)
	}

	if c.Name == "" {
		c.Name = "Pr3thiv"
	}

	asciiDir := filepath.Join(filepath.Dir(path), "assets", "ascii")
	c.ASCIIArts = loadASCIIArts(asciiDir)

	return &c, nil
}

 
func loadASCIIArts(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".txt") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	arts := make([]string, 0, len(names))
	for _, n := range names {
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			continue
		}
		art := strings.TrimRight(string(b), "\r\n")
		if strings.TrimSpace(art) != "" {
			arts = append(arts, art)
		}
	}
	return arts
}
