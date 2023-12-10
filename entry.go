package whale

type Entry struct {
	Key   []byte
	Value []byte
	// 过期时间 0 表示永不过期
	ExpiresAt uint64
}
