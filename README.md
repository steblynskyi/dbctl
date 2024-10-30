# README #

dbctl is a CLI tool for DB permission management. dbctl can perform following tasks in MSSQL server.

- create DB user
- get info of DB Roles (Server Role and DB Role)
- get list of DBs in specific host
- get info of DB users
- Update DB Role to add/remove member


# Table of Contents
[1. Pre-requisites](#markdown-header-pre-requisites)

[2. Usage](#markdown-header-usage)

[3. Examples](#markdown-header-examples)

- [3.1 Get](#markdown-header-1-get)

- [3.2 Update Role](#markdown-header-2-update-role)

- [3.3 Create](#markdown-header-2-update-role)


[4. Run Sanitation Scripts](#markdown-header-4-run-sanitation-scripts)

- [MSSQL DBs](#markdown-header-mssql-dbs)

- [Postgresql DBs](#markdown-header-postgresql-dbs)

[5. Usage and naming convention of DB Roles, Server Roles and DB users
](#markdown-header-usage-and-naming-convention-of-db-roles-server-roles-and-db-users)


## Pre-requisites:

For running dbctl, set DB credentials either by env variables or through flags

- Through env variables

`export DB_HOST="sql-inncenter.devsteblynskyi.com" DB_USER="steblynskyi_admin" DB_PASSWORD="EKVgKT8SDy" DB_PORT=1433 DB_NAME=inncenter01p`

OR

- Pass this flags: --username --password --db-hostname --database-name --port

## Usage:

```
Usage:
  dbctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create DB User
  get         Get details of user, role or DB
  help        Help about any command
  run         Run pre-defined SQL scripts
  update      Update Role or User

Flags:
  -D, --database-name string   DB name
  -H, --db-hostname string     DB Hostname
  -h, --help                   help for dbctl
  -P, --password string        DB Password
  -p, --port string            Port
  -t, --toggle                 Help message for toggle
  -U, --username string        DB Usenname

Use "dbctl [command] --help" for more information about a command.
```

### Examples: ###


### 1. Get

```
Usage:
  dbctl get roles [flags]

Flags:
  -d, --db-name string     DB name
  -h, --help               help for roles
  -l, --list-permissions   List DB/Server Role permissions
  -n, --role-name string   Role name(Server Role/DB Role)
```

Example:

1. Get Server Roles 

`dbctl get roles`

2. Get DB roles in specific DB

`dbctl get roles -d inncenter01p`

3. Fetch member of Server Role

`dbctl get roles -n ReadOnly`

4. Fetching members of DB Role 
Fetch member of DB Role 'ReadOnly' defined in inncenter01p database

`dbctl get roles -n ReadOnly -d inncenter01p`

5. Fetch permission of DB Role
Fetching permissions of DB Role 'ReadOnly' defined in 'inncenter01p' DB

`dbctl get roles -n ReadOnly -d inncenter01p -l`

6. Fetch permission of Server Role

`dbctl get roles -n public -l`

7. Get List of Dbs in specific host

`dbctl get dbs`

8. Get SQL Logins (users) list

`dbctl get users`

### 2. Update Role

```
Usage:
  dbctl update role [flags]

Flags:
  -a, --add-member string    Add member to Role
  -d, --db-name string       DB name
  -r, --drop-member string   Drop member from Role
  -h, --help                 help for role
  -n, --role-name string     Role name(DB Role/Server Role)
```

`Example`:

1. Add User to DB role: Add `RatesAPI` user of `inncenter01p` DB to DB Role `ReadOnly`

`dbctl update role -n ReadOnly -a RatesAPI -d inncenter01p`

2. Remove User from DB Role: Remove `RatesAPI` user of `inncenter01p` DB from DB Role `ReadOnly`

`dbctl update role -n ReadOnly -r RatesAPI -d inncenter01p`

#### 3. Create

```
dbctl create user -h

This command creates DB user.
Usage: dbctl create user --db-username=<username> --db-password=<password>

Usage:
  dbctl create user [flags]

Flags:
      --db-password string   Password for DB user
      --db-username string   DB Username
  -h, --help                 help for user

```

`Example`:

Create DB user with name demo

`dbctl create user --db-username demo --db-password test1`


### 4. Run Sanitation Scripts

sanitation command allows to run sanitation scripts on specific DB or table.
Scripts are stored here: https://bitbucket.org/steblynskyi/dbctl/src/develop/SanitationScripts/

```
dbctl run sanitation --help

Usage:
  dbctl run sanitation [flags]

Flags:
      --all-db               Run Sanitation for All DBs
      --db-name string       Run sanitation for mentioned DB name
  -h, --help                 help for sanitation
      --schema-name string   Run sanitation for specific Schema name (use along with table-name)
      --table-name string    Run sanitation for specific Table name (use along with schema-name)

Global Flags:
  -D, --database-name string   DB name
  -H, --db-hostname string     DB Hostname
  -T, --db-type string         DB Type
  -P, --password string        DB Password
  -p, --port string            Port
  -U, --username string        DB Usenname
```

`Example`:

#### _MSSQL DBs_:

0. Build dbctl binary:
Build dbctl binary, which will generate two binaries:
dbctl-darwin  dbctl-linux


```
make build
```
based on OS, respective binaries can be used in next step.

1. To run sanitation scripts on inncenter01p DB:

- first set details through env var or pass global flags:

To set through env var:

```
export DB_HOST="dbctl-test-prod-backup-22-12-03-cc-fix.cz1vjw8cas2j.us-east-1.rds.amazonaws.com" DB_USER="steblynskyi_admin" DB_PASSWORD="test123" DB_PORT=1433 DB_NAME="inncenter01p"
```


- Then run command like:
```
./dbctl run sanitation --db-type=mssql --db-name=inncenter01p

```

2. To run sanitation for specific table `client_merchantaccount_xref` of inncenter01p DB:

```
./dbctl run sanitation --db-type=mssql --db-name=inncenter --table-name="client_merchantaccount_xref" --schema-name="dbo"
```


#### _Postgresql DBs_:


To run sanitation on `postgres` DB. Make sure DB_NAME flag/env var is set to exact DB name(Case-sensitive) on which we want to run sanitation script. For example:

1. To run sanitation on `reports` DB, set `DB_NAME="reports"` like below:

```
export DB_HOST="dbctl-test-postgres-innsights-prod-22-08-08.cz1vjw8cas2j.us-east-1.rds.amazonaws.com" DB_USER="devops" DB_PASSWORD="test11" DB_PORT=5432 DB_NAME="reports"
```

- Then run sanitation with below command:

```
./dbctl run sanitation --db-type=postgres --db-name=reports
```

2. To run sanitation on `innSights` DB, set `DB_NAME="innSights"` like below:

```
export DB_HOST="dbctl-test-postgres-innsights-prod-22-08-08.cz1vjw8cas2j.us-east-1.rds.amazonaws.com" DB_USER="devops" DB_PASSWORD="test11" DB_PORT=5432 DB_NAME="innSights"
```

- Then run sanitation with below command:

```
./dbctl run sanitation --db-type=postgres --db-name=innsights
```


## Usage and naming convention of DB Roles, Server Roles and DB users

### 1. List of the Role per DB (DB Role):
  - `Usage`: appRole -> In each DB, Separate DB role for each app 
    - Each app that wants to access DB => need to create separate Role for that app
    - Basically, Each Role will need CRUD on all schemas
  - `Naming`:  _app_APPNAME_Role_ e.g. app_RatesApi_Role, app_ReservApi_Role

### 2. Server level Roles:
   - `Usage`: For All server level access/management and broader access to multiple DBs
   - `Naming`: 
      - dbaRole => For DBA team (Permission: `setupadmin`)
      - readOnlyRole => Only for read only access on server and multiple DBs
      - adminRole => Kind of Admin Role, can perform any activity on the server => With All privileges on DB (Similar to RDS Superuser) - (Permission: `sysadmin`)
      - devRole => Read only role with access of multiple DBs, Can be used by DevTeam (for running SELECT SQL queries)
      - userAgentRole => Server role for SQL Server agent Jobs

### 3. Global Users (i.e. SQL Logins)
   - Individual DBA team members => member of dbaRole
   - Individual DevOps team members => member of adminRole
   - readOnlyUser =>  member of readOnlyRole
   - Specific dev team member => member of devRole
   - appUser -> Naming: appAPPNAME e.d. appRatesApi => member of specific appRole
   - Note: Each DB user will have least privilege. So by default `public` role is assigned on user creation.
