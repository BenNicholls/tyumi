package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type Border struct {
	top    gfx.Canvas //top, including upper left corner
	bottom gfx.Canvas //bottom, including bottom right corner
	left   gfx.Canvas //left, including bottom left corner
	right  gfx.Canvas //right, including top right corner

	title   string
	hint    string
	colours col.Pair

	styleFlag borderStyleFlag
	style     *BorderStyle

	//SCROLLBAR STUFF. for now, only vertical scrollbar for lists and the like.
	scrollbar                 bool //whether the scrollbar is enabled. scrollbar will be drawn whenever content doesn't fit
	scrollbarContentHeight    int  //total height of scrolling content
	scrollbarViewportPosition int  //position of the viewed content

	dirty bool //flag indicates when border needs to be re-rendered.
}

func NewBorder(element_size vec.Dims) (b *Border) {
	b = new(Border)
	
	b.top.Init(element_size.W+1, 1)
	b.bottom.Init(element_size.W+1, 1)
	b.left.Init(1, element_size.H+1)
	b.right.Init(1, element_size.H+1)

	b.dirty = true

	return
}

func (b *Border) setColours(col col.Pair) {
	if b.colours == col {
		return
	}

	b.colours = col
	b.dirty = true
}

// renders the border to the internal canvas
func (b *Border) Render() {
	if !b.dirty {
		return
	}

	//determine colours and update internal canavases if necessary.
	colours := b.style.Colours
	if colours.Fore == gfx.COL_DEFAULT {
		colours.Fore = b.colours.Fore
	}
	if colours.Back == gfx.COL_DEFAULT {
		colours.Back = b.colours.Back
	}
	b.top.SetDefaultColours(colours)
	b.bottom.SetDefaultColours(colours)
	b.left.SetDefaultColours(colours)
	b.right.SetDefaultColours(colours)

	for cursor := range vec.EachCoord(b.top.Bounds()) { //top and bottom
		b.top.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_LR])
		b.bottom.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_LR])
	}

	for cursor := range vec.EachCoord(b.left.Bounds()) { //left and right
		b.left.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_UD])
		b.right.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_UD])
	}

	w, h := b.top.Size().W, b.left.Size().H
	b.top.DrawGlyph(vec.Coord{0, 0}, 0, b.style.Glyphs[BORDER_DR])        //upper left corner
	b.right.DrawGlyph(vec.Coord{0, 0}, 0, b.style.Glyphs[BORDER_DL])      //upper right corner
	b.bottom.DrawGlyph(vec.Coord{w - 1, 0}, 0, b.style.Glyphs[BORDER_UL]) //bottom right corner
	b.left.DrawGlyph(vec.Coord{0, h - 1}, 0, b.style.Glyphs[BORDER_UR])   //bottom left corner

	//decorate and draw title
	if b.title != "" {
		decoratedTitle := b.style.DecorateText(b.title)
		if len([]rune(decoratedTitle))%2 == 1 {
			decoratedTitle += string(b.style.TextDecorationPad)
		}
		b.top.DrawText(vec.Coord{1, 0}, 0, decoratedTitle, colours, gfx.DRAW_TEXT_LEFT)
	}

	//decorate and draw hint
	if b.hint != "" {
		decoratedHint := b.style.DecorateText(b.hint)
		if len([]rune(decoratedHint))%2 == 1 {
			decoratedHint = string(b.style.TextDecorationPad) + decoratedHint
		}
		hintOffset := w - len([]rune(decoratedHint))/2 - 1
		b.bottom.DrawText(vec.Coord{hintOffset, 0}, 0, decoratedHint, colours, 0)
	}

	//draw scrollbar if necessary
	if b.scrollbar && b.scrollbarContentHeight > h-1 {
		b.right.DrawGlyph(vec.Coord{0, 1}, 0, gfx.GLYPH_TRIANGLE_UP)
		b.right.DrawGlyph(vec.Coord{0, h - 1}, 0, gfx.GLYPH_TRIANGLE_DOWN)

		barSize := util.Clamp(util.RoundFloatToInt(float64(h-1)/float64(b.scrollbarContentHeight)*float64(h-3)), 1, h-4)

		var barPos int
		if b.scrollbarViewportPosition+h-1 >= b.scrollbarContentHeight {
			barPos = h - 3 - barSize
		} else {
			barPos = util.Clamp(util.RoundFloatToInt(float64(b.scrollbarViewportPosition)/float64(b.scrollbarContentHeight)*float64(h-3)), 0, h-4-barSize)
		}

		for i := range barSize {
			b.right.DrawGlyph(vec.Coord{0, i + barPos + 2}, 0, gfx.GLYPH_FILL)
		}
	}
}

func (b *Border) DrawToCanvas(dst_canvas *gfx.Canvas, offset vec.Coord, depth int) {
	w, h := b.top.Size().W, b.left.Size().H
	offset_top := offset.Add(vec.Coord{-1, -1})
	offset_bottom := offset.Add(vec.Coord{0, h - 1})
	offset_left := offset.Add(vec.Coord{-1, 0})
	offset_right := offset.Add(vec.Coord{w - 1, -1})

	if b.dirty && !b.style.DisableLink {
		b.linkBorderSegment(&b.top, dst_canvas, offset_top, vec.DIR_DOWN, depth)
		b.linkBorderSegment(&b.bottom, dst_canvas, offset_bottom, vec.DIR_UP, depth)
		b.linkBorderSegment(&b.left, dst_canvas, offset_left, vec.DIR_RIGHT, depth)
		b.linkBorderSegment(&b.right, dst_canvas, offset_right, vec.DIR_LEFT, depth)
		b.dirty = false
	}

	b.top.DrawToCanvas(dst_canvas, offset_top, depth)
	b.bottom.DrawToCanvas(dst_canvas, offset_bottom, depth)
	b.left.DrawToCanvas(dst_canvas, offset_left, depth)
	b.right.DrawToCanvas(dst_canvas, offset_right, depth)
}

func (b *Border) linkBorderSegment(border_segment, dst_canvas *gfx.Canvas, dst_offset vec.Coord, inner_dir vec.Direction, depth int) {
	for src_cursor := range vec.EachCoord(border_segment) {
		src_cell := border_segment.GetCell(src_cursor)
		if src_cell.Mode == gfx.DRAW_TEXT {
			continue
		}
		src_glyph := src_cell.Glyph
		if !b.style.glyphIsLinkable(src_glyph) {
			continue
		}

		linkedGlyph := b.getLinkedGlyph(src_glyph, src_cursor.Add(dst_offset), inner_dir, dst_canvas, depth)
		linkedGlyph = b.getLinkedGlyph(linkedGlyph, src_cursor.Add(dst_offset), inner_dir.Inverted(), dst_canvas, depth)
		if src_cursor == vec.ZERO_COORD {
			linkedGlyph = b.getLinkedGlyph(linkedGlyph, src_cursor.Add(dst_offset), inner_dir.RotateCW(), dst_canvas, depth)
		}
		border_segment.DrawGlyph(src_cursor, 0, linkedGlyph)
	}
}

func (b *Border) getLinkedGlyph(src_glyph int, glyph_dst_pos vec.Coord, link_dir vec.Direction, dst_canvas *gfx.Canvas, depth int) int {
	glyph_flags := b.style.getBorderFlags(src_glyph)
	neighbour_pos := glyph_dst_pos.Step(link_dir)
	if dst_canvas.InBounds(neighbour_pos) && dst_canvas.GetDepth(neighbour_pos) == depth {
		neighbour_cell := dst_canvas.GetCell(neighbour_pos)
		if neighbour_cell.Mode == gfx.DRAW_GLYPH {
			if b.style.glyphIsLinkable(neighbour_cell.Glyph) {
				neighbour_cell_flags := b.style.getBorderFlags(neighbour_cell.Glyph)
				if neighbour_cell_flags&getBorderFlagByDirection(link_dir.Inverted()) != 0 {
					glyph_flags |= getBorderFlagByDirection(link_dir)
				}
			}
		} else { // special cases for connecting with title/hint decorations
			if link_dir == vec.DIR_RIGHT {
				if neighbour_cell.Chars[0] == b.style.TextDecorationL || neighbour_cell.Chars[0] == b.style.TextDecorationPad {
					glyph_flags |= BORDER_R
				}
			} else if link_dir == vec.DIR_LEFT {
				if neighbour_cell.Chars[1] == b.style.TextDecorationR || neighbour_cell.Chars[1] == b.style.TextDecorationPad {
					glyph_flags |= BORDER_L
				}
			}
		}
	}

	return b.style.Glyphs[glyph_flags]
}

func (b *Border) EnableScrollbar(height, pos int) {
	b.scrollbar = true
	b.scrollbarContentHeight = height
	b.scrollbarViewportPosition = pos
	b.dirty = true
}

// Updates the position/size of the scrollbar.
// NOTE: this does NOT enable the scrollbar. you have to do that manually during setup.
func (b *Border) UpdateScrollbar(height, pos int) {
	if b.scrollbarContentHeight == height && b.scrollbarViewportPosition == pos {
		return
	}

	b.scrollbarContentHeight = height
	b.scrollbarViewportPosition = pos
	b.dirty = true
}
