local function reserve(key, capacity)
    redis.call('DEL', key)
    return redis.call('PF.RESERVE', key, capacity)
end

assert(#KEYS == 1, 'Expecting key')
assert(#ARGV >= 1, 'Expecting value')

return reserve(KEYS[1], ARGV[1])
