package spiral

type spiral struct {
	x, y, dx, dy int
	offset       int
}

func NewSpiral(offset int) *spiral {
	return &spiral{
		offset: offset,
		dy:     -1,
	}
}

func (s *spiral) Next() (int, int) {
	x, y := s.x, s.y
	if (s.x == s.y) || ((s.x < 0) && (s.x == -s.y)) || ((s.x > 0) && (s.x == 1-s.y)) {
		s.dx, s.dy = -s.dy, s.dx
	}
	s.x += s.dx
	s.y += s.dy
	return x + s.offset, y + s.offset
}
