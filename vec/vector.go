package vec

import "math"
import "github.com/bennicholls/tyumi/util"

//Vec2 is a 2 dimensional vector
type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Set(x, y float64) {
	v.X, v.Y = x, y
}

func (v Vec2) Get() (float64, float64) {
	return v.X, v.Y
}

//Returns (X, Y) as ints, rounding as necessary.
func (v Vec2) GetInt() (int, int) {
	return util.RoundFloatToInt(v.X), util.RoundFloatToInt(v.Y)
}

func (v *Vec2) Mod(dx, dy float64) {
	v.X += dx
	v.Y += dy
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1.X + v2.X, v1.Y + v2.Y}
}

//returns vec2 = v1 - v2
func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v1.X - v2.X, v1.Y - v2.Y}
}

func (v Vec2) Mag() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec2) ToPolar() Vec2Polar {
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

//Add converts to recitlinear components and adds, then converts back to polar.
func (v1 Vec2Polar) Add(v2 Vec2Polar) Vec2Polar {
	return v1.ToRect().Add(v2.ToRect()).ToPolar()
}

//ToRect converts the Polar vector into a rectilinear form
func (v Vec2Polar) ToRect() Vec2 {
	return Vec2{v.R * math.Cos(v.Phi), v.R * math.Sin(v.Phi)}
}

//Reorients vector to ensure R is positive and 0 <= Phi < 2*pi
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

//Returns the shortest anglular distance from v1 to v2. positive for counterclockwise, negative for
//clockwise. NOTE: Do these need to be Pos()'d?? Hmm.
func (v1 Vec2Polar) AngularDistance(v2 Vec2Polar) float64 {
	d := v2.Phi - v1.Phi

	if d > math.Pi {
		d -= 2 * math.Pi
	} else if d < -math.Pi {
		d += 2 * math.Pi
	}

	return d
}