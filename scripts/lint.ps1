function Install-If-Not-Exist {
    param (
        [string]$ToolName,
        [string]$InstallUrl
    )

    if (Get-Command $ToolName -ErrorAction SilentlyContinue) {
        Write-Host "$ToolName is already installed."
    } else {
        Write-Host "$ToolName is not installed. Installing..."
        go install $InstallUrl
    }
}

Install-If-Not-Exist "go-cleanarch" "github.com/roblaszczak/go-cleanarch@latest"

$LINT_VERSION = "1.54.0"
$NEED_INSTALL = $false

if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
    $CURRENT_VERSION = (golangci-lint --version | Select-String -Pattern "\d+\.\d+\.\d+").Matches.Value
    Write-Host "golangci-lint v$CURRENT_VERSION already installed."
    if ($CURRENT_VERSION -eq $LINT_VERSION) {
        $NEED_INSTALL = $false
    } else {
        $NEED_INSTALL = $true
    }
} else {
    $NEED_INSTALL = $true
}

if ($NEED_INSTALL) {
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" -OutFile install.sh
    bash install.sh -b $(go env GOPATH)/bin v$LINT_VERSION
    Remove-Item install.sh
}

go-cleanarch

Write-Host "lint modules:"
Write-Host "$(modules)"

goimports -w -l .

foreach ($module in $(modules)) {
    Set-Location ./internal/$module
    golangci-lint run --config "$ROOT_DIR/.golangci.yaml"
    Set-Location -
}