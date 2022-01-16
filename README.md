# Accessing a relational database with GO
## Prerequisites
- [GO Tutorial: Accessing a relational database](https://go.dev/doc/tutorial/database-access)
- [Cómo instalar MySQL en Ubuntu 20.04](https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-20-04-es)

## Database
### Create a database
```sql
mysql> create database recordings;
```
### Change to the database
```sql
mysql> use recordings;
```
### Create an album table with values
```sql
mysql> source /path/to/create-tables.sql
```

### User and Password
From the command prompt, set the DBUSER and DBPASS environment variables for use by the Go program.

#### On Linux or Mac
```bash
$ export DBUSER=username
```
```bash
$ export DBPASS=password
```
