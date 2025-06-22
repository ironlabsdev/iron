package generate

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ironlabsdev/iron/internal/utils"
	"github.com/spf13/cobra"
)

const (
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
)

//go:embed all:templates/*
var TemplatesFS embed.FS

// GenerateCmd is the base command.
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from templates",
	Long: `Generate code scaffolding from predefined templates.

Available templates:
  oauth   - OAuth authentication implementation`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

// FromTemplate copies and processes template files to the specified full path
func FromTemplate(templateName, fullPath string) error {
	templatePath := filepath.Join("templates", templateName)

	// Check if template exists
	if _, err := TemplatesFS.ReadDir(templatePath); err != nil {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Check if directory exists and validate it
	if err := validateTargetDirectory(fullPath); err != nil {
		return err
	}

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Extract project name from the full path for template processing
	projectName := filepath.Base(fullPath)

	// Walk through template directory
	err := fs.WalkDir(TemplatesFS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root template directory
		if path == templatePath {
			return nil
		}

		// Calculate relative path from template root
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		// Handle special file renames
		destRelPath := handleSpecialFileRenames(relPath)

		// Calculate destination path
		destPath := filepath.Join(fullPath, destRelPath)

		if d.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, 0755)
		}

		// Process file
		return processTemplateFile(path, destPath, projectName)
	})

	if err != nil {
		return fmt.Errorf("failed to generate from template: %w", err)
	}

	// Print success message in green with colored icons
	fmt.Printf("%s✓ Successfully generated %s project in '%s'%s\n", ColorGreen, templateName, fullPath, ColorReset)
	fmt.Printf("%s→ Navigate to your project: %scd %s%s\n", ColorBlue, ColorReset, fullPath, ColorReset)

	return nil
}

// handleSpecialFileRenames renames files back to their original names
func handleSpecialFileRenames(relPath string) string {
	// Handle go.mod and go.sum files
	if strings.HasSuffix(relPath, "go.mod.template") {
		return strings.TrimSuffix(relPath, ".template")
	}
	if strings.HasSuffix(relPath, "go.sum.template") {
		return strings.TrimSuffix(relPath, ".template")
	}

	// You can add more special file handling here if needed
	return relPath
}

// validateTargetDirectory checks if the target directory exists and is empty
func validateTargetDirectory(fullPath string) error {
	// Check if directory exists
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist, which is fine - we'll create it
			return nil
		}
		return fmt.Errorf("failed to check directory: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("path '%s' exists but is not a directory", fullPath)
	}

	// Check if directory is empty
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("directory '%s' is not empty - cannot generate code in a non-empty directory", fullPath)
	}

	return nil
}

// processTemplateFile reads, processes, and writes a template file
func processTemplateFile(srcPath, destPath, projectName string) error {
	// Read template content
	content, err := TemplatesFS.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", srcPath, err)
	}

	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Check if file should be processed as template
	if shouldProcessAsTemplate(srcPath) {
		return processGoTemplate(content, destPath, projectName)
	}

	// Copy file as-is
	return os.WriteFile(destPath, content, 0644)
}

// shouldProcessAsTemplate determines if a file should be processed as a Go template
func shouldProcessAsTemplate(path string) bool {
	// Don't process files that end with .template as Go templates
	// They should be renamed and copied as-is
	if strings.HasSuffix(path, ".template") {
		return false
	}

	// Process specific file types as templates
	ext := filepath.Ext(path)
	templateExts := []string{".go", ".mod", ".yaml", ".yml", ".json", ".md", ".txt", ".env"}

	for _, tExt := range templateExts {
		if ext == tExt {
			return true
		}
	}

	return false
}

// processGoTemplate processes content as a Go template
func processGoTemplate(content []byte, destPath, projectName string) error {
	// Create template data
	data := struct {
		ProjectName      string
		ProjectNameCamel string
		ProjectNameSnake string
		ProjectNameKebab string
	}{
		ProjectName:      projectName,
		ProjectNameCamel: utils.ToCamelCase(projectName),
		ProjectNameSnake: utils.ToSnakeCase(projectName),
		ProjectNameKebab: utils.ToKebabCase(projectName),
	}

	// Parse and execute template
	tmpl, err := template.New("template").Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create destination file
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
