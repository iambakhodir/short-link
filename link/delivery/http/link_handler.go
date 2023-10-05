package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/iambakhodir/short-link/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type LinkHandler struct {
	LUseCae domain.LinkUseCase
}

func NewLinkHandler(e *echo.Echo, us domain.LinkUseCase) {
	handler := &LinkHandler{
		LUseCae: us,
	}

	e.GET("/links", handler.FetchLink)
	e.GET("/links/:id", handler.GetByID)
	e.POST("/links", handler.StoreLink)
	e.DELETE("/links/:id", handler.DeleteLink)
	e.GET("/:alias", handler.RedirectByAlias)
}

func (lh *LinkHandler) RedirectByAlias(c echo.Context) error {
	aliasParam := c.QueryParam("limit")
	ctx := c.Request().Context()

	link, err := lh.LUseCae.GetByAlias(ctx, aliasParam)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, link.Target)
}

func (lh *LinkHandler) GetByID(c echo.Context) error {
	idParam, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()

	link, err := lh.LUseCae.GetById(ctx, int64(idParam))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, link)
}

func (lh *LinkHandler) FetchLink(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitParam)
	//cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listLinks, err := lh.LUseCae.Fetch(ctx, int64(limit))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, listLinks)
}

func (lh *LinkHandler) StoreLink(c echo.Context) error {
	var link domain.Link
	err := c.Bind(&link)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	var ok bool
	if ok, err = isRequestValid(&link); !ok {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = lh.LUseCae.Store(ctx, &link)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, link)
}

func (lh *LinkHandler) DeleteLink(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	id := int64(idParam)
	ctx := c.Request().Context()

	err = lh.LUseCae.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func isRequestValid(m *domain.Link) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
