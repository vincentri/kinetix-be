#/bin/sh

go run github.com/steebchen/prisma-client-go db push --force-reset
go run ./seed/main.go