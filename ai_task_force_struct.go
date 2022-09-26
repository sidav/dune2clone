package main

const (
	AITF_DESIGNATION_ATTACK = iota
	AITF_DESIGNATION_DEFEND = iota
)

type aiTaskForce struct {
	designation                  int
	nextTickToGiveOrders         int
	desiredSize                  int
	maxFullnessPercentForRetreat int
	noRetreatAllowed             bool
	units                        []*unit
	target                       actor
}

func (atf *aiTaskForce) getSize() int {
	return len(atf.units)
}

func (atf *aiTaskForce) getFullnessPercent() int {
	return getPercentInt(len(atf.units), atf.desiredSize)
}

func (atf *aiTaskForce) shouldBeRetreated() bool {
	return atf.maxFullnessPercentForRetreat != 0 && atf.getFullnessPercent() <= atf.maxFullnessPercentForRetreat
}

func (atf *aiTaskForce) isFull() bool {
	return len(atf.units) >= atf.desiredSize
}

func (atf *aiTaskForce) doesHaveUnit(u *unit) bool {
	for _, unt := range atf.units {
		if u == unt {
			return true
		}
	}
	return false
}

func (atf *aiTaskForce) addUnit(u *unit) {
	if atf.doesHaveUnit(u) {
		panic("Duplicated unit in ATF!")
	}
	atf.units = append(atf.units, u)
}

func (atf *aiTaskForce) cleanup() {
	if atf.target != nil && !atf.target.isAlive() {
		atf.target = nil
	}
	for i := len(atf.units) - 1; i >= 0; i-- {
		if !atf.units[i].isAlive() {
			atf.units = append(atf.units[:i], atf.units[i+1:]...)
		}
	}
}
