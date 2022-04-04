# 🎬 asoul-video ![Go](https://github.com/asoul-video/asoul-video/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/asoul-video/asoul-video)](https://goreportcard.com/report/github.com/asoul-video/asoul-video) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/asoul-video/asoul-video)

![](https://screenshotapi-dot-net.storage.googleapis.com/asoul_video__077bf1d6aeee.png)

## 配置开发环境

### 前端

TBD

### 后端

A-SOUL Video 后端二进制文件需要在 Linux 系统上运行，但你可以在 macOS、Windows 等系统上进行开发。

#### 步骤 1: 安装依赖

A-SOUL Video 后端需要安装以下依赖：

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) (v1.8.3 或更高版本)
- [Go](https://golang.org/doc/install) (v1.16 或更高版本)
- [PostgreSQL](https://wiki.postgresql.org/wiki/Detailed_installation_guides) (v12 或更高版本)
- [Golang Migrate](https://github.com/golang-migrate/migrate/) (v4.7.0 或更高版本)

##### macOS

1. 安装 [Homebrew](https://brew.sh/).
1. 安装依赖:

    ```bash
    brew install go postgresql git
    ```

1. 配置 PostgreSQL 数据库自启动:

    ```bash
    brew services start postgresql
    ```

1. 确保在 `$PATH` 环境变量中设置了 PostgreSQL 客户端命令 `psql` 的路径。通过 Homebrew 的安装默认不会设置该环境变量。执行 `brew info postgresql`
   命令，可以在 `Caveats` 段落中看到 Homebrew 提供的安装 `psql` 的方法。除此之外，也可以使用下方的命令进行安装，注意你可能需要根据本机 Homebrew 以及终端环境修改命令中的参数。

   ```bash
   hash psql || { echo 'export PATH="/usr/local/opt/postgresql/bin:$PATH"' >> ~/.bash_profile }
   source ~/.bash_profile
   ```

#### 步骤 2: 数据库初始化

你需要创建一个全新的 Postgres 数据库和一个对该数据库拥有完全操作权限的数据库用户。

1. 为当前 UNIX 用户创建数据库。

    ```bash
    # 对于 Linux 用户，首先需要切换到 postgres 用户终端下
    sudo su - postgres
    ```

    ```bash
    createdb
    ```

2. 创建 A-SOUL Video 用户并设置密码。

    ```bash
    createuser --superuser asoulvideo
    psql -c "ALTER USER asoulvideo WITH PASSWORD 'asoulvideo';"
    ```

3. 创建 A-SOUL Video 数据库。

    ```bash
    createdb --owner=asoulvideo --encoding=UTF8 --template=template0 asoulvideo
    ```

#### 步骤 3: 拉取代码

通常来说，你并不需要拉取所有的历史代码，因此可以设置 `--depth 1` 参数。

```bash
git clone --depth 1 https://github.com/asoul-video/asoul-video
```

**注意** 本仓库开启 Go Module，请将仓库目录防止在你的 `$GOPATH` 目录之外。

#### 步骤 4: 配置数据库连接

A-SOUL Video 后端从以下环境变量中读取数据库连接参数。 [`PG*` 环境变量](http://www.postgresql.org/docs/current/static/libpq-envars.html).

将以下这些环境变量参数添加至 `~/.bashrc` 文件中:

```
export PGHOST=localhost
export PGUSER=asoulvideo
export PGPASSWORD=asoulvideo
export PGDATABASE=asoulvideo
export PGSSLMODE=disable
```

你也可以使用类似 [`direnv`](https://direnv.net/) 这样的工具来管理你的环境变量。

#### 步骤 5: 启动 Web 服务器

启动 Web 服务器前，你需要确保以下环境变量中的参数数正确可用的，并将他们添加至你的 `~/.bashrc` 文件中:

```
# 爬虫上报数据时的鉴权参数
export SOURCE_REPORT_KEY=<REDACTED>
```

```bash
go build . && ./asoul-video
```

## License

MIT
