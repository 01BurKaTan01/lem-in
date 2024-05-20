package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("ERROR: invalid number of arguments: please enter 'go run ./cmd/ filename.txt'")
		fmt.Println("all examples in 'examples' folder")
		return
	}
	rooms, err := ReadContent(args[0])
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	// Print the input file content
	fmt.Println(rooms.AntCount)
	fmt.Println("##start")
	fmt.Printf("%s %d %d\n", rooms.RoomFirst.Name, rooms.RoomFirst.X, rooms.RoomFirst.Y)
	fmt.Println("##end")
	fmt.Printf("%s %d %d\n", rooms.RoomLast.Name, rooms.RoomLast.X, rooms.RoomLast.Y)
	for _, room := range rooms.AllRooms {
		fmt.Printf("%s %d %d\n", room.Name, room.X, room.Y)
	}
	for _, link := range rooms.Links {
		fmt.Println(link)
	}

	fmt.Println()

	mainGraph, cloneGraph := arrangeGraphs(rooms)
	if mainGraph == nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	foundPaths, err := mainGraph.DetectAvailablePaths(cloneGraph, rooms.AntCount)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	antsOnEachPath := PathwiseAntCount(foundPaths, rooms.AntCount)
	AntMovements(antsOnEachPath, foundPaths)
}

type Room struct {
	Name string
	X    int
	Y    int
}

type RoomsAndAnts struct {
	RoomFirst Room
	RoomLast  Room
	AllRooms  []Room
	Links     []string
	AntCount  int
}
type Network struct {
	rooms   map[string]*Point
	Initial *Point
	Finish  *Point
}

type Point struct {
	inverted   bool
	identifier string
	neighbors  []*Point
	prior      *Point
}

func (a *Point) FetchKey() string {
	return a.identifier
}

func SetupGraph() *Network {
	valueTemp := &Network{
		rooms: make(map[string]*Point),
	}
	return valueTemp
}

func (a *Network) AddVertex(identifier string) error {
	valueTemp := &Point{identifier: identifier}
	if _, isHave := a.rooms[identifier]; isHave {
		return fmt.Errorf("already has vertex")
	}
	a.rooms[identifier] = valueTemp
	return nil
}

func (a *Network) InsertEdge(coming, toward *Point) {
	coming.neighbors = append(coming.neighbors, toward)
	toward.neighbors = append(toward.neighbors, coming)
}

func (a *Network) InsertDirectedEdge(coming, toward *Point) {
	coming.neighbors = append(coming.neighbors, toward)
}

func (a *Network) removeEdge(coming, toward *Point) {
	for index, vertex := range coming.neighbors {
		if vertex == toward {
			coming.neighbors = append(coming.neighbors[:index], coming.neighbors[index+1:]...)
			break
		}
	}
	for index, vertex := range toward.neighbors {
		if vertex == coming {
			toward.neighbors = append(toward.neighbors[:index], toward.neighbors[index+1:]...)
			break
		}
	}
}

func (a *Network) TakeVertex(identifier string) *Point {
	if _, have := a.rooms[identifier]; have {
		return a.rooms[identifier]
	}
	return nil
}

func (a *Network) arrangeStart(identifier string) {
	a.Initial = a.TakeVertex(identifier)
}

func (a *Network) arrangeEnd(identifier string) {
	a.Finish = a.TakeVertex(identifier)
}

func (a *Network) BFS(coming, toward *Point) ([]*Point, map[*Point]bool) {
	isVisited := map[*Point]bool{coming: true}
	taskQueue := []*Point{coming}

	for len(taskQueue) > 0 {
		current := taskQueue[0]
		for _, v := range current.neighbors {
			if isVisited[v] {
				continue
			}
			isVisited[v] = true
			v.prior = current
			if v == toward {
				return a.reversePath(v)
			}
			taskQueue = append(taskQueue, v)
		}
		taskQueue = taskQueue[1:]
	}
	return nil, nil
}

func (a *Network) reversePath(finish *Point) ([]*Point, map[*Point]bool) {
	flipped := []*Point{}
	for node := finish; node != nil; node = node.prior {
		node.inverted = true
		flipped = append(flipped, node)
	}
	result := make([]*Point, len(flipped))
	resultOfMap := make(map[*Point]bool)
	for i, j := len(flipped)-1, 0; i >= 0; i, j = i-1, j+1 {
		result[j] = flipped[i]
		resultOfMap[result[j]] = true
	}
	delete(resultOfMap, result[0])
	delete(resultOfMap, result[len(result)-1])
	for i := 1; i < len(result); i++ {
		a.removeEdge(result[i], result[i].prior)
		a.InsertDirectedEdge(result[i], result[i].prior)
	}
	return result, resultOfMap
}

func (a *Network) DetectAvailablePaths(cloneGraph *Network, ants int) ([][]*Point, error) {
	if ants <= 2 {
		path, err := a.SingleDirectionSearch()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return path, nil
	}
	crossings, err := cloneGraph.DetectCrossings()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	a.removeCrossings(crossings)

	foundedPaths := a.ExplorePaths()
	return foundedPaths, nil
}

func (a *Network) SingleDirectionSearch() ([][]*Point, error) {
	path, _ := a.BFS(a.Initial, a.Finish)
	if path == nil {
		return nil, fmt.Errorf("no available paths")
	}
	// without start
	result := [][]*Point{path[1:]}
	return result, nil
}

func (a *Network) ExplorePaths() [][]*Point {
	foundPaths := [][]*Point{}
	visitedVertices := []map[*Point]bool{}
	for {
		path, mapPath := a.BFS(a.Initial, a.Finish)
		if path == nil {
			break
		}
		// if path has only start - end
		if len(path) == 2 {
			// returning path without start vertex
			foundPaths = append(foundPaths, path[1:])
			return foundPaths
		}
		if len(visitedVertices) == 0 {
			visitedVertices = append(visitedVertices, mapPath)
			foundPaths = append(foundPaths, path[1:])
			continue
		}
		if hasVerticesCrossings(visitedVertices, mapPath) {
			continue
		}
		visitedVertices = append(visitedVertices, mapPath)
		foundPaths = append(foundPaths, path[1:])
	}
	return foundPaths
}

func (a *Network) DetectCrossings() ([]string, error) {
	_, pathFound := a.bhandariAlgorithmCrossings(a.Initial, a.Finish)
	if !pathFound {
		return nil, fmt.Errorf("invalid data format")
	}
	crossings := []string{}
	for {
		crossedVertices, pathFound := a.bhandariAlgorithmCrossings(a.Initial, a.Finish)
		if !pathFound {
			break
		}
		crossings = append(crossings, crossedVertices...)
	}
	return crossings, nil
}

func (a *Network) removeCrossings(crossings []string) {
	for _, v := range crossings {
		temp := strings.Split(v, " ")
		a.removeEdge(a.TakeVertex(temp[0]), a.TakeVertex(temp[1]))
	}
}

func hasVerticesCrossings(visitedVertices []map[*Point]bool, path map[*Point]bool) bool {
	for _, v := range visitedVertices {
		if isCrossed(v, path) {
			return true
		}
	}
	return false
}

func isCrossed(path, currentpath map[*Point]bool) bool {
	for vertex := range currentpath {
		if _, have := path[vertex]; have {
			return true
		}
	}
	return false
}

func (a *Network) bhandariAlgorithmCrossings(coming, toward *Point) ([]string, bool) {
	visited := map[*Point]bool{coming: true}
	queue := []*Point{coming}

	for len(queue) > 0 {
		current := queue[0]
		for _, v := range current.neighbors {
			if visited[v] {
				continue
			}
			visited[v] = true
			v.prior = current
			if v == toward {
				cross := a.SwapPathDirections(v)
				return cross, true
			}
			queue = append(queue, v)
		}
		queue = queue[1:]
	}
	return nil, false
}

func (a *Network) SwapPathDirections(finish *Point) []string {
	crossings := []string{}
	for vertex := finish; vertex != nil; vertex = vertex.prior {
		if vertex.prior != nil {
			if vertex.inverted && vertex.prior.inverted {
				crossings = append(crossings, vertex.prior.identifier+" "+vertex.identifier)
			}
		}
		if vertex.prior != nil {
			a.removeEdge(vertex, vertex.prior)
			a.InsertDirectedEdge(vertex, vertex.prior)
		}

		vertex.inverted = true
	}
	return crossings
}

func PathwiseAntCount(locatedPaths [][]*Point, ants int) [][]string {
	antsOnPath := make([]int, len(locatedPaths))
	result := make([][]string, len(locatedPaths))
	antsOnPath[0]++
	antCounter := 1
	result[0] = append(result[0], "L"+strconv.Itoa(antCounter))
	antCounter++
	ants--
	if len(locatedPaths) > 1 {
		for i := 0; ants > 0; {
			if i+1 >= len(locatedPaths) {
				i = 0
			}
			antID := fmt.Sprintf("L%v", antCounter)
			if len(locatedPaths[i])+antsOnPath[i] == len(locatedPaths[i+1])+antsOnPath[i+1] {
				antsOnPath[i]++
				result[i] = append(result[i], antID)
				antCounter++
				ants--
				continue
			} else if len(locatedPaths[i])+antsOnPath[i] < len(locatedPaths[i+1])+antsOnPath[i+1] {
				antsOnPath[i]++
				result[i] = append(result[i], antID)
				antCounter++
				ants--
				i = 0
				continue
			}
			antsOnPath[i+1]++
			result[i+1] = append(result[i+1], antID)
			antCounter++
			ants--
			i++
		}
	} else {
		antsOnPath[0] += ants
		for i := 1; i < antsOnPath[0]; i++ {
			result[0] = append(result[0], "L"+strconv.Itoa(antCounter))
			antCounter++
		}
	}
	return result
}

func AntMovements(antsOnPath [][]string, foundPaths [][]*Point) {
	if len(foundPaths[0]) == 1 {
		count := len(antsOnPath[0])
		for i := 1; count > 0; i++ {
			fmt.Printf("L%d-%s ", i, foundPaths[0][0].identifier)
			count--
		}
		fmt.Println()
		return
	}
	maxLen := len(antsOnPath[0])
	for _, v := range antsOnPath {
		if len(v) > maxLen {
			maxLen = len(v)
		}
	}

	result := make([][]string, 1)
	for index, element, dataStack := 0, 0, 0; index < len(antsOnPath); index++ {
		for resIndex, vertex := range foundPaths[index] {
			if element >= len(antsOnPath[index]) {
				break
			}
			if resIndex+dataStack >= len(result) {
				result = append(result, []string{})
			}
			result[resIndex+dataStack] = append(result[resIndex+dataStack], antsOnPath[index][element]+"-"+vertex.FetchKey())
		}
		if index+1 >= len(antsOnPath) {
			index = -1
			element++
			dataStack++
		}
		if element >= maxLen {
			break
		}
	}
	for _, stack := range result {
		for _, ant := range stack {
			fmt.Printf("%s ", ant)
		}
		fmt.Println()
	}
}

func ReadContent(fileRoute string) (*RoomsAndAnts, error) {
	roomsAndAnts := &RoomsAndAnts{}
	file, err := os.Open(fileRoute)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	defer file.Close()
	for scanner, step := bufio.NewScanner(file), 0; scanner.Scan(); {
		if step == 0 {
			step++
			ants, err := strconv.Atoi(scanner.Text())
			if ants <= 0 || err != nil {
				return nil, fmt.Errorf("invalid data format")
			}
			roomsAndAnts.AntCount = ants
			continue
		}
		if scanner.Text() == "" {
			continue
		}
		if scanner.Text() == "##start" {
			scanner.Scan()
			temp := strings.Split(scanner.Text(), " ")
			if len(temp) != 3 {
				return nil, fmt.Errorf("invalid data format")
			}
			x, _ := strconv.Atoi(temp[1])
			y, _ := strconv.Atoi(temp[2])
			roomsAndAnts.RoomFirst = Room{Name: temp[0], X: x, Y: y}
			continue
		}
		if scanner.Text() == "##end" {
			scanner.Scan()
			temp := strings.Split(scanner.Text(), " ")
			if len(temp) != 3 {
				return nil, fmt.Errorf("invalid data format")
			}
			x, _ := strconv.Atoi(temp[1])
			y, _ := strconv.Atoi(temp[2])
			roomsAndAnts.RoomLast = Room{Name: temp[0], X: x, Y: y}
			continue
		}
		if scanner.Text()[:1] == "#" {
			continue
		}
		if strings.Contains(scanner.Text(), "-") {
			roomsAndAnts.Links = append(roomsAndAnts.Links, scanner.Text())
			continue
		}
		temp := strings.Split(scanner.Text(), " ")
		if len(temp) != 3 {
			return nil, fmt.Errorf("invalid data format")
		}
		x, _ := strconv.Atoi(temp[1])
		y, _ := strconv.Atoi(temp[2])
		roomsAndAnts.AllRooms = append(roomsAndAnts.AllRooms, Room{Name: temp[0], X: x, Y: y})
	}
	if roomsAndAnts.RoomFirst.Name == "" || roomsAndAnts.RoomLast.Name == "" {
		return nil, fmt.Errorf("invalid data format")
	}
	return roomsAndAnts, nil
}

func arrangeGraphs(rooms *RoomsAndAnts) (*Network, *Network) {
	mainGraph := SetupGraph()
	copiedGraph := SetupGraph()

	mainGraph.AddVertex(rooms.RoomFirst.Name)
	mainGraph.AddVertex(rooms.RoomLast.Name)

	copiedGraph.AddVertex(rooms.RoomFirst.Name)
	copiedGraph.AddVertex(rooms.RoomLast.Name)

	mainGraph.arrangeStart(rooms.RoomFirst.Name)
	mainGraph.arrangeEnd(rooms.RoomLast.Name)

	copiedGraph.arrangeStart(rooms.RoomFirst.Name)
	copiedGraph.arrangeEnd(rooms.RoomLast.Name)

	for _, v := range rooms.AllRooms {
		err := mainGraph.AddVertex(v.Name)
		if err != nil {
			return nil, nil
		}
		copiedGraph.AddVertex(v.Name)
	}

	for _, v := range rooms.Links {
		temp := strings.Split(v, "-")
		from := mainGraph.TakeVertex(temp[0])
		to := mainGraph.TakeVertex(temp[1])
		if from == nil || to == nil {
			return nil, nil
		}
		mainGraph.InsertEdge(from, to)
		copiedGraph.InsertEdge(copiedGraph.TakeVertex(temp[0]), copiedGraph.TakeVertex(temp[1]))
	}

	return mainGraph, copiedGraph
}
