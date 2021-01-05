package engine

//This is the gameloop
func Run() {
	
}

//This is the generic update function. Handles input, then calls the current active state's update function.
//This runs at speed determined by 
func update() {

}

//This function updates any UI elements that need updating after the most recent tick in the current active state's window.
//This is a synchronization function, so both update() and render() cannot run while this is happening.
func updateUI() {

}

//builds the frame and renders using whatever the current renderer is (sdl, web, terminal, whatever)
//this runs at speed determined by user-input FPS, defaulting to 60 FPS. this also updates any current animations in the active
//state's ui tree
func render() {

}