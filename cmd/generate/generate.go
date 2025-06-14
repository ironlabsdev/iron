package generate

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ironlabsdev/iron/internal/utils"
	"github.com/spf13/cobra"
)

//go:embed templates/*
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

// FromTemplate copies and processes template files
func FromTemplate(templateName, projectName string) error {
	templatePath := filepath.Join("templates", templateName)

	// Check if template exists
	if _, err := TemplatesFS.ReadDir(templatePath); err != nil {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

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

		// Calculate destination path
		destPath := filepath.Join(projectName, relPath)

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

	fmt.Printf("✅ Successfully generated %s project in '%s/'\n", templateName, projectName)
	fmt.Printf("📁 Navigate to your project: cd %s\n", projectName)

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
