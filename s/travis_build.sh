#!/bin/bash
set -u -e -o pipefail -o verbose

export GO111MODULE=on
go build -o gen-books

./gen-books -no-update-output -analytics UA-113489735-1

if [ -z ${NETLIFY_TOKEN+x} ]
then
    echo "Skipping upload because NETLIFY_TOKEN not set"
else
    netlify deploy --prod --dir=www --site=7df32685-1421-41cf-937a-a92fde6725f4 --auth=${NETLIFY_TOKEN}
fi
