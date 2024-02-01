package helper

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

// Function to evolve a population
func EvolvePopulation(size int, population []Chromosome, fas []Player, free_positions map[int][]string, streamable_players []Player, week string) []Chromosome {

	// Fill cumulative probabilities tracker for each chromosome
	AssignCumProbs(population, size)

	next_generation := make([]Chromosome, size)

	// Implement elitism
	next_generation[size-1] = population[size-1]
	
	for i := 0; i < size-1; i++ {
		
		// Get parents
		parent1 := SelectFirstParent(population)
		parent2 := SelectSecondParent(population)

		// Get children
		child, _ := Crossover(parent1, parent2, fas, free_positions, streamable_players, week)

		// Mutate children
		Mutate(0.10, &child, fas, free_positions, streamable_players, week, streamable_players)

		// Add the total acquisitions to the children chromosomes
		GetTotalAcquisitions(&child)

		// Add children to evolved population
		next_generation[i] = child	
	}

	return next_generation
}

// Function to assign cumulative probabilities to a population for a ranked selection
func AssignCumProbs(population []Chromosome, size int) {

	// Function to get the probability of a chromosome being selected
	get_probability := func(x int) float64 {
		// return ((7 * (50 - float64(x))) + 55) / (2000)
		return math.Pow(float64(x) / float64(size), (1.5)) + 0.02
	}

	cum_prob := get_probability(0)
	population[0].CumProbTracker = cum_prob

	// Loop through each chromosome and assign cumulative probabilities
	for i := 1; i < len(population); i++ {
		cum_prob += get_probability(i)
		population[i].CumProbTracker = cum_prob
	}
}

// Function to select the first parent based on a ranked selection
func SelectFirstParent(population []Chromosome) Chromosome {

	// Get random seed
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	
	// Get the random number to select the chromosome
	rand_num := rng.Float64() * population[len(population)-1].CumProbTracker

	// Loop through each chromosome and return the first one that is greater than the random number
	for _, chromosome := range population {
		if chromosome.CumProbTracker > rand_num {
			return chromosome
		}
	}

	// If no chromosome is returned, return the last one
	return population[len(population)-1]
}

// Function to select the second parent based on a randomized tournament selection
func SelectSecondParent(population []Chromosome) Chromosome {

	// Get random seed
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	tournament := make([][5]Chromosome, 3)

	// Create tournament
	for i := 0; i < 3; i++ {
		
		for j := 0; j < 5; j++ {

			// Insert random chromosome
			tournament[i][j] = population[rng.Intn(len(population))]
		}

		// Sort tournament
		sort.Slice(tournament[i][:], func(k, l int) bool {
			return tournament[i][k].FitnessScore > tournament[i][l].FitnessScore
		})
	}

	return tournament[rng.Intn(3)][0]
}

// Function to get the children of two parents
func Crossover(parent1 Chromosome, parent2 Chromosome, fas []Player, free_positions map[int][]string, streamable_players []Player, week string) (Chromosome, []Player) {

	// Initialize child
	child := Chromosome{Genes: make([]Gene, ScheduleMap[week].GameSpan + 1), FitnessScore: 0, TotalAcquisitions: 0, CumProbTracker: 0.0, DroppedPlayers: make(map[string]DroppedPlayer)}
	for j := 0; j <= ScheduleMap[week].GameSpan; j++ {
		child.Genes[j] = Gene{Roster: make(map[string]Player), NewPlayers: make(map[string]Player), Day: j, Acquisitions: 0, DroppedPlayers: []Player{}, Bench: []Player{}}
	}
	
	// Fill genes with initial streamers
	cur_streamers1 := make([]Player, len(streamable_players))
	sort.Slice(streamable_players, func(i, j int) bool {
		return streamable_players[i].AvgPoints > streamable_players[j].AvgPoints})
	InsertStreamablePlayers(streamable_players, free_positions, week, &child, cur_streamers1)

	for i := 0; i < len(parent1.Genes); i++ {

		MixGenes(parent1.Genes[i], parent2.Genes[i], &child, fas, free_positions, week, cur_streamers1, streamable_players)
	}

	return child, cur_streamers1
}

// Function to mutate a chromosome
func Mutate(rate float64, chromosome *Chromosome, fas []Player, free_positions map[int][]string, cur_streamers []Player, week string, streamable_players []Player) {
	
	// Get random seed
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Get random number to determine if mutation occurs
	rand_num := rng.Float64()

	if rand_num < rate {

		if rand_num < rate * 0.33 {
		// Drop a random player on a random day

			Drop(rng, chromosome, week)

		} else if rand_num < rate * 0.66 {
		// Add a random player in a random position on a random day

			Add(rng, free_positions, chromosome, fas, week, cur_streamers, streamable_players)

		} else {
		// Find a valid swap for a random player on a random day and swap them
			
			Swap(chromosome, free_positions, cur_streamers, streamable_players, week)

		}
	}
}

// Function to insert a player into a chromosome
func InsertPlayer(day int, player Player, free_positions map[int][]string, child *Chromosome, week string, cur_streamers []Player, streamable_players []Player) {

	add := false
	pos_map := GetPosMap(player, child, free_positions, day, week, cur_streamers, streamable_players, day==0, true, &add, false)

	// When added here, counts as a new player
	if _, ok := pos_map[day]; !ok{
		return
	}
	if add {
		child.Genes[day].NewPlayers[player.Name] = player
	}

	for day, pos := range pos_map {
		child.Genes[day].Roster[pos] = player
	}
}

// Function to validate the insertion of a player into a chromosome
func ValidatePlayer(chromosome *Chromosome, player Player, day int) bool {

	// If the player was in the roster or bench in the past, check to see if he was dropped before 3 days ago
	for i := 0; i <= day; i++ {
		if MapContainsValue(chromosome.Genes[i].Roster, player.Name) != "" || Contains(chromosome.Genes[i].Bench, player) {
			if day - i < 2 {
				return false
			}
		}
	}

	return true
}

// Function to find a valid swap for a chromosome mutation
func FindValidSwap(chromosome *Chromosome, free_positions map[int][]string, streamable_players []Player, week string) (Player, int, Player, int) {

	// Get random seed
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	check := func(day1 int, day2 int, player1 Player, player2 Player) bool {

		// Make sure the players are not the same
		if player1.Name == player2.Name {
			return false
		}

		// Check if the players are playing on the days
		if (!Contains(ScheduleMap[week].Games[player1.Team], day2)) || (!Contains(ScheduleMap[week].Games[player2.Team], day1)) {
			return false
		}

		// Make sure neither player is in another day's new players
		for i := 0; i < len(chromosome.Genes); i++ {
			if Contains(chromosome.Genes[i].NewPlayers, player1) {
				if i != day1 {
					return false
				}
			}
			if Contains(chromosome.Genes[i].NewPlayers, player2) {
				if i != day2 {
					return false
				}
			}
		}

		// Check if the players can be put into a position on the days
		opening_for_player1 := false
		opening_for_player2 := false
		for _, pos := range free_positions[day1] {
			if Contains(player2.ValidPositions, pos) {
				if player, ok := chromosome.Genes[day1].Roster[pos]; !ok || player.Name == "" || player.Name == player1.Name {
					opening_for_player2 = true
				}
		}
		for _, pos := range free_positions[day2] {
			if Contains(player1.ValidPositions, pos) {
				if player, ok := chromosome.Genes[day2].Roster[pos]; !ok || player.Name == "" || player.Name == player2.Name {
					opening_for_player1 = true
				}
			}
		}
	}

		return opening_for_player1 && opening_for_player2
	}

	for trials := 0; trials < 100; trials++ {
		// Get two random days that are not the same
		rand_day1 := rng.Intn(len(chromosome.Genes) - 3)
		rand_day2 := rng.Intn(len(chromosome.Genes))
		for rand_day1 >= rand_day2 {
			rand_day2 = rng.Intn(len(chromosome.Genes))
		}

		// Get two random players that were added on that day
		if len(chromosome.Genes[rand_day1].NewPlayers) == 0 || len(chromosome.Genes[rand_day2].NewPlayers) == 0 {
			continue
		}
		rand_index1 := rng.Intn(len(chromosome.Genes[rand_day1].NewPlayers))
		rand_index2 := rng.Intn(len(chromosome.Genes[rand_day2].NewPlayers))
		map_keys1 := make([]string, 0, len(chromosome.Genes[rand_day1].NewPlayers))
		for key := range chromosome.Genes[rand_day1].NewPlayers {
			map_keys1 = append(map_keys1, key)
		}
		map_keys2 := make([]string, 0, len(chromosome.Genes[rand_day2].NewPlayers))
		for key := range chromosome.Genes[rand_day2].NewPlayers {
			map_keys2 = append(map_keys2, key)
		}
		player1 := chromosome.Genes[rand_day1].NewPlayers[map_keys1[rand_index1]]
		player2 := chromosome.Genes[rand_day2].NewPlayers[map_keys2[rand_index2]]

		if check(rand_day1, rand_day2, player1, player2) {
			return player1, rand_day1, player2, rand_day2

		}
	}
	return Player{}, 0, Player{}, 0
}


// Function to get the number of acquisitions in a chromosome
func GetTotalAcquisitions(chromosome *Chromosome) {

	acquisitions := 0

	for _, gene := range chromosome.Genes {

		acquisitions += len(gene.NewPlayers)
	}

	chromosome.TotalAcquisitions = acquisitions
}


// Function to give fitness score to a chromosome
func ScoreFitness(chromosome *Chromosome, week string) {

	fitness_score := 0.0
	penalty_factor := 1.0

	// Loop through each day and add the average points for each player with adjustments for lineups that go over the limit
	if chromosome.TotalAcquisitions > ScheduleMap[week].GameSpan + 1 {
		penalty_factor = 1.0 / math.Pow(1.3, float64(chromosome.TotalAcquisitions - ScheduleMap[week].GameSpan))
	}
	for _, gene := range chromosome.Genes {
		for _, player := range gene.Roster {
			fitness_score += player.AvgPoints
		}
	}

	chromosome.FitnessScore = int(fitness_score * penalty_factor)
}

// // Function to delete all occurences of a player from a chromosome (simplified version)
// func SimpleDeleteAllOccurences(chromosome *Chromosome, player_to_drop Player, week string, start_day int) {

// 	// Remove player from NewPlayers
// 	position := MapContainsValue(chromosome.Genes[start_day].NewPlayers, player_to_drop.Name)
// 	delete(chromosome.Genes[start_day].NewPlayers, position)

// 	// Remove player from Roster
// 	for day := start_day; day < len(chromosome.Genes); day++ {

// 		// If the player is on the bench on a given day remove him
// 		if Contains(chromosome.Genes[day].Bench, player_to_drop) {
// 			chromosome.Genes[day].Bench = Remove(chromosome.Genes[day].Bench, SliceIndexOf(chromosome.Genes[day].Bench, player_to_drop))
// 		}

// 		key := MapContainsValue(chromosome.Genes[day].Roster, player_to_drop.Name)
// 		if key != "" {
// 			delete(chromosome.Genes[day].Roster, key)
// 		}
// 	}

// }


// Function to drop a random player on a random day
func Drop(rng *rand.Rand, chromosome *Chromosome, week string) {

	for not_found := true; not_found; {
		// Until a day with a new player is found, keep generating random days
		day := 0
		test_day := rng.Intn(len(chromosome.Genes))
		for day == 0 {

			if len(chromosome.Genes[test_day].NewPlayers) > 0 {
				day = test_day
				break
			} else {
				test_day = rng.Intn(len(chromosome.Genes))
			}
		}

		// Turn the map of new players into a slice
		new_players := make([]Player, len(chromosome.Genes[day].NewPlayers))
		for _, player := range chromosome.Genes[day].NewPlayers {
			new_players = append(new_players, player)
		}
		rand_index := rng.Intn(len(new_players))
		player_to_drop := new_players[rand_index]

		// Check if the player is ever re-added in the future, if he is, get a new player
		for i := day; i < len(chromosome.Genes); i++ {
			if MapContainsValue(chromosome.Genes[i].Roster, player_to_drop.Name) != "" {
				continue
			}
		}
		not_found = false

		// Delete the player from the roster and decrement the acquisitions
		chromosome.Genes[day].Acquisitions -= 1
		chromosome.TotalAcquisitions -= 1
		RetroDeleteAllOccurrences(chromosome, player_to_drop, week, day)

	}
}

// Function to add a random player in a random position on a random day
func Add(rng *rand.Rand, free_positions map[int][]string, chromosome *Chromosome, fas []Player, week string, cur_streamers []Player, streamable_players []Player) {

	// Functon to check if a player is in a valid position on a day
	check_pos := func (fa Player, rand_day int) bool {
		for _, pos := range fa.ValidPositions {
			if Contains(free_positions[rand_day], pos) {
				if _, ok := chromosome.Genes[rand_day].Roster[pos]; !ok {
					return true
				}
			}
		}
		return false
	}

	// Until a valid player is found, keep generating random players and days
	for not_found := true; not_found; {

		// Generate a random day and player
		rand_day := rng.Intn(len(chromosome.Genes))
		rand_index := rng.Intn(len(fas))
		fa := fas[rand_index]

		// Check if the player is already on the roster or if the player is not playing on the day
		if MapContainsValue(chromosome.Genes[rand_day].Roster, fa.Name) != "" || !Contains(ScheduleMap[week].Games[fa.Team], rand_day) || fa.Injured || !check_pos(fa, rand_day) {
			continue
		}

		// Insert the player into the roster
		// dummy_has_match := false
		add := false
		// matches := GetMatches(fa.ValidPositions, free_positions[rand_day], &dummy_has_match)
		pos_map := GetPosMap(fa, chromosome, free_positions, rand_day, week, cur_streamers, streamable_players, false, true, &add, true)

		for day, pos := range pos_map {
			not_found = false

			chromosome.Genes[day].Roster[pos] = fa
		}

		if !not_found {
			chromosome.Genes[rand_day].NewPlayers[fa.Name] = fa
			chromosome.Genes[rand_day].Acquisitions += 1
			chromosome.TotalAcquisitions += 1
		}
	}
}

// Function to find a valid swap for a random player on a random day and swap them
func Swap(chromosome *Chromosome, free_positions map[int][]string, cur_streamers []Player, streamable_players []Player, week string) {

	player1, day1, player2, day2 := FindValidSwap(chromosome, free_positions, cur_streamers, week)

	// Delete player1 from day1
	RetroDeleteAllOccurrences(chromosome, player1, week, day1)
	// Delete player2 from day2
	RetroDeleteAllOccurrences(chromosome, player2, week, day2)

	if player1.Name != "" && player2.Name != "" {

		InsertPlayer(day2, player1, free_positions, chromosome, week, cur_streamers, streamable_players)
		InsertPlayer(day1, player2, free_positions, chromosome, week, cur_streamers, streamable_players)

	}
}

// Function to copy another chromosome's roster up to a certain index
func CopyUpToIndex(streamable_players []Player, free_positions map[int][]string, week string, parent *Chromosome, child *Chromosome, cur_streamers []Player, crossover_point int) {

	InsertStreamablePlayers(streamable_players, free_positions, week, child, cur_streamers)
	
	for i := 0; i < crossover_point; i++ {

		for _, player := range parent.Genes[i].NewPlayers {

			InsertPlayer(i, player, free_positions, child, week, cur_streamers, streamable_players)
		}

		// After each day, decrement the countdown for dropped players
		for player_name, dropped_player := range child.DroppedPlayers {
			if dropped_player.Countdown > 0 {
				dropped_player.Countdown--
				child.DroppedPlayers[player_name] = dropped_player
			} else {
				delete(child.DroppedPlayers, player_name)
			}
		}

		// Increment the acquisitions
		child.Genes[i].Acquisitions += len(child.Genes[i].NewPlayers)
	}
}

// // Function to cross over a gene from a parent into a child
// func CrossOverGene(parent_gene Gene, child *Chromosome, free_positions map[int][]string, week string, cur_streamers []Player, streamable_players []Player) {

// 	// Loop through each new player in the parent gene and try to insert them into the child
// 	for _, player := range parent_gene.NewPlayers {

// 		if ValidatePlayer(child, player, parent_gene.Day) {

// 			InsertPlayer(parent_gene.Day, player, free_positions, child, week, cur_streamers, streamable_players, false)
// 		}
// 	}
// }

// Function to mix genes from two parents into a child
func MixGenes(parent1_gene Gene, parent2_gene Gene, child *Chromosome, fas []Player, free_positions map[int][]string, week string, cur_streamers []Player, streamable_players []Player) {
	
	// Get random seed
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Create a list of all the new players in the parent genes
	new_players := make([]Player, 0)
	for _, player := range parent1_gene.NewPlayers {
		new_players = append(new_players, player)
	}
	for _, player := range parent2_gene.NewPlayers {
		new_players = append(new_players, player)
	}

	if len(new_players) == 0 {
		return
	}

	// Sort the new players by their average points
	sort.Slice(new_players, func(i, j int) bool {
		return new_players[i].AvgPoints > new_players[j].AvgPoints
	})

	// Get a random number to determine how many players to add to the child
	var rand_num int
	if len(new_players) > 1 {
		rand_num = rng.Intn(len(new_players) - 1) + 1
	} else {
		rand_num = 1
	}



	// Loop through the number of players to add to the child
	for i := 0; i < rand_num; i++ {

		// Insert the player into the child
		if ValidatePlayer(child, new_players[i], parent1_gene.Day) {
			
			InsertPlayer(parent1_gene.Day, new_players[i], free_positions, child, week, cur_streamers, streamable_players)
		}
	}

}