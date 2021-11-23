package repository

import (
	"context"
	"testing"
	"time"

	"github.com/alesr/getground/pkg/database"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	truncateGuestsQuery string = "TRUNCATE guests;"
	truncateTablesQuery string = "TRUNCATE tables;"
)

func TestGetArrivedGuests_INTEGRATION(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	// Arrange

	dbConn := setupDB(t)
	defer dbConn.Close()
	defer truncateHelper(t, dbConn)

	// Create guests
	truncateHelper(t, dbConn)

	guest1 := Guest{
		Name:               "Guest1",
		Table:              1,
		AccompanyingGuests: 3,
	}

	now := time.Now()
	guest2 := Guest{
		Name:               "Guest2",
		Table:              2,
		AccompanyingGuests: 2,
		TimeArrival:        &now,
	}

	createGuestHelper(t, dbConn, &guest1, &guest2)

	repo := New(zap.NewNop(), dbConn)

	// Act

	observed, err := repo.GetArrivedGuests(context.TODO())
	require.NoError(t, err)

	// Assert

	require.Equal(t, 1, len(observed))
	require.Equal(t, guest2.Name, observed[0].Name)
}

func TestGetGuestByName_INTEGRATION(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	// Arrange

	dbConn := setupDB(t)
	defer dbConn.Close()
	defer truncateHelper(t, dbConn)

	// Create guests
	truncateHelper(t, dbConn)

	guest1 := Guest{
		Name:               "Guest1",
		Table:              1,
		AccompanyingGuests: 3,
	}

	createGuestHelper(t, dbConn, &guest1)

	repo := New(zap.NewNop(), dbConn)

	// Act

	observed, err := repo.GetGuestByName(context.TODO(), guest1.Name)
	require.NoError(t, err)

	// Assert

	require.Equal(t, guest1.Name, observed.Name)
}

func TestListGuests_INTEGRATION(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	// Arrange

	dbConn := setupDB(t)
	defer dbConn.Close()
	defer truncateHelper(t, dbConn)

	// Create guests
	truncateHelper(t, dbConn)

	guest1 := Guest{
		Name:               "Guest1",
		Table:              1,
		AccompanyingGuests: 3,
	}

	createGuestHelper(t, dbConn, &guest1)

	repo := New(zap.NewNop(), dbConn)

	// Act

	observed, err := repo.ListGuests(context.TODO())
	require.NoError(t, err)

	// Assert

	require.Equal(t, guest1.Name, observed[0].Name)
}

func TestUpsertGuest_INTEGRATION(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	t.Run("upserts an non existing guest", func(t *testing.T) {
		// Arrange

		dbConn := setupDB(t)
		defer dbConn.Close()
		defer truncateHelper(t, dbConn)

		truncateHelper(t, dbConn)

		guest1 := Guest{
			Name:               "Guest1",
			Table:              1,
			AccompanyingGuests: 3,
		}

		repo := New(zap.NewNop(), dbConn)

		// Act

		err := repo.UpsertGuest(context.TODO(), &guest1)
		require.NoError(t, err)

		// Assert

		observed, err := repo.GetGuestByName(context.TODO(), guest1.Name)
		require.NoError(t, err)

		require.Equal(t, guest1.Name, observed.Name)
	})

	t.Run("upserts an existing guest", func(t *testing.T) {
		// Arrange

		dbConn := setupDB(t)
		defer dbConn.Close()
		defer truncateHelper(t, dbConn)

		truncateHelper(t, dbConn)

		guest1 := Guest{
			Name:               "Guest1",
			Table:              1,
			AccompanyingGuests: 3,
		}

		createGuestHelper(t, dbConn, &guest1)

		repo := New(zap.NewNop(), dbConn)

		// Act

		guest1.Name = "Guest1_UPDATED"

		err := repo.UpsertGuest(context.TODO(), &guest1)
		require.NoError(t, err)

		// Assert

		observed, err := repo.GetGuestByName(context.TODO(), guest1.Name)
		require.NoError(t, err)

		require.Equal(t, guest1.Name, observed.Name)
	})
}

func TestGetTableByNumber_INTEGRATION(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	t.Run("happy path", func(t *testing.T) {
		// Arrange
		dbConn := setupDB(t)
		defer dbConn.Close()
		defer truncateHelper(t, dbConn)

		truncateHelper(t, dbConn)

		table := Table{
			Number: 1,
		}

		createTableHelper(t, dbConn, &table)

		repo := New(zap.NewNop(), dbConn)

		// Act
		observed, err := repo.GetTableByNumber(context.TODO(), table.Number)
		require.NoError(t, err)

		// Assert
		require.Equal(t, table.Number, observed.Number)
	})

	t.Run("tries to get an non existing table", func(t *testing.T) {
		// Arrange
		dbConn := setupDB(t)
		defer dbConn.Close()
		defer truncateHelper(t, dbConn)

		truncateHelper(t, dbConn)

		repo := New(zap.NewNop(), dbConn)

		// Act
		observed, err := repo.GetTableByNumber(context.TODO(), 1)
		require.Error(t, err)
		require.Nil(t, observed)
	})
}

func TestGetTables_INTEGRATION(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	t.Run("happy path", func(t *testing.T) {
		// Arrange
		dbConn := setupDB(t)
		defer dbConn.Close()
		defer truncateHelper(t, dbConn)

		truncateHelper(t, dbConn)

		table1 := Table{
			Number: 1,
		}

		table2 := Table{
			Number: 2,
		}

		createTableHelper(t, dbConn, &table1)
		createTableHelper(t, dbConn, &table2)

		repo := New(zap.NewNop(), dbConn)

		// Act
		observed, err := repo.GetTables(context.TODO())
		require.NoError(t, err)

		// Assert
		require.Equal(t, 2, len(observed))
	})
}

func TestUpsertTable_INTEGRATION(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	t.Run("upserts an non existing table", func(t *testing.T) {
		// Arrange

		dbConn := setupDB(t)
		defer dbConn.Close()
		defer truncateHelper(t, dbConn)

		truncateHelper(t, dbConn)

		table1 := Table{
			Number:         1,
			AvailableSeats: 3,
			Size:           12,
		}

		repo := New(zap.NewNop(), dbConn)

		// Act

		err := repo.UpsertTable(context.TODO(), &table1)
		require.NoError(t, err)

		// Assert

		observed, err := repo.GetTableByNumber(context.TODO(), table1.Number)
		require.NoError(t, err)

		require.Equal(t, table1.Number, observed.Number)
	})

	t.Run("upserts an existing table", func(t *testing.T) {
		// Arrange

		dbConn := setupDB(t)
		defer dbConn.Close()
		defer truncateHelper(t, dbConn)

		truncateHelper(t, dbConn)

		table1 := Table{
			Number:         1,
			AvailableSeats: 3,
			Size:           12,
		}

		createTableHelper(t, dbConn, &table1)

		repo := New(zap.NewNop(), dbConn)

		// Act

		table1.Number = 2

		err := repo.UpsertTable(context.TODO(), &table1)
		require.NoError(t, err)

		// Assert

		observed, err := repo.GetTableByNumber(context.TODO(), table1.Number)
		require.NoError(t, err)

		require.Equal(t, table1.Number, observed.Number)
	})
}

func createGuestHelper(t *testing.T, dbConn *database.DBConn, guest ...*Guest) {
	for _, g := range guest {
		err := dbConn.Table("guests").Create(g).Error
		require.NoError(t, err)
	}
}

func createTableHelper(t *testing.T, dbConn *database.DBConn, table ...*Table) {
	for _, tbl := range table {
		err := dbConn.Table("tables").Create(tbl).Error
		require.NoError(t, err)
	}
}

func setupDB(t *testing.T, tables ...interface{}) *database.DBConn {
	host := "127.0.0.1"
	user := "user"
	password := "password"
	dbName := "party_db"

	dbConn, err := database.Connection(host, user, password, dbName)
	require.NoError(t, err)

	dbConn.Table("guests").AutoMigrate(&Guest{})
	dbConn.Table("tables").AutoMigrate(&Table{})

	return dbConn
}

func truncateHelper(t *testing.T, dbConn *database.DBConn) {
	dbConn.Exec(truncateGuestsQuery)
	dbConn.Exec(truncateTablesQuery)
}
