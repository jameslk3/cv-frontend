package tests

import (
	. "main/functions"
	"math/rand"
	"testing"
	"time"
)

func TestInitPopIntegration(t *testing.T) {

}

func TestCreateChromosome(t *testing.T) {

	LoadSchedule()

	// Get roster and free agent data
	league_id := 424233486
	espn_s2 := ""
	swid := ""
	team_name := "James's Scary Team"
	year := 2024
	week := "13"
	threshold := 35.0

	roster_map, free_agents := GetPlayers(league_id, espn_s2, swid, team_name, year)

	new_optimal_lineup, streamable_players := OptimizeSlotting(roster_map, week, threshold)
	free_positions := GetUnusedPositions(new_optimal_lineup)

	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	chromosome := CreateChromosome(streamable_players, week, free_agents, free_positions, rng)

	// Check to see if on day 0, the roster has all the streamable players who are playing
	for _, player := range streamable_players {
		if Contains(ScheduleMap[week].Games[player.Team], 0) {
			if MapContainsValue(chromosome.Genes[0], player.Name) == "" {
				t.Error("Initial streamer not in roster")
			}
		}
	}

	// Check to see if every player in NewPlayers is in the roster
	for _, gene := range chromosome.Genes {
		for _, player := range gene.NewPlayers {
			if MapContainsValue(gene, player.Name) == "" {
				t.Error("NewPlayer not in roster")
			}
		}
	}
	
}