package nex_match_making_database

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	matchmaking_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
)

// GetDetailedGatheringsByID returns a List of GatheringHolders
func GetDetailedGatheringsByID(manager *common_globals.MatchmakingManager, sourcePID uint64, gatheringIDs []uint32) (types.List[matchmaking_types.GatheringHolder], *nex.Error) {
	gatherings, gatheringTypes, participantArrays, startedTimes, nexError := FindGatheringsByID(manager, gatheringIDs)
	if nexError != nil {
		return types.NewList[matchmaking_types.GatheringHolder](), nexError
	}

	// We're separating each type of Gathering into a map that way we can get it down to 1 DB call per type of
	// Gathering in the ID list. I would prefer to use a map/list combo, but we have to retain the index for
	// later, so a map of maps and converting to a list (where necessary) down the line is more suitable in this
	// scenario.
	gatheringsByType := make(map[string]map[int]matchmaking_types.Gathering)
	for i, gathering := range gatherings {
		gatheringType := gatheringTypes[i]
		if gatheringType != "Gathering" && gatheringType != "MatchmakeSession" && gatheringType != "PersistentGathering" {
			continue
		}

		if gatheringsByType[gatheringType] == nil {
			gatheringsByType[gatheringType] = make(map[int]matchmaking_types.Gathering)
		}

		gatheringsByType[gatheringType][i] = gathering
	}

	gatheringHolders := types.NewList[matchmaking_types.GatheringHolder]()
	for gatheringType, gatheringInfo := range gatheringsByType {
		if gatheringType == "Gathering" {
			for _, gathering := range gatheringInfo {
				gatheringHolder := matchmaking_types.NewGatheringHolder()
				gatheringHolder.Object = gathering.Copy().(matchmaking_types.GatheringInterface)

				gatheringHolders = append(gatheringHolders, gatheringHolder)
			}
		}

		if gatheringType == "MatchmakeSession" {
			matchmakeSessionBases := make([]matchmaking_types.Gathering, 0)
			participantCounts := make([]uint32, 0)
			sessionStartedTimes := make([]types.DateTime, 0)
			for index, gathering := range gatheringInfo {
				matchmakeSessionBases = append(matchmakeSessionBases, gathering)
				participantCounts = append(participantCounts, uint32(len(participantArrays[index])))
				sessionStartedTimes = append(sessionStartedTimes, startedTimes[index])
			}

			matchmakeSessions, nexError := GetMatchmakeSessionsByGathering(manager, manager.Endpoint, matchmakeSessionBases, participantCounts, sessionStartedTimes)
			if nexError != nil {
				return types.NewList[matchmaking_types.GatheringHolder](), nexError
			}

			for _, matchmakeSession := range matchmakeSessions {
				// * Scrap session key and user password
				matchmakeSession.SessionKey = make([]byte, 0)
				matchmakeSession.UserPassword = ""

				gatheringHolder := matchmaking_types.NewGatheringHolder()
				gatheringHolder.Object = matchmakeSession.Copy().(matchmaking_types.GatheringInterface)

				gatheringHolders = append(gatheringHolders, gatheringHolder)
			}
		}

		if gatheringType == "PersistentGathering" {
			persistentGatheringBases := make([]matchmaking_types.Gathering, 0)
			for _, gathering := range gatheringInfo {
				persistentGatheringBases = append(persistentGatheringBases, gathering)
			}

			persistentGatherings, nexError := GetPersistentGatheringsByGathering(manager, persistentGatheringBases, sourcePID)
			if nexError != nil {
				return types.NewList[matchmaking_types.GatheringHolder](), nexError
			}

			for _, persistentGathering := range persistentGatherings {
				// * Scrap persistent gathering password
				persistentGathering.Password = ""

				gatheringHolder := matchmaking_types.NewGatheringHolder()
				gatheringHolder.Object = persistentGathering.Copy().(matchmaking_types.GatheringInterface)

				gatheringHolders = append(gatheringHolders, gatheringHolder)
			}
		}
	}

	return gatheringHolders, nil
}
