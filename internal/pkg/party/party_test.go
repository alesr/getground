package party

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alesr/getground/internal/pkg/party/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	errTestRepo   error = errors.New("test repository error")
	testTableSize int   = 1000
)

func TestNew(t *testing.T) {
	t.Run("returns a new party", func(t *testing.T) {
		party := New(zap.NewNop(), &repository.Mock{}, testTableSize)
		assert.NotNil(t, party)
	})
}

func TestAddGuestToGuestList(t *testing.T) {
	cases := []struct {
		name           string
		given          *AddGuestToGuestListInput
		expectedResult *AddGuestToGuestListOutput
		expectedError  error
	}{
		{
			name: "required name input",
			given: &AddGuestToGuestListInput{
				Name:               "",
				Table:              456,
				AccompanyingGuests: 789,
			},
			expectedResult: nil,
			expectedError:  ErrGuestNameRequired,
		},
		{
			name: "required table input",
			given: &AddGuestToGuestListInput{
				Name:               "123",
				Table:              0,
				AccompanyingGuests: 789,
			},
			expectedResult: nil,
			expectedError:  ErrTableNumberRequired,
		},
		{
			name: "invalid table input",
			given: &AddGuestToGuestListInput{
				Name:               "123",
				Table:              -1,
				AccompanyingGuests: 789,
			},
			expectedResult: nil,
			expectedError:  ErrTableNumberInvalid,
		},
		{
			name: "invalid accompanying guests input",
			given: &AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: -1,
			},
			expectedResult: nil,
			expectedError:  ErrAccompanyingGuestsNumberInvalid,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			party := New(zap.NewNop(), &repository.Mock{}, testTableSize)
			observed, err := party.AddGuestToGuestList(context.TODO(), c.given)

			assert.Equal(t, c.expectedResult, observed)
			assert.Equal(t, c.expectedError, err)
		})
	}

	t.Run("returns and error when GetGuestList fails", func(t *testing.T) {
		repo := repository.Mock{}
		repo.ListGuestsFunc = func(ctx context.Context) ([]repository.Guest, error) {
			return nil, errTestRepo
		}

		party := New(zap.NewNop(), &repo, testTableSize)
		observed, err := party.AddGuestToGuestList(
			context.TODO(),
			&AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 789,
			})

		assert.Nil(t, observed)
		assert.Equal(t, errTestRepo, err)
	})

	t.Run("returns an error when guest is already in guest list", func(t *testing.T) {
		repo := repository.Mock{}
		repo.ListGuestsFunc = func(ctx context.Context) ([]repository.Guest, error) {
			return []repository.Guest{
				{
					Name: "123",
				},
			}, nil
		}

		repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error {
			return ErrGuestAlreadyInList
		}

		party := New(zap.NewNop(), &repo, testTableSize)
		observed, err := party.AddGuestToGuestList(
			context.TODO(),
			&AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 789,
			})

		assert.Nil(t, observed)
		assert.Equal(t, ErrGuestAlreadyInList, err)
	})

	t.Run("returns an error when table has not enough available seats", func(t *testing.T) {
		repo := repository.Mock{}
		repo.ListGuestsFunc = func(ctx context.Context) ([]repository.Guest, error) {
			return nil, nil
		}

		repo.GetTableByNumberFunc = func(ctx context.Context, number int) (*repository.Table, error) {
			return &repository.Table{
				Number:         number,
				AvailableSeats: 1,
			}, nil
		}

		repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error {
			return ErrTableNotEnoughSeats
		}

		party := New(zap.NewNop(), &repo, testTableSize)
		observed, err := party.AddGuestToGuestList(
			context.TODO(),
			&AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 789,
			})

		assert.Nil(t, observed)
		assert.Equal(t, ErrTableNotEnoughSeats, err)
	})

	t.Run("return no error when the table has enough available seats", func(t *testing.T) {
		repo := repository.Mock{}
		repo.ListGuestsFunc = func(ctx context.Context) ([]repository.Guest, error) {
			return nil, nil
		}

		repo.GetTableByNumberFunc = func(ctx context.Context, number int) (*repository.Table, error) {
			return &repository.Table{
				Number:         number,
				AvailableSeats: 2,
			}, nil
		}

		repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error {
			return nil
		}

		repo.UpsertTableFunc = func(ctx context.Context, table *repository.Table) error {
			return nil
		}

		party := New(zap.NewNop(), &repo, testTableSize)
		observed, err := party.AddGuestToGuestList(
			context.TODO(),
			&AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 1,
			})

		assert.Nil(t, err)
		assert.NotNil(t, observed)
	})

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := repository.Mock{}
			repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error {
				return c.expectedError
			}

			party := New(zap.NewNop(), &repo, testTableSize)
			observed, err := party.AddGuestToGuestList(context.TODO(), c.given)

			assert.Equal(t, c.expectedResult, observed)
			assert.Equal(t, c.expectedError, err)
		})
	}
	// TODO(alesr): more tests should be added, but we're going to leave this one for now.
}

func TestGetGuestList(t *testing.T) {
	cases := []struct {
		name            string
		givenMockResult []repository.Guest
		expectedResult  GetGuestListOutput
		expectedError   error
	}{
		{
			name: "returns the guest list",
			givenMockResult: []repository.Guest{
				{
					Name:               "123",
					Table:              456,
					AccompanyingGuests: 789,
				},
			},
			expectedResult: GetGuestListOutput{
				Guests: []Guest{
					{
						Name:               "123",
						Table:              456,
						AccompanyingGuests: 789,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:            "returns an error when GetGuestList fails",
			givenMockResult: nil,
			expectedResult:  GetGuestListOutput{},
			expectedError:   errTestRepo,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := repository.Mock{}
			repo.ListGuestsFunc = func(ctx context.Context) ([]repository.Guest, error) {
				return c.givenMockResult, c.expectedError
			}

			party := New(zap.NewNop(), &repo, testTableSize)
			observed, err := party.GetGuestList(context.TODO())

			assert.Equal(t, c.expectedResult, observed)
			assert.Equal(t, c.expectedError, err)
		})
	}
}

func TestCreateTableAndAddToGuestList(t *testing.T) {
	cases := []struct {
		name                         string
		given                        *AddGuestToGuestListInput
		givenCreateTableMockErr      error
		givenAddToGuestListMockErr   error
		expectedError                error
		expectedCreateTableCalled    bool
		expectedAddToGuestListCalled bool
	}{
		{
			name: "returns an error when table creation fails",
			given: &AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 789,
			},
			givenCreateTableMockErr:      errTestRepo,
			givenAddToGuestListMockErr:   nil,
			expectedError:                errTestRepo,
			expectedCreateTableCalled:    true,
			expectedAddToGuestListCalled: false,
		},
		{
			name: "returns an error when table creation succeeds but guest addition fails",
			given: &AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 789,
			},
			givenCreateTableMockErr:      nil,
			givenAddToGuestListMockErr:   errTestRepo,
			expectedError:                errTestRepo,
			expectedCreateTableCalled:    true,
			expectedAddToGuestListCalled: true,
		},
		{
			name: "returns no error when table creation succeeds and guest addition succeeds",
			given: &AddGuestToGuestListInput{
				Name:               "123",
				Table:              456,
				AccompanyingGuests: 789,
			},
			givenCreateTableMockErr:      nil,
			givenAddToGuestListMockErr:   nil,
			expectedError:                nil,
			expectedCreateTableCalled:    true,
			expectedAddToGuestListCalled: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var createTableCalled, addToGuestListCalled bool

			repo := repository.Mock{}
			repo.UpsertTableFunc = func(ctx context.Context, table *repository.Table) error {
				createTableCalled = true
				return tc.givenCreateTableMockErr
			}

			repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error {
				addToGuestListCalled = true
				return tc.givenAddToGuestListMockErr
			}

			party := New(zap.NewNop(), &repo, testTableSize)
			observedErr := party.createTableAndAddToGuestList(context.TODO(), tc.given)

			assert.Equal(t, tc.expectedError, observedErr)
			assert.Equal(t, tc.expectedCreateTableCalled, createTableCalled)
			assert.Equal(t, tc.expectedAddToGuestListCalled, addToGuestListCalled)
		})
	}
}

func TestWelcomeGuest(t *testing.T) {
	cases := []struct {
		name                        string
		givenGetGuestByNameMockFn   func(ctx context.Context, name string) (*repository.Guest, error)
		givenGetTableByNumberMockFn func(ctx context.Context, number int) (*repository.Table, error)
		givenUpsertTableFn          func(ctx context.Context, table *repository.Table) error
		givenUpsertGuestFn          func(ctx context.Context, guest *repository.Guest) error
		expectedError               error
	}{
		{
			name: "returns an error when get guest fails",
			givenGetGuestByNameMockFn: func(ctx context.Context, name string) (*repository.Guest, error) {
				return nil, errTestRepo
			},
			givenGetTableByNumberMockFn: func(ctx context.Context, number int) (*repository.Table, error) { return nil, nil },
			givenUpsertTableFn:          func(ctx context.Context, table *repository.Table) error { return nil },
			givenUpsertGuestFn:          func(ctx context.Context, guest *repository.Guest) error { return nil },
			expectedError:               errTestRepo,
		},
		{
			name: "returns an error when get table by number fails",
			givenGetGuestByNameMockFn: func(ctx context.Context, name string) (*repository.Guest, error) {
				return &repository.Guest{
					Name: "123",
				}, nil
			},
			givenGetTableByNumberMockFn: func(ctx context.Context, number int) (*repository.Table, error) { return nil, errTestRepo },
			givenUpsertTableFn:          func(ctx context.Context, table *repository.Table) error { return nil },
			givenUpsertGuestFn:          func(ctx context.Context, guest *repository.Guest) error { return nil },
			expectedError:               errTestRepo,
		},
		{
			name: "returns an error when upsert table fails",
			givenGetGuestByNameMockFn: func(ctx context.Context, name string) (*repository.Guest, error) {
				return &repository.Guest{
					Name: "123",
				}, nil
			},
			givenGetTableByNumberMockFn: func(ctx context.Context, number int) (*repository.Table, error) {
				return &repository.Table{
					Number:         456,
					Size:           testTableSize,
					AvailableSeats: 789,
				}, nil
			},
			givenUpsertTableFn: func(ctx context.Context, table *repository.Table) error { return errTestRepo },
			givenUpsertGuestFn: func(ctx context.Context, guest *repository.Guest) error { return nil },
			expectedError:      errTestRepo,
		},
		{
			name: "returns an error when guest upsert fails",
			givenGetGuestByNameMockFn: func(ctx context.Context, name string) (*repository.Guest, error) {
				return &repository.Guest{
					Name: "123",
				}, nil
			},
			givenGetTableByNumberMockFn: func(ctx context.Context, number int) (*repository.Table, error) {
				return &repository.Table{
					Number:         456,
					Size:           testTableSize,
					AvailableSeats: 789,
				}, nil
			},
			givenUpsertTableFn: func(ctx context.Context, table *repository.Table) error { return nil },
			givenUpsertGuestFn: func(ctx context.Context, guest *repository.Guest) error { return errTestRepo },
			expectedError:      errTestRepo,
		},
		{
			name: "returns no error when everything succeeds",
			givenGetGuestByNameMockFn: func(ctx context.Context, name string) (*repository.Guest, error) {
				return &repository.Guest{
					Name: "123",
				}, nil
			},
			givenGetTableByNumberMockFn: func(ctx context.Context, number int) (*repository.Table, error) {
				return &repository.Table{
					Number:         456,
					Size:           testTableSize,
					AvailableSeats: 789,
				}, nil
			},
			givenUpsertTableFn: func(ctx context.Context, table *repository.Table) error { return nil },
			givenUpsertGuestFn: func(ctx context.Context, guest *repository.Guest) error { return nil },
			expectedError:      nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := repository.Mock{}
			repo.GetGuestByNameFunc = tc.givenGetGuestByNameMockFn
			repo.GetTableByNumberFunc = tc.givenGetTableByNumberMockFn
			repo.UpsertTableFunc = tc.givenUpsertTableFn
			repo.UpsertGuestFunc = tc.givenUpsertGuestFn

			party := New(zap.NewNop(), &repo, testTableSize)
			_, observedErr := party.WelcomeGuest(
				context.TODO(),
				&WelcomeGuestInput{
					Name:               "123",
					AccompanyingGuests: 456,
				})

			assert.Equal(t, tc.expectedError, observedErr)
		})
	}

	t.Run("returns an error when guest is not in the list", func(t *testing.T) {
		repo := repository.Mock{}
		repo.GetGuestByNameFunc = func(ctx context.Context, name string) (*repository.Guest, error) {
			return nil, nil
		}
		repo.GetTableByNumberFunc = func(ctx context.Context, number int) (*repository.Table, error) {
			return &repository.Table{
				Number:         456,
				Size:           testTableSize,
				AvailableSeats: 789,
			}, nil
		}
		repo.UpsertTableFunc = func(ctx context.Context, table *repository.Table) error { return nil }
		repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error { return nil }

		party := New(zap.NewNop(), &repo, testTableSize)
		_, observedErr := party.WelcomeGuest(
			context.TODO(),
			&WelcomeGuestInput{
				Name:               "123",
				AccompanyingGuests: 456,
			})

		assert.Equal(t, ErrGuestNotInList, observedErr)
	})

	t.Run("returns an error when the table is not found", func(t *testing.T) {
		repo := repository.Mock{}
		repo.GetGuestByNameFunc = func(ctx context.Context, name string) (*repository.Guest, error) {
			return &repository.Guest{
				Name: "123",
			}, nil
		}
		repo.GetTableByNumberFunc = func(ctx context.Context, number int) (*repository.Table, error) {
			return nil, nil
		}
		repo.UpsertTableFunc = func(ctx context.Context, table *repository.Table) error { return nil }
		repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error { return nil }

		party := New(zap.NewNop(), &repo, testTableSize)
		_, observedErr := party.WelcomeGuest(
			context.TODO(),
			&WelcomeGuestInput{
				Name:               "123",
				AccompanyingGuests: 456,
			})

		assert.Equal(t, ErrTableNumberNotFound, observedErr)
	})

	t.Run("returns an error when the table is full", func(t *testing.T) {
		repo := repository.Mock{}
		repo.GetGuestByNameFunc = func(ctx context.Context, name string) (*repository.Guest, error) {
			return &repository.Guest{
				Name: "123",
			}, nil
		}
		repo.GetTableByNumberFunc = func(ctx context.Context, number int) (*repository.Table, error) {
			return &repository.Table{
				Number:         456,
				Size:           testTableSize,
				AvailableSeats: 0,
			}, nil
		}
		repo.UpsertTableFunc = func(ctx context.Context, table *repository.Table) error { return nil }
		repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error { return nil }

		party := New(zap.NewNop(), &repo, testTableSize)
		_, observedErr := party.WelcomeGuest(
			context.TODO(),
			&WelcomeGuestInput{
				Name:               "123",
				AccompanyingGuests: 456,
			})

		assert.Equal(t, ErrTableNotEnoughSeats, observedErr)
	})

	t.Run("returns no error and the expect output matches", func(t *testing.T) {
		repo := repository.Mock{}
		repo.GetGuestByNameFunc = func(ctx context.Context, name string) (*repository.Guest, error) {
			return &repository.Guest{
				Name: "123",
			}, nil
		}
		repo.GetTableByNumberFunc = func(ctx context.Context, number int) (*repository.Table, error) {
			return &repository.Table{
				Number:         456,
				Size:           testTableSize,
				AvailableSeats: 789,
			}, nil
		}
		repo.UpsertTableFunc = func(ctx context.Context, table *repository.Table) error { return nil }
		repo.UpsertGuestFunc = func(ctx context.Context, guest *repository.Guest) error { return nil }

		party := New(zap.NewNop(), &repo, testTableSize)
		observedOutput, observedErr := party.WelcomeGuest(
			context.TODO(),
			&WelcomeGuestInput{
				Name:               "123",
				AccompanyingGuests: 456,
			})

		assert.Nil(t, observedErr)
		assert.Equal(t, &WelcomeGuestOutput{Name: "123"}, observedOutput)
	})
}

func TestListArrivedGuests(t *testing.T) {
	testTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name                  string
		givenGetArrivedGuests func(ctx context.Context) ([]repository.Guest, error)
		expectedOutput        ListArrivedGuestsOutput
		expectedError         error
	}{
		{
			name: "returns an error when the repository returns an error",
			givenGetArrivedGuests: func(ctx context.Context) ([]repository.Guest, error) {
				return nil, errors.New("some error")
			},
			expectedError: errors.New("some error"),
		},
		{
			name: "returns no error and the expect output matches",
			givenGetArrivedGuests: func(ctx context.Context) ([]repository.Guest, error) {
				return []repository.Guest{
					{
						Name:               "123",
						AccompanyingGuests: 456,
						TimeArrival:        &testTime,
					},
				}, nil
			},
			expectedOutput: ListArrivedGuestsOutput{
				Guests: []GuestArrived{
					{
						Name:               "123",
						AccompanyingGuests: 456,
						TimeArrival:        testTime,
					},
				},
			},
		},
		{
			name: "returns no error and the expect output matches when there are no guests",
			givenGetArrivedGuests: func(ctx context.Context) ([]repository.Guest, error) {
				return []repository.Guest{}, nil
			},
			expectedOutput: ListArrivedGuestsOutput{
				Guests: []GuestArrived{},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := repository.Mock{}
			repo.GetArrivedGuestsFunc = c.givenGetArrivedGuests

			party := New(zap.NewNop(), &repo, testTableSize)
			observedOutput, observedError := party.ListArrivedGuests(context.TODO())

			assert.Equal(t, c.expectedOutput, observedOutput)
			assert.Equal(t, c.expectedError, observedError)
		})
	}
}

func TestGetEmptySeats(t *testing.T) {
	cases := []struct {
		name               string
		givenGetTablesMock func(ctx context.Context) ([]repository.Table, error)
		expectedOutput     GetEmptySeatsOutput
		expectedError      error
	}{
		{
			name: "returns an error when the repository returns an error",
			givenGetTablesMock: func(ctx context.Context) ([]repository.Table, error) {
				return nil, errors.New("some error")
			},
			expectedError: errors.New("some error"),
		},
		{
			name: "returns no error and the expect output matches",
			givenGetTablesMock: func(ctx context.Context) ([]repository.Table, error) {
				return []repository.Table{
					{
						Number:         456,
						Size:           testTableSize,
						AvailableSeats: 789,
					},
				}, nil
			},
			expectedOutput: GetEmptySeatsOutput{
				EmptySeats: 789,
			},
		},
		{
			name: "returns no error and the expect output matches when there are no tables",
			givenGetTablesMock: func(ctx context.Context) ([]repository.Table, error) {
				return []repository.Table{}, nil
			},
			expectedOutput: GetEmptySeatsOutput{
				EmptySeats: 0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := repository.Mock{}
			repo.GetTablesFunc = c.givenGetTablesMock

			party := New(zap.NewNop(), &repo, testTableSize)
			observedOutput, observedError := party.GetEmptySeats(context.TODO())

			assert.Equal(t, c.expectedOutput, observedOutput)
			assert.Equal(t, c.expectedError, observedError)
		})
	}
}
