package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type cartDirection int
type trackDirection int

// implements sort.Interface
type cartSet []cart

const (
	blankTileChar         = ' '
	horizontalTrackChar   = '-'
	verticalTrackChar     = '|'
	curveUpTrackChar      = '/'
	curveDownTrackChar    = '\\'
	intersectionTrackChar = '+'
	upCartChar            = '^'
	leftCartChar          = '<'
	downCartChar          = 'v'
	rightCartChar         = '>'
)

const (
	upDirection cartDirection = iota
	rightDirection
	downDirection
	leftDirection
)

const (
	verticalDirection trackDirection = iota
	horizontalDirection
	intersectionDirection
	// We define up and down for these constants as what will happen when you encounter these types of curves as you are going to the RIGHT
	curveUpDirection
	curveDownDirection
)

type cart struct {
	row, col          int
	currentTrack      *track
	direction         cartDirection
	willTurn          bool
	nextTurnDirection cartDirection
}

type track struct {
	row, col  int
	direction trackDirection
	neighbors []*track
}

func makeTrack(row int, col int, direction trackDirection) *track {
	return &track{
		row:       row,
		col:       col,
		direction: direction,
		// all tracks must have at least two neighbors
		neighbors: make([]*track, 0, 2),
	}
}

func (c *cart) turnAtIntersection() {
	if c.willTurn && c.nextTurnDirection == leftDirection {
		c.nextTurnDirection = rightDirection
		if c.direction == upDirection {
			c.direction = leftDirection
		} else if c.direction == rightDirection {
			c.direction = upDirection
		} else if c.direction == downDirection {
			c.direction = rightDirection
		} else if c.direction == leftDirection {
			c.direction = downDirection
		}
		c.willTurn = false
	} else if c.willTurn && c.nextTurnDirection == rightDirection {
		c.nextTurnDirection = leftDirection
		if c.direction == upDirection {
			c.direction = rightDirection
		} else if c.direction == rightDirection {
			c.direction = downDirection
		} else if c.direction == downDirection {
			c.direction = leftDirection
		} else if c.direction == leftDirection {
			c.direction = upDirection
		}
	} else if !c.willTurn {
		c.willTurn = true
	}
}

// turnAtCurve turns at a curve, and returns true if it could, false otherwise.
func (c *cart) turnAtCurve() {
	if c.direction == upDirection && c.currentTrack.direction == curveUpDirection {
		c.direction = rightDirection
	} else if c.direction == upDirection && c.currentTrack.direction == curveDownDirection {
		c.direction = leftDirection
	} else if c.direction == rightDirection && c.currentTrack.direction == curveUpDirection {
		c.direction = upDirection
	} else if c.direction == rightDirection && c.currentTrack.direction == curveDownDirection {
		c.direction = downDirection
	} else if c.direction == downDirection && c.currentTrack.direction == curveUpDirection {
		c.direction = leftDirection
	} else if c.direction == downDirection && c.currentTrack.direction == curveDownDirection {
		c.direction = rightDirection
	} else if c.direction == leftDirection && c.currentTrack.direction == curveUpDirection {
		c.direction = downDirection
	} else if c.direction == leftDirection && c.currentTrack.direction == curveDownDirection {
		c.direction = upDirection
	}
}

func (c *cart) move() {
	if c.direction == upDirection {
		c.row--
	} else if c.direction == downDirection {
		c.row++
	} else if c.direction == leftDirection {
		c.col--
	} else if c.direction == rightDirection {
		c.col++
	}
	for _, neighborTrack := range c.currentTrack.neighbors {
		if neighborTrack.row == c.row && neighborTrack.col == c.col {
			c.currentTrack = neighborTrack
		}
	}
}

func (set cartSet) Len() int {
	return len(set)
}

func (set cartSet) Less(i int, j int) bool {
	return set[i].row < set[j].row
}

func (set cartSet) Swap(i int, j int) {
	set[i], set[j] = set[j], set[i]
}

func identifyTile(tile rune) (isCart bool, cartTravelDirection cartDirection, direction trackDirection) {
	if tile == horizontalTrackChar {
		direction = horizontalDirection
	} else if tile == verticalTrackChar {
		direction = verticalDirection
	} else if tile == curveUpTrackChar {
		direction = curveUpDirection
	} else if tile == curveDownTrackChar {
		direction = curveDownDirection
	} else if tile == upCartChar {
		isCart = true
		cartTravelDirection = upDirection
		direction = verticalDirection
	} else if tile == intersectionTrackChar {
		direction = intersectionDirection
	} else if tile == downCartChar {
		isCart = true
		cartTravelDirection = downDirection
		direction = verticalDirection
	} else if tile == leftCartChar {
		isCart = true
		cartTravelDirection = leftDirection
		direction = horizontalDirection
	} else if tile == rightCartChar {
		isCart = true
		cartTravelDirection = rightDirection
		direction = horizontalDirection
	}

	return
}

func parseTracks(rawTracks []string) cartSet {
	carts := make(cartSet, 0)
	previousRowTracks := make([]*track, len(rawTracks[0]))
	for row := range rawTracks {
		// The first track on a line cannot possibly be horizontal, unless the track were open.
		var lastTrack *track
		for col, tile := range rawTracks[row] {
			if tile == blankTileChar {
				lastTrack = nil
				previousRowTracks[col] = nil
				continue
			}

			haveCart, cartTravelDrection, direction := identifyTile(tile)
			newTrack := makeTrack(row, col, direction)
			aboveTrack := previousRowTracks[col]
			if direction == horizontalDirection {
				lastTrack.neighbors = append(lastTrack.neighbors, newTrack)
				newTrack.neighbors = append(newTrack.neighbors, lastTrack)
			} else if direction == verticalDirection {
				aboveTrack.neighbors = append(aboveTrack.neighbors, newTrack)
				newTrack.neighbors = append(newTrack.neighbors, aboveTrack)
			} else {
				// handle the \--, --/, and + cases
				if aboveTrack != nil {
					aboveTrack.neighbors = append(aboveTrack.neighbors, newTrack)
					newTrack.neighbors = append(newTrack.neighbors, aboveTrack)
				}
				if lastTrack != nil {
					lastTrack.neighbors = append(lastTrack.neighbors, newTrack)
					newTrack.neighbors = append(newTrack.neighbors, lastTrack)
				}
			}

			lastTrack = newTrack
			previousRowTracks[col] = newTrack

			if haveCart {
				newCart := cart{
					row:               row,
					col:               col,
					direction:         cartTravelDrection,
					currentTrack:      newTrack,
					willTurn:          true,
					nextTurnDirection: leftDirection,
				}
				carts = append(carts, newCart)
			}
		}
	}

	return carts
}

// getCollidedPair returns the indices of the carts that collided
func getCollidedPair(carts cartSet) (int, int) {
	for i := range carts {
		for j := range carts {
			if j != i && carts[i].row == carts[j].row && carts[i].col == carts[j].col {
				return i, j
			}
		}
	}

	return -1, -1
}

func part1(carts cartSet) (int, int) {
	cartsAreCollided := false
	var collidedRow, collidedCol int
	for !cartsAreCollided {
		sort.Sort(carts)
		for i := range carts {
			if carts[i].currentTrack.direction == intersectionDirection {
				carts[i].turnAtIntersection()
			}
			carts[i].move()
			carts[i].turnAtCurve()
			collidedCart1, collidedCart2 := getCollidedPair(carts)
			if collidedCart1 != -1 && collidedCart2 != -1 {
				cartsAreCollided = true
				collidedRow, collidedCol = carts[collidedCart1].row, carts[collidedCart2].col
				break
			}
		}
	}

	return collidedRow, collidedCol
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main in_file")
		return
	}

	inFile := os.Args[1]
	inFileContents, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(err)
	}
	rawTracks := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawTracks = rawTracks[:len(rawTracks)-1]

	carts := parseTracks(rawTracks)
	collidedRow, collidedCol := part1(carts)
	fmt.Printf("%d,%d\n", collidedCol, collidedRow)
}
