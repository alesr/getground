package partyctrl

import (
	"context"
	"net/http"

	"github.com/alesr/getground/internal/pkg/party"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PartyController interface {
	AddGuestToGuestList(c *fiber.Ctx) error
	GetGuestList(c *fiber.Ctx) error
	WelcomeGuest(c *fiber.Ctx) error
	GoodbyeGuest(c *fiber.Ctx) error
	ListArrivedGuests(c *fiber.Ctx) error
	GetEmptySeats(c *fiber.Ctx) error
}

type Controller struct {
	logger  *zap.Logger
	service party.Service
}

func New(logger *zap.Logger, service party.Service) *Controller {
	return &Controller{
		logger:  logger,
		service: service,
	}
}

func (ctrl *Controller) AddGuestToGuestList(c *fiber.Ctx) error {
	var req party.AddGuestToGuestListInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	req.Name = c.Params("name")
	resp, err := ctrl.service.AddGuestToGuestList(context.TODO(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	return c.Status(http.StatusCreated).JSON(resp)
}

func (ctrl *Controller) GetGuestList(c *fiber.Ctx) error {
	resp, err := ctrl.service.GetGuestList(context.TODO())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	return c.JSON(resp)
}

func (ctrl *Controller) WelcomeGuest(c *fiber.Ctx) error {
	var req party.WelcomeGuestInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	resp, err := ctrl.service.WelcomeGuest(context.TODO(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	return c.JSON(resp)
}

func (ctrl *Controller) GoodbyeGuest(c *fiber.Ctx) error {
	var req party.GoodbyeGuestInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	if err := ctrl.service.GoodbyeGuest(context.TODO(), &req); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(nil)
}

func (ctrl *Controller) ListArrivedGuests(c *fiber.Ctx) error {
	resp, err := ctrl.service.ListArrivedGuests(context.TODO())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	return c.JSON(resp)
}

func (ctrl *Controller) GetEmptySeats(c *fiber.Ctx) error {
	resp, err := ctrl.service.GetEmptySeats(context.TODO())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	return c.JSON(resp)
}
