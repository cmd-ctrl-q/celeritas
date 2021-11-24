# Celeritas

### Download Celeritas:
```
go get github.com/cmd-ctrl-q/celeritas
```

### Build Celeritas CLI:
```
make build
```

### Create project:
```
./celeritas new <myapp>
```

### Need help? 
```
./celeritas help
```

### Available Commands
```
help                    - show the help commands
version                 - print application version
migrate                 - run all up migrations that have have yet to run
migrate down            - reverses the most receive migration
migrate reset           - runs all down migrations in reverse order, and then all up migrations
make migration <name>   - creates two new up and down migrations in the migrations folder
make auth               - creates and runs migrations for authentication tables and creates models and middleware
make handler <name>     - creates a stub handler in the handlers directory
make model <name>       - creates a new model in the data directory
make session            - creates a table in the database as a session store
make mail <name>        - creates two starter mail templates in the mail directory
```