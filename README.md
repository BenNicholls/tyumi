# Tyumi

## A Roguelike Engine/Library Thing

Tyumi is a general purpose roguelike library and game engine. While the primary purpose is for making exciting roguelike games, it can be used to make any program that displays visually on a grid. Tyumi's various parts are oranized into separate packages; it's my hope that roguelike developers using Go might be able to get some use out of Tyumi -- even if they don't want to use the whole engine -- by grabbing individual portions of it that they find useful.

### Current State

Tyumi is currently in early stages of development. As it stands the "roguelike" part of the library doesn't even exist yet, apart from I suppose the cell-based renderer. Still, Tyumi has the features to enable making simple games!

The API is liable to change drastically as I flesh out Tyumi's capabilites, but for simple games you can target the latest 0.x release (currently 0.2). Version 0.2 has the following features (in varying states of maturity):

- **Game engine** with simple game loop. Compose your game around a Tyumi.Scene object and Tyumi will run it!
- **Roguelike data structures and algorithms**. Tilemaps, entities, things used for classic roguelikes. The package is currently quite barebones, and lot more is coming here. [package rl]
- **SDL2 based platform** implementation for rendering, audio, and input events [package platform/sdl]
- **2D cell-based canvas with drawing functions**. Canvas cells are square and support both full-width glyph drawing as well as half-width glyph drawing for writing denser text. [package gfx]
- **Animation system**, for making things flash and move and just generally fun to look at.
- **UI system** with a number of predefined elements, which can be composed around to define custom elements. UI elements are then added into a tree structure to build complex UIs. [package gfx/ui]
- **Keyboard and mouse input**. These are very rudimentary right now, but the keyboard support is enough to do simple games. Just as long as you don't need to input a capital letter :P [package input]
- Subscriber-based **Event system**, with support for custom events. [package event]
- Simple **logging system** [package log]
- **Vector** and **Utility** packages with a smattering of useful structures and functions [packages vec and util]
- **Audio system**, with basic support for loading/playing/mixing/unloading sound effects and music

### How To Get It

Get Tyumi in the usual way:

```
go get github.com/bennicholls/tyumi@v0.2.0
```

At the moment the only supported platform for Tyumi is based on [go-sdl2](https://github.com/veandco/go-sdl2), so you'll need to follow the instructions there to set up your dev environment for sdl2 correctly. Eventually other platforms will be added but for now this is what we have.

If you're feeling particularly brave you can track the main branch here instead, but I'm not sure I'd recommend it. Tyumi is something of an organic creature at the moment and I change things at a whim sometimes.

### Examples

Want to see a Tyumi game in action? Check out [Tytris](https://github.com/bennicholls/tytris), a tetris clone I put together to show off Tyumi's features as of v0.1.

More complex is [SpaceShippers](https://github.com/bennicholls/spaceshippers), a roguelike space exploration simulation game. V0.2 of Tyumi consists of additions and changes that were needed while porting SpaceShippers from my old engine, so a lot of the latest features are on display there. SpaceShippers is in early development.

Once the API is more nailed down I'll write up some little example apps, maybe a tutorial?

### Future

There's still lots of work to do. On the horizon are things like:

- **More Helpers for making roguelikes**: This is what Tyumi is supposed to be for, so *coming soon* will be better tile and map structures, procedural generation functions, pathfinding, FOV, actors, AI routines for NPCs and enemies, and much much more! Roguelikes present a huge domain of problems to solve so there's lots of work to do here!
- **More platforms**: At the moment the only platform that has been put together is SDL2 based. SDL2 is nice but Tyumi's platform system is designed so other platforms can be slotted in instead, so we'll have to make some other platform implementations to take advantage of that. In the short term, making an SDL3-based platform seems like a good idea. I also want to make a terminal platform, for making games that run in a terminal just like an old-school roguelike should. Long term I also want to have a WASM platform so people can compile a version of their game for the web.
- **Better Input Handling**: right now input handling is... lacking, to say the least. Mouse clicks don't do anything, keyboard modifier keys are not tracked, gamepad support is non-existent. So there's room for improvement here!
- **More UI Things**: more pre-built UI elements to use as building blocks, with more configuration options, and more ways to interact with them! UI can be a pain so having as much of this stuff done by the engine lets us make games faster. The biggest thing I need to nail down is some kind of consistent Theming Support. The UI package has ways to set styles for borders, default colours for objects, things like that, but it's kind of all over the place at the moment. Need to organize that and make it easier to use for sure.
- **And More!** Tyumi is built and expanded in whatever ways I need at the time while I make games with it, so who knows what features will be added next? If you have any suggestions I'd love to hear them though! Perhaps there will be a time where Tyumi can grow to meet the needs of people other than myself :)

### History

Tyumi is the 2nd generation of the engine, with much of it initially pulled from my previous engine Burl-E. While Burl-E was functional and usable, it was suffering from some longstanding structural issues that were not easy to address. It is my hope that Tyumi will be an evolution on the past, easier to work on and more useful for others who might want to take a dive in.

## Changelog

### V0.2
- [GENERAL] Some performance wins across the board. More efficient walking of the UI tree, more compact data structures, lots of good stuff I promise.
- [GENERAL] Added String() methods to a number of types, making debugging a bit easier. Added string names for Colours, Glyphs, and more.
- [COL] Switched to using a proper Colour type. No more random uint32s all over the codebase.
- [ENGINE] Renamed State to Scene. Other engines call their main game objects Scenes, so we will too.
- [ENGINE] Added multiple scene initializers, for different window configurations (bordered, fullscreen, centered, etc.)
- [ENGINE] Can now change scenes! Games with multiple levels? The possibilities are endless!
- [ENGINE] Added dialogs! Dialogs are scenes that can be opened overtop of another scene. The dialog takes control and the background scene is paused. Once done, the dialog is closed and control returns to the main scene. Good for messages, file pickers, stuff like that.
- [ENGINE] Added timers! Timers call a function after a specified number of ticks. Scene can create timers and will manage them for you, or you can make and manage your own.
- [ENGINE] Added support for Debug/Release build tags. If building a Tyumi project with the Debug build tag some extra functionality is added to the engine, including actions for taking screenshots of the console and dumping the entire UI state to disk.
- [EVENT] Made event streams embeddable. Types with an embedded event stream can now easily listen for and handle events! Wow!
- [EVENT] Added event de-duplication for streams (if you want).
- [EVENT] Event streams can now be temporarily disabled/enabled.
- [GFX] Reworked canvas internals. Canvases now track and can report whether they contain any transparent cells.
- [GFX] Canvases can now be outputted as .xp files.
- [GFX] Added Full Width text drawing mode. You can now also set a default text drawing mode for generic text drawing calls.
- [INPUT] Added Actions! Can now map inputs (key presses, mouse clicks, etc.) to an action, which emits action events that can be responded to. The implementation isn't finished yet but the idea is to allow users to remap inputs, support multiple input devices, all of that.
- [INPUT] Can now suppress all KeyRelease events. Lots of games only care about KeyPresses.
- [RL] The long-awaited Roguelike package! It's barebones for now, but the start is there. Added Tiles, Tilemaps (along with custom ui element for viewing tilemaps), Entities, Stats, and a DataCache container for holding flyweight-style data.
- [UI] Revamped UI compositing! Elements are now depth-sorted and organized into a proper opaque/transparent pipeline. As a result, elements with transparent portions don't need special handling and should "just work" (don't hold me to this, this is a hard problem. seems to work though!)
- [UI] Added an embedded event stream to all elements! Now UI elements can recieve/handle events easily!
- [UI] Added Button! Call Press() to have the button play a pulse animation and call an activation callback.
- [UI] Added ProgressBar! For showing progress via a bar! Good for HP bars, loading bars, etc.
- [UI] Images can now load/reload .xp images from disk after initialization.
- [UI] PageCountainer now accepts custom page types.
- [UI] Added element focus state! Elements can be set as focused, enabling inputs to be received. Also added support for defining a tabbing order to a window, allowing you to press tab to cycle through focused elements.
- [UI] Multiple list improvements, including some extra scrolling functionality (still work to do here) and the ability to set a message to be displayed when the list is empty.
- [UI] Began adding callbacks to common UI occurances. Still could add a lot more.
- [UI] Can now center elements inside their parents (if they have one).
- [UI] More border linking improvements. Elements can now link borders even if they are not siblings.
- [UI] Border titles and hints can now be aligned (left, center, right)
- [VEC] Added distance calc functions for Coords
- [UTIL] Added Set container! Sets contain items, but only 1 of each item. Just like in Math!
- [UTIL] Added State Machines! Define mutually exclusive states and switch between them. Also embedded a state machine into the Scene type, so all scenes can define game states and move between them. It's a common enough pattern, I think it makes sense!



### V0.1
- Initial alpha release.