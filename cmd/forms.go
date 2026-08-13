package main

import (
	"net/http"

	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// GetForms handles retrieval of saved forms.
func (a *App) GetForms(c echo.Context) error {
	out, err := a.core.GetForms()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{out})
}

// GetForm handles retrieval of a single saved form.
func (a *App) GetForm(c echo.Context) error {
	id := getID(c)
	out, err := a.core.GetForm(id, "")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{out})
}

// CreateForm handles saved form creation.
func (a *App) CreateForm(c echo.Context) error {
	var o models.Form
	if err := c.Bind(&o); err != nil {
		return err
	}

	if err := a.validateForm(o); err != nil {
		return err
	}
	o.ListIDs = a.filterPublicListIDs(o.ListIDs)

	out, err := a.core.CreateForm(o)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{out})
}

// UpdateForm handles saved form modification.
func (a *App) UpdateForm(c echo.Context) error {
	var o models.Form
	if err := c.Bind(&o); err != nil {
		return err
	}

	if err := a.validateForm(o); err != nil {
		return err
	}
	o.ListIDs = a.filterPublicListIDs(o.ListIDs)

	id := getID(c)
	out, err := a.core.UpdateForm(id, o)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{out})
}

// DeleteForm handles saved form deletion.
func (a *App) DeleteForm(c echo.Context) error {
	id := getID(c)
	if err := a.core.DeleteForm(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{true})
}

// validateForm validates a saved form's fields.
func (a *App) validateForm(o models.Form) error {
	if !strHasLen(o.Name, 1, stdInputMaxLen) {
		return echo.NewHTTPError(http.StatusBadRequest, a.i18n.Ts("globals.messages.missingFields", "name", "name"))
	}

	return nil
}

// filterPublicListIDs narrows a list of IDs down to only those that
// belong to active, public lists, silently dropping the rest. This keeps
// a form's list_ids safe to expose on its public embed page regardless of
// what the client sent.
func (a *App) filterPublicListIDs(ids pq.Int64Array) pq.Int64Array {
	out := pq.Int64Array{}
	if len(ids) == 0 {
		return out
	}

	idsInt := make([]int, 0, len(ids))
	for _, id := range ids {
		idsInt = append(idsInt, int(id))
	}

	types, err := a.core.GetListTypes(idsInt, nil)
	if err != nil {
		return out
	}

	for _, id := range ids {
		if typ, ok := types[int(id)]; ok && typ == models.ListTypePublic {
			out = append(out, id)
		}
	}

	return out
}
