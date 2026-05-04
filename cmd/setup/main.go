package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
)

var cli struct {
	Apply  bool `help:"Apply changes"`
	DryRun bool `help:"Show changes only (default)"`
}

type Change struct {
	Kind   string `json:"kind"`
	File   string `json:"file"`
	Detail string `json:"detail"`
}

func (c Change) Output() {
	b, _ := json.Marshal(c)
	fmt.Println(string(b))
}

func promptInput(promptMsg string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", promptMsg, defaultVal)
	} else {
		fmt.Printf("%s: ", promptMsg)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" && defaultVal != "" {
		return defaultVal
	}
	return input
}

type Inputs struct {
	AppName     string
	ModulePath  string
	EnvName     string
	XdgPath     string
	PackageName string
}

func gatherInputs() Inputs {
	rawAppName := promptInput("Application name (binary name)", "myapp")
	appName := strings.ToLower(rawAppName)
	re := regexp.MustCompile(`[^a-z0-9-]`)
	appName = re.ReplaceAllString(appName, "-")
	appName = strings.Trim(appName, "-")

	if appName == "" {
		fmt.Fprintln(os.Stderr, "Error: application name is required")
		os.Exit(1)
	}

	githubUser := promptInput("GitHub username", "username")
	defaultModule := fmt.Sprintf("github.com/%s/%s", githubUser, appName)
	modulePath := promptInput("Go module path", defaultModule)

	if modulePath == "" {
		fmt.Fprintln(os.Stderr, "Error: module path is required")
		os.Exit(1)
	}

	envName := strings.ToUpper(appName)
	reEnv := regexp.MustCompile(`[^A-Z0-9]`)
	envName = reEnv.ReplaceAllString(envName, "_")
	envName = strings.Trim(envName, "_")

	xdgPath := appName

	rawPkg := promptInput("Root package name (e.g. config.go)", appName)
	rePkg := regexp.MustCompile(`[^a-z0-9_]`)
	packageName := rePkg.ReplaceAllString(rawPkg, "_")
	packageName = strings.Trim(packageName, "_")

	return Inputs{
		AppName:     appName,
		ModulePath:  modulePath,
		EnvName:     envName,
		XdgPath:     xdgPath,
		PackageName: packageName,
	}
}

func replaceInFile(filePath, old, new, desc string, dryRun bool) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if !strings.Contains(string(content), old) {
		return nil
	}

	Change{Kind: "change", File: filePath, Detail: desc}.Output()

	if !dryRun {
		newContent := strings.ReplaceAll(string(content), old, new)
		if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

func renameDir(oldPath, newPath string, dryRun bool) error {
	Change{Kind: "rename", File: oldPath, Detail: "→ " + newPath}.Output()

	if !dryRun {
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	ctx := kong.Parse(&cli)
	_ = ctx

	dryRun := cli.DryRun || !cli.Apply

	inputs := gatherInputs()

	oldModule := "github.com/nurazon59/go-template"
	oldApp := "go-template"
	oldEnv := "GO_TEMPLATE_CONFIG"
	oldXdg := "go-template"
	oldPkg := "template"

	rootDir, _ := filepath.Abs(".")

	replaceInFile(filepath.Join(rootDir, "go.mod"), oldModule, inputs.ModulePath, "update module path", dryRun)

	oldCmdDir := filepath.Join(rootDir, "cmd", "template")
	newCmdDir := filepath.Join(rootDir, "cmd", inputs.AppName)

	replaceInFile(filepath.Join(oldCmdDir, "main.go"), oldModule+"/cmd", inputs.ModulePath+"/cmd", "update import path in main.go", dryRun)
	replaceInFile(filepath.Join(oldCmdDir, "main_test.go"), oldApp, inputs.AppName, "update binary name in test", dryRun)
	renameDir(oldCmdDir, newCmdDir, dryRun)

	runFile := filepath.Join(rootDir, "cmd", "run.go")
	replaceInFile(runFile, oldModule, inputs.ModulePath, "update import path in run.go", dryRun)
	replaceInFile(runFile, oldPkg+" \"", inputs.PackageName+" \"", "update import alias in run.go", dryRun)
	replaceInFile(runFile, oldEnv, inputs.EnvName, "update env var name", dryRun)
	replaceInFile(runFile, oldXdg, inputs.XdgPath, "update XDG path", dryRun)
	replaceInFile(runFile, fmt.Sprintf("kong.Name(\"%s\")", oldApp), fmt.Sprintf("kong.Name(\"%s\")", inputs.AppName), "update kong.Name", dryRun)
	replaceInFile(runFile, oldPkg+".Load", inputs.PackageName+".Load", "update Load call in run.go", dryRun)

	replaceInFile(filepath.Join(rootDir, "config.go"), "package "+oldPkg, "package "+inputs.PackageName, "update package name in config.go", dryRun)

	replaceInFile(filepath.Join(rootDir, "config_test.go"), "package "+oldPkg, "package "+inputs.PackageName, "update package name in config_test.go", dryRun)

	replaceInFile(filepath.Join(rootDir, "Taskfile.yml"), oldApp, inputs.AppName, "update binary name in Taskfile", dryRun)
	replaceInFile(filepath.Join(rootDir, "Taskfile.yml"), "cmd/template", "cmd/"+inputs.AppName, "update build path in Taskfile", dryRun)

	replaceInFile(filepath.Join(rootDir, ".goreleaser.yaml"), "cmd/template", "cmd/"+inputs.AppName, "update main path in goreleaser", dryRun)

	replaceInFile(filepath.Join(rootDir, "README.md"), oldApp, inputs.AppName, "update app name in README", dryRun)
	replaceInFile(filepath.Join(rootDir, "README.md"), oldEnv, inputs.EnvName, "update env var name in README", dryRun)
	replaceInFile(filepath.Join(rootDir, "README.md"), oldXdg, inputs.XdgPath, "update XDG path in README", dryRun)
	replaceInFile(filepath.Join(rootDir, "README.md"), "cmd/template", "cmd/"+inputs.AppName, "update path in README", dryRun)

	oldBin := filepath.Join(rootDir, "bin", "template")
	newBin := filepath.Join(rootDir, "bin", inputs.AppName)
	if _, err := os.Stat(oldBin); err == nil {
		renameDir(oldBin, newBin, dryRun)
	}

	if dryRun {
		fmt.Println()
		fmt.Println("The above changes will be applied. Run with --apply to actually apply them.")
	} else {
		setupSh := filepath.Join(rootDir, "scripts", "setup.sh")
		if _, err := os.Stat(setupSh); err == nil {
			os.Remove(setupSh)
			fmt.Println("Removed setup.sh")
		}

		zzzDir := filepath.Join(rootDir, "zzz")
		if _, err := os.Stat(zzzDir); err == nil {
			os.RemoveAll(zzzDir)
			fmt.Println("Removed zzz/ directory")
		}

		setupDir := filepath.Join(rootDir, "cmd", "setup")
		if err := os.RemoveAll(setupDir); err == nil {
			fmt.Println("Removed cmd/setup/ directory")
		} else {
			fmt.Printf("Please manually delete cmd/setup/ directory: %v\n", err)
		}

		fmt.Println()
		fmt.Println("Changes applied. Please run the following:")
		fmt.Println("  1. go mod tidy")
		fmt.Println("  2. task ci")
		fmt.Println("  3. git status to review changes")
	}
}
