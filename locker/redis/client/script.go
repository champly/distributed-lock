package client

const (
	// keys: key, value
	// argv: expireTime
	getScript = `
		local key = tostring(KEYS[1])
		local value = tostring(KEYS[2])
		local expireTime = tonumber(ARGV[1])
		
		if redis.call("setnx", key, value) == 0 then
			return 0
		end
		
		redis.call("expire", key, expireTime)
		return 1
`

	// keys: key
	delScript = `
		local key = tostring(KEYS[1])

		if redis.call("del", key) == 0 then
			return 0
		end
		return 1
`

	// keys: key
	// argv: expireTime
	delayScript = `
		local key tostring(KEYS[1])
		local expireTime = tonumber(ARGV[1])

		if redis.call("expire", key, expireTime) == 0 then
			return 0
		end
		return 1
`
)
