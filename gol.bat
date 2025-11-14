@echo off
REM gol.bat - Start IOTA SDK in developer mode with hot-reload in Docker
REM Usage: gol.bat [start|stop|restart|logs|rebuild|shell]

setlocal

if "%1"=="" goto start
if /i "%1"=="start" goto start
if /i "%1"=="stop" goto stop
if /i "%1"=="restart" goto restart
if /i "%1"=="logs" goto logs
if /i "%1"=="rebuild" goto rebuild
if /i "%1"=="shell" goto shell
if /i "%1"=="seed" goto seed
goto usage

:start
echo Starting IOTA SDK in developer mode with hot-reload...
echo.
echo Building dev image (first time only)...
docker compose -f compose.dev.yml build app
echo.
echo Starting all services (DB, Redis, App with hot-reload)...
echo Press Ctrl+C to stop
echo.
docker compose -f compose.dev.yml up app
goto end

:stop
echo Stopping all services...
docker compose -f compose.dev.yml down
echo Services stopped.
goto end

:restart
echo Restarting services...
docker compose -f compose.dev.yml down
timeout /t 2 /nobreak > nul
echo.
echo Starting all services...
docker compose -f compose.dev.yml up
goto end

:logs
echo Viewing application logs...
echo.
if exist "build-errors.log" (
    type build-errors.log
) else (
    echo No build-errors.log file found.
)
goto end

:rebuild
echo Rebuilding application in container...
echo.
echo Step 1: Generating templates...
docker compose -f compose.dev.yml run --rm app templ generate
if errorlevel 1 (
    echo Template generation failed!
    goto end
)
echo.
echo Step 2: Compiling CSS...
docker compose -f compose.dev.yml run --rm app make css
if errorlevel 1 (
    echo CSS compilation failed!
    goto end
)
echo.
echo Step 3: Verifying Go code...
docker compose -f compose.dev.yml run --rm app go vet ./...
if errorlevel 1 (
    echo Go vet found issues!
    goto end
)
echo.
echo Rebuild complete! Air will automatically reload.
goto end

:shell
echo Opening shell in container...
docker compose -f compose.dev.yml run --rm app /bin/sh
goto end

:seed
echo Seeding database with test data...
docker compose -f compose.dev.yml run --rm app go run cmd/command/main.go seed
echo.
echo Database seeded! You can now login with:
echo   Email: test@gmail.com
echo   Password: TestPass123!
goto end

:usage
echo Usage: gol.bat [command]
echo.
echo Commands:
echo   start    - Start database, Redis, and app with hot-reload (default)
echo   stop     - Stop all services
echo   restart  - Restart all services
echo   logs     - View Air build error logs
echo   rebuild  - Regenerate templates and CSS in container
echo   shell    - Open bash shell in container
echo   seed     - Seed database with test data
echo.
echo Examples:
echo   gol.bat           Start in dev mode
echo   gol.bat start     Start in dev mode
echo   gol.bat stop      Stop all services
echo   gol.bat restart   Restart everything
echo   gol.bat rebuild   Full rebuild in container
echo   gol.bat shell     Open container shell
echo   gol.bat seed      Create test user (test@gmail.com / TestPass123!)
goto end

:end
endlocal
