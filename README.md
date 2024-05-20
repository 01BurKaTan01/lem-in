# Lem-in Project

## Overview

The Lem-in project involves simulating an ant colony's movement through a network of rooms and tunnels. The objective is to find the most efficient and fastest path for a group of ants moving from a start room to an end room in the least number of moves. This project tests the ability to handle graph traversal and pathfinding algorithms effectively in Go, adhering to specified constraints and handling various edge cases and input formats.

## Objectives

- Create a program named `lem-in` that reads from a file describing the ants and the colony.
- Display the content of the file and each move the ants make from room to room.
- Output the result in a specific format.

## Input and Output Requirements

### Input File Format

- The first line contains the number of ants.
- Rooms are defined by `name coord_x coord_y` (e.g., `Room 1 2`, `nameoftheroom 1 6`).
- The start room is indicated by `##start` and the end room by `##end`.
- Links are defined by `name1-name2` (e.g., `1-2`).

### Output Format

number_of_ants
the_rooms
the_links
Lx-y Lz-w Lr-o ...


- `Lx` represents an ant, where x is the ant number.
- y, w, o represent room names.

## Movement Rules

1. At the beginning, all ants are in the room `##start`.
2. The goal is to bring all ants to the room `##end` with as few moves as possible.
3. Each room can only contain one ant at a time, except for `##start` and `##end`, which can hold any number of ants.
4. Each tunnel can be used only once per turn.
5. Display only the ants that moved in each turn.
6. Move each ant only once per turn through a tunnel to an empty room.

## Error Handling

- Handle cases where there is no path between `##start` and `##end`.
- Handle rooms that link to themselves.
- Handle invalid inputs (e.g., no `##start` or `##end`, duplicated rooms, links to unknown rooms).
- Display appropriate error messages (e.g., "ERROR: invalid data format").

## Algorithm

- Use a pathfinding algorithm (such as BFS) to determine the fastest way to get n ants across the colony.
- Optimize for the shortest and least congested paths.

## Examples

### Example 1

**Input (test0.txt):**
3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5


**Output:**
3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5

L1-3 L2-2
L1-4 L2-5 L3-3
L1-0 L2-6 L3-4
L2-0 L3-0


### Example 2

**Input (test1.txt):**
3
##start
0 1 0
##end
1 5 0
2 9 0
3 13 0
0-2
2-3
3-1


**Output:**
3
##start
0 1 0
##end
1 5 0
2 9 0
3 13 0
0-2
2-3
3-1

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3
L3-1


### Example 3

**Input (test2.txt):**
3
2 5 0
##start
0 1 2
##end
1 9 2
3 5 4
0-2
0-3
2-1
3-1
2-3

**Output:**
3
2 5 0
##start
0 1 2
##end
1 9 2
3 5 4
0-2
0-3
2-1
3-1
2-3

L1-2 L2-3
L1-1 L2-1 L3-2
L3-1


## Implementation Details

- Use Go's standard packages only.
- Ensure proper error handling and output formatting.
- Follow good coding practices.
- Include unit tests to validate the functionality.

## Usage

To run the program with an input file:

```sh
go run . test0.txt
Additional Notes
Each ant is represented by Lx, where x is its number.
Ensure the solution handles all edge cases and is robust against invalid input.
Consider efficiency and performance, especially for larger colonies.


This README provides a comprehensive overview of the Lem-in project, detailing the objectives, input/output requirements, movement rules, error handling, algorithm, examples, implementation details, and usage instructions.
