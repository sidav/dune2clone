package main

type taskForceMission int

const (
	AITF_MISSION_ATTACK taskForceMission = iota
	AITF_MISSION_DEFEND
	AITF_MISSION_RECON
)

type aiTaskForce struct {
	mission                      taskForceMission
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
