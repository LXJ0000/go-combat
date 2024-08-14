value = redis.call('GET', KEYS[1])
if value == false then -- key doesn't exist
    return redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2])
elseif value == ARGV[1] then -- key exists and lock by this process
    redis.call('EXPIRE', KEYS[1], ARGV[2])
    return "OK"
else -- key exists and lock by another process
    return ""
end