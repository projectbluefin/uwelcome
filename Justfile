version := `go run . version`

version :
    @echo "Current version: {{version}}"

compile :
    @echo "Compiling uwelcome..."

    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o uwelcome_{{version}}_linux_amd64
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o uwelcome_{{version}}_linux_arm64

    @echo "Built uwelcome v{{version}} for Linux on x86_64 and arm64"

translate lang :
    @./translators.sh {{lang}}