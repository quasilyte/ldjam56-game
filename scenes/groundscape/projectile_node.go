package groundscape

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type projectileNode struct {
	data   *gcombat.Projectile
	sprite *graphics.Sprite
}

func newProjectileNode(data *gcombat.Projectile) *projectileNode {
	return &projectileNode{data: data}
}

func (u *projectileNode) IsDisposed() bool {
	return u.data.Disposed
}

func (u *projectileNode) Dispose() {
	u.sprite.Dispose()
}

func (u *projectileNode) Init(scene *gscene.Scene) {
	u.sprite = game.G.NewSprite(u.data.Attacker.Stats.ProjectileImage)
	u.sprite.Pos.Offset = sceneutil.CombatMapOffset(game.G.State.CurrentStage.MapBg)
	u.sprite.Pos.Base = &u.data.Pos
	u.sprite.Rotation = &u.data.Rotation
	scene.AddGraphics(u.sprite, 0)

	if sfx := u.data.Attacker.Stats.FireSound; sfx != assets.AudioNone {
		game.G.PlaySound(sfx)
	}

	u.data.EventDisposed.Connect(nil, func(gsignal.Void) {
		if u.data.Attacker.Stats.SplashDamage {
			effect := newEffectNode(effectNodeConfig{
				Sprite: game.G.NewSprite(assets.ImageExplosion),
				Pos:    u.sprite.Pos.Resolve(),
			})
			scene.AddObject(effect)
		}
		u.Dispose()
	})
}

func (u *projectileNode) Update(delta float64) {}
