package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/iota-uz/iota-sdk/cmd/codegen/generators"
)

func main() {
	var (
		genType    = flag.String("type", "", "Type of code to generate: entity, crud, migration")
		moduleName = flag.String("module", "", "Module name (e.g., fleet, crm)")
		entityName = flag.String("entity", "", "Entity name (e.g., Vehicle, Driver)")
		fields     = flag.String("fields", "", "Comma-separated fields: name:type:tag (e.g., Name:string:required,Age:int:min=0)")
	)

	flag.Parse()

	if *genType == "" {
		fmt.Println("Usage: go run cmd/codegen/main.go -type=<type> [options]")
		fmt.Println("\nTypes:")
		fmt.Println("  entity     - Generate domain aggregate")
		fmt.Println("  crud       - Generate complete CRUD (domain, repo, service, controller, DTOs)")
		fmt.Println("  migration  - Generate migration file template")
		fmt.Println("\nOptions:")
		fmt.Println("  -module    - Module name (required for entity/crud)")
		fmt.Println("  -entity    - Entity name (required for entity/crud)")
		fmt.Println("  -fields    - Field definitions (optional)")
		fmt.Println("\nExamples:")
		fmt.Println("  go run cmd/codegen/main.go -type=crud -module=fleet -entity=Vehicle -fields=\"Make:string:required,Model:string:required,Year:int:min=1900\"")
		fmt.Println("  go run cmd/codegen/main.go -type=migration")
		os.Exit(1)
	}

	switch *genType {
	case "entity":
		if *moduleName == "" || *entityName == "" {
			fmt.Println("Error: -module and -entity are required for entity generation")
			os.Exit(1)
		}
		fieldList := parseFields(*fields)
		if err := generators.GenerateEntity(*moduleName, *entityName, fieldList); err != nil {
			fmt.Printf("Error generating entity: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Generated entity: %s\n", *entityName)

	case "crud":
		if *moduleName == "" || *entityName == "" {
			fmt.Println("Error: -module and -entity are required for CRUD generation")
			os.Exit(1)
		}
		fieldList := parseFields(*fields)
		if err := generators.GenerateCRUD(*moduleName, *entityName, fieldList); err != nil {
			fmt.Printf("Error generating CRUD: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Generated complete CRUD for: %s\n", *entityName)

	case "migration":
		if err := generators.GenerateMigration(); err != nil {
			fmt.Printf("Error generating migration: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Generated migration file")

	default:
		fmt.Printf("Unknown type: %s\n", *genType)
		os.Exit(1)
	}
}

func parseFields(fieldsStr string) []generators.Field {
	if fieldsStr == "" {
		return nil
	}

	var fields []generators.Field
	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		fieldParts := strings.Split(strings.TrimSpace(part), ":")
		if len(fieldParts) < 2 {
			continue
		}

		field := generators.Field{
			Name: fieldParts[0],
			Type: fieldParts[1],
		}

		if len(fieldParts) >= 3 {
			field.Validation = fieldParts[2]
		}

		fields = append(fields, field)
	}

	return fields
}
