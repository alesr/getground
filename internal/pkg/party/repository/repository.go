package repository

import (
	"context"
	"time"
)

type (
	Guest struct {
		Name               string     `gorm:"uniqueIndex,unique,column:name"`
		Table              int        `gorm:"column:table;not null"`
		AccompanyingGuests int        `gorm:"column:accompanying_guests;not null"`
		TimeArrival        *time.Time `db:"time_arrival"`
		TimeDeparture      *time.Time `db:"time_departure"`
	}

	Table struct {
		Number         int `gorm:"uniqueIndex;column:number"`
		AvailableSeats int `gorm:"column:available_seats;not null"`
		Size           int `gorm:"column:size;not null"`
	}

	Repository interface {
		GetArrivedGuests(ctx context.Context) ([]Guest, error)
		GetGuestByName(ctx context.Context, name string) (*Guest, error)
		ListGuests(ctx context.Context) ([]Guest, error)
		UpsertGuest(ctx context.Context, guest *Guest) error

		GetTableByNumber(ctx context.Context, number int) (*Table, error)
		GetTables(ctx context.Context) ([]Table, error)
		UpsertTable(ctx context.Context, table *Table) error
	}
)
