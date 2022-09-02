local function madd(key, values)
    return redis.call('PF.MADD', key, unpack(values))
end

assert(#KEYS == 1, 'Expecting key')

return madd(KEYS[1], ARGV)
