package meroute

import (
	medi "kiraform/src/applications/dependencies/me"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	meschema "kiraform/src/interfaces/rest/schemas/me"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MeHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies medi.MeDependencies
}

func NewMeHandler(DB *gorm.DB, validator *validator.Validate, dependencies medi.MeDependencies) *MeHandler {
	return &MeHandler{
		DB:           DB,
		Validator:    validator,
		Dependencies: dependencies,
	}
}

func NewMeHTTP(g *echo.Group, DB *gorm.DB) {
	h := NewMeHandler(DB, validator.New(), *medi.NewMeDependencies(DB))

	// regist route
	m := g.Group("/me")
	m.GET("", h.Me)
	m.PUT("/user_profile", h.UpdateUserProfile)
	m.PUT("/change_password", h.ChangePassword)
}

// @Security BearerAuth
// @Summary      Get Profile
// @Description  Get your profile detail
// @Tags         Me
// @Accept  	 json
// @Produce  	 json
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/me [get]
func (h *MeHandler) Me(c echo.Context) error {
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// get user id from logged token
	userID, ok := c.Get("user_id").(string)
	if !ok {
		response.Message = "your token is invalid"
		return echo.NewHTTPError(response.Code, response)
	}

	// get data user account
	data, err := h.Dependencies.UC.GetProfile(userID)
	if err != nil {
		response.Message = err.Error()
		return echo.NewHTTPError(response.Code, response)
	}

	// send success response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = data
	return c.JSON(http.StatusOK, response)

}

// @Security BearerAuth
// @Summary      Update Profile
// @Description  Update your profile data
// @Tags         Me
// @Accept       json
// @Produce      json
// @Param        userProfilePayload  body      meschema.UserProfilePayload   true  "User profile payload"
// @Success 	 204  "Profile updated"
// @Failure      400  {object} commonschema.ResponseHTTP "Failure to update"
// @Router       /api/me/user_profile [put]
func (h *MeHandler) UpdateUserProfile(c echo.Context) error {
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	var body meschema.UserProfilePayload
	userID, ok := c.Get("user_id").(string)
	if !ok {
		response.Message = "your token is invalid"
		return c.JSON(response.Code, response)
	}

	if err := c.Bind(&body); err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}

	if err := h.Validator.Struct(body); err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}

	// send into uscease for updating process
	err := h.Dependencies.UC.UpdateProfile(userID, body)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      Change Password
// @Description  Change user password
// @Tags         Me
// @Accept       json
// @Produce      json
// @Param        changePasswordPayload  body      meschema.ChangePasswordPayload   true  "Change password payload"
// @Success 	 204  "Password changed"
// @Failure      400  {object} commonschema.ResponseHTTP "Failure to change password"
// @Router       /api/me/change_password [put]
func (h *MeHandler) ChangePassword(c echo.Context) error {
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	var body meschema.ChangePasswordPayload
	userID, ok := c.Get("user_id").(string)
	if !ok {
		response.Message = "your token is invalid"
		return c.JSON(response.Code, response)
	}

	if err := c.Bind(&body); err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}

	if err := h.Validator.Struct(body); err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}

	// send into uscease for updating process
	err := h.Dependencies.UC.ChangePassword(userID, body)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}
	return c.JSON(http.StatusNoContent, nil)
}
