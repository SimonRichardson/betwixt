package entry

// StringScore is a tuple containing a string and a score
type StringScore struct {
	String string
	Score  Score
}

// String holds a series of string values with associated scores
type String struct {
	values map[string]Score
	total  int
}

// NewString returns a new String
func NewString() *String {
	return &String{make(map[string]Score, 0), 0}
}

// Add adds the method to the string and updates the totals
func (m *String) Add(method string) {
	m.values[method]++
	m.total++
}

// Len returns the total added methods using add
func (m *String) Len() int {
	return m.total
}

// Union returns the most common string
func (m *String) Union() StringScore {
	common := StringScore{"", 0.0}
	for k, v := range m.values {
		if v > common.Score {
			common.String = k
			common.Score = v
		}
	}
	common.Score = common.Score / Score(float64(m.total))
	return common
}

// Difference returns a slice of StringScores that doesn't equal the Union.
func (m *String) Difference() []StringScore {
	var (
		common = m.Union()
		alt    = make([]StringScore, 0)
	)

	for k, v := range m.values {
		if k == common.String {
			continue
		}
		alt = append(alt, StringScore{
			String: k,
			Score:  v / Score(float64(m.total)),
		})
	}

	return alt
}

func (m *String) String() string {
	return m.Union().String
}
