package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/C0d3-5t3w/GTS/internal/config"
	"github.com/C0d3-5t3w/GTS/pkg/compiler"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	verbose    = flag.Bool("verbose", false, "Enable verbose output")
	output     = flag.String("o", "", "Output path for compiled binary")
	skipTS     = flag.Bool("skip-ts", false, "Skip TypeScript compilation")
	skipSCSS   = flag.Bool("skip-scss", false, "Skip SCSS compilation")
	skipPHP    = flag.Bool("skip-php", false, "Skip PHP to HTML conversion")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: No Go source files or commands specified")
		printUsage()
		os.Exit(1)
	}

	cfgPath := *configPath
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath()
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		fmt.Printf("Warning: Could not load config file: %v\nUsing default settings\n", err)
		cfg = &config.CompilerConfig{
			GoCompilerPath: "go",
		}
	}

	if *verbose {
		cfg.OutputOptions.VerboseOutput = true
	}
	if *output != "" {
		args = append(args, "-o", *output)
	}
	if *skipTS {
		cfg.TypeScript.Enabled = false
	}
	if *skipSCSS {
		cfg.SCSS.Enabled = false
	}
	if *skipPHP {
		cfg.PHP.Enabled = false
	}

	compiler := compiler.NewExtendedCompiler(cfg)

	if err := compiler.RunPrePasses(); err != nil {
		fmt.Printf("Error in pre-compilation pass: %v\n", err)
		os.Exit(1)
	}

	if err := compiler.Compile(args); err != nil {
		fmt.Printf("Compilation failed: %v\n", err)
		os.Exit(1)
	}

	if err := compiler.RunPostPasses(); err != nil {
		fmt.Printf("Error in post-compilation pass: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Compilation completed successfully")
}

func printUsage() {
	fmt.Printf(`Extended Go Compiler

Usage: %s [options] [build flags] [packages]

Options:
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Println("\nFor more information, visit: https://github.com/brandonstewart/GTS")
}
