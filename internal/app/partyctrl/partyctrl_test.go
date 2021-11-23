package partyctrl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alesr/getground/internal/pkg/party"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testReqTimeoutMs int = 1000

func TestAddGuestToGuestList(t *testing.T) {
	t.Run("returns status created and the response matching with the mock result", func(t *testing.T) {
		service := party.Mock{}

		var wasCalled bool
		service.AddGuestToGuestListFunc = func(ctx context.Context, in *party.AddGuestToGuestListInput) (*party.AddGuestToGuestListOutput, error) {
			wasCalled = true
			return &party.AddGuestToGuestListOutput{
				Name: in.Name,
			}, nil
		}

		controller := New(zap.NewNop(), &service)

		givenName := "John"
		givenBody := `{"table": 1, "accompanying_guests": 2}`
		expectedStatusCode := http.StatusCreated
		expectedResponse := party.AddGuestToGuestListOutput{
			Name: givenName,
		}

		req := httptest.NewRequest(http.MethodPost, "/guest_list/"+givenName, bytes.NewBufferString(givenBody))
		req.Header.Set("Content-Type", "application/json")

		fiberApp := fiber.New()
		fiberApp.Post("/guest_list/:name", controller.AddGuestToGuestList)

		resp, err := fiberApp.Test(req, testReqTimeoutMs)
		require.NoError(t, err)

		var observed party.AddGuestToGuestListOutput
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&observed))

		assert.True(t, wasCalled)
		assert.Equal(t, expectedStatusCode, resp.StatusCode)
		assert.Equal(t, expectedResponse, observed)
	})
}

func TestGetGuestList(t *testing.T) {
	t.Run("returns status created for successful requests", func(t *testing.T) {
		service := party.Mock{}

		var wasCalled bool
		service.GetGuestListFunc = func(ctx context.Context) (party.GetGuestListOutput, error) {
			wasCalled = true
			return party.GetGuestListOutput{}, nil
		}

		controller := New(zap.NewNop(), &service)

		fiberApp := fiber.New()
		fiberApp.Get("/guest_list", controller.GetGuestList)

		resp, err := fiberApp.Test(httptest.NewRequest(fiber.MethodGet, "/guest_list", nil), testReqTimeoutMs)
		require.NoError(t, err)

		assert.True(t, wasCalled)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("returns the expected result given the party mock returned value", func(t *testing.T) {
		service := party.Mock{}

		service.GetGuestListFunc = func(ctx context.Context) (party.GetGuestListOutput, error) {
			return party.GetGuestListOutput{
				Guests: []party.Guest{
					{
						Name:               "John",
						Table:              1,
						AccompanyingGuests: 2,
					},
				},
			}, nil
		}

		controller := New(zap.NewNop(), &service)

		fiberApp := fiber.New()
		fiberApp.Get("/guest_list", controller.GetGuestList)

		resp, err := fiberApp.Test(httptest.NewRequest(fiber.MethodGet, "/guest_list", nil), testReqTimeoutMs)
		require.NoError(t, err)

		var result party.GetGuestListOutput
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		assert.Equal(t, party.GetGuestListOutput{
			Guests: []party.Guest{
				{
					Name:               "John",
					Table:              1,
					AccompanyingGuests: 2,
				},
			},
		}, result)

	})
}
