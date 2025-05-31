package helpers

import (
	"errors"
	masterrepo "kiraform/src/applications/repos/masters"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CheckAllowedWorkspace(c echo.Context, workspaceID string, DB *gorm.DB) error {
	workspaceRepo := masterrepo.NewWorkspaceRepository(DB)
	notAllowedMessage := "you are not allowed to access this data"

	// get user login
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return errors.New(notAllowedMessage)
	}
	roleName, ok := c.Get("role_name").(string)
	if !ok {
		return errors.New(notAllowedMessage)
	}

	// check valid workspace
	// if user is admin, then allow to access it
	if strings.ToLower(roleName) != "admin" {
		data, err := workspaceRepo.FindWorkspaceUserByUserApproved(workspaceID, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(notAllowedMessage)
			}
			return err
		} else if data == nil {
			return errors.New(notAllowedMessage)
		}
	}

	// set as allowed
	return nil
}
