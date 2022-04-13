# 学习 Backend master class [Golang, Postgres, Docker] 课程笔记
课程在 Udemy 获得 4.8 一致好评，（满分5.0）参与学习人数已达到 8385+

因此课程质量不言而喻

**课程地址 [视频](https://www.youtube.com/c/TECHSCHOOLGURU)**

**课程代码 [代码](https://github.com/techschool/simplebank)**

### 数据库设计
设计工具使用在线工具 [dbdiagram.io/](https://dbdiagram.io/)；
记一个坑：如果 data type 写错，dbdiagram工具不会报错

![效果图](./imgs/Bank-design.png)

### Docker 
- 安装 postgresql
    * docker pull postgres:12-alpine
    * [相关文档](https://hub.docker.com/_/postgres)
- 运行 postgresql 12
    * docker run --name postgresql12 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=rootcc -d postgres:12-alpine
    
    ![runing](./imgs/run-1.jpg)

    * docker stop postgresql12 -> 停止容器

    * docker start postgresql12 -> 启动一个容器

    * docker ps -a 查看所有容器

- 进入 postgresql 容器终端
    * docker exec -it postgresql12 psql -U root    # 因为在本地 psql 不需要密码

    ![runing](./imgs/run-2.jpg)

- 课程数据库管理工具使用的是 TablePlus， 而我个人使用的是`DBeaver`,仅仅是一个数据库管理工具，用那一个都差不多。
    
- 拷贝设计好的sql文件在数据库管理工具执行，生成对应的表。
    ![runing](./imgs/conn.jpg)

### 数据库迁移 （migrate）
[文档/工具库](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

golang-migrate 和 gobuffalo 库的 soda 差不多

- 使用 migrate， Mac 安装
    * brew install golang-migrate
    **直接安装很有可能不成功，报错：Failed to connect to go.googlesource.com port 443: Operation timedout**

    此时需要设置代理：`git config --global http.proxy "localhost:63913"` 具体端口，看个人vpn代码端口；
    安装完还是给禁用代理：git config --global --unset http.proxy

- 生成迁移文件
    - 创建目录 `db/migration`
    * migrate create -ext sql -dir db/migration -seq init_schema
        执行完命令后，会生成一个向上迁移和向下迁移sql文件
    - 升级迁移 migrate up 
    - 版本回退迁移 migrate down
    - 编写好两个迁移 sql

- docker exec -it postgresql12 bash  -> 进入容器交互式终端

- 创建一个新的数据库
    * createdb --username=root --owner=root simple_bank -> 创建数据库
    * psql simple_bank -> use simple_bank
    * \q ->  退出
    * dropdb simple_bank -> 删除
    * exit -> 退出交互式终端
    - 再次创建数据库和进入交互式终端
    * docker exec -it postgresql12 createdb --username=root --owner=root simple_bank
    * docker exec -it postgresql12 psql -U root simple_bank


- 编写 Makefile ，使用 make 构建数据库（工程）
``` makefile
postgres:
    # docker run --name postgres12 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=rootcc -d postgres:12-alpine
	docker run --name postgres12 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=rootcc -d postgres:12-alpine


createdb:
    docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
    docker exec -it postgres12 dropdb simple_bank

.PHONY:postgres createdb dropdb # 类似将命令到处去

## 上面就是定义了三个命令，使用 make 执行；比如： make postgres
```
- docker stop postgresql12 
- docker rm postgresql12
- make postgres
    * 会报错 `Makefile:3: *** missing separator. Stop.`
    * 错误原因是：此条命令从终端拷贝过的，有额外的字符了，空格或者是table不对，颜色不一样，重写了之后颜色变了，就正常了
- 执行 make 命令：
    ![makefile](./imgs/makefile.jpg)
    * 修正数据名字在次执行,停止容器，删除容器再次执行make命令
    * 打开 DBeaver 查看数据库
- 执行迁移数sql脚本命令
    * migrate -path db/migration -database "postgresql://root:rootcc@localhost:5433/simple_bank?sslmode=disable" -verbose up
    * migrate -path db/migration -database "postgresql://root:rootcc@localhost:5433/simple_bank?sslmode=disable" -verbose down
    * 最后将命令添加到 Makefile 里
    - 效果：
    ![migrate](./imgs/migrate.jpg)

**以上就是前三节课程学习内容，做一个总结：**
1. 采用 `dbdiagram.io` 设计数据库
1. 采用 Docker 并创建 postgresql 容器
1. 采用 golang-migrate 迁移 sql 脚本
1. 编写 Makefile，用 make 编译命令
----- 

### 4.Generate CRUD Golang code from SQL | Compare db/sql, gorm, sqlx & sqlc
这一节主要对 db/sql, gorm, sqlx & sqlc 这个几种 CRUD 方式进行对比；
性能(查询速度)最好的无疑是 ` db/sql`,但是缺点最大的也是它，应为需要编写大量的 sql 语句，手动映射struct，容易出错并且枯燥；

`gorm` 不需要编写sql语句，不需要映射，框架已经处理好了，但是在大型项目里，性能比 `db/sql` 慢 3-5 倍；

因此就有了择中选择 `sqlx`, sqlx 实现了映射，但是也要编写sql语句。

再者就是 sqlc 了，自动生成 CRUD sql语句

最后再补充一种 sqluirrel， 暂时没有了解。

- 课程采用 sqlc
[sqlc](https://github.com/kyleconroy/sqlc)

- macOS 安装
    * brew install sqlc
    * sqlc help
    * sqlc version -> 当前 v1.13.0 
        
- sqlc 起步
    * sqlc init 生成sqlc.yaml 配置文件
    ``` yaml
    version: "1"
    packages:
    - name: "db"
        path: "./db/sqlc" # 生成代码路径
        queries: "./db/query/" # 编写 sql 文件语句路径
        schema: "./db/migration/" # sql 迁移脚本
        engine: "postgresql" #  那种中数据库，如：mysql
        emit_prepared_queries: false # 是否预执行sql，可提高性能
        emit_interface: false # 是否生成interface，是否使用interface定义接口
        emit_exact_table_names: false # false -> 表名加复数
        emit_empty_slices: false
        emit_exported_queries: false
        emit_json_tags: true # 是否生成 json tag
        emit_result_struct_pointers: false
        emit_params_struct_pointers: false
        emit_methods_with_db_argument: false
        json_tags_case_style: "none" # json tag 字段命名风格；none -> 跟随table列的名字
        # 生成三个文件的命名规则，分别是db.go、models.go,
        # querier.go 指定的是，sql 语句文件后加 .go 后缀，像下面 account.sql -> account.sql.go
        output_db_file_name: "db.go"
        output_models_file_name: "models.go"
        output_querier_file_name: "querier.go"
    ```
    * 具体配置说明参考文档，不一一细说
    [sqlcConfig](https://docs.sqlc.dev/en/latest/reference/config.html)
    * 将 sqlc generate 添加至 Makefile

    - sqlc 工作流程
        1. sqlc 生成代码是依赖于上面`schema: "./db/migration/"` 配置，里面有建表的sql语句
        1. 第二个很重要的配置`queries: "./db/query/"`, 在配置路径下编写 CRUD SQL 语句
        1. 经过上面两步就可以生成 crud go 代码了  
    - 在 db/query/ 下编写 sql 语句, 例如创建一个 account.sql；
    注释是非常重要的，这个注释是给sqlc识别的
    ``` sql
        -- name: CreateAccount :one
        INSERT INTO accounts (
            owner,
            balance,
            currency
        ) VALUES (
            $1, $2, $3
        ) RETURNING *;
    ```
    - make sqlc 生成代码就在 db/sqlc 目录下

    **简直完美，太方便了**

### 5. Write unit tests for database CRUD with random data in Golang
- 安装 postgres 驱动 `go get github.com/lib/pq`
- 使用 `go get github.com/stretchr/testify` 测试

- `stretchr/testify` 库使用方式
    ``` go
    package db
    import (
        "context"
        "testing"

        "github.com/stretchr/testify/require"
    )
    func TestCreateAccount(t *testing.T) {
        arg := CreateAccountParams{
            Owner:    "张三",
            Balance:  200,
            Currency: "RMB",
        }
        account, err := testQueries.CreateAccount(context.Background(), arg)

        // 检查 err 是否为nil，如果不是nil则测试不通过
        require.NoError(t, err)

        // 检查返回的账号是否为空，空则测试不通过
        require.NotEmpty(t, account)

        // 期望值和实际值对比，不一致则测试不通过
        require.Equal(t, arg.Owner, account.Owner)
        require.Equal(t, arg.Balance, account.Balance)
        require.Equal(t, arg.Currency, account.Currency)
    }

    ```

    - 执行 go test (或者点击：TestCreateAccount 方法上面灰色字体的 run test)
    - 执行整个 package 测试，查看覆盖率 （一图胜千言）如图：
        ![test](./imgs/test.jpg)
    - 执行整个package测试之后，整个包下哪些代码没被测试都会被标颜色，这就很清晰明了

    - 创建 util/random.go 生成随机变量

    - 在 Makefile 添加 test 命令
    - 继续编写测试代码
    
    **本次学习提交代码hash：50fd5b8679c04， msg：feat: #5 编写account crud test**


### 6. A clean way to implement database transaction in Golang
这一节课程讨论在golang里如何实现一个干净的事物（transaction）<br>

- 如：从账户A 向 账户B 转 10块钱，此功能有五以下步
    1. 在transfer表创建一个金额等于10的转账记录
    1. 在entries表创建一条入账记录，账户A amount=-10
    1. 在entries表创建一条入账记录，账户B amount=+10
    1. 从账户A里减去10块钱
    1. 从账户B增加10块钱

**因为涉及多表的操作，要保持可靠一致性，就要使用事物，事物必须满足 `ACID` 属性**
    * Atomicity(A)       -> 原子性，就是最小单位，不可分割，对数据库的修改要么成功执行完成，要不就全部回滚
    * Consistency(C)     -> 一致性，事物执行后，数据库状态应该保持有效
    * Isolation(I)       -> 隔离性，是指一组对数据库的并发修改互相不影响。
    * Durability(D)      -> 持久性，是指对数据的修改，一旦完成，该结果就应当永远不丢失，存盘了。

    **hash:b987b35 msg: feate: 转账事务test**

### 7. DB transaction lock & How to handle deadlock in Golang
- docker exec -it postgres12 psql -U root -d simple_bank

### 8. How to avoid deadlock in DB transaction? Queries order matters!

### 9. Deeply understand transaction isolation levels & read phenomena in MySQL & PostgreSQL


### 10. Setup Github Actions for Golang + Postgres to run automated tests
https://www.ruanyifeng.com/blog/2019/09/getting-started-with-github-actions.html

mkdir -p .github/workflows

touch .github/workflows/ci.yaml