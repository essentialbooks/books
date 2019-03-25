#!/bin/bash
set -u -e -o pipefail
find . -name "*.go" | xargs wc -l
echo ""
wc -l "tmpl/app.js"
echo ""
wc -l "tmpl/main.css"

