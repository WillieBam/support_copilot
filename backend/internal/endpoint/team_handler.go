package endpoint

import (
	"errors"
	"net/http"

	"github.com/WillieBam/support_copilot/backend/internal/service"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/WillieBam/support_copilot/backend/types/requests"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func (h *Handler) getAuthenticatedUser(c *echo.Context) (*models.User, error) {
	uidVal := c.Get("user_uid")
	firebaseUID, ok := uidVal.(string)
	if !ok || firebaseUID == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized session")
	}

	if h.userRepo == nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "user repository unavailable")
	}

	user, err := h.userRepo.GetUserByFirebaseUID(c.Request().Context(), firebaseUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(http.StatusNotFound, "user record not found")
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to retrieve user context")
	}
	return user, nil
}

// CreateTeam handles POST /api/teams
func (h *Handler) CreateTeam(c *echo.Context) error {
	user, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	var req requests.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if h.teamService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "team service unavailable"})
	}

	team, err := h.teamService.CreateTeam(c.Request().Context(), req.TeamName, user.ID)
	if err != nil {
		if errors.Is(err, service.ErrTeamNameRequired) || errors.Is(err, service.ErrTeamNameTooLong) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusConflict, map[string]string{"error": "failed to create team: team name may already exist"})
	}

	return c.JSON(http.StatusCreated, team)
}

// GetTeams handles GET /api/teams/me
func (h *Handler) GetTeams(c *echo.Context) error {
	user, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if h.teamService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "team service unavailable"})
	}

	userWithTeams, err := h.teamService.GetUserTeams(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, userWithTeams)
}

// AddTeamMember handles POST /api/teams/:team_id/members
func (h *Handler) AddTeamMember(c *echo.Context) error {
	user, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	teamIDStr := c.Param("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid team ID"})
	}

	var req requests.AddTeamMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if req.UserID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	if h.teamService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "team service unavailable"})
	}

	err = h.teamService.AddMember(c.Request().Context(), user.ID, teamID, req.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorizedTeamOp) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusConflict, map[string]string{"error": "failed to add member: member may already exist in team"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "member added successfully"})
}

// RemoveTeamMember handles DELETE /api/teams/:team_id/members/:user_id
func (h *Handler) RemoveTeamMember(c *echo.Context) error {
	user, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	teamID, err := uuid.Parse(c.Param("team_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid team ID"})
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	if h.teamService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "team service unavailable"})
	}

	err = h.teamService.RemoveMember(c.Request().Context(), user.ID, teamID, targetUserID)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorizedTeamOp) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, service.ErrUserNotInTeam) || errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "member not found in team"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "member removed successfully"})
}

// DeleteTeam handles DELETE /api/teams/:team_id (super_admin only)
func (h *Handler) DeleteTeam(c *echo.Context) error {
	user, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	teamID, err := uuid.Parse(c.Param("team_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid team ID"})
	}

	if h.teamService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "team service unavailable"})
	}

	err = h.teamService.DeleteTeam(c.Request().Context(), user.Scope, teamID)
	if err != nil {
		if errors.Is(err, service.ErrSuperAdminRequired) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "team not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "team deleted successfully"})
}

// AssignTeamIncident handles POST /api/teams/:team_id/incidents
func (h *Handler) AssignTeamIncident(c *echo.Context) error {
	user, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	teamID, err := uuid.Parse(c.Param("team_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid team ID"})
	}

	var req requests.AssignTeamIncidentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if h.teamService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "team service unavailable"})
	}

	inc, err := h.teamService.AssignIncident(c.Request().Context(), user.ID, teamID, req.IncidentID, req.Title, req.Status, req.Details)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorizedTeamOp) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, inc)
}

// GetTeamIncidents handles GET /api/teams/:team_id/incidents
func (h *Handler) GetTeamIncidents(c *echo.Context) error {
	user, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	teamID, err := uuid.Parse(c.Param("team_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid team ID"})
	}

	if h.teamService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "team service unavailable"})
	}

	incidents, err := h.teamService.ListIncidents(c.Request().Context(), user.ID, teamID)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorizedTeamOp) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, incidents)
}
