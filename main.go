package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	"math"
	"math/rand"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "zack",        // TODO: Your Battlesnake username
		Color:      "#888888", // TODO: Choose color
		Head:       "default", // TODO: Choose head
		Tail:       "default", // TODO: Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	// TODO: Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height
	log.Printf("Board(%v,%v)\nBoard detail %#v\n body %#v\n\n", boardWidth, boardHeight,state.Board, state.You.Body)
	isMoveSafe = boundryCheck(state.You.Head, isMoveSafe, boardHeight, boardWidth) 
	// TODO: Step 2 - Prevent your Battlesnake from colliding with itself
	// mybody := state.You.Body

	coordinateMap := nextMoveCoordinates(state.You.Head,isMoveSafe)
	validateMoves(coordinateMap, isMoveSafe, state.You.Body)

	// TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	// opponents := state.Board.Snakes

	for _, snake := range state.Board.Snakes {
		validateMoves(coordinateMap, isMoveSafe, snake.Body)
	}

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			log.Printf("The moves allowed are %s\n",move)
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	// Choose a random move from the safe ones
	nextMove := safeMoves[rand.Intn(len(safeMoves))]

	// calculating the distance from food
	foodDistanceMap := make(map[string]int)
	for move, isSafe := range isMoveSafe {
		if isSafe {
			foodDistanceMap[move] = minDistance(distanceFromFood(state.You.Head, state.Board.Food))
			log.Printf("The move[%s] == %d distance\n",move, foodDistanceMap[move])
		}
	}

	
	// calculate the move that takes you to the fastest
	first := true
	min := 0
	for move, minDistance := range foodDistanceMap {
		if first {
			nextMove = move
			min = minDistance
			first = false
		} else {
			if minDistance < min {
				nextMove = move
				min = minDistance
			}
		}
	}

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

// minimum distance
func minDistance(distances []int) int{
	min := 0
	for i, value := range distances {
		if i == 0{
			min = value
		} else {
			if value <  min {
				min = value
			}
		}
	}

	return min
}

// Distance to all foodSources from the current position
func distanceFromFood(me Coord,foodSources []Coord) []int{
	distances := make([]int,0)
	for _, food := range foodSources {
		distances = append(distances, calculateDistance(me,food))
	}
	return distances
}

func calculateDistance(here Coord, there Coord) int {
	return int(math.Abs(float64(here.X) - float64(there.X))) + int(math.Abs(float64(here.Y) - float64(there.Y)))
}

// boundry clearence
func boundryCheck(present Coord, moves map[string]bool, height int, width int) map[string]bool {
	if present.X == (width - 1) {
		moves["right"] = false
	}
	if present.X == 0 {
		moves["left"] = false
	}
	if present.Y == 0 {
		moves["down"] = false 
	}
	if present.Y == (height - 1) {
		moves["up"] = false
	}
	log.Printf("Moves post boder check %#v and the head %#v\n",moves, present)
	return moves
}

// function to check if the moves coordinate is invalid coordianates
func validateMoves(coordinateMap map[string]Coord, moves map[string]bool, hazard []Coord) {
	for move, position := range coordinateMap {
		valid := true
		for _, coords := range hazard {
			if (coords.X == position.X) && (coords.Y == position.Y){
				valid = false
				break
			}
		}
		moves[move] = valid
	}
} 

// calculates the position of all the allowed moves.
func nextMoveCoordinates(present Coord, moves map[string]bool) map[string]Coord {
	coordinateMap := make(map[string]Coord)

	for move, isSafe := range moves {
		if isSafe{
			coordinateMap[move] = calculateNextCoordinate(present, move)
		}
	}

	return coordinateMap
}

func calculateNextCoordinate(present Coord, move string) Coord {
	switch move {
	case "up":
		present.Y  += 1
	case "down":
		present.Y -= 1
	case "left":
		present.X -= 1
	case "right":
		present.X += 1
	default:
		break; 
	}

	return present
}

func main() {
	RunServer()
}
