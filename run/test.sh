#!/usr/bin/env sh
mkdir -p build
rm -rf build/coverage*
touch build/coverage.out

go test ./... -cover -coverprofile=build/coverage.part -covermode=count -v ./

if [-$? -ne 0]; then
  exit
fi

echo "mode: count" >> build/coverage.out

grep -h -v "mode: count" build/coverage*.part >>build/coverage.out

go tool cover -func=build/coverage.out
go tool cover -html=build/coverage.out -o=build/coverage.html

echo "[coverage] Report at build/coverage.html"