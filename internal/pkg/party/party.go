package party

import (
	"context"
	"time"

	"github.com/alesr/getground/internal/pkg/party/repository"
	"go.uber.org/zap"
)

type (
	Service interface {
		AddGuestToGuestList(ctx context.Context, in *AddGuestToGuestListInput) (*AddGuestToGuestListOutput, error)
		GetGuestList(ctx context.Context) (GetGuestListOutput, error)
		WelcomeGuest(ctx context.Context, in *WelcomeGuestInput) (*WelcomeGuestOutput, error)
		GoodbyeGuest(ctx context.Context, in *GoodbyeGuestInput) error
		ListArrivedGuests(ctx context.Context) (ListArrivedGuestsOutput, error)
		GetEmptySeats(ctx context.Context) (GetEmptySeatsOutput, error)
	}

	Party struct {
		logger    *zap.Logger
		repo      repository.Repository
		tableSize int
	}
)

var _ Service = (*Party)(nil)

func New(logger *zap.Logger, repo repository.Repository, tableSize int) *Party {
	return &Party{
		logger:    logger.Named("party_service"),
		repo:      repo,
		tableSize: tableSize,
	}
}

func (p *Party) AddGuestToGuestList(ctx context.Context, in *AddGuestToGuestListInput) (*AddGuestToGuestListOutput, error) {
	// Validate input
	if err := in.validate(); err != nil {
		return nil, err
	}

	// Get guest list
	guestlist, err := p.repo.ListGuests(ctx)
	if err != nil {
		return nil, err
	}

	// Iterate over guest list and check if guest is already in the list
	for _, guest := range guestlist {
		if guest.Name == in.Name {
			return nil, ErrGuestAlreadyInList
		}
	}

	// Get table by number
	table, err := p.repo.GetTableByNumber(ctx, in.Table)
	if err != nil {
		return nil, err
	}

	// If table does not exist, create it and add guest to guest list
	if table == nil {
		if err := p.createTableAndAddToGuestList(ctx, in); err != nil {
			return nil, err
		}
		return &AddGuestToGuestListOutput{
			Name: in.Name,
		}, nil
	}

	// Check if table has enough available seats
	requestedSeats := in.AccompanyingGuests + 1
	if table.AvailableSeats < requestedSeats {
		return nil, ErrTableNotEnoughSeats
	}

	// Insert guest into guest list
	guestStore := repository.Guest{
		Name:               in.Name,
		Table:              in.Table,
		AccompanyingGuests: in.AccompanyingGuests,
	}

	if err := p.repo.UpsertGuest(ctx, &guestStore); err != nil {
		return nil, err
	}

	// Update table

	tableStore := repository.Table{
		Number:         table.Number,
		AvailableSeats: table.AvailableSeats - requestedSeats,
		Size:           p.tableSize,
	}

	if err := p.repo.UpsertTable(ctx, &tableStore); err != nil {
		return nil, err
	}

	// Return output
	return &AddGuestToGuestListOutput{
		Name: in.Name,
	}, nil
}

func (p *Party) GetGuestList(ctx context.Context) (GetGuestListOutput, error) {
	guests, err := p.repo.ListGuests(ctx)
	if err != nil {
		return GetGuestListOutput{}, err
	}

	var list = make([]Guest, 0, len(guests))
	for _, guest := range guests {
		list = append(list, Guest{
			Name:               guest.Name,
			Table:              guest.Table,
			AccompanyingGuests: guest.AccompanyingGuests,
		})
	}

	return GetGuestListOutput{
		Guests: list,
	}, nil
}

func (p *Party) WelcomeGuest(ctx context.Context, in *WelcomeGuestInput) (*WelcomeGuestOutput, error) {
	// Get guest in guest list
	guest, err := p.repo.GetGuestByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	// Check if guest is in guest list
	if guest == nil {
		return nil, ErrGuestNotInList
	}

	// Get table by number
	table, err := p.repo.GetTableByNumber(ctx, guest.Table)
	if err != nil {
		return nil, err
	}

	if table == nil {
		return nil, ErrTableNumberNotFound
	}

	// Check if table has enough available seats

	requestedSeats := in.AccompanyingGuests + 1

	if table.AvailableSeats < requestedSeats {
		return nil, ErrTableNotEnoughSeats
	}

	// Update table

	tableStore := repository.Table{
		AvailableSeats: table.AvailableSeats - requestedSeats,
	}

	if err := p.repo.UpsertTable(ctx, &tableStore); err != nil {
		return nil, err
	}

	// Mark guest as present
	arrivalTime := time.Now() // TODO(alesr): Use timezone
	guest.TimeArrival = &arrivalTime

	if err := p.repo.UpsertGuest(ctx, guest); err != nil {
		return nil, err
	}

	// Return output
	return &WelcomeGuestOutput{
		Name: in.Name,
	}, nil
}

func (p *Party) GoodbyeGuest(ctx context.Context, in *GoodbyeGuestInput) error {
	// Get guest in guest list
	guest, err := p.repo.GetGuestByName(ctx, in.Name)
	if err != nil {
		return err
	}

	// Check if guest is in guest list
	if guest == nil {
		return ErrGuestNotInList
	}

	// Get table by number
	table, err := p.repo.GetTableByNumber(ctx, guest.Table)
	if err != nil {
		return err
	}

	if table == nil {
		return ErrTableNumberNotFound
	}

	// Update table

	requestedSeats := guest.AccompanyingGuests + 1
	tableStore := repository.Table{
		AvailableSeats: table.AvailableSeats + requestedSeats,
	}

	if err := p.repo.UpsertTable(ctx, &tableStore); err != nil {
		return err
	}

	// Mark guest as absent

	timeDeparture := time.Now() // TODO(alesr): Use timezone
	guest.TimeDeparture = &timeDeparture

	if err := p.repo.UpsertGuest(ctx, guest); err != nil {
		return err
	}
	return nil
}

func (p *Party) ListArrivedGuests(ctx context.Context) (ListArrivedGuestsOutput, error) {
	guests, err := p.repo.GetArrivedGuests(ctx)
	if err != nil {
		return ListArrivedGuestsOutput{}, err
	}

	list := make([]GuestArrived, 0, len(guests))
	for _, guest := range guests {

		if guest.TimeArrival == nil {
			continue
		}

		list = append(list, GuestArrived{
			Name:               guest.Name,
			TimeArrival:        *guest.TimeArrival,
			AccompanyingGuests: guest.AccompanyingGuests,
		})
	}

	return ListArrivedGuestsOutput{
		Guests: list,
	}, nil
}

func (p *Party) GetEmptySeats(ctx context.Context) (GetEmptySeatsOutput, error) {
	tables, err := p.repo.GetTables(ctx)
	if err != nil {
		return GetEmptySeatsOutput{}, err
	}

	var emptySeats int
	for _, table := range tables {
		emptySeats += table.AvailableSeats
	}

	return GetEmptySeatsOutput{
		EmptySeats: emptySeats,
	}, nil
}

// TODO: Implement operation as repository transaction to prevent broken state
func (p *Party) createTableAndAddToGuestList(ctx context.Context, in *AddGuestToGuestListInput) error {
	// Create table

	requestedSeats := in.AccompanyingGuests + 1

	tableStore := repository.Table{
		Number:         in.Table,
		AvailableSeats: p.tableSize - requestedSeats,
		Size:           p.tableSize,
	}

	if err := p.repo.UpsertTable(ctx, &tableStore); err != nil {
		return err
	}

	// Add guest to guest list

	guestStore := repository.Guest{
		Name:               in.Name,
		Table:              in.Table,
		AccompanyingGuests: in.AccompanyingGuests,
	}

	if err := p.repo.UpsertGuest(ctx, &guestStore); err != nil {
		return err
	}
	return nil
}
