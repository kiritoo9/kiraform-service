package helpers

import (
	"errors"
	masterrepo "kiraform/src/applications/repos/masters"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func baseValidation(c echo.Context) (string, string, error) {
	// get user login
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return "", "", errors.New("your identity is not recognized, please contact our admin")
	}
	roleName, ok := c.Get("role_name").(string)
	if !ok {
		return "", "", errors.New("you have no role registered")
	}

	return userID, roleName, nil
}

func CheckAllowedWorkspace(c echo.Context, workspaceID string, DB *gorm.DB) error {
	workspaceRepo := masterrepo.NewWorkspaceRepository(DB)
	notAllowedMessage := "you are not allowed to access this data"

	userID, roleName, err := baseValidation(c)
	if err != nil {
		return err
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

func CheckAllowedCampaign(c echo.Context, workspaceID string, campaignID string, DB *gorm.DB) error {
	campaignRepo := masterrepo.NewCampaignRepository(DB)
	notAllowedMessage := "you are not allowed to access this data"

	userID, roleName, err := baseValidation(c)
	if err != nil {
		return err
	}

	// check allowed campaign based on workspace and campaign
	// if user is admin, then allow to access it
	if strings.ToLower(roleName) != "admin" {
		data, err := campaignRepo.CheckAllowedUserForCampaign(workspaceID, campaignID, userID)
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
