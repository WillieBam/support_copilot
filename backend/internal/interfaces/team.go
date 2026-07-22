package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/google/uuid"
)

type ITeamRepository interface {
	CreateTeamWithOwner(ctx context.Context, team *models.Team, ownerID uuid.UUID) error
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error)
	GetUserWithTeamsByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	AddTeamMember(ctx context.Context, member *models.TeamMember) error
	RemoveTeamMember(ctx context.Context, teamID, userID uuid.UUID) error
	DeleteTeam(ctx context.Context, teamID uuid.UUID) error
	GetMemberRole(ctx context.Context, teamID, userID uuid.UUID) (string, error)
	ListTeamMembers(ctx context.Context, teamID uuid.UUID) ([]models.TeamMember, error)
	AssignTeamIncident(ctx context.Context, incident *models.TeamIncident) error
	ListTeamIncidents(ctx context.Context, teamID uuid.UUID) ([]models.TeamIncident, error)
}

type ITeamService interface {
	CreateTeam(ctx context.Context, teamName string, creatorID uuid.UUID) (*models.Team, error)
	GetTeam(ctx context.Context, teamID uuid.UUID) (*models.Team, error)
	GetUserTeams(ctx context.Context, userID uuid.UUID) (*models.User, error)
	AddMember(ctx context.Context, requesterID, teamID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, requesterID, teamID, userID uuid.UUID) error
	DeleteTeam(ctx context.Context, userScope string, teamID uuid.UUID) error
	AssignIncident(ctx context.Context, requesterID, teamID, incidentID uuid.UUID, title, status, details string) (*models.TeamIncident, error)
	ListIncidents(ctx context.Context, requesterID, teamID uuid.UUID) ([]models.TeamIncident, error)
}
