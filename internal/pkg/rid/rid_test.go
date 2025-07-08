package rid_test

import (
	"strings"
	"testing"

	"github.com/ArthurWang23/miniblog/internal/pkg/rid"
	"github.com/stretchr/testify/assert"
)

func Salt() string {
	return "staticSalt"
}

func TestResourceID_String(t *testing.T) {
	// 测试UserID转换为字符串
	userID := rid.UserID
	assert.Equal(t, "user", userID.String(), "UserID.String() should return 'user'")

	postID := rid.PostID
	assert.Equal(t, "post", postID.String(), "PostID.String() should return 'post'")
}

func TestResourceID_New(t *testing.T) {
	// 测试生成的ID是否带有正确前缀
	userID := rid.UserID
	uniqueID := userID.New(1)

	assert.True(t, len(uniqueID) > 0, "Generate ID should not be empty")
	assert.Contains(t, uniqueID, "user-", "Generate ID should start with 'user-' prefix")

	anotherID := userID.New(2)
	assert.NotEqual(t, uniqueID, anotherID, "Generate IDs should be unique")
}

func BenchmarkResourceID_New(b *testing.B) {
	// 性能测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := rid.UserID
		_ = userID.New(uint64(i))
	}
}

func FuzzResourceID_New(f *testing.F) {
	f.Add(uint64(1))      // 添加一个种子值counter为1
	f.Add(uint64(123456)) // 添加一个较大的种子值

	f.Fuzz(func(t *testing.T, counter uint64) {
		result := rid.UserID.New(counter)

		assert.NotEmpty(t, result, "The generated unique identifier should not be empty")

		assert.Contains(t, result, rid.UserID.String()+"-", "The generated unique identifier should contain the correct prefix")

		// 断言前缀不会与uniqueStr部分重叠
		splitParts := strings.SplitN(result, "-", 2)
		assert.Equal(t, rid.UserID.String(), splitParts[0], "The prefix part of the result should correctly match the UserID")

		if len(splitParts) == 2 {
			assert.Equal(t, 6, len(splitParts[1]), "The unique identifier part should have a length of 6")
		} else {
			t.Errorf("The format of the generated unique identifier dose not meet expectation")
		}
	})
}
