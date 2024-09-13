package main

// Imports
import (
    "net/http"
    "github.com/labstack/echo"
    "github.com/google/uuid"
    "time"
)

// Models
type CreateJournalEntryCommand struct {
    Name string `json:"name"`
    Details *string `json:"details"`
    Date *time.Time `json:"date"`
}

type JournalEntry struct {
    Id uuid.UUID `json:"id"`
    Name string `json:"name"`
    Details *string `json:"details"`
    Date *time.Time `json:"date"`
    Tags string `json:"tags"`
}

// Service Interface
type Service interface {
    CreateJournalEntry(request *CreateJournalEntryCommand) (JournalEntry, error)
}

// Handlers
type Handler struct {
    service Service
}

func (h *Handler) CreateJournalEntry(c echo.Context) error {
    m := CreateJournalEntryCommand{}
	err := c.Bind(&m)
	if err != nil {
		return err
	}
    response, err := h.service.CreateJournalEntry(&m)
    if err != nil {
        return err
    }

    return c.JSON(http.StatusOK, response)
}

func (h *Handler) RegisterHandlers(e *echo.Echo, middleware echo.MiddlewareFunc) {
    e.POST("/create_journal_entry", h.CreateJournalEntry, middleware)
}