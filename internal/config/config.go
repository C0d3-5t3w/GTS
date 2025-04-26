package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CompilerConfig struct {
	GoCompilerPath string            `yaml:"goCompilerPath"`
	DefaultFlags   []string          `yaml:"defaultFlags"`
	Extensions     []string          `yaml:"extensions"`
	CustomPasses   []CompilerPass    `yaml:"customPasses"`
	OutputOptions  OutputOptions     `yaml:"outputOptions"`
	DebugMode      bool              `yaml:"debugMode"`
	EnvVars        map[string]string `yaml:"envVars"`
	TypeScript     TypeScriptConfig  `yaml:"typescript"`
	SCSS           SCSSConfig        `yaml:"scss"`
	PHP            PHPConfig         `yaml:"php"`
}

type CompilerPass struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
	Enabled bool     `yaml:"enabled"`
}

type OutputOptions struct {
	Directory       string `yaml:"directory"`
	VerboseOutput   bool   `yaml:"verboseOutput"`
	GenerateReports bool   `yaml:"generateReports"`
}

type TypeScriptConfig struct {
	Enabled bool     `yaml:"enabled"`
	TscPath string   `yaml:"tscPath"`
	SrcDir  string   `yaml:"srcDir"`
	OutDir  string   `yaml:"outDir"`
	Options []string `yaml:"options"`
}

type SCSSConfig struct {
	Enabled  bool     `yaml:"enabled"`
	SassPath string   `yaml:"sassPath"`
	SrcDir   string   `yaml:"srcDir"`
	OutDir   string   `yaml:"outDir"`
	Options  []string `yaml:"options"`
}

type PHPConfig struct {
	Enabled bool     `yaml:"enabled"`
	PhpPath string   `yaml:"phpPath"`
	SrcDirs []string `yaml:"srcDirs"`
	Options []string `yaml:"options"`
}

func LoadConfig(configPath string) (*CompilerConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config CompilerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.GoCompilerPath == "" {
		config.GoCompilerPath = "go"
	}

	if config.TypeScript.TscPath == "" {
		config.TypeScript.TscPath = "tsc"
	}

	if config.SCSS.SassPath == "" {
		config.SCSS.SassPath = "sass"
	}

	if config.PHP.PhpPath == "" {
		config.PHP.PhpPath = "php"
	}

	return &config, nil
}

func DefaultConfigPath() string {

	if _, err := os.Stat("./pkg/config/config.yaml"); err == nil {
		return "./pkg/config/config.yaml"
	}

	defaultPath := filepath.Join(os.Getenv("HOME"), ".GTS", "config.yaml")
	return defaultPath
}
