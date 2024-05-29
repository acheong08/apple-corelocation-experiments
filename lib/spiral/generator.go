package spiral

type spiral struct {
	x, y, dx, dy int
	offsetX      int
	offsetY      int
}

func NewSpiral(offsetX, offsetY int) *spiral {
	return &spiral{
		offsetX: offsetX,
		offsetY: offsetY,
		dy:      -1,
	}
}

func (s *spiral) Next() (int, int) {
	x, y := s.x, s.y
	if (s.x == s.y) || ((s.x < 0) && (s.x == -s.y)) || ((s.x > 0) && (s.x == 1-s.y)) {
		s.dx, s.dy = -s.dy, s.dx
	}
	s.x += s.dx
	s.y += s.dy
	return x + s.offsetX, y + s.offsetY
}
