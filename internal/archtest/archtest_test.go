package archtest_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const modulePrefix = "github.com/aognio/webknife"

var featurePackages = []string{
	"pkg/static",
	"pkg/proxy",
	"pkg/echo",
	"pkg/respond",
	"pkg/redirect",
	"pkg/auth",
	"pkg/headers",
	"pkg/observe",
	"pkg/redact",
	"pkg/server",
}

func TestFeaturePackagesAreIndependent(t *testing.T) {
	for _, pkg := range featurePackages {
		imports := importsInPackage(t, pkg)
		for _, imp := range imports {
			if strings.HasPrefix(imp, modulePrefix+"/pkg/") && imp != modulePrefix+"/"+pkg && imp != modulePrefix+"/pkg/webknife" {
				t.Errorf("%s imports sibling package %s", pkg, imp)
			}
		}
	}
}

func TestFeaturePackagesDoNotImportInternal(t *testing.T) {
	for _, pkg := range featurePackages {
		imports := importsInPackage(t, pkg)
		for _, imp := range imports {
			if strings.HasPrefix(imp, modulePrefix+"/internal/") {
				t.Errorf("%s imports internal package %s", pkg, imp)
			}
		}
	}
}

func TestPortsDoNotImportInternal(t *testing.T) {
	imports := importsInPackage(t, "pkg/webknife")
	for _, imp := range imports {
		if strings.HasPrefix(imp, modulePrefix+"/internal/") {
			t.Errorf("pkg/webknife imports internal package %s", imp)
		}
	}
}

func TestCLIDoesNotImportFeaturePackages(t *testing.T) {
	imports := importsInPackage(t, "internal/adapters/cli")
	for _, imp := range imports {
		if strings.HasPrefix(imp, modulePrefix+"/pkg/") {
			t.Errorf("CLI adapter imports feature package %s", imp)
		}
	}
}

func importsInPackage(t *testing.T, relPkgDir string) []string {
	t.Helper()
	// Resolve relative to project root (two levels up from internal/archtest/)
	projectRoot, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("resolving project root: %v", err)
	}
	absDir := filepath.Join(projectRoot, relPkgDir)
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		t.Skipf("package directory %s does not exist", absDir)
	}

	fset := token.NewFileSet()
	var imports []string

	entries, err := os.ReadDir(absDir)
	if err != nil {
		t.Fatalf("reading %s: %v", absDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		filePath := filepath.Join(absDir, entry.Name())
		f, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parsing %s: %v", filePath, err)
		}
		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			imports = append(imports, path)
		}
	}

	return imports
}
