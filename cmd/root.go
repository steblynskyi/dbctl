package cmd

import (
	"os"

	"bitbucket.org/steblynskyi/dbpermissionmanagement/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbctl",
	Short: "DB Permission Management Tool",
	Long:  `dbctl is a CLI tool for creating and updating DB user, DB Roles (Server Role and DB Role) in MSSQL Server.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	//err := rootCmd.Execute()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	host, usrname, passwd, port, dbname, dbtype string
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dbpermissionmanagement.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&dbname, "database-name", "D", "", "DB name")
	rootCmd.PersistentFlags().StringVarP(&usrname, "username", "U", "", "DB Usenname")
	rootCmd.PersistentFlags().StringVarP(&passwd, "password", "P", "", "DB Password")
	rootCmd.PersistentFlags().StringVarP(&host, "db-hostname", "H", "", "DB Hostname")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "Port")
	rootCmd.PersistentFlags().StringVarP(&dbtype, "db-type", "T", "", "DB Type")

	rootCmd.MarkPersistentFlagRequired("db-type")
}

// SetGlobalFlags sets the DB credentials/connection related flags
func SetGlobalFlags(cmd *cobra.Command) {

	utils.Cfg.DbName, _ = cmd.Flags().GetString("database-name")
	utils.Cfg.Username, _ = cmd.Flags().GetString("username")
	utils.Cfg.Password, _ = cmd.Flags().GetString("password")
	utils.Cfg.Host, _ = cmd.Flags().GetString("db-hostname")
	utils.Cfg.Port, _ = cmd.Flags().GetString("port")
	utils.Cfg.DbType, _ = cmd.Flags().GetString("db-type")

	//fmt.Printf("utils.Cfg is: %+v\n", utils.Cfg)
}
