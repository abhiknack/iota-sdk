package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const controllerTemplate = `package controllers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/components/base/pagination"
	"github.com/iota-uz/iota-sdk/modules/{{.ModuleName}}/domain/aggregates/{{.EntityLower}}"
	"github.com/iota-uz/iota-sdk/modules/{{.ModuleName}}/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/{{.ModuleName}}/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/htmx"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
	"github.com/iota-uz/iota-sdk/pkg/shared"
)

type {{.EntityName}}Controller struct {
	app     application.Application
	service *services.{{.EntityName}}Service
	basePath string
}

func New{{.EntityName}}Controller(app application.Application) application.Controller {
	return &{{.EntityName}}Controller{
		app:      app,
		service:  app.Service(services.{{.EntityName}}Service{}).(*services.{{.EntityName}}Service),
		basePath: "/{{.ModuleName}}/{{.EntityPlural}}",
	}
}

func (c *{{.EntityName}}Controller) Key() string {
	return c.basePath
}

func (c *{{.EntityName}}Controller) Register(r *mux.Router) {
	commonMiddleware := []mux.MiddlewareFunc{
		middleware.Authorize(),
		middleware.RedirectNotAuthenticated(),
		middleware.ProvideUser(),
		middleware.ProvideDynamicLogo(c.app),
		middleware.ProvideLocalizer(c.app.Bundle()),
		middleware.NavItems(),
		middleware.WithPageContext(),
	}

	getRouter := r.PathPrefix(c.basePath).Subrouter()
	getRouter.Use(commonMiddleware...)
	getRouter.HandleFunc("", c.List).Methods(http.MethodGet)
	getRouter.HandleFunc("/new", c.New).Methods(http.MethodGet)
	getRouter.HandleFunc("/{id}", c.Edit).Methods(http.MethodGet)

	setRouter := r.PathPrefix(c.basePath).Subrouter()
	setRouter.Use(commonMiddleware...)
	setRouter.Use(middleware.WithTransaction())
	setRouter.HandleFunc("", c.Create).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}", c.Update).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}/delete", c.Delete).Methods(http.MethodPost)
}

func (c *{{.EntityName}}Controller) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	paginationParams := composables.UsePaginated(r)
	filterDTO, err := composables.UseQuery(&dtos.{{.EntityName}}FilterDTO{
		Page:     paginationParams.Page,
		PageSize: paginationParams.Limit,
	}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &{{.EntityLower}}.FindParams{
		TenantID: tenantID,
		Limit:    filterDTO.PageSize,
		Offset:   (filterDTO.Page - 1) * filterDTO.PageSize,
		SortBy:   {{.EntityLower}}.FieldCreatedAt,
		SortDesc: true,
	}

	entities, err := c.service.GetPaginated(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := c.service.Count(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	paginationState := pagination.New(r.URL.Path, filterDTO.Page, int(total), filterDTO.PageSize)

	if htmx.IsHxRequest(r) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<div>TODO: Render {{.EntityPlural}} table</div>")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<div>TODO: Render {{.EntityPlural}} page with pagination: %v</div>", paginationState)
}

func (c *{{.EntityName}}Controller) New(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<div>TODO: Render new {{.EntityLower}} form</div>")
}

func (c *{{.EntityName}}Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dto, err := composables.UseForm(&dtos.{{.EntityName}}CreateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity := {{.EntityLower}}.New{{.EntityName}}(
		uuid.New(),
		tenantID,
{{- range .Fields}}
		dto.{{.Name}},
{{- end}}
	)

	created, err := c.service.Create(ctx, entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.SetTrigger(w, "{{.EntityLower}}Created", map[string]interface{}{"id": created.ID()})
		htmx.Redirect(w, c.basePath)
		return
	}

	shared.Redirect(w, r, c.basePath)
}

func (c *{{.EntityName}}Controller) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entity, err := c.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<div>TODO: Render edit form for {{.EntityLower}}: %v</div>", entity.ID())
}

func (c *{{.EntityName}}Controller) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dto, err := composables.UseForm(&dtos.{{.EntityName}}UpdateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existing, err := c.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	entity := {{.EntityLower}}.New{{.EntityName}}(
		existing.ID(),
		existing.TenantID(),
{{- range .Fields}}
		dto.{{.Name}},
{{- end}}
		{{.EntityLower}}.WithTimestamps(existing.CreatedAt(), existing.UpdatedAt()),
	)

	updated, err := c.service.Update(ctx, entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.SetTrigger(w, "{{.EntityLower}}Updated", map[string]interface{}{"id": updated.ID()})
		htmx.Redirect(w, c.basePath)
		return
	}

	shared.Redirect(w, r, c.basePath)
}

func (c *{{.EntityName}}Controller) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := c.service.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.SetTrigger(w, "{{.EntityLower}}Deleted", map[string]interface{}{"id": id})
		htmx.Redirect(w, c.basePath)
		return
	}

	shared.Redirect(w, r, c.basePath)
}
`

type controllerTemplateData struct {
	ModuleName   string
	EntityName   string
	EntityLower  string
	EntityPlural string
	Fields       []fieldTemplateData
}

func GenerateController(moduleName, entityName string, fields []Field) error {
	entityLower := strings.ToLower(entityName[:1]) + entityName[1:]
	entityPlural := strings.ToLower(entityName) + "s"

	data := controllerTemplateData{
		ModuleName:   moduleName,
		EntityName:   entityName,
		EntityLower:  entityLower,
		EntityPlural: entityPlural,
		Fields:       make([]fieldTemplateData, len(fields)),
	}

	for i, f := range fields {
		data.Fields[i] = fieldTemplateData{
			Name:      f.Name,
			NameLower: strings.ToLower(f.Name[:1]) + f.Name[1:],
			Type:      f.Type,
		}
	}

	basePath := filepath.Join("modules", moduleName, "presentation", "controllers")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	outputPath := filepath.Join(basePath, entityLower+"_controller.go")
	return generateFromTemplate(controllerTemplate, outputPath, data)
}
