package party

import "time"

type (

	// AddGuestToGuestListOutput defines the output struct for adding guests to the guestlist.
	AddGuestToGuestListOutput struct {
		Name string
	}

	// Guest defines the guest struct.
	Guest struct {
		Name               string     `json:"name"`
		Table              int        `json:"table"`
		AccompanyingGuests int        `json:"accompanying_guests"`
		TimeArrival        *time.Time `json:"time_arrived,omitempty"`
	}

	// GuestList defines the guestlist output struct.
	GetGuestListOutput struct {
		Guests []Guest `json:"guests"`
	}

	// WelcomeGuestInput defines the input struct for welcoming guests.
	WelcomeGuestInput struct {
		Name               string `json:"name"`
		AccompanyingGuests int    `json:"accompanying_guests"`
	}

	// WelcomeGuestOutput defines the output struct for welcoming guests.
	WelcomeGuestOutput struct {
		Name string `json:"name"`
	}

	// GoodbyeGuestInput defines the input struct for leaving guests.
	GoodbyeGuestInput struct {
		Name string
	}

	GuestArrived struct {
		Name               string    `json:"name"`
		AccompanyingGuests int       `json:"accompanying_guests"`
		TimeArrival        time.Time `json:"time_arrived"`
	}

	ListArrivedGuestsOutput struct {
		Guests []GuestArrived `json:"guests"`
	}

	GetEmptySeatsOutput struct {
		EmptySeats int `json:"empty_seats"`
	}
)

// AddGuestToGuestListInput defines the input struct for adding guests to the guestlist.
type AddGuestToGuestListInput struct {
	Name               string `json:"-"`
	Table              int    `json:"table"` // Table Number
	AccompanyingGuests int    `json:"accompanying_guests"`
}

func (r *AddGuestToGuestListInput) validate() error {
	if r.Name == "" {
		return ErrGuestNameRequired
	}

	if r.Table == 0 {
		return ErrTableNumberRequired
	}

	if r.Table < 0 {
		return ErrTableNumberInvalid
	}

	if r.AccompanyingGuests < 0 {
		return ErrAccompanyingGuestsNumberInvalid
	}

	// TODO(:alesr): Validate table size
	return nil
}
