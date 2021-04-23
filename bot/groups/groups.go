package groups

import (
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-discord/app"
)

func VerifyGroupExistenceAndCreateIfNot(groupId string) (*linker.Group, error) {
	group, err := linker.GetGroupByID(groupId)
	if err == nil {
		return group, nil
	}

	group, err = linker.CreateGroup(groupId, app.Config.ServiceName)
	if err != nil {
		return nil, err
	}

	return group, nil
}
