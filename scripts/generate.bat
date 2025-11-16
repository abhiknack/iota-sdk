@echo off
REM IOTA SDK Code Generator Helper Script (Windows)
REM Usage: scripts\generate.bat [command] [options]

setlocal enabledelayedexpansion

if "%1"=="" goto :help
if "%1"=="help" goto :help
if "%1"=="--help" goto :help
if "%1"=="-h" goto :help

set COMMAND=%1
shift

set MODULE=
set ENTITY=
set FIELDS=

:parse_args
if "%1"=="" goto :execute
if "%1"=="-m" (
    set MODULE=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="--module" (
    set MODULE=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="-e" (
    set ENTITY=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="--entity" (
    set ENTITY=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="-f" (
    set FIELDS=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="--fields" (
    set FIELDS=%2
    shift
    shift
    goto :parse_args
)
echo Unknown option: %1
goto :help

:execute
if "%COMMAND%"=="crud" goto :crud
if "%COMMAND%"=="entity" goto :entity
if "%COMMAND%"=="migration" goto :migration
echo Unknown command: %COMMAND%
goto :help

:crud
if "%MODULE%"=="" (
    echo Error: -m ^(module^) and -e ^(entity^) are required for crud generation
    exit /b 1
)
if "%ENTITY%"=="" (
    echo Error: -m ^(module^) and -e ^(entity^) are required for crud generation
    exit /b 1
)

echo Generating CRUD for %ENTITY% in %MODULE% module...
if "%FIELDS%"=="" (
    go run cmd/codegen/main.go -type=crud -module=%MODULE% -entity=%ENTITY%
) else (
    go run cmd/codegen/main.go -type=crud -module=%MODULE% -entity=%ENTITY% -fields="%FIELDS%"
)
goto :done

:entity
if "%MODULE%"=="" (
    echo Error: -m ^(module^) and -e ^(entity^) are required for entity generation
    exit /b 1
)
if "%ENTITY%"=="" (
    echo Error: -m ^(module^) and -e ^(entity^) are required for entity generation
    exit /b 1
)

echo Generating entity %ENTITY% in %MODULE% module...
if "%FIELDS%"=="" (
    go run cmd/codegen/main.go -type=entity -module=%MODULE% -entity=%ENTITY%
) else (
    go run cmd/codegen/main.go -type=entity -module=%MODULE% -entity=%ENTITY% -fields="%FIELDS%"
)
goto :done

:migration
echo Generating migration file...
go run cmd/codegen/main.go -type=migration
goto :done

:help
echo IOTA SDK Code Generator
echo.
echo Usage: scripts\generate.bat [command] [options]
echo.
echo Commands:
echo     crud        Generate complete CRUD (domain, repo, service, controller, DTOs)
echo     entity      Generate domain aggregate only
echo     migration   Generate migration file template
echo     help        Show this help message
echo.
echo Options:
echo     -m, --module     Module name (required for crud/entity)
echo     -e, --entity     Entity name (required for crud/entity)
echo     -f, --fields     Field definitions (optional)
echo.
echo Examples:
echo     # Generate complete CRUD
echo     scripts\generate.bat crud -m fleet -e Vehicle -f "Make:string:required,Model:string:required"
echo.
echo     # Generate entity only
echo     scripts\generate.bat entity -m crm -e Contact -f "Name:string:required,Email:string:email"
echo.
echo     # Generate migration
echo     scripts\generate.bat migration
echo.
echo Field Format:
echo     FieldName:Type:Validation
echo.
echo     Types: string, int, int64, float64, bool, time.Time, uuid.UUID
echo     Validation: required, min=N, max=N, len=N, email, url
exit /b 0

:done
echo.
echo ✓ Generation complete!
exit /b 0
