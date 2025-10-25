package page

import "github.com/rivo/tview"

const (
	// Auth страница аутентификации.
	Auth = "Auth"
	// UserDataList страница пользовательских данных.
	UserDataList = "UserDataList"
)

// Page страница.
type Page interface {
	Component() tview.Primitive
	Render() tview.Primitive
}

// Pages страницы.
type Pages struct {
	TviewPages *tview.Pages
	pages      map[string]Page
}

// AddPage добавляет страницу.
func (p *Pages) AddPage(key string, page Page, resize, visible bool) {
	p.TviewPages.AddPage(key, page.Component(), resize, visible)
	if visible {
		page.Render()
	}
	p.pages[key] = page
}

// SwitchToPage переключает страницы.
func (p *Pages) SwitchToPage(page string) {
	e, ok := p.pages[page]
	if ok {
		e.Render()
		p.TviewPages.SwitchToPage(page)
	}
}

// NewPages конструктор.
func NewPages() *Pages {
	return &Pages{
		TviewPages: tview.NewPages(),
		pages:      make(map[string]Page),
	}
}
