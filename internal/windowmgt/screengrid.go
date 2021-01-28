package windowmgt

import (
	"github.com/rotisserie/eris"

	"github.com/aacebedo/snapi3/internal"
)

type Position struct {
	x int32
	y int32
}

type ScreenGrid struct {
	rows       uint32
	cols       uint32
	cellWidth  uint32
	cellHeight uint32
	positions  [][]*Position
}

func NewPosition(x, y int32) (res *Position) {
	res = new(Position)
	res.x = x
	res.y = y

	return
}

func (p *Position) X() int32 {
	return p.x
}

func (p *Position) Y() int32 {
	return p.y
}

func NewScreenGrid(rows, cols, screenWidth, screenHeight uint32) (res *ScreenGrid, err error) {
	if rows == 0 || cols == 0 || screenWidth == 0 || screenHeight == 0 {
		err = eris.Wrapf(internal.InvalidArgumentError, "Unable to create a screen grid as "+
			"one of the arguments is equal to '0': rows='%d', cols='%d', screenWidth='%d', screenHeight='%d'",
			rows, cols, screenWidth, screenWidth)

		return
	}

	res = new(ScreenGrid)
	res.cols = cols
	res.rows = rows
	res.cellWidth = screenWidth / cols
	res.cellHeight = screenHeight / rows
	res.positions = make([][]*Position, 0, rows)
	currentY := int32(0)

	for i := uint32(0); i < rows; i++ {
		res.positions = append(res.positions, make([]*Position, 0, cols))
		currentX := int32(0)

		for j := uint32(0); j < cols; j++ {
			res.positions[i] = append(res.positions[i], NewPosition(currentX, currentY))
			currentX += int32(res.cellWidth)
		}

		currentY += int32(res.cellHeight)
	}

	return
}

func (sg *ScreenGrid) CellWidth() uint32 {
	return sg.cellWidth
}

func (sg *ScreenGrid) CellHeight() uint32 {
	return sg.cellHeight
}

func (sg *ScreenGrid) GetPosition(row, col uint32) (res *Position, err error) {
	if row < sg.rows && col < sg.cols {
		res = sg.positions[row][col]
	} else {
		if row >= sg.rows {
			err = eris.Wrapf(internal.OutOfRangeArgumentError, "Given row '%d' is more than the grid rows ('%d')", row, sg.rows)
		} else {
			err = eris.Wrapf(internal.OutOfRangeArgumentError, "Given column '%d' is more than the grid columns ('%d')", col, sg.cols)
		}
	}

	return
}
