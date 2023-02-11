package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type Card struct {
	ID      int
	Counter map[string]int
}

type Request struct {
	Cards        []*Card
	Requirements map[string]int
}

func NewCard(id int, counter map[string]int) *Card {
	return &Card{
		ID:      id,
		Counter: counter,
	}
}

type Team []*Card

func (t Team) Len() int {
	return len(t)
}

func (t Team) Less(i, j int) bool {
	for skill := range t[i].Counter {
		if t[j].Counter[skill] > t[i].Counter[skill] {
			return true
		}
	}
	return false
}

func (t Team) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func createTeam(cards []*Card, reqs map[string]int) [][]int {
	fmt.Println("Start searching", reqs)
	res := [][]int{}
	q := [][]*Card{}
	q = append(q, []*Card{})
	last_t := 0
	q_count := 0
	for len(q) > 0 {
		t := q[0]
		q = q[1:]
		m := map[string]int{}
		q_count++
		if len(t) > last_t {
			fmt.Println("Team size:", len(t))
			last_t = len(t)
			q_count = 0
			sort.Slice(q, func(i, j int) bool {
				i_score := 0
				j_score := 0
				im := map[string]int{}
				jm := map[string]int{}
				for _, card := range q[i] {
					for skill, count := range card.Counter {
						im[skill] += count
					}
				}
				for _, card := range q[j] {
					for skill, count := range card.Counter {
						jm[skill] += count
					}
				}

				for skill, count := range reqs {
					if im[skill] >= count {
						i_score += count
					} else if cards[i].Counter[skill] >= 1 {
						i_score += cards[i].Counter[skill]
					}
					if jm[skill] >= count {
						j_score += count
					} else if cards[j].Counter[skill] >= 1 {
						j_score += cards[j].Counter[skill]
					}
				}
				return i_score > j_score
			})

		}
		if q_count > 250 {
			continue
		}
		if len(t) > 4 {
			break
		}
		if len(res) > 40 {
			break
		}
		for _, card := range t {
			for skill, count := range card.Counter {
				m[skill] += count
			}
		}

		flag := true
		for skill, count := range reqs {
			if m[skill] < count {
				flag = false
				break
			}
		}
		if flag {
			team := make([]int, len(t))
			for i, card := range t {
				team[i] = card.ID
			}
			res = append(res, team)
			if len(res) >= 20 {
				return res
			}
			continue
		}
		sort.Slice(cards, func(i, j int) bool {
			i_score := 0
			j_score := 0
			for skill, count := range reqs {
				if m[skill] < count {
					skill_missing := count - m[skill]
					if cards[i].Counter[skill] >= skill_missing {
						i_score += skill_missing
					} else if cards[i].Counter[skill] >= 1 {
						i_score += cards[i].Counter[skill]
					}
					if cards[j].Counter[skill] >= skill_missing {
						j_score += skill_missing
					} else if cards[j].Counter[skill] >= 1 {
						j_score += cards[j].Counter[skill]
					}
				}
			}
			return i_score > j_score
		})
		for idx, card := range cards {
			nt := make([]*Card, len(t)+1)
			copy(nt, t)
			nt[len(t)] = card
			q = append(q, nt)
			if idx > 20 {
				break
			}
			if len(q) > 500 {
				break
			}
		}
	}
	return res
}

func main() {
	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var cards []*Card
	err = json.Unmarshal(file, &cards)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		var request Request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		results := createTeam(cards, request.Requirements)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	err = http.ListenAndServe(":28080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
