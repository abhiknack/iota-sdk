# Kiro Commands Reference

This file documents workflows you can ask Kiro to perform directly (no hooks needed).

## Git & PR Workflows

### Commit and Push
**Ask:** "Commit and push my changes"

**What happens:**
- Analyzes changed files
- Runs formatting (`make check fmt`) if Go/Templ files changed
- Validates translations (`make check tr`) if TOML files changed
- Creates logical commits with conventional commit messages
- Pulls latest changes to avoid conflicts
- Pushes to current branch

### Create Pull Request
**Ask:** "Create a PR" or "Commit and create PR"

**What happens:**
- Commits and pushes changes (as above)
- **On staging branch:** Creates new feature branch + PR automatically
- **On feature branch:** Checks if PR exists, creates only if needed
- Generates multilingual PR description (English + Russian)
- Includes comprehensive test plan
- Returns PR URL

### Address PR Comments
**Ask:** "Address PR comments for #123" or provide PR URL

**What happens:**
- Fetches unresolved review comments from GitHub
- Checks CI status and gets failure logs if failing
- Fixes CI failures first (tests, linting, build errors)
- Addresses each review comment systematically
- Commits and pushes fixes
- Verifies CI passes after changes

### Review PR
**Ask:** "Review PR #123" or provide PR URL

**What happens:**
- Fetches PR files and diffs from GitHub
- Performs comprehensive code review
- Identifies Critical (❌), Minor (🟡), and Style (🟢) issues
- Posts review comments directly on GitHub
- Submits overall review (APPROVE/REQUEST_CHANGES/COMMENT)

## Testing & Quality

### Fix Failing Tests
**Ask:** "Fix failing tests"

**What happens:**
- Runs `go vet ./...` to catch compilation errors
- Runs `make test` to identify failures
- Fixes each test incrementally
- Validates fixes with targeted test runs
- Ensures no regressions

### Fix Linting Issues
**Ask:** "Fix linting issues"

**What happens:**
- Runs `make check lint`
- Removes unused code (variables, functions, imports)
- Fixes formatting issues
- Re-validates until clean

### Fix E2E Tests
**Ask:** "Fix E2E tests"

**What happens:**
- Resets E2E database (`make e2e reset`)
- Checks for port conflicts and environment issues
- Runs E2E tests to identify failures
- Categorizes failures (database, timing, forms, navigation)
- Fixes systematically using Playwright debugging tools
- Validates with full E2E suite

## Code Quality

### Refactor & Review
**Ask:** "Refactor my changes" or "Review my code"

**What happens:**
- Analyzes uncommitted changes
- Identifies security issues, DRY violations, hard-coded values
- Fixes by priority (security → API misuse → style)
- Runs validation (`go vet`, `make check tr`)
- Provides detailed report

### Holistic Refactor
**Ask:** "Holistic refactor" or "Refactor with breaking changes allowed"

**What happens:**
- Deep analysis of code architecture
- Identifies layered hacks and tight coupling
- Proposes design options (A/B/C)
- Implements simplest option that removes complexity
- May break compatibility for better design
- Includes migration notes

## Database

### Create Migration
**Ask:** "Create a migration for [description]"

**What happens:**
- Generates timestamp for filename
- Creates migration with Up and Down sections
- Includes proper tenant_id and indexes
- Validates reversibility

### Connect to Database
**Ask:** "Connect to staging database" or "Connect to local database"

**What happens:**
- Provides connection command with credentials
- For staging: Uses Railway connection details
- For local: Uses local PostgreSQL credentials

## Translations

### Sync Translations
**Ask:** "Sync translations" or use the manual hook

**What happens:**
- Runs `make check tr` to find inconsistencies
- Updates all three files (en.toml, ru.toml, uz.toml)
- Ensures proper key naming
- Validates with `make check tr`

## Issue Management

### Fix GitHub Issue
**Ask:** "Fix issue #123" or provide issue URL

**What happens:**
- Fetches issue details from GitHub
- Moves issue to "In Progress" in project board
- Creates feature branch
- Implements fix with TDD approach (when applicable)
- Runs tests and validation
- Ready for commit/PR (use separate commands)

## Context & Documentation

### Dump Project Context
**Ask:** "Dump context" or "Create context dump"

**What happens:**
- Captures git state, work in progress, TODOs
- Analyzes module status and active development areas
- Extracts action items from code
- Provides resumption guidance
- Saves to `CONTEXT_DUMP_[timestamp].md`

## Tips

- **Be specific:** "Fix tests in modules/finance" is better than "fix tests"
- **Combine workflows:** "Commit, create PR, and address any CI failures"
- **Use hooks:** For repetitive tasks, enable the appropriate hook
- **Check steering:** The `.kiro/steering/` files provide context for all operations
