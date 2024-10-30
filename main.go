package main

import (
	"embed"

	"bitbucket.org/steblynskyi/dbpermissionmanagement/cmd"
	"bitbucket.org/steblynskyi/dbpermissionmanagement/utils"
)

//go:embed templates/* SanitationScripts/*
var templateFs embed.FS

func main() {

	utils.TemplateFs = templateFs
	cmd.Execute()

}
