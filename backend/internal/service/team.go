package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/google/uuid"
)

var (
	ErrTeamNameRequired      = errors.New("team name is required")
	ErrTeamNameTooLong       = errors.New("team name must be 20 characters or less")
	ErrUnauthorizedTeamOp    = errors.New("unauthorized team operation: owner permission required")
	ErrSuperAdminRequired    = errors.New("unauthorized operation: super_admin scope required to delete a team")
	ErrUserNotInTeam         = errors.New("user is not a member of this team")
)

type teamService struct {
	teamRepo interfaces.ITeamRepository
}

func NewTeamService(teamRepo interfaces.ITeamRepository) interfaces.ITeamService {
	return &teamService{teamRepo: teamRepo}
}

func (s *teamService) CreateTeam(ctx context.Context, teamName string, creatorID uuid.UUID) (*models.Team, error) {
	name := strings.TrimSpace(teamName)
	if name == "" {
		return nil, ErrTeamNameRequired
	}
	if len(name) > 20 {
		return nil, ErrTeamNameTooLong
	}

	team := &models.Team{
		ID:        uuid.New(),
		TeamName:  name,
		CreatedAt: time.Now(),
	}

	err := s.teamRepo.CreateTeamWithOwner(ctx, team, creatorID)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamID uuid.UUID) (*models.Team, error) {
	return s.teamRepo.GetTeamByID(ctx, teamID)
}

func (s *teamService) GetUserTeams(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.teamRepo.GetUserWithTeamsByID(ctx, userID)
}

func (s *teamService) AddMember(ctx context.Context, requesterID, teamID, userID uuid.UUID) error {
	reqRole, err := s.teamRepo.GetMemberRole(ctx, teamID, requesterID)
	if err != nil || reqRole != "owner" {
		return ErrUnauthorizedTeamOp
	}

	member := &models.TeamMember{
		ID:     uuid.New(),
		TeamID: teamID,
		UserID: userID,
		Role:   "member",
	}
	return s.teamRepo.AddTeamMember(ctx, member)
}

func (s *teamService) RemoveMember(ctx context.Context, requesterID, teamID, userID uuid.UUID) error {
	reqRole, err := s.teamRepo.GetMemberRole(ctx, teamID, requesterID)
	if err != nil || reqRole != "owner" {
		return ErrUnauthorizedTeamOp
	}

	_, err = s.teamRepo.GetMemberRole(ctx, teamID, userID)
	if err != nil {
		return ErrUserNotInTeam
	}

	return s.teamRepo.RemoveTeamMember(ctx, teamID, userID)
}

func (s *teamService) DeleteTeam(ctx context.Context, userScope string, teamID uuid.UUID) error {
	if userScope != "super_admin" {
		return ErrSuperAdminRequired
	}
	return s.teamRepo.DeleteTeam(ctx, teamID)
}

func (s *teamService) AssignIncident(ctx context.Context, requesterID, teamID, incidentID uuid.UUID, title, status, details string) (*models.TeamIncident, error) {
	_, err := s.teamRepo.GetMemberRole(ctx, teamID, requesterID)
	if err != nil {
		return nil, ErrUnauthorizedTeamOp
	}

	inc := &models.TeamIncident{
		ID:         uuid.New(),
		IncidentID: incidentID,
		TeamID:     teamID,
		AssignedBy: requesterID,
		Title:      title,
		Status:     status,
		Details:    details,
		AssignedAt: time.Now(),
	}

	err = s.teamRepo.AssignTeamIncident(ctx, inc)
	if err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *teamService) ListIncidents(ctx context.Context, requesterID, teamID uuid.UUID) ([]models.TeamIncident, error) {
	_, err := s.teamRepo.GetMemberRole(ctx, teamID, requesterID)
	if err != nil {
		return nil, ErrUnauthorizedTeamOp
	}
	return s.teamRepo.ListTeamIncidents(ctx, teamID)
}
