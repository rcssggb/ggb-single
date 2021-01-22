package player

import (
	"github.com/aunum/goro/pkg/v1/layer"
	m "github.com/aunum/goro/pkg/v1/model"
)

var qLearningModel *m.Sequential

// InitQLearning instantiates models for actor-critic learning of
// a discrete policy
func (p *Player) InitQLearning() error {
	qModel, err := m.NewSequential("discretePolicy")
	if err != nil {
		return err
	}

	inSize := len(p.State()) + len(p.SampleDiscreteActionVector())

	xShape := []int{1, inSize}
	yShape := []int{1, 1}

	qModel.AddLayers(
		layer.FC{Input: len(p.State()) + len(p.SampleDiscreteActionVector()), Output: 64},
		layer.FC{Input: 256, Output: 64},
		layer.FC{Input: 64, Output: 16},
		layer.FC{Input: 4, Output: 1, Activation: layer.Linear},
	)

	qModel.Compile(
		m.NewInput("state", xShape),
		m.NewInput("action-value", yShape),
		m.WithBatchSize(1),
	)

	qLearningModel = qModel

	return nil
}
