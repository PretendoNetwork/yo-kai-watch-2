// i very loosely stole this code from https://github.com/PretendoNetwork/nex-protocols-common-go/blob/main/match-making/database/find_gathering_by_id.go - Trace
package nex_match_making_database

import (
	"database/sql"
	"math"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
)

// FindGatheringIDsByOwner finds gathering IDs on a database using the given owner's PID and a ResultRange. Returns a list of uint32
func FindGatheringIDsByOwner(manager *common_globals.MatchmakingManager, connection *nex.PRUDPConnection, id types.PID, resultRange types.ResultRange) ([]uint32, *nex.Error) {
	if resultRange.Length == 0 {
		// TODO: change default max length, this is likely inaccurate
		resultRange.Length = 50
	}

	if resultRange.Offset == types.NewUInt32(math.MaxUint32) {
		resultRange.Offset = 0
	}

	rows, err := manager.Database.Query(`SELECT
		id FROM matchmaking.gatherings WHERE
		registered=true AND
		host_pid <> 0 AND
		owner_pid=$1 AND
		array_length(participants, 1) < max_participants
		LIMIT $2 OFFSET $3`,
		id,
		resultRange.Length,
		resultRange.Offset,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nex.NewError(nex.ResultCodes.RendezVous.SessionVoid, err.Error())
		} else {
			return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
		}
	}

	var gatheringIDs []uint32

	for rows.Next() {
		var gatheringID uint32

		err = rows.Scan(
			&gatheringID,
		)

		if err != nil {
			common_globals.Logger.Critical(err.Error())
			continue
		}

		gatheringIDs = append(gatheringIDs, gatheringID)
	}

	rows.Close()

	return gatheringIDs, nil
}
