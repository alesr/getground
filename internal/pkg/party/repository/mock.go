package repository

import "context"

var _ Repository = (*Mock)(nil)

type Mock struct {
	GetArrivedGuestsFunc func(ctx context.Context) ([]Guest, error)
	GetGuestByNameFunc   func(ctx context.Context, name string) (*Guest, error)
	ListGuestsFunc       func(ctx context.Context) ([]Guest, error)
	UpsertGuestFunc      func(ctx context.Context, guest *Guest) error
	GetTableByNumberFunc func(ctx context.Context, number int) (*Table, error)
	GetTablesFunc        func(ctx context.Context) ([]Table, error)
	UpsertTableFunc      func(ctx context.Context, table *Table) error
}

func (m *Mock) GetArrivedGuests(ctx context.Context) ([]Guest, error) {
	return m.GetArrivedGuestsFunc(ctx)
}

func (m *Mock) GetGuestByName(ctx context.Context, name string) (*Guest, error) {
	return m.GetGuestByNameFunc(ctx, name)
}

func (m *Mock) ListGuests(ctx context.Context) ([]Guest, error) {
	return m.ListGuestsFunc(ctx)
}

func (m *Mock) UpsertGuest(ctx context.Context, guest *Guest) error {
	return m.UpsertGuestFunc(ctx, guest)
}

func (m *Mock) GetTableByNumber(ctx context.Context, number int) (*Table, error) {
	return m.GetTableByNumberFunc(ctx, number)
}

func (m *Mock) GetTables(ctx context.Context) ([]Table, error) {
	return m.GetTablesFunc(ctx)
}

func (m *Mock) UpsertTable(ctx context.Context, table *Table) error {
	return m.UpsertTableFunc(ctx, table)
}
