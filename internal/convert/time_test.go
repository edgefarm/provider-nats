package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeToRFC3339(t *testing.T) {
	assert := assert.New(t)
	// Test data
	// UNIX Timestamp    ISO 8601                     RFC 3339
	// 1673275712     == 2023-01-09T14:48:32+00:00 == 2023-01-09T14:48:32Z

	// positive tests
	// RFC3339 format
	dRfc, err := RFC3339ToTime("2023-01-09T14:48:32Z")
	assert.Equal(dRfc.Unix(), int64(1673275712))
	assert.Nil(err)

	// ISO 8601 format
	dRfc, err = RFC3339ToTime("2023-01-09T14:48:32+00:00")
	assert.Equal(dRfc.Unix(), int64(1673275712))
	assert.Nil(err)

	dStr, err := TimeToRFC3339(dRfc)
	assert.Equal(dStr, "2023-01-09T14:48:32Z")
	assert.Nil(err)
	// negative tests
	// wrong input format
	// UTC format
	_, err = RFC3339ToTime("01/09/2023 @ 2:48pm")
	assert.NotNil(err)

	// RFC 822, 1036, 1123, 2822 formats
	_, err = RFC3339ToTime("Mon, 09 Jan 2023 14:48:32 +0000")
	assert.NotNil(err)

	// RFC 2822 format
	_, err = RFC3339ToTime("Monday, 09-Jan-23 14:48:32 UTC")
	assert.NotNil(err)

	// wrong user input
	// wrong input format
	_, err = RFC3339ToTime("2023-01-09T14:48:32+00:00")
	assert.Nil(err)

	// timezone missing
	_, err = RFC3339ToTime("2023-01-09T14:48:32")
	assert.NotNil(err)
}
