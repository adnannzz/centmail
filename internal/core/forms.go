package core

import (
	"database/sql"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
)

// GetForms retrieves all saved forms.
func (c *Core) GetForms() ([]models.Form, error) {
	out := []models.Form{}
	if err := c.q.GetForms.Select(&out); err != nil {
		c.log.Printf("error fetching forms: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.forms}", "error", pqErrMsg(err)))
	}

	return out, nil
}

// GetForm retrieves a form by its ID or UUID.
func (c *Core) GetForm(id int, formUUID string) (models.Form, error) {
	var out models.Form
	if err := c.q.GetForm.Get(&out, id, formUUID); err != nil {
		if err == sql.ErrNoRows {
			return models.Form{}, echo.NewHTTPError(http.StatusBadRequest,
				c.i18n.Ts("globals.messages.notFound", "name", "{globals.terms.form}"))
		}

		c.log.Printf("error fetching form: %v", err)
		return models.Form{}, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.form}", "error", pqErrMsg(err)))
	}

	return out, nil
}

// CreateForm creates a new saved form.
func (c *Core) CreateForm(f models.Form) (models.Form, error) {
	uu, err := uuid.NewV4()
	if err != nil {
		c.log.Printf("error generating UUID: %v", err)
		return models.Form{}, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorUUID", "error", err.Error()))
	}

	var newID int
	if err := c.q.CreateForm.Get(&newID, uu.String(), f.Name, f.ListIDs, f.RedirectURL); err != nil {
		c.log.Printf("error creating form: %v", err)
		return models.Form{}, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorCreating", "name", "{globals.terms.form}", "error", pqErrMsg(err)))
	}

	return c.GetForm(newID, "")
}

// UpdateForm updates a given form.
func (c *Core) UpdateForm(id int, f models.Form) (models.Form, error) {
	res, err := c.q.UpdateForm.Exec(id, f.Name, f.ListIDs, f.RedirectURL)
	if err != nil {
		c.log.Printf("error updating form: %v", err)
		return models.Form{}, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorUpdating", "name", "{globals.terms.form}", "error", pqErrMsg(err)))
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return models.Form{}, echo.NewHTTPError(http.StatusBadRequest,
			c.i18n.Ts("globals.messages.notFound", "name", "{globals.terms.form}"))
	}

	return c.GetForm(id, "")
}

// DeleteForm deletes a saved form.
func (c *Core) DeleteForm(id int) error {
	if _, err := c.q.DeleteForm.Exec(id); err != nil {
		c.log.Printf("error deleting form: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorDeleting", "name", "{globals.terms.form}", "error", pqErrMsg(err)))
	}

	return nil
}
