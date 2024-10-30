package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:       "get users|roles|dbs",
	Short:     "Get details of user, role or DB",
	ValidArgs: []string{"users", "roles", "dbs"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Long: `Display details of one or many DB user/Roles/DBs.
Prints a table of the most important information about the specified resources.
Usage: dbctl get roles --name=<role_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// getCmd.Flags().StringVarP(&name, "name", "n", "", "DB User or Role name")
	// getCmd.Flags().StringVarP(&listPer, "list-permissions", "l", "", "List DB user/Role permissions")
	//getCmd.MarkFlagsRequiredTogether("type", "name")

}
