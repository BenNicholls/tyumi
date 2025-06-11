package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var ACTION_PAGE_NEXT = input.RegisterAction("Next Page")

func init() {
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_PAGE_NEXT, input.K_TAB)
}

// PageElement is the interface defining the elements that can act as pages in the PageContainer. It's satisfied by
// the Page element, which you can embed into a custom element.
type PageElement interface {
	element
	initPage(size vec.Dims, title string)
	activate()
	deactivate()
	getTab() *Textbox
}

// PageContainer contains multiple pages and displays them one at a time, with a familiar tab interface at the top
// for swapping pages.
type PageContainer struct {
	Element

	OnPageChanged func()

	pages            []PageElement
	currentPageIndex int //this is set to -1 on container creation, indicating no pages are selected (since they don't exist yet)
}

func NewPageContainer(size vec.Dims, pos vec.Coord, depth int) (pc *PageContainer) {
	pc = new(PageContainer)
	pc.Init(size, pos, depth)

	return
}

func (pc *PageContainer) Init(size vec.Dims, pos vec.Coord, depth int) {
	pc.Element.Init(size, pos, depth)
	pc.TreeNode.Init(pc)

	pc.pages = make([]PageElement, 0)
	pc.currentPageIndex = -1 //no pages in container, so no selection
}

// creates and adds a new page to the pagecontainer, and returns a reference to the new page for the user to populate
// with other ui stuff
func (pc *PageContainer) CreatePage(title string) *Page {
	newpage := newPage(pc.size.Shrink(0, 3), title)
	pc.addPage(newpage)

	return newpage
}

// Adds a supplied PageElement to the PageContainer. The provided page will be initialized, resized, and repositioned
// here, so don't bother initializing it beforehand. Once added, the page element can be used just like a normal page.
func (pc *PageContainer) AddPage(title string, page PageElement) {
	page.initPage(pc.size.Shrink(0, 3), title)
	pc.addPage(page)
}

// Selects the next page in the container. If at the end, wraps around to the first tab.
func (pc *PageContainer) NextPage() {
	if len(pc.pages) < 2 {
		return
	}

	pc.selectPage(util.CycleClamp(pc.currentPageIndex+1, 0, len(pc.pages)-1))
}

// Selects the previous page in the container. If at the start, wraps around to the last tab.
func (pc *PageContainer) PrevPage() {
	if len(pc.pages) < 2 {
		return
	}

	pc.selectPage(util.CycleClamp(pc.currentPageIndex-1, 0, len(pc.pages)-1))
}

func (pc *PageContainer) addPage(page PageElement) {
	//find position for next tab
	x := 1
	for _, page := range pc.pages {
		x += page.getTab().Size().W + 1
	}
	page.getTab().MoveTo(vec.Coord{x, 1})
	pc.AddChild(page.getTab())
	pc.AddChild(page)
	pc.pages = append(pc.pages, page)

	if len(pc.pages) == 1 { //first page added
		pc.selectPage(0)
	}
}

func (pc *PageContainer) selectPage(page_index int) {
	if page_index == pc.currentPageIndex {
		return
	}

	if page_index < 0 || page_index >= len(pc.pages) {
		log.Error("Bad Page Select! got ", page_index, " number of pages is ", len(pc.pages))
		return
	}

	//remove previous selected page if there is one (index -1 means no page selected)
	if pc.currentPageIndex >= 0 {
		pc.getSelectedPage().deactivate()
	}

	pc.currentPageIndex = page_index
	pc.getSelectedPage().activate()
	fireCallbacks(pc.OnPageChanged)
}

func (pc *PageContainer) getSelectedPage() PageElement {
	if len(pc.pages) == 0 {
		return nil
	}

	return pc.pages[pc.currentPageIndex]
}

// Retrives the index for the current page. If there are no pages, returns -1.
func (pc PageContainer) GetPageIndex() int {
	return pc.currentPageIndex
}

func (pc *PageContainer) renderIfDirty() {
	if len(pc.pages) == 0 {
		return
	}

	//blank out border below selected page's tab
	selectedTab := pc.getSelectedPage().getTab()
	cursor := selectedTab.Bounds().Coord
	cursor.Move(0, 2)
	brush := gfx.NewGlyphVisuals(selectedTab.getBorderStyle().GetGlyph(gfx.LINK_UL), selectedTab.Border.colours)
	pc.DrawVisuals(cursor, BorderDepth, brush)
	brush.Glyph = gfx.GLYPH_NONE
	for range selectedTab.Size().W {
		cursor.Move(1, 0)
		pc.DrawVisuals(cursor, BorderDepth, brush)
	}
	cursor.Move(1, 0)
	brush.Glyph = selectedTab.getBorderStyle().GetGlyph(gfx.LINK_UR)
	pc.DrawVisuals(cursor, BorderDepth, brush)
}

func (pc *PageContainer) HandleAction(action input.ActionID) (action_handled bool) {
	switch action {
	case ACTION_PAGE_NEXT:
		pc.NextPage()
	default:
		return
	}

	return true
}

// Page is the content for a tab in a PageContainer. Size is defined and controlled by the PageContainer.
// Pages are initialized as deactivated (hidden), and will be activated when selected by the page container.
type Page struct {
	Element

	OnActivate   func() //callback called when the page is activated (switched to)
	OnDeactivate func() //callback called when the page is deactivated (switched away from)

	tab *Textbox //textbox for the tab in the pagecontainer
}

func (p *Page) initPage(size vec.Dims, title string) {
	if title == "" {
		title = " "
	}

	p.Element.Init(size, vec.Coord{0, 3}, BorderDepth)
	p.EnableBorder()
	p.Border.SetStyle(BORDER_STYLE_INHERIT)

	p.tab = NewTextbox(vec.Dims{FIT_TEXT, 1}, vec.Coord{1, 1}, 5, title, ALIGN_LEFT)
	p.tab.EnableBorder()
	p.tab.Border.SetStyle(BORDER_STYLE_INHERIT)

	p.deactivate()
}

func newPage(size vec.Dims, title string) (p *Page) {
	p = new(Page)
	p.initPage(size, title)

	return
}

func (p *Page) getTab() *Textbox {
	return p.tab
}

func (p *Page) activate() {
	p.tab.EnableBorder()
	p.Show()
	fireCallbacks(p.OnActivate)
}

func (p *Page) deactivate() {
	p.tab.DisableBorder()
	p.Hide()
	fireCallbacks(p.OnDeactivate)
}

// No-op. Pages cannot be moved relative to their container.
func (p *Page) Move(dx, dy int) {}

// No-op. Pages cannot be moved relative to their container.
func (p *Page) MoveTo(pos vec.Coord) {}
