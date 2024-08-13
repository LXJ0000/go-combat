-- // 以下步骤必须为原子操作 这里采用 lua 脚本实现
-- // 1. 检查是否为自己加的锁
-- // 2. 解锁

if redis.call('get', KEYS[1]) == ARGV[1] then
    -- 锁属于自己 or 锁不存在
    return redis.call('del', KEYS[1])
else
    -- 锁不属于自己
    return 0
end