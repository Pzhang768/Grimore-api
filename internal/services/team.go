package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Pzhang768/Grimore-api/internal/models"
)

var validAgentTypes = map[string]struct{}{
	"fetcher":     {},
	"analyser":    {},
	"tailor":      {},
	"coordinator": {},
}

var (
	ErrTeamNameRequired   = errors.New("name is required")
	ErrAgentsRequired     = errors.New("at least one agent is required")
	ErrTooManyAgents      = errors.New("a team can have at most 4 agents")
	ErrInvalidAgentType   = errors.New("invalid agent type")
	ErrDuplicatePosition  = errors.New("agent positions must be unique")
	ErrDuplicateAgentType = errors.New("each agent type may appear at most once per team")
	ErrTeamNotFound       = errors.New("team not found")
	ErrTeamForbidden      = errors.New("team belongs to another user")
)

type CreateAgentInput struct {
	AgentType string
	Position  int
	Context   map[string]any
}

type TeamService struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{db: db}
}

func (s *TeamService) GetTeam(ctx context.Context, userID, teamID uuid.UUID) (*models.Team, []models.TeamAgent, error) {
	var team models.Team
	if err := s.db.WithContext(ctx).First(&team, "id = ?", teamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrTeamNotFound
		}
		return nil, nil, err
	}
	if team.UserID != userID {
		return nil, nil, ErrTeamForbidden
	}
	var agents []models.TeamAgent
	if err := s.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&agents).Error; err != nil {
		return nil, nil, err
	}
	return &team, agents, nil
}

func (s *TeamService) ListTeams(ctx context.Context, userID uuid.UUID) ([]models.Team, error) {
	var teams []models.Team
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *TeamService) CreateTeam(ctx context.Context, userID uuid.UUID, name string, agents []CreateAgentInput) (*models.Team, []models.TeamAgent, error) {
	if name == "" {
		return nil, nil, ErrTeamNameRequired
	}
	if len(agents) == 0 {
		return nil, nil, ErrAgentsRequired
	}
	if len(agents) > 4 {
		return nil, nil, ErrTooManyAgents
	}

	positions := make(map[int]struct{})
	agentTypes := make(map[string]struct{})
	for _, a := range agents {
		if _, ok := validAgentTypes[a.AgentType]; !ok {
			return nil, nil, fmt.Errorf("%w: %q", ErrInvalidAgentType, a.AgentType)
		}
		if _, dup := agentTypes[a.AgentType]; dup {
			return nil, nil, ErrDuplicateAgentType
		}
		agentTypes[a.AgentType] = struct{}{}
		if _, dup := positions[a.Position]; dup {
			return nil, nil, ErrDuplicatePosition
		}
		positions[a.Position] = struct{}{}
	}

	team := &models.Team{
		UserID: userID,
		Name:   name,
	}

	var teamAgents []models.TeamAgent
	for _, a := range agents {
		agentCtx := a.Context
		if agentCtx == nil {
			agentCtx = map[string]any{}
		}
		teamAgents = append(teamAgents, models.TeamAgent{
			AgentType: a.AgentType,
			Position:  a.Position,
			Context:   agentCtx,
		})
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		for i := range teamAgents {
			teamAgents[i].TeamID = team.ID
		}
		return tx.Create(&teamAgents).Error
	})
	if err != nil {
		return nil, nil, err
	}

	return team, teamAgents, nil
}
