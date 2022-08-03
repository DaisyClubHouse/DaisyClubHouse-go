package room

type ChessMatrix struct {
	matrix    [][]bool      // 排列矩阵 15*15
	history   []*PieceCoord // 历史落子记录
	latestPos *PieceCoord   // 最近一次落子的位置
}

// PieceCoord 落子坐标
type PieceCoord struct {
	x, y int
}

func NewChessMatrix(size int) *ChessMatrix {
	matrix := make([][]bool, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]bool, size)
	}

	return &ChessMatrix{
		matrix:    matrix,
		history:   make([]*PieceCoord, 0),
		latestPos: nil,
	}
}

func (matrix *ChessMatrix) Put(x, y int) {
	matrix.matrix[x][y] = true

	// 更新最近一次落子的位置
	matrix.updatePieceCoord(x, y)
}

func (matrix *ChessMatrix) updatePieceCoord(x, y int) {
	pos := PieceCoord{x, y}

	matrix.latestPos = &pos
	matrix.history = append(matrix.history, &pos)
}

func (matrix *ChessMatrix) Existed(x, y int) bool {
	return matrix.matrix[x][y]
}

func (matrix *ChessMatrix) IsWin() bool {

	return false
}
