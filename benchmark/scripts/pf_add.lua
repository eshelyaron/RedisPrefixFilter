local function add(key, value)
    return redis.call('PF.ADD', key, value)
end

assert(#KEYS == 1, 'Expecting key')
assert(#ARGV == 1, 'Expecting value')

return add(KEYS[1], ARGV[1])
