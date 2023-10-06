package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/iambakhodir/short-link/domain"
	"github.com/iambakhodir/short-link/domain/random"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseSuccessObject struct {
	Message string      `json:"message"`
	Data    domain.Link `json:"data"`
}

type ResponseSuccessArray struct {
	Message string        `json:"message"`
	Data    []domain.Link `json:"data"`
}

type LinkHandler struct {
	LUseCae domain.LinkUseCase
}

func NewLinkHandler(e *echo.Echo, us domain.LinkUseCase) {
	handler := &LinkHandler{
		LUseCae: us,
	}

	e.GET("/links", handler.FetchLinks)
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

	return c.JSON(http.StatusOK, ResponseSuccessObject{Message: "ok", Data: link})
}

func (lh *LinkHandler) FetchLinks(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitParam)
	//cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listLinks, err := lh.LUseCae.Fetch(ctx, int64(limit))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ResponseSuccessArray{Message: "ok", Data: listLinks})
}

func (lh *LinkHandler) StoreLink(c echo.Context) error {
	var req domain.LinkRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	var ok bool
	if ok, err = isRequestValid(&req); !ok {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	alias := req.Alias
	length := req.Length

	if alias == "" {
		if length > 0 {
			alias = random.NewRandomString(length) //TODO improve generator
		} else {
			alias = random.NewRandomString(viper.GetInt("alias_length")) //TODO improve generator
		}
	}

	ctx := c.Request().Context()

	id, err := lh.LUseCae.Store(ctx, domain.Link{
		Target: req.Target,
		Alias:  alias,
	})
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	link, err := lh.LUseCae.GetById(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ResponseSuccessObject{Message: "ok", Data: link})
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
	case domain.ErrLinkIsExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func isRequestValid(m *domain.LinkRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
