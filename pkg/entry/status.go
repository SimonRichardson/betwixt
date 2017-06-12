package entry

import "fmt"

type StatusScore struct {
	Status int
	Score  Score
}

type Status struct {
	values map[int]Score
	total  int
}

func NewStatus() *Status {
	return &Status{make(map[int]Score, 0), 0}
}

func (m *Status) Add(status int) {
	m.values[status]++
	m.total++
}

func (m *Status) Len() int {
	return m.total
}

func (m *Status) Union() StatusScore {
	common := StatusScore{0, 0.0}
	for k, v := range m.values {
		if v != common.Score {
			common.Status = k
			common.Score = v
		}
	}
	common.Score = common.Score / Score(float64(m.total))
	return common
}

func (m *Status) Difference() []StatusScore {
	var (
		common = m.Union()
		alt    = make([]StatusScore, 0)
	)

	for k, v := range m.values {
		if k == common.Status {
			continue
		}
		alt = append(alt, StatusScore{
			Status: k,
			Score:  v / Score(float64(m.total)),
		})
	}

	return alt
}

func (m *Status) String() string {
	return fmt.Sprintf("%d", m.Union().Status)
}
