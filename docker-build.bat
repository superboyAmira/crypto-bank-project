@echo off
REM Docker build and run script for Windows
echo Crypto Bank Project - Docker Management

:menu
echo.
echo ================================
echo   Docker Management Menu
echo ================================
echo 1. Build all services
echo 2. Start all services
echo 3. Stop all services
echo 4. View logs
echo 5. Restart services
echo 6. Clean all (remove containers and volumes)
echo 7. Exit
echo ================================
echo.

set /p choice="Enter your choice (1-7): "

if "%choice%"=="1" goto build
if "%choice%"=="2" goto start
if "%choice%"=="3" goto stop
if "%choice%"=="4" goto logs
if "%choice%"=="5" goto restart
if "%choice%"=="6" goto clean
if "%choice%"=="7" goto end

echo Invalid choice. Please try again.
goto menu

:build
echo Building Docker images...
docker-compose build
echo Build completed!
pause
goto menu

:start
echo Starting all services...
docker-compose up -d
echo Services started!
docker-compose ps
pause
goto menu

:stop
echo Stopping all services...
docker-compose down
echo Services stopped!
pause
goto menu

:logs
echo Showing logs (Ctrl+C to stop)...
docker-compose logs -f
goto menu

:restart
echo Restarting services...
docker-compose restart
echo Services restarted!
pause
goto menu

:clean
echo WARNING: This will remove all containers, volumes, and images!
set /p confirm="Are you sure? (yes/no): "
if /i "%confirm%"=="yes" (
    docker-compose down -v --rmi all
    echo Cleanup completed!
) else (
    echo Cleanup cancelled.
)
pause
goto menu

:end
echo Goodbye!
exit /b 0

