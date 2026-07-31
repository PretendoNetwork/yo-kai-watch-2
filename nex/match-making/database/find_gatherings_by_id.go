package nex_match_making_database

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	pqextended "github.com/PretendoNetwork/pq-extended"
)

// FindGatheringsByID finds gatherings on a database with the given IDs. Returns the gatherings, their types, the participant lists and the started times
func FindGatheringsByID(manager *common_globals.MatchmakingManager, ids []uint32) ([]match_making_types.Gathering, []string, [][]uint64, []types.DateTime, *nex.Error) {
	rows, err := manager.Database.Query(`SELECT
		id,
		owner_pid, 
		host_pid, 
		min_participants, 
		max_participants, 
		participation_policy, 
		policy_argument, 
		flags, 
		state, 
		description, 
		type, 
		participants, 
		started_time FROM matchmaking.gatherings 
		WHERE id=ANY($1) AND registered=true`,
		pqextended.Array(ids),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil, nex.NewError(nex.ResultCodes.RendezVous.SessionVoid, err.Error())
		} else {
			return nil, nil, nil, nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
		}
	}

	var gatherings []match_making_types.Gathering
	var gatheringTypes []string
	var participantArrays [][]uint64
	var startedTimes []types.DateTime

	for rows.Next() {
		gathering := match_making_types.NewGathering()
		var gatheringType string
		var participants []uint64
		var startedTime types.DateTime

		err := rows.Scan(
			&gathering.ID,
			&gathering.OwnerPID,
			&gathering.HostPID,
			&gathering.MinimumParticipants,
			&gathering.MaximumParticipants,
			&gathering.ParticipationPolicy,
			&gathering.PolicyArgument,
			&gathering.Flags,
			&gathering.State,
			&gathering.Description,
			&gatheringType,
			pqextended.Array(&participants),
			&startedTime,
		)

		if err != nil {
			common_globals.Logger.Critical(err.Error())
			continue
		}

		gatherings = append(gatherings, gathering)
		gatheringTypes = append(gatheringTypes, gatheringType)
		participantArrays = append(participantArrays, participants)
		startedTimes = append(startedTimes, startedTime)
	}

	return gatherings, gatheringTypes, participantArrays, startedTimes, nil
}
