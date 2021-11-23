package party

import "context"

type Mock struct {
	AddGuestToGuestListFunc func(ctx context.Context, in *AddGuestToGuestListInput) (*AddGuestToGuestListOutput, error)
	GetGuestListFunc        func(ctx context.Context) (GetGuestListOutput, error)
	WelcomeGuestFunc        func(ctx context.Context, in *WelcomeGuestInput) (*WelcomeGuestOutput, error)
	GoodbyeGuestFunc        func(ctx context.Context, in *GoodbyeGuestInput) error
	ListArrivedGuestsFunc   func(ctx context.Context) (ListArrivedGuestsOutput, error)
	GetEmptySeatsFunc       func(ctx context.Context) (GetEmptySeatsOutput, error)
}

func (m *Mock) AddGuestToGuestList(ctx context.Context, in *AddGuestToGuestListInput) (*AddGuestToGuestListOutput, error) {
	return m.AddGuestToGuestListFunc(ctx, in)
}

func (m *Mock) GetGuestList(ctx context.Context) (GetGuestListOutput, error) {
	return m.GetGuestListFunc(ctx)
}

func (m *Mock) WelcomeGuest(ctx context.Context, in *WelcomeGuestInput) (*WelcomeGuestOutput, error) {
	return m.WelcomeGuestFunc(ctx, in)
}

func (m *Mock) GoodbyeGuest(ctx context.Context, in *GoodbyeGuestInput) error {
	return m.GoodbyeGuestFunc(ctx, in)
}

func (m *Mock) ListArrivedGuests(ctx context.Context) (ListArrivedGuestsOutput, error) {
	return m.ListArrivedGuestsFunc(ctx)
}

func (m *Mock) GetEmptySeats(ctx context.Context) (GetEmptySeatsOutput, error) {
	return m.GetEmptySeatsFunc(ctx)
}
