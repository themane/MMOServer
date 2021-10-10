package schedulers

import (
	"context"
	"github.com/procyon-projects/chrono"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"time"
)

type ScheduledMissionManager struct {
	userRepository     models.UserRepository
	universeRepository models.UniverseRepository
	taskScheduler      chrono.TaskScheduler
	scheduledTasks     map[string]chrono.ScheduledTask
	logger             *constants.LoggingUtils
}

func NewScheduledMissionManager(userRepository models.UserRepository, universeRepository models.UniverseRepository,
	logLevel string,
) *ScheduledMissionManager {
	return &ScheduledMissionManager{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		taskScheduler:      chrono.NewDefaultTaskScheduler(),
		scheduledTasks:     map[string]chrono.ScheduledTask{},
		logger:             constants.NewLoggingUtils("SCHEDULER_MISSION_JOBS", logLevel),
	}
}

func (m *ScheduledMissionManager) ScheduleAttackMission(attackMission models.AttackMission, missionTime time.Time, returnTime time.Time) {
	task := func(ctx context.Context) {
		m.attack(attackMission)
	}
	scheduledTask, err := m.taskScheduler.Schedule(task,
		chrono.WithStartTime(missionTime.Year(), missionTime.Month(), missionTime.Day(),
			missionTime.Hour(), missionTime.Minute(), missionTime.Second()))
	m.scheduledTasks[attackMission.Id] = scheduledTask
	if err != nil {
		m.logger.Error("error in scheduled task for attack mission", err)
	}
}

func (m *ScheduledMissionManager) ScheduleSpyMission(spyMission models.SpyMission, missionTime time.Time, returnTime time.Time) {
	task := func(ctx context.Context) {
		m.spy(spyMission)
	}
	scheduledTask, err := m.taskScheduler.Schedule(task,
		chrono.WithStartTime(missionTime.Year(), missionTime.Month(), missionTime.Day(),
			missionTime.Hour(), missionTime.Minute(), missionTime.Second()))
	m.scheduledTasks[spyMission.Id] = scheduledTask
	if err != nil {
		m.logger.Error("error in scheduled task for spy mission", err)
	}
}

func (m *ScheduledMissionManager) attack(attackMission models.AttackMission) {
	m.logger.Println("Attacking", attackMission)
}

func (m *ScheduledMissionManager) spy(spyMission models.SpyMission) {
	m.logger.Println("Spying", spyMission)
}
