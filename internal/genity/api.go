package genity

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/wondenge/go-genity/internal/errors"
	"github.com/wondenge/go-genity/pkg/log"
	"github.com/wondenge/go-genity/pkg/pagination"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/genitys/<id>", res.get)
	r.Get("/genitys", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/genitys", res.create)
	r.Put("/genitys/<id>", res.update)
	r.Delete("/genitys/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	genity, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(genity)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	genitys, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = genitys
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateGenityRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	genity, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(genity, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateGenityRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	genity, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(genity)
}

func (r resource) delete(c *routing.Context) error {
	genity, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(genity)
}
