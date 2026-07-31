package nex_match_making_database

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	pqextended "github.com/PretendoNetwork/pq-extended"
)

// GetPersistentGatheringsByGathering gets a persistent gathering with the given gathering data
func GetPersistentGatheringsByGathering(manager *common_globals.MatchmakingManager, gatherings []match_making_types.Gathering, sourcePID uint64) ([]match_making_types.PersistentGathering, *nex.Error) {
	var gatheringIDs []uint32
	for _, gathering := range gatherings {
		gatheringIDs = append(gatheringIDs, uint32(gathering.ID))
	}

	rows, err := manager.Database.Query(`SELECT
		community_type,
		password,
		attribs,
		application_buffer,
		participation_start_date,
		participation_end_date,
		(SELECT COUNT(ms.id)
			FROM matchmaking.matchmake_sessions AS ms
			INNER JOIN matchmaking.gatherings AS gms ON ms.id = gms.id
			WHERE gms.registered=true
			AND ms.matchmake_system_type=5 -- matchmake_system_type=5 is only used in matchmake sessions attached to a persistent gathering
			AND ms.attribs[1]=g.id) AS matchmake_session_count,
		COALESCE((SELECT cp.participation_count
			FROM matchmaking.community_participations AS cp
			WHERE cp.user_pid=$2
			AND cp.gathering_id=g.id), 0) AS participation_count
		FROM matchmaking.persistent_gatherings
		WHERE id=ANY($1)`,
		gatheringIDs,
		sourcePID,
	)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
	}

	persistentGatherings := make([]match_making_types.PersistentGathering, 0)
	for rows.Next() {
		resultPersistentGathering := match_making_types.NewPersistentGathering()
		var resultAttribs []uint32

		err = rows.Scan(
			&resultPersistentGathering.CommunityType,
			&resultPersistentGathering.Password,
			pqextended.Array(&resultAttribs),
			&resultPersistentGathering.ApplicationBuffer,
			&resultPersistentGathering.ParticipationStartDate,
			&resultPersistentGathering.ParticipationEndDate,
			&resultPersistentGathering.MatchmakeSessionCount,
			&resultPersistentGathering.ParticipationCount,
		)
		if err != nil {
			common_globals.Logger.Critical(err.Error())
			continue
		}

		attributesSlice := make([]types.UInt32, len(resultAttribs))
		for i, value := range resultAttribs {
			attributesSlice[i] = types.NewUInt32(value)
		}
		resultPersistentGathering.Attribs = attributesSlice

		persistentGatherings = append(persistentGatherings, resultPersistentGathering)
	}

	if len(persistentGatherings) != len(gatherings) {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, "change_error")
	}

	for i := range persistentGatherings {
		persistentGatherings[i].Gathering = gatherings[i]
	}

	return persistentGatherings, nil
}
