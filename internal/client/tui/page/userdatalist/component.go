package userdatalist

import (
	"github.com/rivo/tview"

	"github.com/ktigay/goph-keeper/internal/entity"
)

// ComponentData компонент полей данных.
type ComponentData interface {
	AddToForm(form *tview.Form)
}

// StringComponent текстовый компонент  [entity.DataTypeText].
type StringComponent struct {
	data        string
	fieldWidth  int
	fieldHeight int
	maxLength   int
	changed     func(text any)
}

// AddToForm добавляет компонент в форму.
func (s *StringComponent) AddToForm(form *tview.Form) {
	form.AddTextArea("Data (string)", s.data, s.fieldWidth, s.fieldHeight, s.maxLength, func(text string) {
		s.changed(text)
	})
}

// NewStringComponent конструктор.
func NewStringComponent(data string, fieldWidth, fieldHeight, maxLength int, changed func(text any)) *StringComponent {
	return &StringComponent{
		data:        data,
		fieldWidth:  fieldWidth,
		fieldHeight: fieldHeight,
		maxLength:   maxLength,
		changed:     changed,
	}
}

// BinaryComponent компонент для бинарных данных [entity.DataTypeBinary].
type BinaryComponent struct {
	data        string
	fieldWidth  int
	fieldHeight int
	maxLength   int
	changed     func(text any)
}

// AddToForm добавляет компонент в форму.
func (b *BinaryComponent) AddToForm(form *tview.Form) {
	form.AddTextArea("Data (base64)", b.data, b.fieldWidth, b.fieldHeight, b.maxLength, func(text string) {
		b.changed(text)
	})
}

// NewBinaryComponent конструктор.
func NewBinaryComponent(data string, fieldWidth, fieldHeight, maxLength int, changed func(text any)) *BinaryComponent {
	return &BinaryComponent{
		data:        data,
		fieldWidth:  fieldWidth,
		fieldHeight: fieldHeight,
		maxLength:   maxLength,
		changed:     changed,
	}
}

// CardComponent компонент для [entity.DataTypeCard].
type CardComponent struct {
	data       entity.UserDataCard
	fieldWidth int
	changed    func(any)
}

// AddToForm добавляет компонент в форму.
func (c *CardComponent) AddToForm(form *tview.Form) {
	var number, expMonth, expYear, cvc *tview.InputField

	number = tview.NewInputField()
	expMonth = tview.NewInputField()
	expYear = tview.NewInputField()
	cvc = tview.NewInputField()

	changed := func(_ string) {
		c.data.Number = number.GetText()
		c.data.ExpMonth = expMonth.GetText()
		c.data.ExpYear = expYear.GetText()
		c.data.CVC = cvc.GetText()
		c.changed(c.data)
	}

	number.
		SetLabel("Card Number").
		SetFieldWidth(c.fieldWidth).
		SetText(c.data.Number).
		SetChangedFunc(changed)
	form.AddFormItem(number)

	expMonth.
		SetLabel("ExpMonth").
		SetFieldWidth(c.fieldWidth).
		SetText(c.data.ExpMonth).
		SetChangedFunc(changed)
	form.AddFormItem(expMonth)

	expYear.
		SetLabel("ExpYear").
		SetFieldWidth(c.fieldWidth).
		SetText(c.data.ExpYear).
		SetChangedFunc(changed)
	form.AddFormItem(expYear)

	cvc.
		SetLabel("CVC").
		SetFieldWidth(c.fieldWidth).
		SetText(c.data.CVC).
		SetChangedFunc(changed)
	form.AddFormItem(cvc)
}

// NewCardComponent конструктор.
func NewCardComponent(data entity.UserDataCard, fieldWidth int, changed func(any)) *CardComponent {
	return &CardComponent{
		data:       data,
		fieldWidth: fieldWidth,
		changed:    changed,
	}
}

// ComponentDataFactory фабрика компонентов.
func ComponentDataFactory(data *entity.UserData, fieldWidth, fieldHeight, maxLength int, changed func(any)) ComponentData {
	switch data.Type {
	case entity.DataTypeText:
		return NewStringComponent(data.GetData().(string), fieldWidth, fieldHeight, maxLength, changed)
	case entity.DataTypeBinary:
		return NewBinaryComponent(data.GetData().(string), fieldWidth, fieldHeight, maxLength, changed)
	case entity.DataTypeCard:
		return NewCardComponent(data.GetData().(entity.UserDataCard), fieldWidth, changed)
	}
	return nil
}
