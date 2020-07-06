# apacheLogToSQLite

**apacheLogToSQLite** read an apache access log file and fill a SQLite database, to allow you to ease stat building

# Installation 

```bash
git clone https://github.com/everycheck/apacheLogToSQLite.git
go build
```

# usage 

Cli option : 
 - `-sqlite` :  to change path of sqlite database
 - `-log` :  to change path of apache access log file
 - `-clearDb` :  if db already exist we clear it

```bash
./app -log /var/log/apache2/access.log 
```

Have fun