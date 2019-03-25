#!/usr/bin/env pwsh

# you can pass additional args like:
# -update-go-deps
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

go build -o ./gen.exe
exitIfFailed

./gen.exe -analytics UA-113489735-1 $args
exitIfFailed

Remove-Item -Force -ErrorAction SilentlyContinue ./gen.exe
netlify deploy --prod --dir=www --site=7df32685-1421-41cf-937a-a92fde6725f4 --open
