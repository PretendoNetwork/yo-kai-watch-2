// i very loosely stole this code from https://github.com/PretendoNetwork/nex-protocols-common-go/blob/main/match-making/database/find_gathering_by_id.go - Trace
package nex_match_making_database

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	matchmaking_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	pqextended "github.com/PretendoNetwork/pq-extended"
)

// FindMatchmakeSessionsByOwner finds MatchmakeSessions on a database using the given owner's pid and a ResultRange. Returns a list of GatheringHolders
func FindMatchmakeSessionsByOwner(manager *common_globals.MatchmakingManager, connection *nex.PRUDPConnection, id types.PID, resultRange types.ResultRange) (types.List[matchmaking_types.GatheringHolder], *nex.Error) {
	if resultRange.Length == 0 {
		// TODO: change default max length, this is likely highly inaccurate
		resultRange.Length = 50
	}

	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	rows, err := manager.Database.Query(`SELECT
		g.id,
		g.owner_pid,
		g.host_pid,
		g.min_participants,
		g.max_participants,
		g.participation_policy,
		g.policy_argument,
		g.flags,
		g.state,
		g.description,
		array_length(g.participants, 1),
		g.started_time,
		ms.game_mode,
		ms.attribs,
		ms.open_participation,
		ms.matchmake_system_type,
		ms.application_buffer,
		ms.progress_score,
		ms.session_key,
		ms.option_zero,
		ms.matchmake_param,
		ms.user_password,
		ms.refer_gid,
		ms.user_password_enabled,
		ms.system_password_enabled,
		ms.codeword
		FROM matchmaking.gatherings AS g
		INNER JOIN matchmaking.matchmake_sessions AS ms ON ms.id = g.id
		WHERE
		g.registered=true AND
		g.type='MatchmakeSession' AND
		g.host_pid <> 0 AND
		g.owner_pid=$1 AND
		ms.open_participation=true AND
		array_length(g.participants, 1) < g.max_participants AND
		ms.user_password_enabled=false AND
		ms.system_password_enabled=false
		LIMIT $2 OFFSET $3`,
		id,
		resultRange.Length,
		resultRange.Offset,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.NewList[matchmaking_types.GatheringHolder](), nex.NewError(nex.ResultCodes.RendezVous.SessionVoid, err.Error())
		} else {
			return types.NewList[matchmaking_types.GatheringHolder](), nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
		}
	}

	var gatheringHolders types.List[matchmaking_types.GatheringHolder]

	for rows.Next() {
		resultMatchmakeSession := matchmaking_types.NewMatchmakeSession()
		var startedTime time.Time
		var resultAttribs []uint32
		var resultMatchmakeParam []byte

		err = rows.Scan(
			&resultMatchmakeSession.Gathering.ID,
			&resultMatchmakeSession.Gathering.OwnerPID,
			&resultMatchmakeSession.Gathering.HostPID,
			&resultMatchmakeSession.Gathering.MinimumParticipants,
			&resultMatchmakeSession.Gathering.MaximumParticipants,
			&resultMatchmakeSession.Gathering.ParticipationPolicy,
			&resultMatchmakeSession.Gathering.PolicyArgument,
			&resultMatchmakeSession.Gathering.Flags,
			&resultMatchmakeSession.Gathering.State,
			&resultMatchmakeSession.Gathering.Description,
			&resultMatchmakeSession.ParticipationCount,
			&startedTime,
			&resultMatchmakeSession.GameMode,
			pqextended.Array(&resultAttribs),
			&resultMatchmakeSession.OpenParticipation,
			&resultMatchmakeSession.MatchmakeSystemType,
			&resultMatchmakeSession.ApplicationBuffer,
			&resultMatchmakeSession.ProgressScore,
			&resultMatchmakeSession.SessionKey,
			&resultMatchmakeSession.Option,
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

		resultMatchmakeSession.StartedTime = resultMatchmakeSession.StartedTime.FromTimestamp(startedTime)

		attributesSlice := make([]types.UInt32, len(resultAttribs))
		for i, value := range resultAttribs {
			attributesSlice[i] = types.NewUInt32(value)
		}
		resultMatchmakeSession.Attributes = attributesSlice

		matchmakeParamBytes := nex.NewByteStreamIn(resultMatchmakeParam, endpoint.LibraryVersions(), endpoint.ByteStreamSettings())
		resultMatchmakeSession.MatchmakeParam.ExtractFrom(matchmakeParamBytes)

		resultGatheringHolder := matchmaking_types.NewGatheringHolder()
		resultGatheringHolder.Object = resultMatchmakeSession.Copy().(matchmaking_types.MatchmakeSession)

		gatheringHolders = append(gatheringHolders, resultGatheringHolder)
	}

	rows.Close()

	return gatheringHolders, nil
}
