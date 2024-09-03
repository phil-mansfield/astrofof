package symfof

// Particle represents a single particle. Most of the time particles will
// be in "code units" where positions are equal to the size of a single
// internal FoF cell, and velocities are in units that lead to G being 1.
type Particle struct {
	// ID is a 64-bit integer which uniquely identifies the particle.
	ID uint64
	// X and V are the position and velocity of the particle.
	X, V [3]float32
}

func ParticleXCmp(p1, p2 Particle) int {
	if p1.X[0] < p2.X[0] {
		return -1
	} else if p1.X[0] > p2.X[0] {
		return +1
	}
	return 0
}

func ParticleYCmp(p1, p2 Particle) int {
	if p1.X[1] < p2.X[1] {
		return -1
	} else if p1.X[1] > p2.X[1] {
		return +1
	}
	return 0
}

func ParticleZCmp(p1, p2 Particle) int {
	if p1.X[2] < p2.X[2] {
		return -1
	} else if p1.X[2] > p2.X[2] {
		return +1
	}
	return 0
}