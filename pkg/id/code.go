package id

func NewCode(id uint64, options ...func(*CodeOptions)) string {
	ops := getCodeOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}

	id = id*uint64(ops.n1) + ops.salt

	var code []rune
	slIdx := make([]byte, ops.l)

	charLen := len(ops.chars)
	charLenUI := uint64(charLen)

	for i := 0; i < ops.l; i++ {
		slIdx[i] = byte(id % charLenUI)
		slIdx[i] = (slIdx[i] + byte(i)*slIdx[0]) % byte(charLen)
		id /= charLenUI
	}

	for i := 0; i < ops.l; i++ {
		idx := (byte(i) * byte(ops.n2)) % byte(ops.l)
		code = append(code, ops.chars[slIdx[idx]])
	}
	return string(code)
}
