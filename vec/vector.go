package vec

import (
	"math"

	"github.com/bennicholls/tyumi/util"
)

// Vec2f is a 2 dimensional vector of floats
type Vec2f struct {
	X, Y float64
}

func (v *Vec2f) Set(x, y float64) {
	v.X, v.Y = x, y
}

func (v *Vec2f) Mod(dx, dy float64) {
	v.X += dx
	v.Y += dy
}

func (v1 Vec2f) Add(v2 Vec2f) Vec2f {
	return Vec2f{v1.X + v2.X, v1.Y + v2.Y}
}

// returns vec2f = v1 - v2
func (v1 Vec2f) Sub(v2 Vec2f) Vec2f {
	return Vec2f{v1.X - v2.X, v1.Y - v2.Y}
}

func (v Vec2f) Mag() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec2f) ToVec2i() Vec2i {
	return Vec2i{util.RoundFloatToInt(v.X), util.RoundFloatToInt(v.Y)}
}

func (v Vec2f) ToPolar() Vec2Polar {
	return Vec2Polar{v.Mag(), math.Atan2(v.Y, v.X)}
}

type Vec2Polar struct {
	R, Phi float64
}

func (v *Vec2Polar) Set(r, phi float64) {
	v.R, v.Phi = r, phi
}

func (v Vec2Polar) Get() (float64, float64) {
	return v.R, v.Phi
}

// Add converts to recitlinear components and adds, then converts back to polar.
func (v1 Vec2Polar) Add(v2 Vec2Polar) Vec2Polar {
	return v1.ToRect().Add(v2.ToRect()).ToPolar()
}

// ToRect converts the Polar vector into a rectilinear form
func (v Vec2Polar) ToRect() Vec2f {
	return Vec2f{v.R * math.Cos(v.Phi), v.R * math.Sin(v.Phi)}
}

// Reorients vector to ensure R is positive and 0 <= Phi < 2*pi
func (v *Vec2Polar) Pos() {
	if v.R < 0 {
		v.Phi += math.Pi
		v.R = -v.R
	}

	for v.Phi < 0 {
		v.Phi += 2 * math.Pi
	}

	for v.Phi > 2*math.Pi {
		v.Phi -= 2 * math.Pi
	}
}

// Returns the shortest anglular distance from v1 to v2. positive for counterclockwise, negative for clockwise.
// NOTE: Do these need to be Pos()'d?? Hmm.
func (v1 Vec2Polar) AngularDistance(v2 Vec2Polar) float64 {
	d := v2.Phi - v1.Phi

	if d > math.Pi {
		d -= 2 * math.Pi
	} else if d < -math.Pi {
		d += 2 * math.Pi
	}

	return d
}

type Vec2i struct {
	X, Y int
}

func (v1 Vec2i) Add(v2 Vec2i) Vec2i {
	return Vec2i{v1.X + v2.X, v1.Y + v2.Y}
}

func (v1 Vec2i) Sub(v2 Vec2i) Vec2i {
	return Vec2i{v1.X - v2.X, v1.Y - v2.Y}
}

// Returns the magnitude of the vector. Note that this is a float and not an int.
func (v Vec2i) Mag() float64 {
	return v.ToVec2f().Mag()
}

func (v Vec2i) ToVec2f() Vec2f {
	return Vec2f{X: float64(v.X), Y: float64(v.Y)}
}
