package nex_match_making_database

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	pqextended "github.com/PretendoNetwork/pq-extended"
)

// GetMatchmakeSessionsByGathering gets matchmake sessions with the given gathering data
func GetMatchmakeSessionsByGathering(manager *common_globals.MatchmakingManager, endpoint *nex.PRUDPEndPoint, gatherings []match_making_types.Gathering, participantCounts []uint32, startedTimes []types.DateTime) ([]match_making_types.MatchmakeSession, *nex.Error) {
	var gatheringIDs []uint32
	for _, gathering := range gatherings {
		gatheringIDs = append(gatheringIDs, uint32(gathering.ID))
	}

	rows, err := manager.Database.Query(`SELECT
		game_mode,
		attribs,
		open_participation,
		matchmake_system_type,
		application_buffer,
		progress_score,
		session_key,
		option_zero,
		matchmake_param,
		user_password,
		refer_gid,
		user_password_enabled,
		system_password_enabled,
		codeword
		FROM matchmaking.matchmake_sessions WHERE id=ANY($1)`,
		pqextended.Array(gatheringIDs),
	)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
	}

	matchmakeSessions := make([]match_making_types.MatchmakeSession, 0)
	for rows.Next() {
		resultMatchmakeSession := match_making_types.NewMatchmakeSession()
		var resultMatchmakeParam []byte
		var resultAttribs []uint32

		err = rows.Scan(
			&resultMatchmakeSession.GameMode,
			pqextended.Array(&resultAttribs),
			&resultMatchmakeSession.OpenParticipation,
			&resultMatchmakeSession.MatchmakeSystemType,
			&resultMatchmakeSession.ApplicationBuffer,
			&resultMatchmakeSession.ProgressScore,
			&resultMatchmakeSession.SessionKey,
			&resultMatchmakeSession.Option0,
			&resultMatchmakeParam,
			&resultMatchmakeSession.UserPassword,
			&resultMatchmakeSession.ReferGID,
			&resultMatchmakeSession.UserPasswordEnabled,
			&resultMatchmakeSession.SystemPasswordEnabled,
			&resultMatchmakeSession.CodeWord,
		)
		if err != nil {
			common_globals.Logger.Critical(err.Error())
			continue
		}

		attributesSlice := make([]types.UInt32, len(resultAttribs))
		for i, value := range resultAttribs {
			attributesSlice[i] = types.NewUInt32(value)
		}
		resultMatchmakeSession.Attributes = attributesSlice

		matchmakeParamBytes := nex.NewByteStreamIn(resultMatchmakeParam, endpoint.LibraryVersions(), endpoint.ByteStreamSettings())
		resultMatchmakeSession.MatchmakeParam.ExtractFrom(matchmakeParamBytes)

		matchmakeSessions = append(matchmakeSessions, resultMatchmakeSession)
	}

	if len(matchmakeSessions) != len(gatherings) {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, "change_error")
	}

	for i := range matchmakeSessions {
		matchmakeSessions[i].Gathering = gatherings[i]
		matchmakeSessions[i].ParticipationCount = types.NewUInt32(participantCounts[i])
		matchmakeSessions[i].StartedTime = startedTimes[i]
	}

	return matchmakeSessions, nil
}
