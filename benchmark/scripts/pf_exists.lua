local function exists(key, value)
    return redis.call('PF.EXISTS', key, value)
end

assert(#KEYS == 1, 'Expecting key')
assert(#ARGV == 1, 'Expecting value')

return exists(KEYS[1], ARGV[1])
