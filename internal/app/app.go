package app

import (
	"fmt"
	"net"

	"github.com/alesr/getground/internal/app/partyctrl"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Create app struct
type App struct {
	logger    *zap.Logger
	fiberApp  *fiber.App
	partyCtrl partyctrl.PartyController
}

func New(logger *zap.Logger, fiberApp *fiber.App, partyCtrl partyctrl.PartyController) *App {
	return &App{
		logger:    logger.Named("party_app"),
		fiberApp:  fiberApp,
		partyCtrl: partyCtrl,
	}
}

func (a *App) Run(port string) error {
	a.fiberApp.Post("/guest_list/:name", a.partyCtrl.AddGuestToGuestList)
	a.fiberApp.Get("/guest_list", a.partyCtrl.GetGuestList)
	a.fiberApp.Put("/guests/:name", a.partyCtrl.WelcomeGuest)
	a.fiberApp.Delete("/guests/:name", a.partyCtrl.GoodbyeGuest)
	a.fiberApp.Get("/guests", a.partyCtrl.ListArrivedGuests)
	a.fiberApp.Get("/seats_empty", a.partyCtrl.GetEmptySeats)

	if err := a.fiberApp.Listen(net.JoinHostPort("", port)); err != nil {
		return fmt.Errorf("failed to serve http request: %w", err)
	}
	return nil
}
