@echo off
REM Build script for Windows
echo Building Crypto Bank Project for Windows...

echo.
echo Building bank-service...
cd bank-service
set CGO_ENABLED=0
set GOOS=windows
go build -ldflags="-w -s" -o main.exe ./cmd/server
cd ..

echo.
echo Building exchange-service...
cd exchange-service
set CGO_ENABLED=0
set GOOS=windows
go build -ldflags="-w -s" -o main.exe ./cmd/server
cd ..

echo.
echo Building analytics-service...
cd analytics-service
set CGO_ENABLED=0
set GOOS=windows
go build -ldflags="-w -s" -o main.exe ./cmd/server
cd ..

echo.
echo Building notification-service...
cd notification-service
set CGO_ENABLED=0
set GOOS=windows
go build -ldflags="-w -s" -o main.exe ./cmd/server
cd ..

echo.
echo Build completed successfully!
echo Executables created:
echo   - bank-service\main.exe
echo   - exchange-service\main.exe
echo   - analytics-service\main.exe
echo   - notification-service\main.exe
pause

