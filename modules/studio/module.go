package studio

import (
	"github.com/iota-uz/iota-sdk/modules/studio/infrastructure/persistence"
	"github.com/iota-uz/iota-sdk/modules/studio/presentation/controllers"
	"github.com/iota-uz/iota-sdk/modules/studio/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Register(app application.Application) error {
	moduleDefRepo := persistence.NewModuleDefinitionRepository()

	moduleDefService := services.NewModuleDefinitionService(moduleDefRepo, app.EventPublisher())
	codeGenService := services.NewCodeGeneratorService(moduleDefService)

	app.RegisterServices(
		moduleDefService,
		codeGenService,
	)

	app.RegisterControllers(
		controllers.NewModuleDefinitionController(app),
	)

	return nil
}

func (m *Module) Name() string {
	return "Studio"
}
