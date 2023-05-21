# Nextlinux FP Management Corrections Tool

This is meant as a temporary tool to facilitate the adding of multiple
package corrections to Nextlinux.

Long term, Nextlinux is working on an Enteprise CLI called `nextlinuxctl` that will replace
the functionality of this tool. Once that is released (coming soon), this tool will become deprecated.

## Prerequisites
This script has two dependencies:
* go version 1.16
* dependencies download. This can be done by running `go mod download` in the root of the repository
* edit the [config](./fpmanagement.yaml) and add your Nextlinux deployment information

Execute this script with:
```shell script
go run main.go
```

## Correction Content

The [corrections.json](./corrections.json) contains a list of Nextlinux False-Positive Management Package Corrections for various
common java packages. This can be customized (but the schema should be maintained) as needed.
