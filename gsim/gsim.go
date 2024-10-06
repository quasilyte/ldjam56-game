package gsim

import (
	"math"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
)

type Runner struct {
	stage *gcombat.Stage

	phase       int
	untilPhase  float64
	activeCards []*activeCard
	cardsSlice  []gcombat.Card

	allUnits    []*gcombat.Unit
	projectiles []*gcombat.Projectile

	finished bool

	EventUpdateCards       gsignal.Event[[]gcombat.Card]
	EventFinished          gsignal.Event[*gcombat.Team]
	EventProjectileCreated gsignal.Event[*gcombat.Projectile]
}

type activeCard struct {
	data       gcombat.Card
	phasesLeft int
}

func NewRunner(stage *gcombat.Stage) *Runner {
	r := &Runner{
		stage:      stage,
		cardsSlice: make([]gcombat.Card, 0, 16),
	}

	for _, team := range stage.Teams {
		for _, u := range team.Units {
			u.Reload = u.Stats.Reload * game.G.Rand.FloatRange(0.5, 1.25)
			r.allUnits = append(r.allUnits, u)
		}
	}
	gmath.Shuffle(&game.G.Rand, r.allUnits)

	return r
}

func (r *Runner) Update(delta float64) {
	if r.finished {
		return
	}

	if r.untilPhase == 0 {
		r.untilPhase = 5
		r.updatePhase()
	}
	r.untilPhase = gmath.ClampMin(r.untilPhase-delta, 0)
	r.stage.Time += delta

	{
		live := r.projectiles[:0]
		for _, p := range r.projectiles {
			r.updateProjectile(p, delta)
			if !p.Disposed {
				live = append(live, p)
			}
		}
		r.projectiles = live
	}

	unitsByTeam := [2]int{}
	live := r.allUnits[:0]
	for _, u := range r.allUnits {
		if u.IsDisposed() {
			u.Team.Casualties++
			u.EventDisposed.Emit(gsignal.Void{})
			continue
		}
		unitsByTeam[u.Team.Index]++
		live = append(live, u)
	}
	r.allUnits = live

	if unitsByTeam[0] == 0 {
		r.finished = true
		r.EventFinished.Emit(r.stage.Teams[1])
		return
	}
	if unitsByTeam[1] == 0 {
		r.finished = true
		r.EventFinished.Emit(r.stage.Teams[0])
		return
	}

	gmath.RandIterate(&game.G.Rand, r.allUnits, func(u *gcombat.Unit) bool {
		r.updateUnit(u, delta)
		return false
	})
}

func (r *Runner) updatePhase() {
	{
		live := r.activeCards[:0]
		for _, c := range r.activeCards {
			c.phasesLeft--
			if c.phasesLeft <= 0 {
				continue
			}
			live = append(live, c)
		}
		r.activeCards = live
	}

	if r.phase < len(r.stage.Teams[0].Cards) {
		c := r.stage.Teams[0].Cards[r.phase]
		r.activeCards = append(r.activeCards, &activeCard{
			data:       c,
			phasesLeft: c.Kind.Info().Duration,
		})
	}
	if r.phase < len(r.stage.Teams[1].Cards) {
		c := r.stage.Teams[1].Cards[r.phase]
		r.activeCards = append(r.activeCards, &activeCard{
			data:       c,
			phasesLeft: c.Kind.Info().Duration,
		})
	}

	{
		r.cardsSlice = r.cardsSlice[:0]
		for _, c := range r.activeCards {
			r.cardsSlice = append(r.cardsSlice, c.data)
		}
		r.EventUpdateCards.Emit(r.cardsSlice)
	}

	r.phase++
}

func (r *Runner) updateProjectile(p *gcombat.Projectile, delta float64) {
	speed := 500.0
	if p.GoodAim {
		speed *= 1.5
	}
	p.Pos = p.Pos.MoveTowards(p.AimPos, speed*delta)
	if p.Pos == p.AimPos {
		r.detonateProjectile(p)
		p.EventDisposed.Emit(gsignal.Void{})
		p.Disposed = true
	}
}

func (r *Runner) detonateProjectile(p *gcombat.Projectile) {
	if p.Target.IsDisposed() {
		return
	}

	maxDist := p.Attacker.Stats.ProjectileHitRadius
	if p.GoodAim {
		maxDist += 4
	}
	if p.Pos.DistanceTo(p.Target.Pos) > maxDist {
		return
	}

	damage := p.Attacker.Stats.Damage * game.G.Rand.FloatRange(0.8, 1.2)
	def := r.unitDefense(p.Target)
	damage *= 1 - def
	p.Target.HP = gmath.ClampMin(p.Target.HP-damage, 0)
}

func (r *Runner) unitDefense(u *gcombat.Unit) float64 {
	cx := int(u.Pos.X) / 64
	cy := int(u.Pos.Y) / 64
	tile := r.stage.Level.Tiles[cy][cx]
	return u.Stats.TerrainDefense[tile]
}

func (r *Runner) unitSpeed(u *gcombat.Unit) float64 {
	cx := int(u.Pos.X) / 64
	cy := int(u.Pos.Y) / 64
	tile := r.stage.Level.Tiles[cy][cx]
	multiplier := u.Stats.TerrainSpeed[tile]
	if u.Stats.Infantry {
		if r.cardIsActive(u.Team.Index, gcombat.CardInfatryCharge) {
			multiplier *= 1.33
		}
	}
	return u.Stats.Speed * multiplier
}

func (r *Runner) updateUnit(u *gcombat.Unit, delta float64) {
	u.Reload = gmath.ClampMin(u.Reload-delta, 0)
	if u.Reload == 0 {
		r.maybeOpenFire(u)
	}

	if !u.Waypoint.IsZero() {
		dist := r.unitSpeed(u) * delta
		u.Pos = u.Pos.MoveTowards(u.Waypoint, dist)
		if u.Pos == u.Waypoint {
			u.Waypoint = gmath.Vec{}
			return
		}
		return
	}

	switch u.Stats.Kind {
	case gcombat.UnitRifle:
		r.updateRifleUnit(u, delta)
	case gcombat.UnitLaser:
		r.updateLaserUnit(u, delta)
	}
}

func (r *Runner) updateTakeCover(u *gcombat.Unit, delta float64) {
	var bestPos gmath.Vec
	minDist := math.MaxFloat64
	found := false
	for row, rowTiles := range r.stage.Level.Tiles {
		for col, t := range rowTiles {
			if t != gcombat.TileForest {
				continue
			}
			pos := gmath.Vec{
				X: float64(col*64) + 32,
				Y: float64(row*64) + 32,
			}
			dist := u.Pos.DistanceTo(pos)
			if dist < minDist {
				found = true
				minDist = dist
				bestPos = pos
			}
		}
	}
	if !found {
		return
	}
	desiredPos := bestPos.Add(game.G.Rand.Offset(-14, 14))
	u.Waypoint = u.Pos.MoveTowards(desiredPos, 64)
}

func (r *Runner) updateLaserUnit(u *gcombat.Unit, delta float64) {
	if r.cardIsActive(u.Team.Index, gcombat.CardTakeCover) {
		r.updateTakeCover(u, delta)
		return
	}
}

func (r *Runner) updateRifleUnit(u *gcombat.Unit, delta float64) {
	if r.cardIsActive(u.Team.Index, gcombat.CardTakeCover) {
		r.updateTakeCover(u, delta)
		return
	}

	effectiveIndex := u.Team.Index
	if r.cardIsActive(u.Team.EnemyIndex(), gcombat.CardSuppressiveFire) {
		if !r.cardIsActive(u.Team.Index, gcombat.CardInfatryCharge) {
			effectiveIndex = u.Team.EnemyIndex()
		}
	}

	var desiredPos gmath.Vec
	if effectiveIndex == 0 {
		threshold := (r.stage.Width - float64(r.stage.Level.DeployWidth*64)) - 64
		if u.Pos.X >= threshold {
			return
		}
		desiredPos = gmath.Vec{
			Y: u.SpawnPos.Y + game.G.Rand.FloatRange(-14, 14),
			X: u.Pos.X + 64 + game.G.Rand.FloatRange(-14, 14),
		}
	} else {
		threshold := float64(r.stage.Level.DeployWidth*64) + 64
		if u.Pos.X <= threshold {
			return
		}
		desiredPos = gmath.Vec{
			Y: u.SpawnPos.Y + game.G.Rand.FloatRange(-14, 14),
			X: u.Pos.X - 64 + game.G.Rand.FloatRange(-14, 14),
		}
	}

	u.Waypoint = u.Pos.MoveTowards(desiredPos, 64)
}

func (r *Runner) cardIsActive(teamIndex int, k gcombat.CardKind) bool {
	for _, a := range r.activeCards {
		if a.data.Kind == k && a.data.TeamIndex == teamIndex {
			return true
		}
	}
	return false
}

func (r *Runner) maybeOpenFire(u *gcombat.Unit) {
	u.Reload = u.Stats.Reload * game.G.Rand.FloatRange(0.7, 1.5)
	if game.G.Rand.Chance(0.1) && u.Reload < 1.25 {
		u.Reload *= 0.5
	}
	if u.Stats.Infantry && r.cardIsActive(u.Team.Index, gcombat.CardSuppressiveFire) {
		u.Reload *= game.G.Rand.FloatRange(0.2, 0.4)
	}

	var target *gcombat.Unit
	{
		bestScore := 0.0
		for _, u2 := range r.allUnits {
			if u2.Team == u.Team {
				continue
			}
			score := 1000.0 - u.Pos.DistanceTo(u2.Pos)
			if score > bestScore {
				bestScore = score
				target = u2
			}
		}
	}
	if target == nil {
		return
	}

	var aimPos gmath.Vec
	goodAim := false
	{
		chanceToHit := u.Stats.BaseAccuracy
		if u.Stats.Infantry && r.cardIsActive(u.Team.Index, gcombat.CardSuppressiveFire) {
			chanceToHit *= 0.2
		}
		if r.cardIsActive(u.Team.Index, gcombat.CardLuckyShot) {
			chanceToHit *= 2.0
		}
		dist := u.Pos.DistanceTo(target.Pos)
		if dist > u.Stats.AccuracyDist {
			chanceToHit *= 0.25
		}
		if game.G.Rand.Chance(chanceToHit) {
			aimPos = target.Pos
			goodAim = true
		} else {
			spread := math.Pow(dist, 0.75) * 0.5
			aimPos = target.Pos.Add(game.G.Rand.Offset(-spread, +spread))
		}
	}

	p := &gcombat.Projectile{
		Attacker: u,
		Target:   target,
		Pos:      u.Pos,
		AimPos:   aimPos,
		GoodAim:  goodAim,
		Rotation: u.Pos.AngleToPoint(aimPos),
	}
	r.projectiles = append(r.projectiles, p)
	r.EventProjectileCreated.Emit(p)
}
