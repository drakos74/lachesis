package lb

import (
	"github.com/drakos74/lachesis/store/network"
	"github.com/rs/zerolog/log"
)

func LeaderFollowerPartition() network.Switch {
	return &CaptainCrewSwitch{}
}

type CaptainCrewSwitch struct {
	parallelism int
}

func (r *CaptainCrewSwitch) Register(id int) {
	// this is always the leader e.g. the last node added to the cluster
	// this makes our failure scenarios easier to build
	r.parallelism = id
}

func (r *CaptainCrewSwitch) DeRegister(id int) {
	// delegate leadership to the next in line
	log.Info().Int("index", id).Msg("De-Register From Network")
	r.parallelism--
}

func (r CaptainCrewSwitch) Route(key network.Key) ([]int, error) {
	return []int{r.parallelism}, nil
}
