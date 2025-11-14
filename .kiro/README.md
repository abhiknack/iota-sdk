# Kiro Configuration for IOTA SDK

This directory contains Kiro-specific configuration migrated from your `.claude` setup.

## Structure

```
.kiro/
├── steering/              # Always-included context (like Claude agents)
│   ├── go-development.md      # Go/DDD/ITF patterns
│   ├── database-patterns.md   # PostgreSQL/migrations
│   └── ui-development.md      # HTMX/Templ/translations
├── hooks/                 # Automated workflows (like Claude commands)
│   ├── manual-*.json          # Manual trigger hooks
│   └── on-save-*.json         # Auto-trigger hooks (disabled by default)
├── COMMANDS.md           # Quick reference for asking Kiro
└── README.md             # This file
```

## Steering Files (Always Active)

These provide context for every Kiro interaction:

- **go-development.md**: DDD architecture, ITF testing, dependency injection
- **database-patterns.md**: Multi-tenant patterns, migrations, SQL security
- **ui-development.md**: HTMX integration, IOTA components, translations
- **railway-deployment.md**: Railway platform operations, staging/production deployment
- **e2e-testing.md**: Playwright E2E testing patterns, database isolation
- **configuration-management.md**: Environment files, Docker, Makefile management
- **debugging.md**: Systematic debugging workflow, root cause analysis
- **refactoring-standards.md**: Code quality standards, security patterns
- **excalidraw-design.md**: UI mockup design system (manual inclusion only)

## Hooks (Automated Workflows)

### Manual Hooks (Enabled)
Click the hook button or ask Kiro to run:

- **Sync Translations**: Synchronize en/ru/uz.toml files
- **Refactor & Review**: Comprehensive code review of changes
- **Fix Linting**: Run and fix `make check lint` issues
- **Fix Failing Tests**: Systematically fix broken tests

### Auto Hooks (Disabled by Default)
Enable in Kiro's Agent Hooks view:

- **Run Tests on Save**: Auto-run tests when `*_test.go` files saved
- **Format on Save**: Auto-format Go/Templ files on save

## Quick Start

### Just Ask Kiro
Instead of slash commands, just ask naturally:

```
"Commit and push my changes"
"Create a PR for this feature"
"Fix failing tests"
"Review PR #123"
"Sync translations"
"Create a migration for adding user roles"
```

See `COMMANDS.md` for full list of workflows.

### Enable Hooks
1. Open Command Palette (Ctrl/Cmd+Shift+P)
2. Search "Open Kiro Hook UI"
3. Enable desired hooks

### View Steering Rules
Steering files are automatically loaded. To modify:
- Edit files in `.kiro/steering/`
- Changes take effect immediately

## Migration from Claude

Your `.claude` setup has been converted:

| Claude | Kiro | Status |
|--------|------|--------|
| `.claude/agents/*.md` (9 agents) | `.kiro/steering/*.md` | ✅ Converted (9 steering files) |
| `.claude/commands/*.md` (10 commands) | Ask Kiro directly | ✅ Available as natural commands |
| `.claude/commands/*.md` | `.kiro/hooks/*.json` | ✅ Created (6 hooks) |

## What's Different

### Claude Code
- Used "agents" for specialized contexts
- Used "slash commands" for workflows
- Required explicit agent delegation

### Kiro
- Uses "steering" for always-on context
- Uses natural language for workflows
- Uses "hooks" for automation
- No explicit delegation needed

## Examples

### Before (Claude)
```
/commit-pr --base staging
```

### Now (Kiro)
```
"Create a PR to staging"
```

### Before (Claude)
```
Task(agent: go-editor) "Fix the user service tests"
```

### Now (Kiro)
```
"Fix the tests in modules/users/services"
```

The steering files ensure Kiro knows Go/DDD/ITF patterns automatically.

## Tips

1. **Steering is powerful**: All 3 steering files load on every interaction
2. **Hooks are optional**: Use for repetitive tasks
3. **Natural language works**: Just describe what you want
4. **Check COMMANDS.md**: For complex workflows
5. **Keep .claude folder**: As reference/backup

## Support

- Steering files: Edit `.kiro/steering/*.md`
- Hooks: Use "Open Kiro Hook UI" command
- Commands: See `.kiro/COMMANDS.md`
- Original setup: See `.claude/` folder
