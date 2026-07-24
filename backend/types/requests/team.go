package requests

import "github.com/google/uuid"

type CreateTeamRequest struct {
	TeamName string `json:"team_name" binding:"required"`
}

type AddTeamMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type AssignTeamIncidentRequest struct {
	IncidentID uuid.UUID `json:"incident_id" binding:"required"`
	Title      string    `json:"title" binding:"required"`
	Status     string    `json:"status"`
	Details    string    `json:"details"`
}
