package auth

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ktigay/goph-keeper/internal/client/entity"
)

// Callbacks callback события.
type Callbacks struct {
	OnSignIn func(entity.Credentials) error
	OnSignUp func(entity.Credentials) error
}

// Page структура страницы.
type Page struct {
	callbacks Callbacks
	cmp       *tview.Grid
}

// Component компонент страницы.
func (l *Page) Component() tview.Primitive {
	return l.cmp
}

// Render рендер.
func (l *Page) Render() tview.Primitive {
	// Create empty Box to pad each side of appGrid
	bx := tview.NewBox()

	e := entity.Credentials{}

	loginLabel := tview.NewTextView().SetText("Credentials:")
	login := tview.NewInputField()
	login.SetChangedFunc(func(text string) {
		e.Login = text
	})

	passLabel := tview.NewTextView().SetText("Password:")
	pass := tview.NewInputField().SetMaskCharacter('*')
	pass.SetChangedFunc(func(text string) {
		e.Password = text
	})

	noticeTxt := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	// style := tcell.Style{}.Background(tcell.ColorNone)

	signIn := tview.NewButton("Sign-in")
	// signIn.SetStyle(style).SetActivatedStyle(style)
	signIn.SetSelectedFunc(func() {
		err := l.callbacks.OnSignIn(e)
		if err != nil {
			noticeTxt.
				SetTextColor(tcell.ColorRed).
				SetText(fmt.Errorf("sign-in failed: %w", err).Error())
		}
		signIn.Blur()
	})

	signUp := tview.NewButton("Sign-up")
	// signUp.SetStyle(style).SetActivatedStyle(style)
	signUp.SetSelectedFunc(func() {
		err := l.callbacks.OnSignUp(e)
		if err != nil {
			noticeTxt.SetTextColor(tcell.ColorRed).SetText(fmt.Errorf("sign-up failed: %w", err).Error())
			return
		}
		noticeTxt.
			SetTextColor(tcell.ColorGreen).
			SetText("sign-up succeeded")
		signUp.Blur()
	})

	// Create Grid containing the application's widgets
	appGrid := l.cmp.
		SetColumns(-1, 16, 26, -1).
		SetRows(-1, 2, 2, 3, 3, -1).
		AddItem(bx, 0, 0, 3, 1, 0, 0, false). // Left - 3 rows
		AddItem(bx, 0, 1, 1, 1, 0, 0, false). // Top - 1 row
		AddItem(bx, 0, 3, 3, 1, 0, 0, false). // Right - 3 rows
		AddItem(bx, 4, 1, 1, 1, 0, 0, false). // Bottom - 1 row
		AddItem(loginLabel, 1, 1, 1, 1, 0, 0, false).
		AddItem(login, 1, 2, 1, 1, 0, 0, false).
		AddItem(passLabel, 2, 1, 1, 1, 0, 0, false).
		AddItem(pass, 2, 2, 1, 1, 0, 0, false).
		AddItem(noticeTxt, 3, 1, 1, 2, 0, 0, false).
		AddItem(signIn, 4, 1, 1, 1, 1, 0, false).
		AddItem(signUp, 4, 2, 1, 1, 1, 0, false)

	appGrid.SetGap(0, 1)

	return appGrid
}

// New конструктор.
func New(c Callbacks) *Page {
	return &Page{
		callbacks: c,
		cmp:       tview.NewGrid(),
	}
}
