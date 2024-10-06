package groundscape

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/ebitengine-graphics/particle"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type projectileNode struct {
	data    *gcombat.Projectile
	sprite  *graphics.Sprite
	emitter *particle.Emitter
	state   *sceneState
}

func newProjectileNode(data *gcombat.Projectile, state *sceneState) *projectileNode {
	return &projectileNode{
		data:  data,
		state: state,
	}
}

func (p *projectileNode) IsDisposed() bool {
	return p.data.Disposed
}

func (p *projectileNode) Dispose() {
	p.sprite.Dispose()
	if p.emitter != nil {
		p.emitter.Pos.Base = nil
		p.emitter.SetEmitting(false)
	}
}

func (p *projectileNode) Init(scene *gscene.Scene) {
	p.sprite = game.G.NewSprite(p.data.Attacker.Stats.ProjectileImage)
	p.sprite.Pos.Offset = sceneutil.CombatMapOffset(game.G.State.CurrentStage.MapBg)
	p.sprite.Pos.Base = &p.data.Pos
	p.sprite.Rotation = &p.data.Rotation
	scene.AddGraphics(p.sprite, 0)

	if sfx := p.data.Attacker.Stats.FireSound; sfx != assets.AudioNone {
		game.G.PlaySound(sfx)
	}

	switch p.data.Attacker.Stats.Kind {
	case gcombat.UnitMissile:
		emitter := particle.NewEmitter(particleTemplateMissileTrail)
		emitter.Pos.Offset = sceneutil.CombatMapOffset(game.G.State.CurrentStage.MapBg)
		emitter.Pos.Base = &p.data.Pos
		emitter.Rotation = &p.data.Rotation
		emitter.PivotOffset.X = -6
		p.state.renderer.AddEmitter(emitter)
		emitter.SetEmitting(true)
		p.state.emitters = append(p.state.emitters, emitter)
		p.emitter = emitter
	}

	p.data.EventDisposed.Connect(nil, func(gsignal.Void) {
		if p.data.Attacker.Stats.SplashDamage {
			effect := newEffectNode(effectNodeConfig{
				Sprite: game.G.NewSprite(assets.ImageExplosion),
				Pos:    p.sprite.Pos.Resolve(),
			})
			scene.AddObject(effect)
		}
		p.Dispose()
	})
}

func (p *projectileNode) Update(delta float64) {}
