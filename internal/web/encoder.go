// +build js,wasm

package web

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/takeshixx/deen/internal/plugins"
	"github.com/takeshixx/deen/pkg/types"
)

type EncoderWidget struct {
	vecty.Core
	Parent     *DeenWeb
	Content    string
	ContentLen vecty.MarkupOrChild
	Plugin     *types.DeenPlugin
	Border     bool
}

func (e *EncoderWidget) Render() vecty.ComponentOrHTML {
	var m vecty.Applyer
	if e.Border {
		m = vecty.Style("border", "3px solid blue")
	} else {
		m = vecty.Style("border", "1px dotted black")
	}
	w := elem.Div(
		vecty.Markup(
			vecty.Style("margin-bottom", "15px"),
		),
		e.CreatePluginSelects(),
		elem.TextArea(
			vecty.Markup(
				m,
				vecty.Style("font-family", "monospace"),
				vecty.Style("width", "80%"),
				vecty.Style("display", "inline-block"),
				vecty.Property("rows", 20),
				event.Input(func(event *vecty.Event) {
					e.Content = event.Target.Get("value").String()
					e.Parent.RunChainFrom(e)
					vecty.Rerender(e)
				}),
				event.Click(func(event *vecty.Event) {
					e.Parent.SetCurrentEncoder(e)
				}),
			),
			vecty.Text(e.Content),
		),
	)
	return w
}

func (e *EncoderWidget) SetContent(data string) {
	e.Content = data
	e.Render()
}

func (e *EncoderWidget) ClearContent() {
	e.Content = ""
	e.Render()
}

func (e *EncoderWidget) Process() (processed []byte, err error) {
	if e.Plugin == nil {
		err = fmt.Errorf("No plugin set")
		return
	}
	var reader io.Reader
	if len(e.Content) > 1 {
		reader = strings.NewReader(e.Content)
	}
	if e.Plugin.ProcessDeenTaskFunc != nil {
		var outWriter bytes.Buffer
		task := types.NewDeenTask(&outWriter)
		task.Reader = reader
		if e.Plugin.Unprocess {
			e.Plugin.UnprocessDeenTaskFunc(task)
		} else {
			e.Plugin.ProcessDeenTaskFunc(task)
		}
		select {
		case err = <-task.ErrChan:
		case <-task.DoneChan:
		}
		processed = outWriter.Bytes()
	} else {
		if e.Plugin.Unprocess {
			processed, err = e.Plugin.UnprocessStreamFunc(reader)
		} else {
			processed, err = e.Plugin.ProcessStreamFunc(reader)
		}
	}
	return
}

func (e *EncoderWidget) CreatePluginSelects() vecty.ComponentOrHTML {
	var selectOptions []vecty.MarkupOrChild
	selectOptions = append(selectOptions, vecty.Markup(
		vecty.Style("display", "inline-block"),
		vecty.Style("list-style-type", "none"),
		vecty.Style("padding", "15px"),
		vecty.Style("vertical-align", "top"),
	))
	for _, c := range plugins.PluginCategories {
		filteredPlugins := plugins.GetForCategory(c, false)
		var options []vecty.MarkupOrChild
		options = append(options, elem.Option(
			vecty.Markup(
				vecty.Attribute("disabled", nil),
			),
			vecty.Text(c),
		))
		for _, p := range filteredPlugins {
			pluginName := p
			options = append(options, elem.Option(
				vecty.Markup(
					event.Click(func(v *vecty.Event) {
						e.Parent.RunPlugin(pluginName)
					}),
				),
				vecty.Text(pluginName),
			))
		}
		selectOptions = append(selectOptions, elem.ListItem(
			vecty.Markup(
				vecty.Style("display", "block"),
				vecty.Style("margin", "10px 0 10px 0"),
			),
			elem.Select(options...)),
		)
	}
	selectOptions = append(selectOptions, elem.ListItem(e.CreateEncoderInfo()))
	w := elem.UnorderedList(selectOptions...)
	return w
}

func (e *EncoderWidget) CreateEncoderInfo() vecty.MarkupOrChild {
	e.ContentLen = elem.Label(
		vecty.Text(fmt.Sprintf("Len: %d", len(e.Content))),
	)
	return e.ContentLen
}

func NewEncoderWidget(parent *DeenWeb) (e *EncoderWidget) {
	e = &EncoderWidget{
		Parent: parent,
		Border: false,
	}
	return
}
