local function mexists(key, values)
    return redis.call('CF.MEXISTS', key, unpack(values))
end

assert(#KEYS == 1, 'Expecting key')

return mexists(KEYS[1], ARGV)
