---
inclusion: manual
---

# Excalidraw Design System for IOTA ERP

This steering file is only included when explicitly referenced with #excalidraw-design.

## IOTA ERP Design Principles

### Frame-Based Layout
- **Frame 1**: Main table/list view
- **Frame 2**: Create form (separate view)
- **Frame 3**: Edit form (separate view)
- **Frame 4**: Detail/view page (separate view)

### Infinity Scroll
- **Never use pagination** - always infinity scroll
- Include loading indicators
- Show "Load more" states
- Smooth scrolling transitions

## Color Palette
```
Primary: #ffffff (background), #e5e7eb/#d1d5db (borders)
Text: #374151/#1f2937 (primary), #6b7280 (secondary)
Accent: #3b82f6/#06b6d4 (blue/teal)
States: #f3f4f6 (disabled), #10b981 (success), #ef4444 (error)
```

## Typography & Sizing
```
Headers: 24px (page), 18px (section)
Body: 16px (standard), 14px (secondary)
Elements: 40px (input/button), 6px (border radius)
Layout: 64px (header), 280px (sidebar)
Spacing: 8px/12px/16px/24px/32px
```

## Component Templates

### Header
```json
{"type": "rectangle", "width": 1200, "height": 64, "backgroundColor": "#ffffff"}
{"type": "text", "text": "IOTA ERP", "fontSize": 18}
```

### Table with Infinity Scroll
```json
{"type": "rectangle", "width": 800, "height": 40, "backgroundColor": "#f9fafb"}
{"type": "text", "text": "# ↑↓", "fontSize": 14}
```

### Infinity Scroll Loader
```json
{"type": "text", "text": "⟳", "fontSize": 20, "strokeColor": "#3b82f6"}
{"type": "text", "text": "Loading more...", "fontSize": 14}
```

## Usage
When creating Excalidraw schemas:
1. Use frame-based layout (separate views)
2. Implement infinity scroll (no pagination)
3. Apply IOTA color palette
4. Follow typography standards
5. Generate valid Excalidraw JSON
