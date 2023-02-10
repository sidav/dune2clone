DATE=$(date +%d_%m_%y)
EXE_NAME=sidavRTS_${DATE}.exe
ZIP_NAME=sidavRTS_${DATE}.zip

echo "Building ${EXE_NAME}..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w" -o ${EXE_NAME} *.go

echo "Packing the data to ${ZIP_NAME}..."
zip -qr ${ZIP_NAME} ${EXE_NAME} resources

echo "Removing residual ${EXE_NAME}..."
rm ${EXE_NAME}

echo "Build completed."
