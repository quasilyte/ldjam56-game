package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"

	_ "image/png"
)

func registerAudioResources(loader *resource.Loader) {
	audioResources := map[resource.AudioID]resource.AudioInfo{
		AudioUnitDestroyed: {Path: "audio/unit_destroyed.wav", Volume: +0.1},

		AudioLaser1: {Path: "audio/laser1.wav", Volume: -0.3},
		AudioLaser2: {Path: "audio/laser2.wav", Volume: -0.3},
		AudioLaser3: {Path: "audio/laser3.wav", Volume: -0.3},

		AudioRifle1: {Path: "audio/rifle1.wav", Volume: -0.4},
		AudioRifle2: {Path: "audio/rifle2.wav", Volume: -0.4},
		AudioRifle3: {Path: "audio/rifle3.wav", Volume: -0.4},
	}

	for id, res := range audioResources {
		loader.AudioRegistry.Set(id, res)
		loader.LoadAudio(id)
	}
}

func NumSamples(a resource.AudioID) int {
	switch a {
	case AudioLaser1:
		return 3
	case AudioRifle1:
		return 3
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioUnitDestroyed

	AudioLaser1
	AudioLaser2
	AudioLaser3

	AudioRifle1
	AudioRifle2
	AudioRifle3
)
