package room

type ChessMatrix [][]bool

func NewChessMatrix(size int) ChessMatrix {
	matrix := make([][]bool, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]bool, size)
	}
	return matrix
}

func (matrix ChessMatrix) Put(x, y int) {
	matrix[x][y] = true
}

func (matrix ChessMatrix) Existed(x, y int) bool {
	return matrix[x][y]
}
