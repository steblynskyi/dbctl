package utils

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/sync/errgroup"
)

// ExecutionSummary used to print final execution summary of sanitation scripts
type ExecutionSummary struct {
	DbName             string
	SuccessTablesCount int
	FailedTablesCount  int
	SuccessTableNames  []string
	FailedTableNames   []string
	ErrorList          []string
	TotalExecutionTime time.Duration
}

var Result ExecutionSummary

// RunConcurrentSanityScripts runs sanitation scripts for multiple tables in parallel with ErrGroup goroutine
func (d Database) RunConcurrentSanityScripts(DbName string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	d.Namespace = "SanitationScripts/" + Cfg.DbType + "/" + DbName
	templateNames, err := GetTemplateNames(d.Namespace)
	if err != nil {
		return fmt.Errorf("failed to get templateNames: %w", err)
	}

	//Kept for testing purpose, Uncomment while testing with SELECT query templates
	//templateNames = []string{"getServerRoles", "getAllUsers", "getAllDbNames", "getServerRolePermissions"}
	//for testing on postgresql
	//templateNames = []string{"getPgVersion", "getPgRoles", "getAllPgDbNames"}

	fmt.Printf("\nRunning sanitation script for tables: %v\n", strings.Join(templateNames, ", "))

	for _, templateName := range templateNames {
		templateName := templateName
		eg.Go(func() error {

			d.TemplateName = templateName
			d.Parameters = nil
			Result.DbName = DbName

			return d.SanitationQueryExecutionTemplate()
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to execute template: %v %w", d.TemplateName, err)
	}

	return nil
}

// GetTemplateNames returns list of template names in specific directory
// Used to get list of sanitation script templates for specific DB
func GetTemplateNames(dir string) ([]string, error) {

	// Parse all the templates in the directory
	templates, err := template.ParseGlob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse all template names: %v", err)
	}

	// Create a slice of template names
	var templateNames []string
	for _, t := range templates.Templates() {

		//as Templates() returns filename.sql as well, ignore those names while adding in templateNames
		if strings.Contains(t.Name(), ".sql") {
			continue
		}
		templateNames = append(templateNames, t.Name())
	}

	return templateNames, nil
}

// PrintFinalSummary prints Sanitation script Execution Result Summary in table format
func PrintFinalSummary() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"Parameter", "Value"}, rowConfigAutoMerge)
	t.SetAutoIndex(true)

	t.SetTitle("Sanitation Script Execution Summary: " + Result.DbName)

	t.AppendRow(table.Row{"Successful Execution Tables Count", Result.SuccessTablesCount})
	t.AppendRow(table.Row{"Failed Execution Tables Count", Result.FailedTablesCount})
	t.AppendRow(table.Row{"Successful Execution Tablenames", strings.Join(Result.SuccessTableNames, ", ")})
	t.AppendRow(table.Row{"Failed Execution Tablenames", strings.Join(Result.FailedTableNames, ", ")})
	t.AppendRow(table.Row{"Total Execution Time", Result.TotalExecutionTime.String()})
	for _, errMsg := range Result.ErrorList {
		t.AppendRow(table.Row{"Errors", errMsg}, rowConfigAutoMerge)
	}

	//t.AppendRow(table.Row{"Errors", utils.Result.ErrorList}, rowConfigAutoMerge)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 90, AutoMerge: true},
		{Number: 2, WidthMax: 100},
	})

	t.Render()

	//fmt.Printf("PerQueryTime: %v\n", PerQueryTime)
}

// AskForConfirmation prints message and ask for user confirmation, before running
// actual sanitation script on the host
// Added to prevent sanitation on incorrect host
func AskForConfirmation() string {
	var confirmation string

	cc := text.Colors{text.FgGreen}
	fmt.Println(cc.Sprintf("\nRunning Query against DB: %v@%v\n", Cfg.Username, Cfg.Host))
	fmt.Printf("Do you want to run sanitation against this host? \nOnly 'yes' will be accepted to approve. \n\nEnter a value:")

	fmt.Scanln(&confirmation)

	return confirmation
}

func PrintPerTableExecutionTime() {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"Table Name", "Query Execution Time"}, rowConfigAutoMerge)
	t.SetAutoIndex(true)

	t.SetTitle("Per table exectution time summary: " + Result.DbName)

	PerQueryTime.Range(func(key, value any) bool {
		t.AppendRow(table.Row{key, value})
		return true
	})

	// for tableName, time := range PerQueryTime {
	// 	t.AppendRow(table.Row{tableName, time})
	// }

	t.Render()
}
