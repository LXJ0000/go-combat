-- 刷新过期时间
if redis.call('get', KEYS[1]) == ARGV[1] then
    -- 锁属于自己
    return redis.call('expire', KEYS[1], ARGV[2])
else
    -- 锁不属于自己
    return 0
end