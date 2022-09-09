local function mexists(key, values)
    return redis.call('PF.MEXISTS', key, unpack(values))
end

assert(#KEYS == 1, 'Expecting key')

return mexists(KEYS[1], ARGV)
