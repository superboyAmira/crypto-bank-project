@echo off
REM Clean rebuild script for Windows
echo ====================================
echo   Docker Clean Rebuild
echo ====================================
echo.
echo This will:
echo   1. Stop all containers
echo   2. Remove all images
echo   3. Clear build cache
echo   4. Rebuild from scratch
echo.
set /p confirm="Continue? (yes/no): "

if /i not "%confirm%"=="yes" (
    echo Cancelled.
    exit /b 0
)

echo.
echo Step 1: Stopping containers...
docker-compose down

echo.
echo Step 2: Removing images...
docker-compose down --rmi all

echo.
echo Step 3: Clearing build cache...
docker builder prune -af

echo.
echo Step 4: Rebuilding without cache...
docker-compose build --no-cache

echo.
echo ====================================
echo   Rebuild Complete!
echo ====================================
echo.
echo To start services, run:
echo   docker-compose up -d
echo.
pause

