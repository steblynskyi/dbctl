package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/steblynskyi/dbpermissionmanagement/utils"
	"github.com/spf13/cobra"
)

// roleCmd represents the role command
var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Update Role",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		SetGlobalFlags(cmd)

		roleflags.RoleName, _ = cmd.Flags().GetString("role-name")
		roleflags.DbName, _ = cmd.Flags().GetString("db-name")
		roleflags.AddMember, _ = cmd.Flags().GetString("add-member")
		roleflags.DropMember, _ = cmd.Flags().GetString("drop-member")

		//fmt.Printf("flags struct: %+v\n", roleflags)

		if roleflags.DbName != "" && roleflags.AddMember != "" {
			// If a and d are set -> Call addMemberDbRole
			addMemberDbRole()
		} else if roleflags.DbName != "" && roleflags.DropMember != "" {
			removeMemberDbRole()
		}
	},
}

var (
	addMember, dropMember string
	roleflags             UpdateRolesFlags
)

type UpdateRolesFlags struct {
	RoleName   string
	DbName     string
	AddMember  string
	DropMember string
}

func init() {
	updateCmd.AddCommand(roleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// roleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// roleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	roleCmd.Flags().StringVarP(&dbName, "db-name", "d", "", "DB name")
	roleCmd.Flags().StringVarP(&roleName, "role-name", "n", "", "Role name(DB Role/Server Role)")
	roleCmd.Flags().StringVarP(&addMember, "add-member", "a", "", "Add member to Role")
	roleCmd.Flags().StringVarP(&dropMember, "drop-member", "r", "", "Drop member from Role")

	roleCmd.MarkFlagRequired("role-name")
	roleCmd.MarkFlagsMutuallyExclusive("add-member", "drop-member")

}

func addMemberDbRole() {
	fmt.Printf("Adding '%v' in DB Role '%v'.\n", roleflags.AddMember, roleflags.RoleName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	RoleFlagsMap := map[string]interface{}{
		"RoleName": roleflags.RoleName,
		"DbName":   roleflags.DbName,
		"Username": roleflags.AddMember,
	}

	db.Namespace = "mssql/user"
	db.TemplateName = "addMemberToDBRole"
	db.Parameters = RoleFlagsMap
	db.RowReturn = false

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to add member to DB Role:  %v\n", err)
	}

	fmt.Printf("Successfully added '%v' in DB Role '%v'.\n", roleflags.AddMember, roleflags.RoleName)
}

func removeMemberDbRole() {

	fmt.Printf("Removing '%v' from DB Role '%v'.\n", roleflags.DropMember, roleflags.RoleName)

	db, err := utils.InitDB()
	if err != nil {
		fmt.Printf("failed to connect with DB %v\n", err)
		os.Exit(1)
	}

	RoleFlagsMap := map[string]interface{}{
		"RoleName": roleflags.RoleName,
		"DbName":   roleflags.DbName,
		"Username": roleflags.DropMember,
	}

	db.Namespace = "mssql/user"
	db.TemplateName = "dropMemberFromDBRole"
	db.Parameters = RoleFlagsMap
	db.RowReturn = false

	err = db.CommonQueryExecutionTemplate()
	if err != nil {
		fmt.Printf("Failed to add member to DB Role:  %v\n", err)
	}

	fmt.Printf("\nSuccessfully removed '%v' in DB Role '%v'.\n", roleflags.DropMember, roleflags.RoleName)
}
