package utils

import (
	"fmt"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func QParams(c echo.Context) *commonschema.QueryParams {
	q := commonschema.QueryParams{
		Page:    1,  // default
		Limit:   10, // default
		Search:  "",
		OrderBy: "",
	}

	if c.QueryParam("page") != "" {
		p, err := strconv.Atoi(c.QueryParam("page"))
		if err == nil {
			q.Page = p
		}
	}

	if c.QueryParam("limit") != "" {
		l, err := strconv.Atoi(c.QueryParam("limit"))
		if err == nil {
			q.Limit = l
		}
	}

	if c.QueryParam("search") != "" {
		q.Search = c.QueryParam("search")
	}

	if c.QueryParam("orderBy") != "" {
		o := strings.Split(c.QueryParam("orderBy"), ":")
		if len(o) > 1 {
			q.OrderBy = fmt.Sprintf("%s %s", o[0], o[1])
		}
	}

	return &q
}
