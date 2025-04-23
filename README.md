# dbtools

This is a small library with common Go/DB patterns I've extracted from some of my personal projects.
Note that these functions are hard-coded for Postgres, although they could potentially be made generic.


## migration.Run

`migration.Run` is a helper that will use the golang-mirgate/migrate package to run all
migrations in a given folder.

The main use I have for this is for *tests*.
For example I might want have a fresh database state for my tests, or I might want
to "namespace" different packages' tests into different databases, so that they can
run in parallel.

To make this more ergonomic for tests, there's the `dbtest` package.

## dbtest

This provides a quick helper to call on tests to:
1. Create a database
2. Migrate the database
3. Return a *sql.DB connection

The `dbtest.DB` will also truncate the database when
called multiple times, ensuring that tests have always a
clear state to run.

The intention is that the target DB is running on the
local machine, at default ports.
Customizations around username, password and suffixing can be done
by changing `dbtest.Config`.

Example:

```
func TestDB1(t *testing.T) {
    db := dbtest.DB(t, "./migrations", "iam_service")
    ...
    // My testing
}

func TestDB2(t *testing.T) {
    db := dbtest.DB(t, "./migrations", "iam_service")
    ...
    // My testing
}
```

`gotest.DB` will ensure that the test database prefixed with
"iam_service" will be truncated between tests.

Alternatively, to be able to run them in parallel:


```
func TestDB1(t *testing.T) {
    t.Parallel()
    db := dbtest.DB(t, "./migrations", "iam_service_1")
    ...
    // My testing
}

func TestDB2(t *testing.T) {
    t.Parallel()
    db := dbtest.DB(t, "./migrations", "iam_service_2")
    ...
    // My testing
}
```

This will create 2 distinct test databases for the tests.

Also to make `DB` easier to use, the migrationPath can be
the migrations path relative to the root of the project instead
of the package that it's being called. `DB` walk up the folder tree
looking for it.

### dbtest/dbgorm

`gormtest.DB()` is just a wrapper on `dbtest.DB()` that returns
a gorm.DB instead, for convenience.

### dberrors

`dberrors.ToStatusErr()` takes a "database error" plus a custom message
and tries to map that into a `status.Error` (gRPC).
