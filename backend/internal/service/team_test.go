package service_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/internal/mocks"
	"github.com/WillieBam/support_copilot/backend/internal/service"
	"github.com/WillieBam/support_copilot/backend/types/models"
)

var _ = Describe("TeamService", func() {
	var (
		teamSvc  interfaces.ITeamService
		teamRepo *mocks.ITeamRepository
		ctx      context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		teamRepo = &mocks.ITeamRepository{}
		teamSvc = service.NewTeamService(teamRepo)
	})

	AfterEach(func() {
		teamRepo.AssertExpectations(GinkgoT())
	})

	Context("CreateTeam", func() {
		It("should fail if team name is empty", func() {
			team, err := teamSvc.CreateTeam(ctx, "   ", uuid.New())
			Expect(err).To(Equal(service.ErrTeamNameRequired))
			Expect(team).To(BeNil())
		})

		It("should fail if team name is longer than 20 characters", func() {
			longName := "ThisTeamNameIsWayTooLongForConstraint"
			team, err := teamSvc.CreateTeam(ctx, longName, uuid.New())
			Expect(err).To(Equal(service.ErrTeamNameTooLong))
			Expect(team).To(BeNil())
		})

		It("should succeed when team name is valid", func() {
			creatorID := uuid.New()
			teamRepo.On("CreateTeamWithOwner", ctx, mock.AnythingOfType("*models.Team"), creatorID).Return(nil)

			team, err := teamSvc.CreateTeam(ctx, "DevOps", creatorID)
			Expect(err).NotTo(HaveOccurred())
			Expect(team).NotTo(BeNil())
			Expect(team.TeamName).To(Equal("DevOps"))
		})

		It("should return error if repo fails to create team", func() {
			creatorID := uuid.New()
			teamRepo.On("CreateTeamWithOwner", ctx, mock.AnythingOfType("*models.Team"), creatorID).Return(errors.New("db error"))

			team, err := teamSvc.CreateTeam(ctx, "DevOps", creatorID)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("db error"))
			Expect(team).To(BeNil())
		})
	})

	Context("GetTeam", func() {
		It("should fetch team by ID", func() {
			teamID := uuid.New()
			expectedTeam := &models.Team{ID: teamID, TeamName: "SRE Core"}
			teamRepo.On("GetTeamByID", ctx, teamID).Return(expectedTeam, nil)

			team, err := teamSvc.GetTeam(ctx, teamID)
			Expect(err).NotTo(HaveOccurred())
			Expect(team).To(Equal(expectedTeam))
		})
	})

	Context("GetUserTeams", func() {
		It("should fetch user with teams", func() {
			userID := uuid.New()
			expectedUser := &models.User{ID: userID, Email: "user@test.com"}
			teamRepo.On("GetUserWithTeamsByID", ctx, userID).Return(expectedUser, nil)

			user, err := teamSvc.GetUserTeams(ctx, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).To(Equal(expectedUser))
		})
	})

	Context("AddMember", func() {
		var (
			teamID      uuid.UUID
			requesterID uuid.UUID
			targetID    uuid.UUID
		)

		BeforeEach(func() {
			teamID = uuid.New()
			requesterID = uuid.New()
			targetID = uuid.New()
		})

		It("should fail if requester is not team owner", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, requesterID).Return("member", nil)

			err := teamSvc.AddMember(ctx, requesterID, teamID, targetID)
			Expect(err).To(Equal(service.ErrUnauthorizedTeamOp))
		})

		It("should succeed and assign member role when requester is team owner", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, requesterID).Return("owner", nil)
			teamRepo.On("AddTeamMember", ctx, mock.MatchedBy(func(m *models.TeamMember) bool {
				return m.TeamID == teamID && m.UserID == targetID && m.Role == "member"
			})).Return(nil)

			err := teamSvc.AddMember(ctx, requesterID, teamID, targetID)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("RemoveMember", func() {
		var (
			teamID      uuid.UUID
			ownerID     uuid.UUID
			memberID    uuid.UUID
			nonMemberID uuid.UUID
		)

		BeforeEach(func() {
			teamID = uuid.New()
			ownerID = uuid.New()
			memberID = uuid.New()
			nonMemberID = uuid.New()
		})

		It("should fail if requester is not owner", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, memberID).Return("member", nil)

			err := teamSvc.RemoveMember(ctx, memberID, teamID, uuid.New())
			Expect(err).To(Equal(service.ErrUnauthorizedTeamOp))
		})

		It("should fail if target user is not in team", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, ownerID).Return("owner", nil)
			teamRepo.On("GetMemberRole", ctx, teamID, nonMemberID).Return("", gorm.ErrRecordNotFound)

			err := teamSvc.RemoveMember(ctx, ownerID, teamID, nonMemberID)
			Expect(err).To(Equal(service.ErrUserNotInTeam))
		})

		It("should succeed when owner removes a member", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, ownerID).Return("owner", nil)
			teamRepo.On("GetMemberRole", ctx, teamID, memberID).Return("member", nil)
			teamRepo.On("RemoveTeamMember", ctx, teamID, memberID).Return(nil)

			err := teamSvc.RemoveMember(ctx, ownerID, teamID, memberID)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("DeleteTeam", func() {
		var teamID uuid.UUID

		BeforeEach(func() {
			teamID = uuid.New()
		})

		It("should fail if user scope is not super_admin", func() {
			err := teamSvc.DeleteTeam(ctx, "engineer", teamID)
			Expect(err).To(Equal(service.ErrSuperAdminRequired))
		})

		It("should succeed if user scope is super_admin", func() {
			teamRepo.On("DeleteTeam", ctx, teamID).Return(nil)

			err := teamSvc.DeleteTeam(ctx, "super_admin", teamID)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("AssignIncident & ListIncidents", func() {
		var (
			teamID     uuid.UUID
			userID     uuid.UUID
			incidentID uuid.UUID
		)

		BeforeEach(func() {
			teamID = uuid.New()
			userID = uuid.New()
			incidentID = uuid.New()
		})

		It("should fail AssignIncident if user is not in team", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, userID).Return("", errors.New("not member"))

			inc, err := teamSvc.AssignIncident(ctx, userID, teamID, incidentID, "High Latency", "OPEN", "Details")
			Expect(err).To(Equal(service.ErrUnauthorizedTeamOp))
			Expect(inc).To(BeNil())
		})

		It("should succeed AssignIncident when user is in team", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, userID).Return("member", nil)
			teamRepo.On("AssignTeamIncident", ctx, mock.MatchedBy(func(inc *models.TeamIncident) bool {
				return inc.TeamID == teamID && inc.IncidentID == incidentID && inc.Title == "High Latency"
			})).Return(nil)

			inc, err := teamSvc.AssignIncident(ctx, userID, teamID, incidentID, "High Latency", "OPEN", "Details")
			Expect(err).NotTo(HaveOccurred())
			Expect(inc).NotTo(BeNil())
			Expect(inc.Title).To(Equal("High Latency"))
		})

		It("should ListIncidents when user is in team", func() {
			teamRepo.On("GetMemberRole", ctx, teamID, userID).Return("member", nil)
			expectedIncidents := []models.TeamIncident{
				{ID: uuid.New(), TeamID: teamID, Title: "Incident 1"},
			}
			teamRepo.On("ListTeamIncidents", ctx, teamID).Return(expectedIncidents, nil)

			incidents, err := teamSvc.ListIncidents(ctx, userID, teamID)
			Expect(err).NotTo(HaveOccurred())
			Expect(incidents).To(Equal(expectedIncidents))
		})
	})
})
