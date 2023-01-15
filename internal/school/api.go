package school

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ujianmu/go-skeleton/internal/errors"
	"github.com/ujianmu/go-skeleton/pkg/log"
	"github.com/ujianmu/go-skeleton/pkg/pagination"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/school/<id>", res.get)
	r.Get("/school", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/school", res.create)
	r.Put("/school/<id>", res.update)
	r.Delete("/school/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	school, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(school)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	schools, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = schools
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateSchoolRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	school, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(school, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateSchoolRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	school, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(school)
}

func (r resource) delete(c *routing.Context) error {
	school, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(school)
}
