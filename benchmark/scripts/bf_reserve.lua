local function reserve(key, error_rate, capacity)
    redis.call('DEL', key)
    return redis.call('BF.RESERVE', key, error_rate, capacity, 'NONSCALING')
end

assert(#KEYS == 1, 'Expecting key')
assert(#ARGV >= 2, 'Expecting value')

return reserve(KEYS[1], ARGV[1], ARGV[2])
