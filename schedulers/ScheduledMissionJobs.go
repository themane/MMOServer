package schedulers

import (
	"context"
	"github.com/procyon-projects/chrono"
	"github.com/themane/MMOServer/mongoRepository/models"
	"log"
	"time"
)

type ScheduledMissionManager struct {
	userRepository     models.UserRepository
	universeRepository models.UniverseRepository
	taskScheduler      chrono.TaskScheduler
	scheduledTasks     map[string]chrono.ScheduledTask
}

func NewScheduledMissionManager(userRepository models.UserRepository, universeRepository models.UniverseRepository,
) *ScheduledMissionManager {
	return &ScheduledMissionManager{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		taskScheduler:      chrono.NewDefaultTaskScheduler(),
		scheduledTasks:     map[string]chrono.ScheduledTask{},
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
		log.Println("Error in scheduling mission")
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
		log.Println("Error in scheduling mission")
	}
}

func (m *ScheduledMissionManager) attack(attackMission models.AttackMission) {
	log.Println("Attacking: ", attackMission)
}

func (m *ScheduledMissionManager) spy(spyMission models.SpyMission) {
	log.Println("Spying: ", spyMission)
}
