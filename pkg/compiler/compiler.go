package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/C0d3-5t3w/GTS/internal/config"
)

type ExtendedCompiler struct {
	config *config.CompilerConfig
}

func NewExtendedCompiler(cfg *config.CompilerConfig) *ExtendedCompiler {
	return &ExtendedCompiler{
		config: cfg,
	}
}

func (c *ExtendedCompiler) RunPrePasses() error {
	for _, pass := range c.config.CustomPasses {
		if pass.Type == "pre" && pass.Enabled {
			if c.config.OutputOptions.VerboseOutput {
				fmt.Printf("Running pre-compilation pass: %s\n", pass.Name)
			}

			cmd := exec.Command(pass.Command, pass.Args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Env = os.Environ()
			for k, v := range c.config.EnvVars {
				cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
			}

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("pre-pass %s failed: %w", pass.Name, err)
			}
		}
	}
	return nil
}

func (c *ExtendedCompiler) Compile(args []string) error {

	if c.config.TypeScript.Enabled {
		if err := c.CompileTypeScript(); err != nil {
			return fmt.Errorf("TypeScript compilation failed: %w", err)
		}
	}

	if c.config.SCSS.Enabled {
		if err := c.CompileSCSS(); err != nil {
			return fmt.Errorf("SCSS compilation failed: %w", err)
		}
	}

	if c.config.PHP.Enabled {
		if err := c.ConvertPHPToHTML(); err != nil {
			return fmt.Errorf("PHP to HTML conversion failed: %w", err)
		}
	}

	extendedArgs := c.applyExtensions(args)

	for _, flag := range c.config.DefaultFlags {
		if !containsFlag(extendedArgs, flag) {
			extendedArgs = append([]string{flag}, extendedArgs...)
		}
	}

	if !isCommand(extendedArgs) {
		extendedArgs = append([]string{"build"}, extendedArgs...)
	}

	if c.config.OutputOptions.VerboseOutput {
		fmt.Printf("Executing: %s %s\n", c.config.GoCompilerPath, strings.Join(extendedArgs, " "))
	}

	cmd := exec.Command(c.config.GoCompilerPath, extendedArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()
	for k, v := range c.config.EnvVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return cmd.Run()
}

func (c *ExtendedCompiler) CompileTypeScript() error {
	if c.config.OutputOptions.VerboseOutput {
		fmt.Println("Compiling TypeScript files...")
	}

	if err := os.MkdirAll(c.config.TypeScript.OutDir, 0755); err != nil {
		return fmt.Errorf("failed to create TypeScript output directory: %w", err)
	}

	args := append(c.config.TypeScript.Options,
		"--outDir", c.config.TypeScript.OutDir,
		c.config.TypeScript.SrcDir)

	cmd := exec.Command(c.config.TypeScript.TscPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if c.config.OutputOptions.VerboseOutput {
		fmt.Printf("Executing: %s %s\n", c.config.TypeScript.TscPath, strings.Join(args, " "))
	}

	return cmd.Run()
}

func (c *ExtendedCompiler) CompileSCSS() error {
	if c.config.OutputOptions.VerboseOutput {
		fmt.Println("Compiling SCSS files...")
	}

	if err := os.MkdirAll(c.config.SCSS.OutDir, 0755); err != nil {
		return fmt.Errorf("failed to create SCSS output directory: %w", err)
	}

	srcPath := filepath.Join(c.config.SCSS.SrcDir, "*.scss")

	args := append(c.config.SCSS.Options, srcPath, c.config.SCSS.OutDir)

	cmd := exec.Command(c.config.SCSS.SassPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if c.config.OutputOptions.VerboseOutput {
		fmt.Printf("Executing: %s %s\n", c.config.SCSS.SassPath, strings.Join(args, " "))
	}

	return cmd.Run()
}

func (c *ExtendedCompiler) RunPostPasses() error {
	for _, pass := range c.config.CustomPasses {
		if pass.Type == "post" && pass.Enabled {
			if c.config.OutputOptions.VerboseOutput {
				fmt.Printf("Running post-compilation pass: %s\n", pass.Name)
			}

			cmd := exec.Command(pass.Command, pass.Args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Env = os.Environ()
			for k, v := range c.config.EnvVars {
				cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
			}

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("post-pass %s failed: %w", pass.Name, err)
			}
		}
	}
	return nil
}

func (c *ExtendedCompiler) applyExtensions(args []string) []string {
	extendedArgs := make([]string, len(args))
	copy(extendedArgs, args)

	for _, ext := range c.config.Extensions {
		switch ext {
		case "custom-import":

			extendedArgs = append(extendedArgs, "-gcflags=all=-importmap=oldpkg=newpkg")
		case "enhanced-generics":

			extendedArgs = append(extendedArgs, "-tags=enhancedgenerics")
		case "code-contracts":

			extendedArgs = append(extendedArgs, "-tags=contracts")
		}
	}

	return extendedArgs
}

func containsFlag(args []string, flag string) bool {
	flagName := strings.Split(flag, "=")[0]
	for _, arg := range args {
		if strings.HasPrefix(arg, flagName) {
			return true
		}
	}
	return false
}

func isCommand(args []string) bool {
	commands := []string{"build", "clean", "doc", "env", "bug", "fix", "fmt",
		"generate", "get", "install", "list", "mod", "run", "test", "tool", "version", "vet"}

	if len(args) == 0 {
		return false
	}

	for _, cmd := range commands {
		if args[0] == cmd {
			return true
		}
	}
	return false
}

func (c *ExtendedCompiler) ConvertPHPToHTML() error {
	if c.config.OutputOptions.VerboseOutput {
		fmt.Println("Converting PHP files to HTML...")
	}

	for _, dir := range c.config.PHP.SrcDirs {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".php") {
				outputPath := strings.TrimSuffix(path, ".php") + ".html"

				if c.config.OutputOptions.VerboseOutput {
					fmt.Printf("Converting %s to %s\n", path, outputPath)
				}

				cmd := exec.Command(c.config.PHP.PhpPath, append(c.config.PHP.Options, path)...)
				output, err := cmd.Output()
				if err != nil {
					return fmt.Errorf("PHP conversion failed for %s: %w", path, err)
				}

				if err := os.WriteFile(outputPath, output, 0644); err != nil {
					return fmt.Errorf("failed to write HTML file %s: %w", outputPath, err)
				}
			}
			return nil
		}); err != nil {
			return fmt.Errorf("PHP to HTML conversion failed: %w", err)
		}
	}

	return nil
}
