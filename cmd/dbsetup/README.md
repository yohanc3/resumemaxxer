# DB setup helpers

## How to create migration up and down files

Run: 

```sh
$ cd cmd/dbsetup && sh create_migration_files.sh <file_name>
```

The output will be something like:

```
/.../cmd/dbsetup/migration_files/000001_<file_name>.up
.sql
/.../cmd/dbsetup/migration_files/000001_<file_name>.do
wn.sql
Succesfully created both files
```


Example: 
```sh
$ cd cmd/dbsetup && create_migration_files.sh create_users_table
```

Output:
```
/.../cmd/dbsetup/migration_files/000001_create_users_table.up
.sql
/.../cmd/dbsetup/migration_files/000001_create_users_table.do
wn.sql
Succesfully created both files
```

## How to run migration up and down files

Migrations are handled behind the scenes. See [cmd/dbsetup/main.go](main.go)

## How to install migrate

Run:

``` sh
$ curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz
```
