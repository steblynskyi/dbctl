package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/steblynskyi/dbpermissionmanagement/utils"
	"github.com/spf13/cobra"
)

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Get Roles related info",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// flags := RolesFlags{
		// 	RoleType: cmd.Flags().Lookup("role-type").Value.String(),
		// 	RoleName: cmd.Flags().Lookup("name").Value.String(),
		// 	ListPermissions, _: cmd.Flags().GetBool("list-permissions"),

		// }
		SetGlobalFlags(cmd)
		flags.RoleName, _ = cmd.Flags().GetString("role-name")
		flags.DbName, _ = cmd.Flags().GetString("db-name")
		flags.ListPermissions, _ = cmd.Flags().GetBool("list-permissions")

		//fmt.Printf("flags struct: %+v\n", flags)
		if flags.DbName == "" && flags.RoleName == "" && !flags.ListPermissions {
			// If all flags are NOT set -> list server roles
			getServerRoles(flags)
		} else if flags.DbName != "" && flags.RoleName != "" && flags.ListPermissions {
			//If n,l, d ALL are set -> listDBRolePermissions
			getDBRolePermissions(flags)
		} else if flags.DbName == "" && flags.RoleName != "" && !flags.ListPermissions {
			//If only n is set -> get members of server Role
			getServerRoleMembers(flags)
		} else if flags.DbName == "" && flags.RoleName != "" && flags.ListPermissions {
			//If n & l is set -> listServerRolePermissions
			getServerRolePermissions(flags)
		} else if flags.DbName != "" && flags.RoleName == "" && !flags.ListPermissions {
			//If ONLY d is set -> listDBRoles
			getDBRoles(flags)
		} else if flags.DbName != "" && flags.RoleName != "" && !flags.ListPermissions {
			//If ONLY n and d is set -> get members of DB Role
			getDBRoleMembers(flags)
		}

	},
}

var (
	listPer          bool
	roleName, dbName string
	flags            RolesFlags
)

type RolesFlags struct {
	RoleName        string
	DbName          string
	ListPermissions bool
}

func init() {
	getCmd.AddCommand(rolesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rolesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rolesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rolesCmd.Flags().StringVarP(&dbName, "db-name", "d", "", "DB name")
	rolesCmd.Flags().StringVarP(&roleName, "role-name", "n", "", "Role name(Server Role/DB Role)")
	rolesCmd.Flags().BoolVarP(&listPer, "list-permissions", "l", false, "List DB/Server Role permissions")

}

// getDBRoles returns list of DB Roles in specific DB
func getDBRoles(flags RolesFlags) {
	fmt.Printf("Fetching list of DB Roles in %v DB...\n", flags.DbName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	RoleFlagsMap := map[string]interface{}{
		"DbName": flags.DbName,
	}

	db.Namespace = "mssql/user"
	db.TemplateName = "getDBRoles"
	db.Parameters = RoleFlagsMap
	db.RowReturn = true

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to get DB Roles:  %v\n", err)
	}
}

// getServerRoles returns list of server roles in the DB Host
func getServerRoles(flags RolesFlags) {
	fmt.Printf("Fetching list of Server Roles...\n")

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB => %v\n", err)
		os.Exit(1)
	}
	if utils.Cfg.DbType == "mssql" {
		db.Namespace = "mssql/user"
		db.TemplateName = "getServerRoles"
	} else if utils.Cfg.DbType == "postgres" {
		db.Namespace = "postgres/user"
		db.TemplateName = "getPgRoles"
	}

	// db.Namespace = "mssql/user"
	// db.TemplateName = "getServerRoles"
	db.Parameters = nil
	db.RowReturn = true

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to get Server Roles:  %v\n", err)
	}
}

// getDBRolePermissions returns specific DB Role permissions
func getDBRolePermissions(flags RolesFlags) {
	fmt.Printf("Fetching permissions of DB Role '%v' defined in '%v' DB\n", flags.RoleName, flags.DbName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	RoleFlagsMap := map[string]interface{}{
		"RoleName": flags.RoleName,
		"DbName":   flags.DbName,
	}

	db.Namespace = "mssql/user"
	db.TemplateName = "getDBRolePermissions"
	db.Parameters = RoleFlagsMap
	db.RowReturn = true

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to get DB Role permissions:  %v\n", err)
	}

}

// getServerRolePermissions returns specific Server Role permissions
func getServerRolePermissions(flags RolesFlags) {
	fmt.Printf("Fetching permissions of Server Role '%v'.\n", flags.RoleName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	RoleFlagsMap := map[string]interface{}{
		"RoleName": flags.RoleName,
	}

	db.Namespace = "mssql/user"
	db.TemplateName = "getServerRolePermissions"
	db.Parameters = RoleFlagsMap
	db.RowReturn = true

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to get permission of Server Role:  %v\n", err)
	}

}

// getServerRoleMembers returns specific Server Role members
func getServerRoleMembers(flags RolesFlags) {
	fmt.Printf("Fetching members of Server Role '%v'...\n", flags.RoleName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	RoleFlagsMap := map[string]interface{}{
		"RoleName": flags.RoleName,
	}

	if utils.Cfg.DbType == "mssql" {
		db.Namespace = "mssql/user"
		db.TemplateName = "getServerRoleMembers"
	} else if utils.Cfg.DbType == "postgres" {
		db.Namespace = "postgres/user"
		db.TemplateName = "getPgServerRoleMembers"
	}

	// db.Namespace = "mssql/user"
	// db.TemplateName = "getServerRoleMembers"
	db.Parameters = RoleFlagsMap
	db.RowReturn = true

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to get Server Role members:  %v\n", err)
	}

}

// getDBRoleMembers returns specific DB Role members
func getDBRoleMembers(flags RolesFlags) {
	fmt.Printf("Fetching members of DB Role '%v' defined in '%v' DB...\n", flags.RoleName, flags.DbName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	RoleFlagsMap := map[string]interface{}{
		"RoleName": flags.RoleName,
		"DbName":   flags.DbName,
	}

	db.Namespace = "mssql/user"
	db.TemplateName = "getDBRoleMembers"
	db.Parameters = RoleFlagsMap
	db.RowReturn = true

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to get DB Role members:  %v\n", err)
	}

}
