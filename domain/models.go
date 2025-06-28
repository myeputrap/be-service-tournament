package domain

import (
	"time"
)

// --- User Related ---
type User struct {
	ID             int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Email          string     `json:"email" gorm:"column:email;unique;not null"`
	Username       string     `json:"username" gorm:"column:username;unique;not null"`
	HashedPassword string     `json:"hashed_password" gorm:"column:hashed_password;not null"`
	PhoneNumber    string     `json:"phone_number" gorm:"column:phone_number"`
	FullName       string     `json:"full_name" gorm:"column:full_name"`
	Gender         string     `json:"gender" gorm:"column:gender"`
	PhotoProfile   *string    `json:"photo_profile" gorm:"column:photo_profile"`
	BirthDate      *time.Time `json:"birth_date" gorm:"column:birth_date"`
	City           *string    `json:"city" gorm:"column:city"`
	Nik            *string    `json:"nik" gorm:"column:nik"`
	Tier           *Tier      `json:"tier" gorm:"foreignKey:TierID;references:id"`
	TierID         *int64     `json:"tier_id" gorm:"column:tier_id"`
	Role           *Role      `json:"role" gorm:"foreignKey:RoleID;references:id"`
	RoleID         int64      `json:"role_id" gorm:"column:role_id"`
	ReferalCode    *string    `json:"referal_code" gorm:"column:referal_code"`
	CreatedAt      time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	//UpdatedAt      time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string { return "web.users" }

type Tier struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TierName  string    `json:"tier_name" gorm:"column:tier_name;unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Tier) TableName() string { return "web.tier" }

type Role struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RoleName string `json:"role_name" gorm:"column:name;unique;not null"`
}

func (Role) TableName() string { return "web.role" }

type PlayerProgress struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID      int64     `json:"user_id" gorm:"column:user_id"`
	Level       int       `json:"level" gorm:"column:level"`
	CurrentXP   int       `json:"current_experience_level" gorm:"column:current_experience_level"`
	TotalXP     int       `json:"total_experience_level" gorm:"column:total_experience_level"`
	LastLevelUp time.Time `json:"last_levelup" gorm:"column:last_levelup"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (PlayerProgress) TableName() string { return "player_progress" }

type Level struct {
	Level      int `json:"level" gorm:"column:level;primaryKey"`
	XPRequired int `json:"xp_required" gorm:"column:xp_required"`
}

func (Level) TableName() string { return "levels" }

type PointsLog struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"column:user_id"`
	Change    int       `json:"points_change" gorm:"column:points_change"`
	Reason    string    `json:"reason" gorm:"column:reason"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (PointsLog) TableName() string { return "points_log" }

type Tournament struct {
	ID                int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name              string    `json:"name" gorm:"column:name" validate:"required"`
	Description       string    `json:"description" gorm:"column:description"`
	Type              string    `json:"type" gorm:"column:type" validate:"required"`
	Location          string    `json:"location" gorm:"column:location" validate:"required"`
	Fee               float64   `json:"fee" gorm:"column:fee" validate:"required"`
	Latitude          float64   `json:"latitude" gorm:"column:latitude"`
	Longitude         float64   `json:"longitude" gorm:"column:longitude"`
	FeeInfo           string    `json:"fee_info" gorm:"column:fee_info" validate:"required"`
	StartDate         time.Time `json:"start_date" gorm:"column:start_date" validate:"required"`
	EndDate           time.Time `json:"end_date" gorm:"column:end_date" validate:"required"`
	GenderEligibility string    `json:"gender_eligibility" gorm:"column:gender_eligibility" validate:"required"`
	Quota             int       `json:"quota" gorm:"column:quota" validate:"required"`
	Organizer         string    `json:"organizer" gorm:"column:organizer" validate:"required"`
	CreatedBy         int64     `json:"created_by" gorm:"column:created_by;autoCreateTime"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Tournament) TableName() string { return "web.tournaments" }

type Participant struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TournamentID int64     `json:"tournament_id" gorm:"column:tournament_id"`
	UserAID      int64     `json:"player1" gorm:"column:user_a_id"`
	UserA        User      `json:"player1_detail" gorm:"foreignKey:UserAID;references:id"`
	UserBID      int64     `json:"player2" gorm:"column:user_b_id"`
	UserB        User      `json:"player2_detail" gorm:"foreignKey:UserBID;references:id"`
	State        string    `json:"state" gorm:"column:state"`
	PaymentProof *string   `json:"payment_proof" gorm:"column:payment_proof"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Participant) TableName() string { return "web.participants" }

type Phase struct {
	ID            int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TournamentID  int64     `json:"tournament_id" gorm:"column:tournament_id"`
	Name          string    `json:"name" gorm:"column:name"`
	Type          string    `json:"type" gorm:"column:type"`
	GroupCount    int       `json:"group_count" gorm:"column:group_count"`
	PairsPerGroup int       `json:"pairs_per_group" gorm:"column:pairs_per_group"`
	Round         int       `json:"round" gorm:"column:round"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Phase) TableName() string { return "phases" }

type TournamentGroup struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TournamentID int64     `json:"tournament_id" gorm:"column:tournament_id"`
	Name         string    `json:"name" gorm:"column:name"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (TournamentGroup) TableName() string { return "tournament_groups" }

type GroupParticipant struct {
	ID            int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	GroupID       int64     `json:"group_id" gorm:"column:group_id"`
	ParticipantID int64     `json:"participant_id" gorm:"column:participant_id"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (GroupParticipant) TableName() string { return "group_participants" }

type MatchHeader struct {
	ID            int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	PhaseID       int64     `json:"phase_id" gorm:"column:phase_id"`
	GroupID       int64     `json:"group_id" gorm:"column:group_id"`
	MatchNumber   int       `json:"match_number" gorm:"column:match_number"`
	ScheduledTime time.Time `json:"scheduled_time" gorm:"column:scheduled_time"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (MatchHeader) TableName() string { return "matches_header" }

type MatchDetail struct {
	ID            int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MatchID       int64     `json:"match_id" gorm:"column:match_id"`
	ParticipantID int64     `json:"participant_id" gorm:"column:participant_id"`
	Score         int       `json:"score" gorm:"column:score"`
	IsWinner      bool      `json:"is_winner" gorm:"column:is_winner"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (MatchDetail) TableName() string { return "matches_detail" }

type GroupStanding struct {
	ID             int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ParticipantID  int64     `json:"participant_id" gorm:"column:participant_id"`
	GroupID        int64     `json:"group_id" gorm:"column:group_id"`
	MatchesPlayed  int       `json:"matches_played" gorm:"column:matches_played"`
	Wins           int       `json:"wins" gorm:"column:wins"`
	Draws          int       `json:"draws" gorm:"column:draws"`
	Losses         int       `json:"losses" gorm:"column:losses"`
	Points         int       `json:"points" gorm:"column:points"`
	GoalsFor       int       `json:"goals_for" gorm:"column:goals_for"`
	GoalsAgainst   int       `json:"goals_against" gorm:"column:goals_against"`
	GoalDifference int       `json:"goal_difference" gorm:"column:goal_difference"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (GroupStanding) TableName() string { return "group_standings" }

type Banner struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Title     string    `json:"title" gorm:"column:title"`
	Content   string    `json:"content" gorm:"column:content"`
	ImageURL  string    `json:"image_url" gorm:"column:image_url"`
	CtaLabel  string    `json:"cta_label" gorm:"column:cta_label"`
	CtaURL    string    `json:"cta_url" gorm:"column:cta_url"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Banner) TableName() string { return "banners" }

type LeaderboardScore struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64     `json:"user_id" gorm:"column:user_id"`
	TournamentID int64     `json:"tournament_id" gorm:"column:tournament_id"`
	Score        int       `json:"score" gorm:"column:score"`
	Month        int       `json:"month" gorm:"column:month"`
	Year         int       `json:"year" gorm:"column:year"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (LeaderboardScore) TableName() string { return "leaderboard_scores" }
