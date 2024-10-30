package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/steblynskyi/dbpermissionmanagement/utils"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// sanitationCmd represents the sanitation command
var sanitationCmd = &cobra.Command{
	Use:   "sanitation",
	Short: "Run sanitation scripts for specific DB/Tables",
	Long: `
sanitation command allows to run sanitation scripts on specific DB or table.
Scripts are stored here: https://bitbucket.org/steblynskyi/dbctl/src/develop/SanitationScripts/

For example:
1. To run sanitation scripts on inncenter01p DB:
./dbctl run sanitation --db-type=mssql --db-name=inncenter01p
	
2. To run sanitation for specific table "client_merchantaccount_xref" of inncenter01p DB:
./dbctl run sanitation --db-type=mssql --db-name=inncenter --table-name="client_merchantaccount_xref" --schema-name="dbo"`,
	Run: func(cmd *cobra.Command, args []string) {
		SetGlobalFlags(cmd)

		sflags.DbName, _ = cmd.Flags().GetString("db-name")
		sflags.SchemaName, _ = cmd.Flags().GetString("schema-name")
		sflags.TableName, _ = cmd.Flags().GetString("table-name")
		sflags.AllDB, _ = cmd.Flags().GetBool("all-db")

		if sflags == (SanitationFlags{}) {
			fmt.Printf("To run sanitation, Please set either of the below flags:\n 1) --db-name OR \n2) --all-db OR \n3) --table-name and --schema-name flags\n")
		} else if sflags.DbName == "" && sflags.TableName != "" && sflags.SchemaName != "" {
			fmt.Printf("Please set --db-name flag as well")
		} else if sflags.DbName != "" && sflags.TableName == "" && sflags.SchemaName == "" {
			RunSanityForSpecificDB(sflags)
		} else if sflags.DbName != "" && sflags.TableName != "" && sflags.SchemaName != "" {
			RunSanityForSpecificTable(sflags)
		} else if sflags.AllDB {
			RunSanityForAllDB(sflags)
		}
		//runSanitationScript(DbName)
	},
}

var (
	schemaName, tableName string
	allDB                 bool
	sflags                SanitationFlags
)

type SanitationFlags struct {
	DbName     string
	TableName  string
	SchemaName string
	AllDB      bool
}

func init() {
	runCmd.AddCommand(sanitationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sanitationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sanitationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sanitationCmd.Flags().StringVar(&dbName, "db-name", "", "Run sanitation for mentioned DB name")
	sanitationCmd.Flags().StringVar(&schemaName, "schema-name", "", "Run sanitation for specific Schema name (use along with table-name)")
	sanitationCmd.Flags().StringVar(&tableName, "table-name", "", "Run sanitation for specific Table name (use along with schema-name)")
	sanitationCmd.Flags().BoolVar(&allDB, "all-db", false, "Run Sanitation for All DBs")

	//sanitationCmd.Flags().StringVar(&dbName, "script-name", "", "Sanitation script name")
	//sanitationCmd.MarkFlagRequired("db-name")
	sanitationCmd.MarkFlagsRequiredTogether("table-name", "schema-name")
}

// RunSanityForSpecificTable runs santitation for specific table in specific DB
func RunSanityForSpecificTable(flags SanitationFlags) {
	fmt.Printf("Running Sanitation for table %v inside DB %v\n", flags.SchemaName+"."+flags.TableName, flags.DbName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	if utils.Cfg.DbType == "mssql" {
		db.Namespace = "mssql/" + flags.DbName
		db.TemplateName = flags.SchemaName + "." + flags.TableName
	} else if utils.Cfg.DbType == "postgres" {
		db.Namespace = "postgres/" + flags.DbName
		db.TemplateName = flags.SchemaName + "." + flags.TableName
	}

	db.Parameters = nil
	db.RowReturn = false

	//fmt.Printf("db is %+v\n", db)
	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to run Sanitation script for DB:  %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully completed sanitation for for table %v inside DB %v\n", flags.SchemaName+"."+flags.TableName, flags.DbName)

}

// RunSanityForSpecificDB runs sanity for specific DB(includes all the required tables)
func RunSanityForSpecificDB(flags SanitationFlags) {
	fmt.Printf("Runing sanitation for tables in DB %v\n", flags.DbName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	//Only run go ahead with sanitation if user enters 'yes'
	//otherwise cancel sanitation
	confirmation := utils.AskForConfirmation()
	if strings.ToLower(confirmation) != "yes" {
		fmt.Printf("query apply cancelled")
		os.Exit(1)
	}

	//Measure sanity execution run time
	startTime := time.Now()

	//Run sanity for all tables in parallel
	err = db.RunConcurrentSanityScripts(flags.DbName)

	endTime := time.Now()
	utils.Result.TotalExecutionTime = endTime.Sub(startTime)

	// For printing Sanitation script Execution Result in table format
	utils.PrintFinalSummary()
	fmt.Println()
	utils.PrintPerTableExecutionTime()
	if err != nil {
		cc := text.Colors{text.FgRed}

		fmt.Println(cc.Sprintf("\nFAIL - Sanitation FAILED for %v, Error: %v \n\n", flags.DbName, err))
		os.Exit(1)
	}
	cc := text.Colors{text.FgGreen}
	fmt.Println(cc.Sprintf("\n SUCCESS! - Sanitation successful for %v ", flags.DbName))

}

// RunSanityForAllDB runs sanity for all DB in specific host
func RunSanityForAllDB(flags SanitationFlags) {
	fmt.Println("Inside All DBs")
}
