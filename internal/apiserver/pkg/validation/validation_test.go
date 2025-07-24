package validation

import "testing"

func BenchmarkIsValidUsername(b *testing.B) {
	testUsernames := []string{
		"valid_user123",         // 合法，常规输入
		"user_too_long_example", // 长度超过 20
		"sh",                    // 长度不足 3
		"in*valid",              // 包含非法字符
		"12345678901234567890",  // 合法，刚好 20 个字符
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, username := range testUsernames {
			isValidUsername(username)
		}
	}
}
