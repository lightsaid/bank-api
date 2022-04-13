-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id=$1 LIMIT 1;   -- 普通 查询 语句，当并发操作同一个表时，并不会阻止另一个事务读取 (Read Commited) 下一个 SELECT 语句解决这个问题

-- name: GetAccountForUpdate :one
SELECT * FROM accounts -- for update 查询语句，当两个事务同时操作一张表时，另一个事务会等待前一个事务commited了，才能读取
WHERE id=$1 LIMIT 1
FOR NO KEY UPDATE; -- 告诉查询，udpate accounts 表不会更新主键ID，就不会产生死锁了
-- FOR UPDATE;  -- 此时，因为外键约束，当多事务并发时，会发生死锁
--（场景：因为accounts中的id 在entries、transfer表中是外建，当一个事务在插入/更新entries｜transfer时，另一个事务在更新accounts 就会发生死锁）
-- 因为插入/更新entries｜transfer的事务会担心，accounts 修改了id，想等另一个事务提交了才操作， 而更新accounts的事务也担心会更新了id，影响了第一个是事务
-- 因此双方在等待对方commit，就产生了死锁。

-- name: ListAccounts :many
SELECT * FROM accounts 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAcount :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *; 

-- name: AddUpdateAcountBlance :one
UPDATE accounts 
SET balance = balance + sqlc.arg(amount) -- 把参数名字设置为 amount
WHERE id = sqlc.arg(id)
RETURNING *; 


-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id=$1;