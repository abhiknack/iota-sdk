package generators

import (
	"fmt"
)

func GenerateCRUD(moduleName, entityName string, fields []Field) error {
	fmt.Println("Generating domain aggregate...")
	if err := GenerateEntity(moduleName, entityName, fields); err != nil {
		return fmt.Errorf("failed to generate entity: %w", err)
	}

	fmt.Println("Generating repository...")
	if err := GenerateRepository(moduleName, entityName, fields); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	fmt.Println("Generating service...")
	if err := GenerateService(moduleName, entityName); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	fmt.Println("Generating DTOs...")
	if err := GenerateDTO(moduleName, entityName, fields); err != nil {
		return fmt.Errorf("failed to generate DTOs: %w", err)
	}

	fmt.Println("Generating controller...")
	if err := GenerateController(moduleName, entityName, fields); err != nil {
		return fmt.Errorf("failed to generate controller: %w", err)
	}

	fmt.Println("\n✓ CRUD generation complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Create database migration for the table")
	fmt.Println("2. Register service in module.go: app.RegisterServices()")
	fmt.Println("3. Register controller in module.go: app.RegisterControllers()")
	fmt.Println("4. Add permissions to permissions/constants.go")
	fmt.Println("5. Create templates in presentation/templates/pages/")
	fmt.Println("6. Add translations to presentation/locales/")
	fmt.Println("7. Run: templ generate && make css")

	return nil
}
