---
inclusion: always
---

# UI Development Standards for IOTA SDK

## HTMX Integration (CRITICAL)
**ALWAYS use pkg/htmx package functions** - NEVER access headers directly

### Response Headers
```go
htmx.Redirect(w, path)
htmx.SetTrigger(w, event, data)
htmx.Refresh(w)
htmx.PushUrl(w, url)
```

### Request Headers
```go
htmx.IsHxRequest(r)
htmx.Target(r)
htmx.CurrentUrl(r)
```

## IOTA SDK Components

### Button Components
```go
@button.Primary(button.Props{
    Size: button.SizeNormal,
    Icon: icons.Plus(icons.Props{Size: "16"}),
    Attrs: templ.Attributes{
        "hx-post": "/api/create",
        "hx-target": "#content",
    },
}) {
    { pageCtx.T("CreateNew") }
}
```

### Input Components
```go
@input.Text(&input.Props{
    Label: pageCtx.T("Username"),
    Placeholder: pageCtx.T("EnterUsername"),
    Error: validationErrors["username"],
    Attrs: templ.Attributes{
        "name": "username",
        "required": true,
    },
})
```

### Badge Components
```go
@badge.New(badge.Props{
    Variant: badge.VariantGreen,
    Size: badge.SizeNormal,
}) {
    Active
}
```

## Security Patterns

### Template Security
```go
// ALWAYS use for dynamic URLs
<a href={ templ.URL(dynamicURL) }>

// NEVER use @templ.Raw() with user content
@templ.Raw(sanitizedHTML)  // Only for trusted content

// CSRF tokens in forms
<input type="hidden" name="csrf_token" value={ pageCtx.CSRFToken } />
```

## Translation Management
- **Always edit all three files**: en.toml, ru.toml, uz.toml
- **Avoid TOML reserved keys**: Never use `OTHER`, `ID`, `DESCRIPTION`
- **Follow enum pattern**: `Module.Enums.EnumType.VALUE`
- **Always validate**: Run `make check tr` after changes

## Form Field Naming (CRITICAL)
**ALWAYS use CamelCase for HTML form field names**:
- ✅ CORRECT: `name="FirstName"`, `name="EmailAddress"`
- ❌ INCORRECT: `name="first_name"`, `name="first-name"`

## Composables Usage
```go
// Parse form data
formData, err := composables.UseForm(&DTO{}, r)

// Get page context (for translations)
pageCtx := composables.UsePageCtx(ctx)

// Handle HTMX vs regular requests
if htmx.IsHxRequest(r) {
    htmx.SetTrigger(w, "dataUpdated", data)
    // Return partial HTML
} else {
    shared.Redirect(w, r, "/success")
}
```

## Tailwind CSS Colors (OKLCH)
- Surface: `bg-surface-100` to `bg-surface-600`
- Primary: `bg-primary-100` to `bg-primary-700`
- Semantic: `bg-green-*`, `bg-red-*`, `bg-yellow-*`, `bg-blue-*`
- Badge: `bg-badge-pink`, `bg-badge-yellow`, `bg-badge-green`

## Build Commands
- Generate templates: `templ generate`
- Compile CSS: `make css`
- Combined: `templ generate && make css`
- Format templates: `templ fmt`
