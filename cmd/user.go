package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/steblynskyi/dbpermissionmanagement/utils"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Create DB User",
	// ValidArgs: []string{"create", "update"},
	Long: `
This command creates DB user.
Usage: dbctl create user --db-username=<username> --db-password=<password>`,
	//Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		SetGlobalFlags(cmd)

		username, _ = cmd.Flags().GetString("db-username")
		password, _ = cmd.Flags().GetString("db-password")

		CreateUser(username, password)

	},
}

var username, password string

type User struct {
	Username  string
	Password  string
	DefaultDb string
}

func init() {
	createCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	userCmd.Flags().StringVar(&username, "db-username", "", "DB Username")
	userCmd.Flags().StringVar(&password, "db-password", "", "Password for DB user")

	userCmd.MarkFlagRequired("db-username")
	userCmd.MarkFlagRequired("db-password")

}

func CreateUser(username string, password string) {

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	//fill up query params with struct
	userDetails := User{
		Username:  username,
		Password:  password,
		DefaultDb: "inncenter01p",
	}

	db.Namespace = "mssql/user"
	db.TemplateName = "createUser"
	db.Parameters = userDetails
	db.RowReturn = false

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to create user:  %v\n", err)
	}
}
