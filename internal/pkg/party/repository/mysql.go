package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alesr/getground/pkg/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var _ Repository = (*MySQL)(nil)

type MySQL struct {
	logger *zap.Logger
	dbConn *database.DBConn
}

func New(logger *zap.Logger, dbConn *database.DBConn) *MySQL {
	return &MySQL{
		logger: logger.Named("mysql_repository"),
		dbConn: dbConn,
	}
}

func (m *MySQL) GetArrivedGuests(ctx context.Context) ([]Guest, error) {
	var guests []Guest
	result := m.dbConn.Table("guests").Where("time_arrival IS NOT NULL").Find(&guests)
	if result.Error != nil {
		return nil, result.Error
	}

	return guests, nil
}

func (m *MySQL) GetGuestByName(ctx context.Context, name string) (*Guest, error) {
	var guest Guest
	result := m.dbConn.Table("guests").Find(&guest, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &guest, nil
}

func (m *MySQL) ListGuests(ctx context.Context) ([]Guest, error) {
	var guests []Guest
	result := m.dbConn.Table("guests").Find(&guests)
	if result.Error != nil {
		return nil, result.Error
	}
	return guests, nil
}

func (m *MySQL) UpsertGuest(ctx context.Context, guest *Guest) error {
	var g Guest
	result := m.dbConn.Table("guests").Where("name = ?", guest.Name).Find(&g)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("could not find guest: %w", result.Error)
	}

	// if user not found, create new user
	if result.RowsAffected == 0 {
		result = m.dbConn.Table("guests").Create(guest)
		if result.Error != nil {
			return fmt.Errorf("could not create guest: %w", result.Error)
		}
		return nil
	}

	// update user
	result = m.dbConn.Table("guests").Where("name = ?", g.Name).Update(guest)
	if result.Error != nil {
		return fmt.Errorf("could not update guest: %w", result.Error)
	}
	return nil
}

func (m *MySQL) GetTableByNumber(ctx context.Context, number int) (*Table, error) {
	var table Table
	result := m.dbConn.Table("tables").Where("number = ?", number).Find(&table)
	if result.Error != nil {
		return nil, result.Error
	}
	return &table, nil
}

func (m *MySQL) GetTables(ctx context.Context) ([]Table, error) {
	var tables []Table
	result := m.dbConn.Table("tables").Find(&tables)
	if result.Error != nil {
		return nil, result.Error
	}
	return tables, nil
}

func (m *MySQL) UpsertTable(ctx context.Context, tbl *Table) error {
	var table Table
	result := m.dbConn.Table("tables").Where("number = ?", tbl.Number).Find(&table)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("could not find table: %w", result.Error)
	}

	// if user not found, create new user
	if result.RowsAffected == 0 {
		result = m.dbConn.Table("tables").Create(tbl)
		if result.Error != nil {
			return fmt.Errorf("could not create table: %w", result.Error)
		}
		return nil
	}

	// update user
	result = m.dbConn.Table("tables").Where("number = ?", table.Number).Update(tbl)
	if result.Error != nil {
		return fmt.Errorf("could not update table: %w", result.Error)
	}
	return nil
}
