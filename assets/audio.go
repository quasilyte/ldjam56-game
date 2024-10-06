package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"

	_ "image/png"
)

func registerAudioResources(loader *resource.Loader) {
	audioResources := map[resource.AudioID]resource.AudioInfo{
		AudioButtonClick:     {Path: "audio/button_click.wav", Volume: -0.5},
		AudioButtonClickSoft: {Path: "audio/button_click_soft.wav", Volume: -0.25},
		AudioVictory:         {Path: "audio/victory.wav", Volume: +0},

		AudioUnitDestroyed: {Path: "audio/unit_destroyed.wav", Volume: +0.1},

		AudioLaser1: {Path: "audio/laser1.wav", Volume: -0.8},
		AudioLaser2: {Path: "audio/laser2.wav", Volume: -0.8},
		AudioLaser3: {Path: "audio/laser3.wav", Volume: -0.8},

		AudioRifle1: {Path: "audio/rifle1.wav", Volume: -0.4},
		AudioRifle2: {Path: "audio/rifle2.wav", Volume: -0.4},

		AudioMissile1: {Path: "audio/missile1.wav", Volume: -0.6},
		AudioMissile2: {Path: "audio/missile2.wav", Volume: -0.6},

		AudioHunter1: {Path: "audio/hunter1.wav", Volume: -0.3},
		AudioHunter2: {Path: "audio/hunter2.wav", Volume: -0.3},
		AudioHunter3: {Path: "audio/hunter3.wav", Volume: -0.3},

		AudioTank1: {Path: "audio/tank1.wav", Volume: -0.6},
		AudioTank2: {Path: "audio/tank2.wav", Volume: -0.6},
		AudioTank3: {Path: "audio/tank3.wav", Volume: -0.6},
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
		return 2
	case AudioHunter1:
		return 3
	case AudioTank1:
		return 3
	case AudioMissile1:
		return 2
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioButtonClick
	AudioButtonClickSoft
	AudioVictory

	AudioUnitDestroyed

	AudioLaser1
	AudioLaser2
	AudioLaser3

	AudioRifle1
	AudioRifle2

	AudioMissile1
	AudioMissile2

	AudioHunter1
	AudioHunter2
	AudioHunter3

	AudioTank1
	AudioTank2
	AudioTank3
)
