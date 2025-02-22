package extractor

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = `abc 123
def 245
qqq 123
xxx`

func TestBasicExtractor(t *testing.T) {
	input := ConvertReaderToStringChan(ioutil.NopCloser(strings.NewReader(testData)), 1)
	ex, err := New(input, &Config{
		Regex:   `(\d+)`,
		Extract: "val:{1}",
		Workers: 1,
	})
	assert.NoError(t, err)

	vals := unbatchMatches(ex.ReadChan())
	assert.Equal(t, "abc 123", vals[0].Line)
	assert.Equal(t, 2, len(vals[0].Groups))
	assert.Equal(t, 4, len(vals[0].Indices))
	assert.Equal(t, "123", vals[0].Groups[0])
	assert.Equal(t, "val:123", vals[0].Extracted)
	assert.Equal(t, uint64(1), vals[0].LineNumber)
	assert.Equal(t, uint64(1), vals[0].MatchNumber)

	assert.Equal(t, 3, len(vals))

	assert.Equal(t, uint64(0), ex.IgnoredLines())
	assert.Equal(t, uint64(3), ex.MatchedLines())
	assert.Equal(t, uint64(4), ex.ReadLines())
}
