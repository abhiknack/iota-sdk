package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/modules/studio/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/studio/presentation/mappers"
	"github.com/iota-uz/iota-sdk/modules/studio/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
	"github.com/iota-uz/iota-sdk/pkg/shared"
)

type ModuleDefinitionController struct {
	app              application.Application
	service          *services.ModuleDefinitionService
	generatorService *services.CodeGeneratorService
	basePath         string
}

func NewModuleDefinitionController(app application.Application) application.Controller {
	return &ModuleDefinitionController{
		app:              app,
		service:          app.Service(services.ModuleDefinitionService{}).(*services.ModuleDefinitionService),
		generatorService: app.Service(services.CodeGeneratorService{}).(*services.CodeGeneratorService),
		basePath:         "/studio/modules",
	}
}

func (c *ModuleDefinitionController) Key() string {
	return c.basePath
}

func (c *ModuleDefinitionController) Register(r *mux.Router) {
	commonMiddleware := []mux.MiddlewareFunc{
		middleware.Authorize(),
		middleware.RedirectNotAuthenticated(),
		middleware.ProvideUser(),
		middleware.ProvideDynamicLogo(c.app),
		middleware.ProvideLocalizer(c.app.Bundle()),
		middleware.NavItems(),
		middleware.WithPageContext(),
	}

	sub := r.PathPrefix("/studio/modules").Subrouter()
	sub.Use(commonMiddleware...)

	sub.HandleFunc("", c.List).Methods(http.MethodGet)
	sub.HandleFunc("/new", c.New).Methods(http.MethodGet)
	sub.HandleFunc("", c.Create).Methods(http.MethodPost)
	sub.HandleFunc("/{id}", c.Detail).Methods(http.MethodGet)
	sub.HandleFunc("/{id}/edit", c.Edit).Methods(http.MethodGet)
	sub.HandleFunc("/{id}", c.Update).Methods(http.MethodPut)
	sub.HandleFunc("/{id}", c.Delete).Methods(http.MethodDelete)
	sub.HandleFunc("/{id}/entities", c.AddEntity).Methods(http.MethodPost)
	sub.HandleFunc("/{id}/entities/{entityId}", c.RemoveEntity).Methods(http.MethodDelete)
	sub.HandleFunc("/{id}/generate", c.Generate).Methods(http.MethodPost)
	sub.HandleFunc("/{id}/preview", c.Preview).Methods(http.MethodGet)
}

func (c *ModuleDefinitionController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	modules, err := c.service.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

func (c *ModuleDefinitionController) New(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Module creation form"))
}

func (c *ModuleDefinitionController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dto, err := composables.UseForm(&dtos.ModuleDefinitionCreateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mod, err := c.service.Create(ctx, dto.Name, dto.DisplayName, dto.Description, dto.Icon, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.Redirect(w, r, "/studio/modules/"+mod.ID().String())
}

func (c *ModuleDefinitionController) Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mod, err := c.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mod)
}

func (c *ModuleDefinitionController) Edit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Module edit form"))
}

func (c *ModuleDefinitionController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dto, err := composables.UseForm(&dtos.ModuleDefinitionUpdateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.service.Update(ctx, id, dto.DisplayName, dto.Description, dto.Icon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.Redirect(w, r, "/studio/modules/"+id.String())
}

func (c *ModuleDefinitionController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.Redirect(w, r, "/studio/modules")
}

func (c *ModuleDefinitionController) AddEntity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var dto dtos.EntityDefinitionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity := mappers.ToEntityDefinition(&dto)

	_, err = c.service.AddEntity(ctx, id, entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (c *ModuleDefinitionController) RemoveEntity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entityID, err := uuid.Parse(vars["entityId"])
	if err != nil {
		http.Error(w, "Invalid entity ID", http.StatusBadRequest)
		return
	}

	_, err = c.service.RemoveEntity(ctx, id, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ModuleDefinitionController) Generate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = c.generatorService.GenerateModule(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Module generated successfully"})
}

func (c *ModuleDefinitionController) Preview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	preview, err := c.generatorService.PreviewCode(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preview)
}
