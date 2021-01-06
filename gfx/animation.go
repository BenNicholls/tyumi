package gfx

//Anything that can do animations on a Canvas
type Animator interface {
	Update()
	Render(*Canvas)
	Done() bool
}