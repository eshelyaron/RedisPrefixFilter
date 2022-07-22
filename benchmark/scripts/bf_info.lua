local function info(key)
    return redis.call('BF.INFO', key)
end

assert(#KEYS == 1, 'Expecting key')

return info(KEYS[1])
