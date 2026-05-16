package apitest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSCN_AIPSM_001_008_ProviderSettingsContractExposesOnlyRealUserFacingProviders(t *testing.T) {
	sources := loadProviderSettingsContractSources(t)

	combined := strings.ToLower(strings.Join(sources.allText(), "\n"))
	for _, providerID := range []string{"gemini", "lm_studio", "xai"} {
		if !strings.Contains(combined, providerID) {
			t.Fatalf("SCN-AIPSM-001/008 expected user-facing provider id %q in provider settings contract", providerID)
		}
	}
	for _, forbidden := range []string{`"fake"`, "`fake`", "fakeprovider", "fake_provider"} {
		if strings.Contains(combined, forbidden) {
			t.Fatalf("SCN-AIPSM-001/008 expected fake provider to stay out of user-facing provider contract, found %q", forbidden)
		}
	}
}

func TestSCN_AIPSM_002_005_009_ProviderSettingsSecretBoundaryIsRedactedInPublicContract(t *testing.T) {
	sources := loadProviderSettingsContractSources(t)

	for _, source := range sources.goSources {
		for _, field := range providerSettingsStructFields(source.file) {
			if source.path == "internal/controller/wails/provider_settings_controller.go" &&
				field.structName == "SaveProviderSettingsRequestDTO" &&
				normalizeContractName(field.name) == "credentialinput" {
				continue
			}
			if isForbiddenSecretField(field) {
				t.Fatalf(
					"SCN-AIPSM-002/005/009 expected no API key or raw payload field in DTO/read model contract, got %s.%s %s",
					source.path,
					field.name,
					field.typeName,
				)
			}
		}
	}
	for path, text := range sources.tsSources {
		assertTextHasNoProviderSettingsSecretOutputFields(t, "SCN-AIPSM-002/005/009", path, text)
	}
}

func TestSCN_AIPSM_002_006_ProviderSettingsPublicSeamsAreFrozen(t *testing.T) {
	sources := loadProviderSettingsContractSources(t)

	requiredWailsSeams := []string{
		"ListProviderSettings",
		"SaveProviderSettings",
		"ResetProviderSettings",
		"ValidateProviderSettings",
	}
	for _, seam := range requiredWailsSeams {
		if !sources.hasGoSymbol(seam) {
			t.Fatalf("SCN-AIPSM-002/006 expected public Wails seam %s in backend usecase/controller contract", seam)
		}
		if !sources.hasText(seam) {
			t.Fatalf("SCN-AIPSM-002/006 expected public seam %s in frontend gateway contract", seam)
		}
	}

	for _, seam := range []string{"ListProviderModels", "ResolveProviderExecutionSettings"} {
		if !sources.hasGoSymbol(seam) {
			t.Fatalf("SCN-AIPSM-002/006 expected backend internal consumer seam %s in usecase contract", seam)
		}
		if strings.Contains(sources.textForPath("internal/controller/wails/provider_settings_controller.go"), seam+"(") {
			t.Fatalf("SCN-AIPSM-002/006 expected %s to stay out of provider settings Wails public controller", seam)
		}
		if sources.hasText(seam) {
			t.Fatalf("SCN-AIPSM-002/006 expected backend internal seam %s to stay out of provider settings frontend gateway contract", seam)
		}
	}
}

func TestSCN_AIPSM_002_006_SaveProviderSettingsRequestAllowsOnlyTransientCredentialInput(t *testing.T) {
	sources := loadProviderSettingsContractSources(t)

	saveFields := sources.goSaveRequestFields()
	if !hasAnyNormalizedField(saveFields, "credentialinput") {
		t.Fatalf("SCN-AIPSM-002/006 expected SaveProviderSettings input DTO to allow transient credentialInput, got %#v", saveFields)
	}
	for path, text := range sources.tsSources {
		if strings.Contains(text, "credentialInput") && !strings.Contains(path, "gateway-contract/provider-settings/provider-settings-gateway-contract.ts") {
			t.Fatalf("SCN-AIPSM-002/006 expected credentialInput only in save gateway input contract, got %s", path)
		}
	}
}

func TestSCN_AIPSM_004_006_SaveProviderSettingsRequestKeepsReferenceSettingsOut(t *testing.T) {
	sources := loadProviderSettingsContractSources(t)

	saveFields := sources.goSaveRequestFields()
	if len(saveFields) == 0 {
		t.Fatalf("SCN-AIPSM-004/006 expected SaveProviderSettings request DTO fields")
	}
	if !hasAnyNormalizedField(saveFields, "providerid", "provider") {
		t.Fatalf("SCN-AIPSM-004/006 expected SaveProviderSettings request to include provider id, got %#v", saveFields)
	}
	if !hasAnyNormalizedField(saveFields, "endpoint", "endpointurl") {
		t.Fatalf("SCN-AIPSM-004/006 expected SaveProviderSettings request to include endpoint, got %#v", saveFields)
	}
	if !hasAnyNormalizedField(saveFields, "apikeyinputpresent", "apikeyprovided", "hasapikeyinput", "credentialinputpresent", "credentialprovided") {
		t.Fatalf("SCN-AIPSM-004/006 expected SaveProviderSettings request to include API key input presence, got %#v", saveFields)
	}
	for _, field := range saveFields {
		if isForbiddenReferenceSettingField(field.name) {
			t.Fatalf("SCN-AIPSM-004/006 expected SaveProviderSettings request not to include reference-side setting field %q", field.name)
		}
	}
	saveRequestContractCount := 0
	for path, text := range sources.tsSources {
		if assertSaveProviderSettingsTextKeepsReferenceSettingsOut(t, path, text) {
			saveRequestContractCount++
		}
	}
	if saveRequestContractCount == 0 {
		t.Fatal("SCN-AIPSM-004/006 expected SaveProviderSettings request contract in TypeScript sources")
	}
}

func TestSCN_DFSS_006_ProviderSettingsSecretBackendRejectsUnsupportedBackend(t *testing.T) {
	source, err := parseProviderSettingsGoSource(
		providerSettingsRepoRoot(t),
		"internal/bootstrap/app_controller.go",
	)
	if err != nil {
		t.Fatalf("SCN-DFSS-006 expected bootstrap provider settings backend source: %v", err)
	}

	requiredSnippets := []string{
		"providerSettingsSecretBackendEnv",
		"providerSettingsSecretBackendInMemory",
		"case providerSettingsSecretBackendInMemory:",
		`case "":`,
		"unsupported",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(source.text, snippet) {
			t.Fatalf("SCN-DFSS-006 expected provider settings backend selector to expose safe unsupported-backend handling")
		}
	}
	if strings.Contains(source.text, "default:\n\t\tstore, err := repository.NewProviderSettingsKeyringSecretStore()") {
		t.Fatal("SCN-DFSS-006 expected unsupported backend values not to silently fall back to keyring")
	}
}

type providerSettingsContractSources struct {
	goSources []providerSettingsGoSource
	tsSources map[string]string
}

type providerSettingsGoSource struct {
	path string
	file *ast.File
	text string
}

type providerSettingsField struct {
	structName string
	name       string
	typeName   string
}

func loadProviderSettingsContractSources(t *testing.T) providerSettingsContractSources {
	t.Helper()

	repoRoot := providerSettingsRepoRoot(t)
	goPaths := []string{
		"internal/usecase/provider_settings_contract.go",
		"internal/controller/wails/provider_settings_controller.go",
	}
	tsPaths := []string{
		"frontend/src/application/gateway-contract/provider-settings/provider-settings-gateway-contract.ts",
		"frontend/src/application/contract/provider-settings/provider-settings-screen-contract.ts",
	}

	result := providerSettingsContractSources{tsSources: map[string]string{}}
	for _, path := range goPaths {
		source, err := parseProviderSettingsGoSource(repoRoot, path)
		if err != nil {
			t.Fatalf("provider-settings-contract-freeze expected %s: %v", path, err)
		}
		result.goSources = append(result.goSources, source)
	}
	for _, path := range tsPaths {
		text, err := os.ReadFile(filepath.Join(repoRoot, filepath.Clean(path))) // #nosec G304 -- fixed repo-local contract source.
		if err != nil {
			t.Fatalf("provider-settings-contract-freeze expected %s: %v", path, err)
		}
		result.tsSources[path] = string(text)
	}
	return result
}

func providerSettingsRepoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("provider-settings-contract-freeze could not resolve test file path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
}

func parseProviderSettingsGoSource(repoRoot string, path string) (providerSettingsGoSource, error) {
	text, err := os.ReadFile(filepath.Join(repoRoot, filepath.Clean(path))) // #nosec G304 -- fixed repo-local contract source.
	if err != nil {
		return providerSettingsGoSource{}, fmt.Errorf("read provider settings Go source: %w", err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), path, text, parser.ParseComments)
	if err != nil {
		return providerSettingsGoSource{}, fmt.Errorf("parse provider settings Go source: %w", err)
	}
	return providerSettingsGoSource{path: path, file: file, text: string(text)}, nil
}

func (sources providerSettingsContractSources) allText() []string {
	var texts []string
	for _, source := range sources.goSources {
		texts = append(texts, source.text)
	}
	for _, text := range sources.tsSources {
		texts = append(texts, text)
	}
	return texts
}

func (sources providerSettingsContractSources) textForPath(path string) string {
	for _, source := range sources.goSources {
		if source.path == path {
			return source.text
		}
	}
	return sources.tsSources[path]
}

func (sources providerSettingsContractSources) hasGoSymbol(name string) bool {
	for _, source := range sources.goSources {
		for _, declaration := range source.file.Decls {
			if hasProviderSettingsGoSymbol(declaration, name) {
				return true
			}
		}
	}
	return false
}

func (sources providerSettingsContractSources) hasText(text string) bool {
	for _, sourceText := range sources.tsSources {
		if strings.Contains(sourceText, text) {
			return true
		}
	}
	return false
}

func (sources providerSettingsContractSources) goSaveRequestFields() []providerSettingsField {
	var fields []providerSettingsField
	for _, source := range sources.goSources {
		for _, field := range providerSettingsStructFields(source.file) {
			normalizedStructName := normalizeContractName(field.structName)
			if strings.Contains(normalizedStructName, "saveprovidersettings") &&
				(strings.Contains(normalizedStructName, "request") || strings.Contains(normalizedStructName, "input") || strings.Contains(normalizedStructName, "command")) {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

func hasProviderSettingsGoSymbol(declaration ast.Decl, name string) bool {
	switch declaration := declaration.(type) {
	case *ast.FuncDecl:
		return declaration.Name != nil && declaration.Name.Name == name
	case *ast.GenDecl:
		return hasProviderSettingsGoTypeSymbol(declaration, name)
	}
	return false
}

func hasProviderSettingsGoTypeSymbol(declaration *ast.GenDecl, name string) bool {
	for _, spec := range declaration.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		if typeSpec.Name.Name == name || providerSettingsInterfaceHasMethod(typeSpec, name) {
			return true
		}
	}
	return false
}

func providerSettingsInterfaceHasMethod(typeSpec *ast.TypeSpec, name string) bool {
	interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return false
	}
	for _, method := range interfaceType.Methods.List {
		if providerSettingsFieldHasName(method, name) {
			return true
		}
	}
	return false
}

func providerSettingsFieldHasName(field *ast.Field, name string) bool {
	for _, methodName := range field.Names {
		if methodName.Name == name {
			return true
		}
	}
	return false
}

func providerSettingsStructFields(file *ast.File) []providerSettingsField {
	var fields []providerSettingsField
	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok {
			continue
		}
		fields = append(fields, providerSettingsStructFieldsFromGenDecl(genDecl)...)
	}
	return fields
}

func providerSettingsStructFieldsFromGenDecl(genDecl *ast.GenDecl) []providerSettingsField {
	var fields []providerSettingsField
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok || !strings.Contains(normalizeContractName(typeSpec.Name.Name), "provider") {
			continue
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		fields = append(fields, providerSettingsStructFieldsFromType(typeSpec.Name.Name, structType)...)
	}
	return fields
}

func providerSettingsStructFieldsFromType(structName string, structType *ast.StructType) []providerSettingsField {
	var fields []providerSettingsField
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fields = append(fields, providerSettingsField{
				structName: structName,
				name:       name.Name,
				typeName:   exprText(field.Type),
			})
		}
	}
	return fields
}

func isForbiddenSecretField(field providerSettingsField) bool {
	normalizedName := normalizeContractName(field.name)
	normalizedType := normalizeContractName(field.typeName)
	if strings.Contains(normalizedName, "requesttoken") ||
		strings.Contains(normalizedName, "credentialreference") ||
		strings.Contains(normalizedName, "credentialstate") ||
		strings.Contains(normalizedName, "apikeyinputpresent") {
		return false
	}
	for _, forbidden := range []string{
		"apikey",
		"secret",
		"authorization",
		"providertoken",
		"rawrequest",
		"rawresponse",
		"rawprompt",
	} {
		if strings.Contains(normalizedName, forbidden) && (normalizedType == "string" || normalizedType == "byte" || strings.Contains(normalizedType, "string")) {
			return true
		}
	}
	return false
}

func isForbiddenReferenceSettingField(name string) bool {
	normalized := normalizeContractName(name)
	for _, forbidden := range []string{
		"model",
		"batchapi",
		"batch",
		"processingmethod",
		"processmethod",
		"executionmethod",
		"selectedprovider",
		"providerselection",
	} {
		if strings.Contains(normalized, forbidden) {
			return true
		}
	}
	return false
}

func hasAnyNormalizedField(fields []providerSettingsField, names ...string) bool {
	expected := map[string]struct{}{}
	for _, name := range names {
		expected[normalizeContractName(name)] = struct{}{}
	}
	for _, field := range fields {
		if _, ok := expected[normalizeContractName(field.name)]; ok {
			return true
		}
	}
	return false
}

func assertTextHasNoProviderSettingsSecretOutputFields(t *testing.T, scenarioID string, path string, text string) {
	t.Helper()

	normalized := normalizeContractName(text)
	for _, forbidden := range []string{
		"apikey:string",
		"apikey?:string",
		"secret:string",
		"secret?:string",
		"authorization:string",
		"rawrequest",
		"rawresponse",
		"rawprompt",
	} {
		if strings.Contains(normalized, forbidden) {
			t.Fatalf("%s expected no API key or raw payload field in %s, found %q", scenarioID, path, forbidden)
		}
	}
}

func assertSaveProviderSettingsTextKeepsReferenceSettingsOut(t *testing.T, path string, text string) bool {
	t.Helper()

	saveRequestText, ok := providerSettingsTypeScriptInterfaceText(text, "SaveProviderSettingsRequest")
	if !ok {
		return false
	}
	normalized := normalizeContractName(saveRequestText)
	saveIndex := strings.Index(normalized, "saveprovidersettingsrequest")
	if saveIndex < 0 {
		t.Fatalf("SCN-AIPSM-004/006 expected SaveProviderSettings request contract in %s", path)
	}
	saveSection := normalized[saveIndex:]
	for _, forbidden := range []string{"model", "batchapi", "processingmethod", "selectedprovider", "providerselection"} {
		if strings.Contains(saveSection, forbidden) {
			t.Fatalf("SCN-AIPSM-004/006 expected SaveProviderSettings contract in %s not to include %q", path, forbidden)
		}
	}
	return true
}

func providerSettingsTypeScriptInterfaceText(text string, name string) (string, bool) {
	interfaceIndex := strings.Index(text, "interface "+name)
	if interfaceIndex < 0 {
		return "", false
	}
	openIndex := strings.Index(text[interfaceIndex:], "{")
	if openIndex < 0 {
		return "", false
	}
	start := interfaceIndex + openIndex
	depth := 0
	for index := start; index < len(text); index++ {
		switch text[index] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[interfaceIndex : index+1], true
			}
		}
	}
	return "", false
}

func exprText(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.ArrayType:
		return exprText(expr.Elt)
	case *ast.StarExpr:
		return exprText(expr.X)
	case *ast.SelectorExpr:
		return exprText(expr.X) + "." + expr.Sel.Name
	default:
		return ""
	}
}

func normalizeContractName(value string) string {
	value = strings.ToLower(value)
	replacer := strings.NewReplacer("_", "", "-", "", " ", "", "\t", "", "\n", "", "\r", "")
	return replacer.Replace(value)
}
