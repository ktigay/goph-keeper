package userdatalist

import (
	"time"

	"github.com/rivo/tview"

	"github.com/ktigay/goph-keeper/internal/entity"
)

// Callbacks callback события.
type Callbacks struct {
	OnItemUpdate  func(entity.UserData) error
	OnItemAdd     func(entity.UserData) error
	OnItemDelete  func(entity.UserData) error
	OnRefreshData func()
	OnQuit        func()
}

// Page структура страницы пользовательских данных.
type Page struct {
	cmp          *tview.Flex
	callbacks    Callbacks
	dataSourceFn func() ([]entity.UserData, error)
	userData     []*entity.UserData
	form         *tview.Form
	metaForm     *tview.Form
	list         *tview.List
	notice       *tview.TextView
	activeIdx    int
}

// RefreshData обновляет данные.
func (u *Page) RefreshData() {
	data, _ := u.dataSourceFn()
	u.userData = make([]*entity.UserData, 0, len(data))
	for _, d := range data {
		u.userData = append(u.userData, &d)
	}
}

// Component возвращает компонент страницы.
func (u *Page) Component() tview.Primitive {
	return u.cmp
}

// Render рендер.
func (u *Page) Render() tview.Primitive {
	u.RefreshData()

	list := u.list

	changed := func(i int, text, secondary string, r rune) {
		u.activeIdx = i
		u.renderForm(u.saveBtnChanged, u.delBtnChanged)
	}
	list.SetSelectedFunc(changed)
	list.SetChangedFunc(changed)

	u.renderMetaForm()
	u.renderList()

	return u.cmp
}

func (u *Page) renderList() *tview.List {
	list := u.list
	list.Clear()
	for _, v := range u.userData {
		list.AddItem(v.Title, "", 0, nil)
	}
	return list
}

func (u *Page) renderForm(save, delete func(userData entity.UserData)) *tview.Form {
	data := u.userData[u.activeIdx]
	form := u.form
	form.Clear(true)

	if data.IsNew && data.Title == "" {
		form.SetTitle("New Record")
	} else {
		form.SetTitle(data.Title)
	}

	form.AddInputField("Title", data.Title, 40, nil, func(text string) {
		data.Title = text
	})

	u.notice = tview.NewTextView().
		SetSize(2, 0)
	form.AddFormItem(u.notice)

	form.AddDropDown("Type", []string{
		string(entity.DataTypeText),
		string(entity.DataTypeBinary),
		string(entity.DataTypeCard),
	},
		func() int {
			switch data.Type {
			case entity.DataTypeBinary:
				return 1
			case entity.DataTypeCard:
				return 2
			default:
				return 0
			}
		}(),
		func(option string, optionIndex int) {
			if string(data.Type) != option {
				data.Type = entity.UserDataType(option)
				u.renderForm(save, delete)
			}
		})

	dataCmp := ComponentDataFactory(data, 80, 10, 0, func(d any) {
		_ = data.SetData(d)
	})
	dataCmp.AddToForm(form)

	if !data.UpdatedAt.IsZero() {
		form.AddTextView("Created At", data.CreatedAt.In(time.Local).Format(time.RFC822), 40, 1, false, false)
		form.AddTextView("Updated At", data.UpdatedAt.In(time.Local).Format(time.RFC822), 40, 1, false, false)
	}

	if data.IsSynced {
		form.AddTextView("Is Synced", "synced", 40, 1, false, false)
	} else {
		form.AddTextView("Is Synced", "not sync", 40, 1, false, false)
	}

	if len(data.MetaData) > 0 {
		form.AddTextView("", "Metadata", 0, 1, false, false)
	}
	for i, m := range data.MetaData {
		field := tview.NewInputField()
		field.SetLabel(m.Title).
			SetFieldWidth(80).
			SetText(m.Value).
			SetChangedFunc(func(val string) {
				data.MetaData[i].Value = val
			})
		form.AddFormItem(field)
	}

	form.AddButton("Remove metadata", func() {
		mLen := len(data.MetaData)
		if mLen > 0 {
			cnt := form.GetFormItemCount()
			data.MetaData = data.MetaData[:len(data.MetaData)-1]
			form.RemoveFormItem(cnt - 1)
		}
	})

	form.AddButton("Save", func() {
		save(*data)
	})
	if delete != nil {
		form.AddButton("Delete", func() {
			delete(*data)
		})
	}

	return form
}

func (u *Page) renderMetaForm() *tview.Form {
	form := u.form
	metaForm := u.metaForm
	metaForm.Clear(true)

	metaForm.AddTextView("", "Add metadata", 0, 1, false, false)
	metaForm.AddInputField("Title", "", 20, nil, func(text string) {})
	metaForm.AddInputField("Value", "", 20, nil, func(text string) {})

	metaForm.AddButton("Add metadata", func() {
		data := u.userData[u.activeIdx]

		cnt := metaForm.GetFormItemCount()
		title := metaForm.GetFormItem(cnt - 2).(*tview.InputField)
		val := metaForm.GetFormItem(cnt - 1).(*tview.InputField)
		if title.GetText() == "" {
			return
		}

		data.MetaData = append(data.MetaData, entity.MetaData{
			Title: title.GetText(),
			Value: val.GetText(),
		})
		i := len(data.MetaData) - 1
		v := &data.MetaData[i]

		if i == 0 {
			form.AddTextView("", "Metadata", 0, 1, false, false)
		}
		form.AddInputField(title.GetText(), val.GetText(), 80, nil, func(val string) {
			v.Value = val
		})
		title.SetText("")
		val.SetText("")
	})
	return metaForm
}

func (u *Page) addNewDataItem(add func(userData entity.UserData)) {
	data := &entity.UserData{
		Title:    "",
		Type:     entity.DataTypeText,
		MetaData: make([]entity.MetaData, 0),
		IsNew:    true,
		IsSynced: false,
	}
	u.userData = append(u.userData, data)
	u.activeIdx = len(u.userData) - 1
	u.renderForm(add, nil)
}

func (u *Page) saveBtnChanged(data entity.UserData) {
	if err := u.callbacks.OnItemUpdate(data); err != nil {
		u.notice.SetText(err.Error())
		return
	}

	u.RefreshData()
	u.renderList().SetCurrentItem(0)
}

func (u *Page) delBtnChanged(data entity.UserData) {
	_ = u.callbacks.OnItemDelete(data)

	u.RefreshData()
	u.renderList().SetCurrentItem(0)
}

// New конструктор.
func New(c Callbacks, dataSource func() ([]entity.UserData, error)) *Page {
	cmp := tview.NewFlex()
	cmp.SetDirection(tview.FlexRow)

	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	cmp.AddItem(flex, 0, 4, false)

	list := tview.NewList().ShowSecondaryText(false)
	list.SetTitle("User Data List").SetBorder(true).SetTitleAlign(tview.AlignCenter)
	flex.AddItem(list, 0, 1, false)

	formFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(formFlex, 0, 4, false)

	form := tview.NewForm()
	form.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	formFlex.AddItem(form, 0, 4, false)

	metaForm := tview.NewForm()
	metaForm.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	formFlex.AddItem(metaForm, 0, 1, false)

	page := Page{
		cmp:          cmp,
		callbacks:    c,
		dataSourceFn: dataSource,
		form:         form,
		metaForm:     metaForm,
		list:         list,
	}

	addBtn := tview.NewButton("New record")
	addBtn.SetSelectedFunc(func() {
		page.addNewDataItem(func(data entity.UserData) {
			if err := page.callbacks.OnItemAdd(data); err != nil {
				page.notice.SetText(err.Error())
				addBtn.Blur()
				return
			}

			page.RefreshData()
			page.renderList().SetCurrentItem(0)
			addBtn.Blur()
		})
	})

	refreshBtn := tview.NewButton("Refresh Data")
	refreshBtn.SetSelectedFunc(func() {
		page.callbacks.OnRefreshData()
		refreshBtn.Blur()
	})

	quitBtn := tview.NewButton("Quit")
	quitBtn.SetSelectedFunc(func() {
		page.callbacks.OnQuit()
		quitBtn.Blur()
	})

	cmp.AddItem(
		tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(addBtn, 20, 1, false).
			AddItem(tview.NewBox(), 1, 1, false).
			AddItem(refreshBtn, 20, 1, false).
			AddItem(tview.NewBox(), 1, 1, false).
			AddItem(quitBtn, 20, 1, false),

		1, 1, false)

	return &page
}
