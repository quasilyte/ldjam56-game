package groundscape

import "github.com/quasilyte/ebitengine-graphics/particle"

type sceneState struct {
	emitters []*particle.Emitter
	renderer *particle.Renderer
}

func (state *sceneState) UpdateEmitters(delta float64) {
	active := state.emitters[:0]
	for _, e := range state.emitters {
		if e.Pos.Base == nil && e.NumParticles() == 0 {
			e.Dispose()
			continue
		}
		e.UpdateWithDelta(delta)
		active = append(active, e)
	}
	state.emitters = active
}
