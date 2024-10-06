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
			u.Team.CasualtyRefunds += gmath.Iround(float64(u.Stats.Cost) * 0.2)
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
		ac := activeCard{
			data:       c,
			phasesLeft: c.Kind.Info().Duration,
		}
		r.activeCards = append(r.activeCards, &ac)
		if ac.data.Kind.Info().Category == gcombat.CardCategoryModifier {
			ac2 := ac
			ac2.data.TeamIndex = r.stage.Teams[ac.data.TeamIndex].EnemyIndex()
			r.activeCards = append(r.activeCards, &ac2)
		}
	}
	if r.phase < len(r.stage.Teams[1].Cards) {
		c := r.stage.Teams[1].Cards[r.phase]
		ac := activeCard{
			data:       c,
			phasesLeft: c.Kind.Info().Duration,
		}
		r.activeCards = append(r.activeCards, &ac)
		if ac.data.Kind.Info().Category == gcombat.CardCategoryModifier {
			ac2 := ac
			ac2.data.TeamIndex = r.stage.Teams[ac.data.TeamIndex].EnemyIndex()
			r.activeCards = append(r.activeCards, &ac2)
		}
	}

	{
		r.cardsSlice = r.cardsSlice[:0]
		for _, c := range r.activeCards {
			if c.data.Kind == gcombat.CardFirstAid {
				for _, u := range r.allUnits {
					if u.Team.Index != c.data.TeamIndex {
						continue
					}
					if !u.Stats.Infantry {
						continue
					}
					if u.HP < u.Stats.MaxHP {
						u.HP += (u.Stats.MaxHP - u.HP) * 0.75
					}
				}
			}
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

func (r *Runner) dealSplashDamage(p *gcombat.Projectile) {
	for _, u := range r.allUnits {
		if u == p.Target {
			continue
		}
		dist := u.Pos.DistanceTo(p.Pos)
		if dist > 12 {
			continue
		}
		r.dealDamage(p.Attacker.Stats.Damage*0.75, p.Attacker, u)
	}
}

func (r *Runner) detonateProjectile(p *gcombat.Projectile) {
	if p.Attacker.Stats.SplashDamage {
		r.dealSplashDamage(p)
	}

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

	r.dealDamage(p.Attacker.Stats.Damage, p.Attacker, p.Target)
}

func (r *Runner) dealDamage(value float64, attacker, target *gcombat.Unit) {
	damage := value * game.G.Rand.FloatRange(0.8, 1.2)
	def := r.unitDefense(target)
	if def > 0 && target.Stats.Infantry && r.cardIsActive(target.Team.Index, gcombat.CardBadCover) {
		def *= 0.25
	}
	damage *= 1 - def
	if !target.Stats.Infantry {
		damage *= attacker.Stats.AntiArmorDamage
	}
	target.HP = gmath.ClampMin(target.HP-damage, 0)
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
		if r.cardIsActive(u.Team.Index, gcombat.CardInfantryCharge) {
			multiplier *= 1.6
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
		if u.Stats.Kind == gcombat.UnitTank {
			r.maybeRunOverInfantry(u)
		}
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
	case gcombat.UnitMissile:
		r.updateMissileUnit(u, delta)
	case gcombat.UnitHunter:
		r.updateHunterUnit(u, delta)
	case gcombat.UnitTank:
		r.updateTankUnit(u, delta)
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
			if dist > 96 {
				continue
			}
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

func (r *Runner) updateHunterUnit(u *gcombat.Unit, delta float64) {
	if r.cardIsActive(u.Team.Index, gcombat.CardTakeCover) {
		r.updateTakeCover(u, delta)
		return
	}

	if r.cardIsActive(u.Team.Index, gcombat.CardStandGround) {
		cx := int(u.Pos.X) / 64
		cy := int(u.Pos.Y) / 64
		u.Waypoint = gmath.Vec{
			X: (float64(cx*64) + 32) + game.G.Rand.FloatRange(-14.0, 14.0),
			Y: (float64(cy*64) + 32) + game.G.Rand.FloatRange(-14.0, 14.0),
		}
		return
	}

	u.Waypoint = r.rushWaypoint(u, u.Team.Index != 0, 160)
}

func (r *Runner) updateTankUnit(u *gcombat.Unit, delta float64) {
	if r.cardIsActive(u.Team.Index, gcombat.CardTankRush) {
		u.Waypoint = r.rushWaypoint(u, u.Team.Index != 0, 64)
		return
	}
}

func (r *Runner) updateLaserUnit(u *gcombat.Unit, delta float64) {
	if r.cardIsActive(u.Team.Index, gcombat.CardTakeCover) {
		r.updateTakeCover(u, delta)
		return
	}
}

func (r *Runner) updateMissileUnit(u *gcombat.Unit, delta float64) {
	if r.cardIsActive(u.Team.Index, gcombat.CardTakeCover) {
		r.updateTakeCover(u, delta)
		return
	}
}

func (r *Runner) rushWaypoint(u *gcombat.Unit, inverse bool, thresholdConstant float64) gmath.Vec {
	var desiredPos gmath.Vec
	if inverse {
		threshold := float64(r.stage.Level.DeployWidth*64) + thresholdConstant
		if u.Pos.X <= threshold {
			return gmath.Vec{}
		}
		desiredPos = gmath.Vec{
			Y: u.SpawnPos.Y + game.G.Rand.FloatRange(-14, 14),
			X: u.Pos.X - 64 + game.G.Rand.FloatRange(-14, 14),
		}
	} else {
		threshold := (r.stage.Width - float64(r.stage.Level.DeployWidth*64)) - thresholdConstant
		if u.Pos.X >= threshold {
			return gmath.Vec{}
		}
		desiredPos = gmath.Vec{
			Y: u.SpawnPos.Y + game.G.Rand.FloatRange(-14, 14),
			X: u.Pos.X + 64 + game.G.Rand.FloatRange(-14, 14),
		}
	}

	return u.Pos.MoveTowards(desiredPos, 64)
}

func (r *Runner) updateRifleUnit(u *gcombat.Unit, delta float64) {
	if r.cardIsActive(u.Team.Index, gcombat.CardTakeCover) {
		r.updateTakeCover(u, delta)
		return
	}

	if r.cardIsActive(u.Team.Index, gcombat.CardStandGround) {
		cx := int(u.Pos.X) / 64
		cy := int(u.Pos.Y) / 64
		u.Waypoint = gmath.Vec{
			X: (float64(cx*64) + 32) + game.G.Rand.FloatRange(-14.0, 14.0),
			Y: (float64(cy*64) + 32) + game.G.Rand.FloatRange(-14.0, 14.0),
		}
		return
	}

	inverse := u.Team.Index != 0
	if r.cardIsActive(u.Team.EnemyIndex(), gcombat.CardSuppressiveFire) {
		if !r.cardIsActive(u.Team.Index, gcombat.CardInfantryCharge) {
			inverse = !inverse
		}
	}

	u.Waypoint = r.rushWaypoint(u, inverse, 64)
}

func (r *Runner) cardIsActive(teamIndex int, k gcombat.CardKind) bool {
	for _, a := range r.activeCards {
		if a.data.Kind == k && a.data.TeamIndex == teamIndex {
			return true
		}
	}
	return false
}

func (r *Runner) maybeRunOverInfantry(u *gcombat.Unit) {
	for _, u2 := range r.allUnits {
		if u2.Team == u.Team {
			continue
		}
		if !u2.Stats.Infantry {
			continue
		}
		dist := u2.Pos.DistanceTo(u.Pos)
		if dist <= 14 {
			r.dealDamage(100, u, u2)
		}
	}
}

func (r *Runner) maybeOpenFire(u *gcombat.Unit) {
	u.Reload = u.Stats.Reload * game.G.Rand.FloatRange(0.7, 1.5)
	if game.G.Rand.Chance(0.1) && u.Reload < 1.25 {
		u.Reload *= 0.5
	}
	if u.Stats.SuppressiveROF && r.cardIsActive(u.Team.Index, gcombat.CardSuppressiveFire) {
		u.Reload *= game.G.Rand.FloatRange(0.2, 0.4)
	}
	if u.Stats.IonStorm && r.cardIsActive(u.Team.Index, gcombat.CardIonStorm) {
		return
	}

	aimedDist := u.Stats.AccuracyDist
	var target *gcombat.Unit
	focused := false
	{
		bestScore := 0.0
		focusFire := r.cardIsActive(u.Team.Index, gcombat.CardFocusFire)
		for i, u2 := range r.allUnits {
			if u2.Team == u.Team {
				continue
			}
			dist := u.Pos.DistanceTo(u2.Pos)
			score := 1000.0 - dist
			if !u2.Stats.Infantry {
				if dist <= u.Stats.AccuracyDist {
					score *= 0.8 + (u.Stats.AntiArmorDamage * 0.2)
				} else {
					score *= 0.5 + (u.Stats.AntiArmorDamage * 0.5)
				}
			}
			wasFocused := false
			if focusFire && dist <= (1.5*u.Stats.AccuracyDist) {
				score += float64(i * 30)
				if u2.HP <= u2.Stats.MaxHP*0.8 {
					wasFocused = true
					score *= (u2.HP + u2.Stats.MaxHP*0.2) / u2.Stats.MaxHP
				}
			}
			if score > bestScore {
				focused = wasFocused
				bestScore = score
				target = u2
			}
		}
	}
	if target == nil {
		return
	}
	if focused {
		aimedDist *= 1.25
	}

	var aimPos gmath.Vec
	goodAim := false
	{
		chanceToHit := u.Stats.BaseAccuracy
		if u.Stats.SuppressiveROF && r.cardIsActive(u.Team.Index, gcombat.CardSuppressiveFire) {
			chanceToHit *= 0.2
		}
		if r.cardIsActive(u.Team.Index, gcombat.CardLuckyShot) {
			chanceToHit *= 2.0
			if r.cardIsActive(target.Team.Index, gcombat.CardBadCover) {
				chanceToHit *= 2.0
			}
		}
		dist := u.Pos.DistanceTo(target.Pos)
		if dist > aimedDist {
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
