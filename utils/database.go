package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	sqlTemplate "github.com/NicklasWallgren/sqlTemplate/pkg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sethvargo/go-envconfig"

	"embed"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	Cfg        DbConfig
	TemplateFs embed.FS
	dbContext  = context.Background()
	ErrorList  []string
	//PerQueryTime = map[string]string{}
	//Mutex        = &sync.Mutex{}
	PerQueryTime = sync.Map{}
)

type Database struct {
	DB           *sql.DB
	Namespace    string
	TemplateName string
	Parameters   interface{}
	RowReturn    bool
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		DB: db,
	}
}

type DbConfig struct {
	Host     string `env:"DB_HOST"`
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DbName   string `env:"DB_NAME"`
	Port     string `env:"DB_PORT"`
	DbType   string `env:"DB_TYPE"`
	// ConnectionString string `env:"ConnectionString__Inncenter,required"`
}

// GetEnvVar parses env variables
func GetEnvVar() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &Cfg); err != nil {
		log.Fatal(err)
	}
}

// Open a database connection based on the configuration.
func Open(DbType string) (*sql.DB, error) {

	var ConnectionString string
	if DbType == "mssql" {
		ConnectionString = fmt.Sprintf("server=%v;user id=%v;password=%v;port=%v;database=%v", Cfg.Host, Cfg.Username, Cfg.Password, Cfg.Port, Cfg.DbName)
	} else if DbType == "postgres" {
		DbType = "pgx"
		ConnectionString = fmt.Sprintf("host=%v user=%v password=%v port=%v dbname=%v", Cfg.Host, Cfg.Username, Cfg.Password, Cfg.Port, Cfg.DbName)
	}
	//ConnectionString := fmt.Sprintf("server=%v;user id=%v;password=%v;port=%v;database=%v", Cfg.Host, Cfg.Username, Cfg.Password, Cfg.Port, Cfg.DbName)
	fmt.Printf("Connection string is: %s\n", ConnectionString)
	db, connectionError := sql.Open(DbType, ConnectionString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return nil, connectionError
	}

	return db, connectionError
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sql.DB) error {

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

// InitDB connects with DB and on successful connection returns DB object
func InitDB() (*Database, error) {
	//set env variables
	GetEnvVar()
	// var cfg DbConfig

	// if err := envconfig.Process(context.Background(), &cfg); err != nil {
	// 	log.Fatal(err)
	// }

	//Check if DbConfig is empty(if no credentials are provided)
	if Cfg == (DbConfig{}) {
		return nil, fmt.Errorf(`unable to locate DB credentials.
You can configure credentials by setting
1. env variable: DB_HOST DB_USER DB_PASSWORD DB_NAME DB_PORT
		OR
2. this flags: --username --password --db-hostname --database-name --port`)
	}

	//Connect with DB
	dbConn, err := Open(Cfg.DbType)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with DB %v", err)
	}

	//return db object
	db := NewDatabase(dbConn)

	return db, nil

}

// RunTransaction executes query as part of SQL transaction. It rollbacks the transaction in case of error
// otherwise commits the transaction
func (d Database) RunTransaction(query string, rowReturn bool) error {
	var rows *sql.Rows
	//Start the transaction
	txn, err := d.DB.BeginTx(dbContext, nil)
	if err != nil {
		return fmt.Errorf("failed to start the transaction: %w", err)
	}

	//Execute the query, based on return value
	//If query retruns row, execute it with QueryContext otherwise use ExecContext
	if rowReturn {
		rows, err = txn.QueryContext(dbContext, query)
		if err != nil {
			//Rollback the transaction if query returns an error
			txErr := txn.Rollback()
			return fmt.Errorf("failed to run query :%v %v", err, txErr)
		}
		defer rows.Close()

		// Iterate through the result set.
		var count int
		fmt.Printf("\nResults: TableName: %v\n========================\n", d.TemplateName)

		//get Column names
		columnNames, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to get column names: %v", err)
		}

		columnHeaders := []interface{}{}

		for _, col := range columnNames {
			columnHeaders = append(columnHeaders, col)
		}

		// For printing Query Result in table format
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleColoredBright)
		t.AppendHeader(columnHeaders)
		t.SetAutoIndex(true)
		t.SetTitle(d.TemplateName)
		for rows.Next() {
			currentRowColumns := make([]interface{}, len(columnNames))
			for i := range currentRowColumns {
				currentRowColumns[i] = new(interface{})
			}

			//Get Current Row Column values and copy to currentRowColumns
			err := rows.Scan(currentRowColumns...)
			if err != nil {
				return fmt.Errorf("failed to get column values with scan: %v", err)
			}

			//Get values for row in string format
			rowData := []interface{}{}
			for _, val := range currentRowColumns {
				rowData = append(rowData, fmt.Sprintf("%v", reflect.ValueOf(val).Elem()))
			}
			//Append current row values into table for display
			t.AppendRow(rowData)

			count++
		}
		//Render will print the final result table
		t.Render()

		fmt.Printf("DB Returned %v rows\n", count)

	} else {
		// Get the current time
		startTime := time.Now()
		_, err = txn.ExecContext(dbContext, query)
		endTime := time.Now()
		totalTime := endTime.Sub(startTime).String()

		// Mutex.Lock()
		// PerQueryTime[d.TemplateName] = totalTime
		// Mutex.Unlock()
		PerQueryTime.Store(d.TemplateName, totalTime)

		if err != nil {
			//Rollback the transaction if query returns an error
			txErr := txn.Rollback()
			return fmt.Errorf("failed to run query :%v %v", err, txErr)
		}
		fmt.Printf("Query Execution completed successfully: %v Time Took: %v\n", d.TemplateName, totalTime)
	}

	//Commit the transaction if query is executed successfully.
	if err = txn.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}

	return nil
}

func printExecutionPlan(query string) string {
	var confirmation string

	cc := text.Colors{text.FgGreen}
	fmt.Println(cc.Sprintf("\nRunning Query against DB: %v@%v\n", Cfg.Username, Cfg.Host))
	fmt.Println(cc.Sprintf("\n==== Query Execution Plan ==== %v\n", query))
	//fmt.Printf("\n==== Query Execution Plan ==== %v\n", query)

	fmt.Printf("Do you want to run these query? \ndbctl will run these query described above. \nOnly 'yes' will be accepted to approve. \n\nEnter a value:")

	fmt.Scanln(&confirmation)

	return confirmation
}

func ParseTemplate(namespace string, templateName string, parameters interface{}) (sqlTemplate.QueryTemplate, error) {
	//Templates parsing implementation using sqlt

	sqlt := sqlTemplate.NewQueryTemplateEngine()
	if err := sqlt.Register(namespace, TemplateFs, ".sql"); err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	var tmpl sqlTemplate.QueryTemplate
	var err error

	//Parse template & update named values with Struct/Map
	if reflect.ValueOf(parameters).Kind().String() == "struct" {
		tmpl, err = sqlt.ParseWithValuesFromStruct(namespace, templateName, parameters)
		if err != nil {
			return nil, fmt.Errorf("failed to parse templates with struct values: %w", err)
		}
	} else if reflect.ValueOf(parameters).Kind().String() == "map" {

		param := parameters.(map[string]interface{}) //Type assertion
		tmpl, err = sqlt.ParseWithValuesFromMap(namespace, templateName, param)
		if err != nil {
			return nil, fmt.Errorf("failed to parse templates with map values: %w", err)
		}
	} else {
		//Execute queries without any parameters (i.e. Execute query as it is)
		tmpl, err = sqlt.Parse(namespace, templateName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}
	}

	return tmpl, nil
}

// CommonQueryExecutionTemplate runs on each of the commands execution
func (d Database) CommonQueryExecutionTemplate() error {

	//Check if DB connection is still alive
	if err := d.DB.PingContext(dbContext); err != nil {
		return fmt.Errorf("error checking db connection: %w", err)
	}
	fmt.Printf("Successfully Connected with DB %v\n", Cfg.Host)

	//Parse template
	template, err := ParseTemplate(d.Namespace, d.TemplateName, d.Parameters)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	//Based on parsed template, print SQL query execution plan
	query := template.GetQuery()
	confirmation := printExecutionPlan(query)

	//Only run query(start the transaction) if user enters 'yes'
	//otherwise cancel query execution
	if strings.ToLower(confirmation) != "yes" {
		return fmt.Errorf("query apply cancelled")
	}

	//Start the transaction
	if err = d.RunTransaction(query, d.RowReturn); err != nil {
		return fmt.Errorf("failed to Run the transaction: %w", err)
	}

	return nil
}

// CommonQueryExecutionTemplate runs on each of the commands execution
func (d Database) SanitationQueryExecutionTemplate() error {
	//Check if DB connection is still alive
	if err := d.DB.PingContext(dbContext); err != nil {
		//errCh <- fmt.Sprintf("Template name %v: \n Error: %v",d.TemplateName, err)
		ErrorList = append(ErrorList, fmt.Sprintf("Template name %v: \n Error: %v", d.TemplateName, err))
		Result.ErrorList = append(Result.ErrorList, fmt.Sprintf("Template name %v: \n Error: %v", d.TemplateName, err))
		Result.FailedTableNames = append(Result.FailedTableNames, d.TemplateName)
		Result.FailedTablesCount++

		return fmt.Errorf("error checking db connection: %w", err)
	}
	//fmt.Printf("Successfully Connected with DB %v\n", Cfg.Host)

	//Parse template
	template, err := ParseTemplate(d.Namespace, d.TemplateName, d.Parameters)
	if err != nil {
		ErrorList = append(ErrorList, fmt.Sprintf("Template name %v: \n Error: %v", d.TemplateName, err))
		Result.ErrorList = append(Result.ErrorList, fmt.Sprintf("Template name %v: \n Error: %v", d.TemplateName, err))
		Result.FailedTableNames = append(Result.FailedTableNames, d.TemplateName)
		Result.FailedTablesCount++
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	//Based on parsed template, print SQL query execution plan
	query := template.GetQuery()

	//If SQL query contains SELECT statement, set rowReturn to true. Which uses txn.QueryContext
	//otherwise, in case of other statements like INSERT, UPDATE, DELETE or TRUNCATE set rowReturn to False. Which uses txn.ExecContext
	// if strings.Contains(query, "SELECT") {
	// 	d.RowReturn = true
	// } else {
	// 	d.RowReturn = false
	// }
	d.RowReturn = false
	d.PrintSanitationExecutionPlan(query)

	//Start the transaction
	//fmt.Printf("DB object: %+v", d)
	if err = d.RunTransaction(query, d.RowReturn); err != nil {
		//errCh <- fmt.Sprintf("%v", err)
		ErrorList = append(ErrorList, fmt.Sprintf("Template name %v: \n Error: %v", d.TemplateName, err))
		Result.ErrorList = append(Result.ErrorList, fmt.Sprintf("Template name %v: \n Error: %v", d.TemplateName, err))
		Result.FailedTableNames = append(Result.FailedTableNames, d.TemplateName)
		Result.FailedTablesCount++
		return fmt.Errorf("failed to Run the transaction: %w", err)
	}

	Result.SuccessTableNames = append(Result.SuccessTableNames, d.TemplateName)
	Result.SuccessTablesCount++

	return nil
}

func (d Database) PrintSanitationExecutionPlan(query string) {
	cc := text.Colors{text.FgGreen}
	fmt.Println(cc.Sprintf("\nRunning Query against DB: %v@%v\n", Cfg.Username, Cfg.Host))
	fmt.Println(cc.Sprintf("\n==== Query Execution Plan for table: %v ==== %v\n", d.TemplateName, query))

}
