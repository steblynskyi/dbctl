package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/steblynskyi/dbpermissionmanagement/utils"
	"github.com/spf13/cobra"
)

// dbsCmd represents the dbs command
var dbsCmd = &cobra.Command{
	Use:   "dbs",
	Short: "Get DB related info",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		SetGlobalFlags(cmd)
		GetAllDBNames()
	},
}

func init() {
	getCmd.AddCommand(dbsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetAllDBNames() {
	fmt.Printf("Fetching list of DB...\n")

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	if utils.Cfg.DbType == "mssql" {
		db.Namespace = "mssql/user"
		db.TemplateName = "getAllDbNames"
	} else if utils.Cfg.DbType == "postgres" {
		db.Namespace = "postgres/user"
		db.TemplateName = "getAllPgDbNames"
	}

	db.Parameters = nil
	db.RowReturn = true

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to get Server Roles:  %v\n", err)
	}
}
