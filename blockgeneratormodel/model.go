package blockgeneratormodel

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

type BlockGenerator struct {
	TargetTotalEvents    int64
	BlocksPerRealization int
	TargetEventsPerBlock int
	Seed                 int64
}
type Block struct {
	RealizationIndex int   `json:"realization_index"`
	BlockIndex       int   `json:"block_index"`
	BlockEventCount  int   `json:"block_event_count"`
	BlockEventStart  int64 `json:"block_event_start"` //inclusive - will be one greater than previous event end
	BlockEventEnd    int64 `json:"block_event_end"`   //inclusive - will be one less than event start if event count is 0.
}

func (b Block) ContainsEvent(eventIndex int) bool {
	if b.BlockEventStart <= int64(eventIndex) {
		if b.BlockEventEnd >= int64(eventIndex) {
			return true
		}
	}
	return false
}
func (bg BlockGenerator) GenerateBlocks() []Block {
	blocks := make([]Block, 0)
	var EventStart int64 = 1
	var EventEnd int64 = 1
	poisson := distuv.Poisson{}
	poisson.Lambda = float64(bg.TargetEventsPerBlock)
	poisson.Src = rand.NewSource(uint64(bg.Seed))
	Index := 1
	Realization := 1
	for {
		events := int(poisson.Rand())
		EventEnd = EventStart + (int64(events) - 1)
		block := Block{BlockIndex: Index, RealizationIndex: Realization, BlockEventCount: events, BlockEventStart: EventStart, BlockEventEnd: EventEnd}
		blocks = append(blocks, block)
		if Index == bg.BlocksPerRealization {
			Realization++
			Index = 0 //always will be adding at the last line of the method.
			if EventEnd >= bg.TargetTotalEvents {
				return blocks
			}
		}
		EventStart = EventEnd + 1
		Index++
	}
}
