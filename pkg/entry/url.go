package entry

import "fmt"

// HostPath is a tuple of both the Host and the Path
type HostPath struct {
	Host, Path string
}

func (h HostPath) String() string {
	if len(h.Host) < 1 {
		return h.Path
	}
	return fmt.Sprintf("%s%s", h.Host, h.Path)
}

// HostPathScore is a tuple of both the HostPath and a Score for sorting
type HostPathScore struct {
	HostPath HostPath
	Score    Score
}

// URL holds various HostPath and scores, along with the total number of urls
type URL struct {
	values map[HostPath]Score
	total  int
}

// NewURL creates a URL
func NewURL() *URL {
	return &URL{make(map[HostPath]Score, 0), 0}
}

// Add adds a HostPath to the URL for counting.
func (u *URL) Add(host HostPath) {
	u.values[host]++
	u.total++
}

// Len returns the total amount of HostPath's added
func (u *URL) Len() int {
	return u.total
}

// Union returns the most common HostPathScore
func (u *URL) Union() HostPathScore {
	common := HostPathScore{HostPath{}, 0}
	for k, v := range u.values {
		if v > common.Score {
			common.HostPath = k
			common.Score = v
		}
	}

	common.Score = common.Score / Score(float64(u.total))
	return common
}

// Difference returns all the HostPathScore's that aren't the most common.
func (u *URL) Difference() []HostPathScore {
	var (
		common = u.Union()
		alt    = make([]HostPathScore, 0)
	)

	for k, v := range u.values {
		if k == common.HostPath {
			continue
		}
		alt = append(alt, HostPathScore{
			HostPath: k,
			Score:    v / Score(float64(u.total)),
		})
	}

	return alt
}

func (u *URL) String() string {
	return u.Union().HostPath.String()
}
