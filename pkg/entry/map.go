package entry

import (
	"encoding/json"
	"sort"
)

// Value defines a Key Value tuple pairing.
type Value struct {
	Key, Value string
}

// ValuePromoted defines a value that can be promoted
type ValuePromoted struct {
	Value    string
	Promoted bool
}

// ScorePromoted defines a score that can be promoted
type ScorePromoted struct {
	Score    Score
	Promoted bool
}

// Values is a type alias for a map which has a key and value of string
type Values map[string]string

// Each loops through each value in a predictable way.
// Note: this can be inefficient
func (v Values) Each(fn func(string, interface{})) {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		var m interface{}
		json.Unmarshal([]byte(v[k]), &m)
		fn(k, m)
	}
}

// ValuesPromoted is a type alias for a map with key as a string and value as a
// ValuePromoted
type ValuesPromoted map[string]ValuePromoted

// ValueScore is a tuple of a Value and a Score
type ValueScore struct {
	Value Value
	Score Score
}

// ValuesScore is a tuple of a series of Values and Score
type ValuesScore struct {
	Values Values
	Score  Score
}

// Map defines a structure of values with associated scores and a threshold
// score, which states what should be worked upon.
type Map struct {
	values    map[Value]*ScorePromoted
	total     int
	threshold Score
	merge     bool
}

// NewMap creates a Map with some default sane values.
func NewMap() *Map {
	return &Map{make(map[Value]*ScorePromoted, 0), 0, 1.0, false}
}

// Add adds ValuesPromoted to the map
func (p *Map) Add(params ValuesPromoted) {
	for k, v := range params {
		key := Value{Key: k, Value: v.Value}
		if _, ok := p.values[key]; ok {
			p.values[key].Score++
			continue
		}
		p.values[key] = &ScorePromoted{
			Score:    0,
			Promoted: v.Promoted,
		}
	}
	p.total++
}

// Len returns the total number of ValuesPromoted created
func (p *Map) Len() int {
	return p.total
}

// Union returns the common ValuesScore of the map.
func (p *Map) Union() ValuesScore {
	scored := make(map[Value]Score, 0)
	for k, v := range p.values {
		var (
			score = v.Score
			total = Score(float64(p.total - 1))
		)
		if v.Promoted {
			score = total
		}
		scored[k] = score / total
	}

	values := make(Values, 0)
	for k, v := range scored {
		if v >= p.threshold {
			values[k.Key] = k.Value
		}
	}

	return ValuesScore{
		Values: values,
		Score:  p.threshold,
	}
}

// Difference returns a slice of ValuesScore that are not common.
func (p *Map) Difference() []ValuesScore {
	var (
		res []ValuesScore

		common = p.Union()
		values = make(map[Score]Values, 0)
	)

	for k, v := range p.values {
		// Remove if it's already found in unique
		if _, ok := common.Values[k.Key]; ok {
			continue
		}

		score := v.Score
		if _, ok := values[score]; !ok {
			values[score] = make(Values, 0)
		}
		values[score][k.Key] = k.Value
	}

	for k, v := range values {
		res = append(res, ValuesScore{
			Values: v,
			Score:  k,
		})
	}

	sort.Sort(ValuesScores(res))

	// Remove duplicates

	var highest Values
	for k, v := range res {
		if k == 0 {
			highest = v.Values
			continue
		}

		for k, w := range v.Values {
			if _, ok := highest[k]; ok {
				if p.merge {
					var (
						x = extract(res[0].Values[k])
						y = extract(w)
					)
					res[0].Values[k] = concat(x, y)
				}

				delete(v.Values, k)
			}
		}
	}

	return res
}

// ValuesScores is a type alias for sorting a slice of ValuesScore
type ValuesScores []ValuesScore

func (v ValuesScores) Len() int {
	return len(v)
}

func (v ValuesScores) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ValuesScores) Less(i, j int) bool {
	return v[i].Score < v[j].Score
}

func extract(v string) []string {
	var x []string
	json.Unmarshal([]byte(v), &x)
	return x
}

func concat(a, b []string) string {
	bytes, _ := json.Marshal(append(a, b...))
	return string(bytes)
}
