# 🎬 asoul-video ![Go](https://github.com/asoul-video/asoul-video/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/asoul-video/asoul-video)](https://goreportcard.com/report/github.com/asoul-video/asoul-video) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/asoul-video/asoul-video)

The backend server of https://asoul.video/

## Set up development environment

The ASOUL-Video backend server binary is meant to be run on Linux system, but you can also develop it on macOS.

### Step 1: Install dependencies

ASOUL-Video backend has the following dependencies:

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) (v1.8.3 or higher)
- [Go](https://golang.org/doc/install) (v1.16 or higher)
- [PostgreSQL](https://wiki.postgresql.org/wiki/Detailed_installation_guides) (v12 or higher)
- [Golang Migrate](https://github.com/golang-migrate/migrate/) (v4.7.0 or higher)

#### macOS

1. Install [Homebrew](https://brew.sh/).
1. Install dependencies:

    ```bash
    brew install go postgresql git
    ```

1. Configure PostgreSQL to start automatically:

    ```bash
    brew services start postgresql
    ```

1. Ensure `psql`, the PostgreSQL command line client, is on your `$PATH`. Homebrew does not put it there by default.
   Homebrew gives you the command to run to insert `psql` in your path in the "Caveats" section
   of `brew info postgresql`. Alternatively, you can use the command below. It might need to be adjusted depending on
   your Homebrew prefix (`/usr/local` below) and shell (bash below).

   ```bash
   hash psql || { echo 'export PATH="/usr/local/opt/postgresql/bin:$PATH"' >> ~/.bash_profile }
   source ~/.bash_profile
   ```

### Step 2: Initialize your database

You need a fresh Postgres database and a database user that has full ownership of that database.

1. Create a database for the current Unix user:

    ```bash
    # For Linux users, first access the postgres user shell
    sudo su - postgres
    ```

    ```bash
    createdb
    ```

2. Create the ASOUL-Video user and password:

    ```bash
    createuser --superuser asoulvideo
    psql -c "ALTER USER asoulvideo WITH PASSWORD 'asoulvideo';"
    ```

3. Create the ASOUL-Video database:

    ```bash
    createdb --owner=asoulvideo --encoding=UTF8 --template=template0 asoulvideo
    ```

### Step 3: Get the code

Generally, you don't need a full clone, so set `--depth` to `1`:

```bash
git clone --depth 1 https://github.com/asoul-video/asoul-video
```

**NOTE** The repository has Go Modules enabled, please clone to somewhere outside your `$GOPATH`.

### Step 4: Configure database settings

The Fork AI backend reads PostgreSQL connection configuration from
the [`PG*` environment variables](http://www.postgresql.org/docs/current/static/libpq-envars.html).

Add these, for example, in your `~/.bashrc`:

```
export PGPORT=5432
export PGHOST=localhost
export PGUSER=asoulvideo
export PGPASSWORD=asoulvideo
export PGDATABASE=asoulvideo
export PGSSLMODE=disable
```

You can also use a tool like [`direnv`](https://direnv.net/) to source these env vars on demand when you start the
backend.

### Step 5: Start the web server

The web server requires few environment variables to make it fully working, add them to your `~/.bashrc`:

```
export SOURCE_REPORT_KEY=<REDACTED>
```

```bash
go build . && ./asoul-video
```

## License

MIT
