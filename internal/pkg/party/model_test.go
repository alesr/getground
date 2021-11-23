package party

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name     string
		given    AddGuestToGuestListInput
		expected error
	}{
		{
			name: "valid input",
			given: AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 789,
			},
			expected: nil,
		},
		{
			name: "required name input",
			given: AddGuestToGuestListInput{
				Name:               "",
				Table:              456,
				AccompanyingGuests: 789,
			},
			expected: ErrGuestNameRequired,
		},
		{
			name: "required table input",
			given: AddGuestToGuestListInput{
				Name:               "123",
				Table:              0,
				AccompanyingGuests: 789,
			},
			expected: ErrTableNumberRequired,
		},
		{
			name: "invalid table input",
			given: AddGuestToGuestListInput{
				Name:               "123",
				Table:              -1,
				AccompanyingGuests: 789,
			},
			expected: ErrTableNumberInvalid,
		},
		{
			name: "invalid accompanying guests input",
			given: AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: -1,
			},
			expected: ErrAccompanyingGuestsNumberInvalid,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			observed := c.given.validate()
			assert.Equal(t, c.expected, observed)
		})
	}
}
