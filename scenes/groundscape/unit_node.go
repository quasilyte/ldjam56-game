package groundscape

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type unitNode struct {
	data   *gcombat.Unit
	sprite *graphics.Sprite
}

func newUnitNode(data *gcombat.Unit) *unitNode {
	return &unitNode{data: data}
}

func (u *unitNode) IsDisposed() bool {
	return u.data.HP <= 0
}

func (u *unitNode) Init(scene *gscene.Scene) {
	u.sprite = game.G.NewSprite(u.data.Stats.Image)
	u.sprite.Pos.Offset = sceneutil.CombatMapOffset(game.G.State.CurrentStage.MapBg)
	u.sprite.Pos.Base = &u.data.Pos
	if u.data.Team.Index != 0 {
		u.sprite.SetHorizontalFlip(true)
	}
	scene.AddGraphics(u.sprite, 0)

	u.data.EventDisposed.Connect(nil, func(gsignal.Void) {
		u.Dispose()
	})
}

func (u *unitNode) Dispose() {
	u.sprite.Dispose()
}

func (u *unitNode) Update(delta float64) {}
