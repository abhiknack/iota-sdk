#!/bin/bash

# IOTA SDK Code Generator Helper Script
# Usage: ./scripts/generate.sh [command] [options]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

show_help() {
    cat << EOF
IOTA SDK Code Generator

Usage: ./scripts/generate.sh [command] [options]

Commands:
    crud        Generate complete CRUD (domain, repo, service, controller, DTOs)
    entity      Generate domain aggregate only
    migration   Generate migration file template
    help        Show this help message

Options:
    -m, --module     Module name (required for crud/entity)
    -e, --entity     Entity name (required for crud/entity)
    -f, --fields     Field definitions (optional)

Examples:
    # Generate complete CRUD
    ./scripts/generate.sh crud -m fleet -e Vehicle -f "Make:string:required,Model:string:required"

    # Generate entity only
    ./scripts/generate.sh entity -m crm -e Contact -f "Name:string:required,Email:string:email"

    # Generate migration
    ./scripts/generate.sh migration

Field Format:
    FieldName:Type:Validation

    Types: string, int, int64, float64, bool, time.Time, uuid.UUID
    Validation: required, min=N, max=N, len=N, email, url

EOF
}

if [ $# -eq 0 ] || [ "$1" = "help" ] || [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    show_help
    exit 0
fi

COMMAND=$1
shift

MODULE=""
ENTITY=""
FIELDS=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--module)
            MODULE="$2"
            shift 2
            ;;
        -e|--entity)
            ENTITY="$2"
            shift 2
            ;;
        -f|--fields)
            FIELDS="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

case $COMMAND in
    crud)
        if [ -z "$MODULE" ] || [ -z "$ENTITY" ]; then
            echo "Error: -m (module) and -e (entity) are required for crud generation"
            exit 1
        fi
        
        CMD="go run cmd/codegen/main.go -type=crud -module=$MODULE -entity=$ENTITY"
        if [ -n "$FIELDS" ]; then
            CMD="$CMD -fields=\"$FIELDS\""
        fi
        
        echo "Generating CRUD for $ENTITY in $MODULE module..."
        eval $CMD
        ;;
        
    entity)
        if [ -z "$MODULE" ] || [ -z "$ENTITY" ]; then
            echo "Error: -m (module) and -e (entity) are required for entity generation"
            exit 1
        fi
        
        CMD="go run cmd/codegen/main.go -type=entity -module=$MODULE -entity=$ENTITY"
        if [ -n "$FIELDS" ]; then
            CMD="$CMD -fields=\"$FIELDS\""
        fi
        
        echo "Generating entity $ENTITY in $MODULE module..."
        eval $CMD
        ;;
        
    migration)
        echo "Generating migration file..."
        go run cmd/codegen/main.go -type=migration
        ;;
        
    *)
        echo "Unknown command: $COMMAND"
        show_help
        exit 1
        ;;
esac

echo ""
echo "✓ Generation complete!"
