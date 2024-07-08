


## To get started
To run the API and MySQL DB.
```sh
make start
```

To stop the API and the MySQL DB.
```sh
make stop
```

To exec into the MySQL DB container.
```sh
make exec
```

To restart the MySQL DB container and exec into it.
```sh
make db-restart
```

To run the tests using a Docker container.
```sh
make test
```


## Setup
The MySQL DB container uses the sql file `sql/init.sql` as the entrypoint to populate the db with test data.

If the connection to the mysql db fails, an empty inmemory noop DB will be used.

The DB contains three tables:
- users
- gender
- decisions