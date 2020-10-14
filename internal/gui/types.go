package gui

import (
	"io/ioutil"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/takeshixx/deen/internal/plugins"
)

// DeenGUI represents a GUI instance.
type DeenGUI struct {
	App                  fyne.App
	MainWindow           fyne.Window
	MainLayout           *fyne.Container
	PluginList           *widget.ScrollContainer
	Plugins              []string
	EncoderWidgetsScroll *widget.ScrollContainer
	EncoderWidgets       *widget.Box
	Encoders             []*DeenEncoder
	HistoryList          *widget.Group
	History              []string
	CurrentFocus         int // The index of the encoder widget in Encoders
}

// NewDeenGUI initializes a new DeenGUI instance.
func NewDeenGUI() (dg *DeenGUI, err error) {
	dg = &DeenGUI{
		App: app.NewWithID("io.deen.app"),
		PluginList: widget.NewScrollContainer(
			widget.NewAccordionContainer(),
		),
		EncoderWidgetsScroll: widget.NewScrollContainer(widget.NewVBox()),
		EncoderWidgets:       widget.NewVBox(),
		HistoryList:          widget.NewGroup("History"),
	}
	dg.MainWindow = dg.App.NewWindow("deen")
	dg.newMainLayout()
	dg.newMainMenu()
	dg.loadPluginList()

	// Create the root encoder widget (must always exist)
	if _, err = dg.AddEncoder(); err != nil {
		return
	}

	// Setup the theme
	if dg.App.Preferences().String("theme") == "light" {
		dg.App.Settings().SetTheme(theme.LightTheme())
	} else {
		dg.App.Settings().SetTheme(theme.DarkTheme())
	}

	dg.MainWindow.SetMaster()
	dg.MainWindow.SetContent(dg.MainLayout)
	dg.MainWindow.Resize(fyne.NewSize(640, 480))
	dg.addCustomShortcuts()
	dg.updateGUI()
	return
}

// Run is the main function that should
// be called to run the GUI. This will
// block until the GUI is closed.
func (dg *DeenGUI) Run() {
	dg.MainWindow.ShowAndRun()
}

func (dg *DeenGUI) newMainLayout() {
	dg.MainLayout = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, dg.PluginList, dg.HistoryList),
		dg.PluginList,           // left
		dg.HistoryList,          // right
		dg.EncoderWidgetsScroll, // middle
	)
}

func (dg *DeenGUI) newMainMenu() {
	dg.MainWindow.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("File",
				fyne.NewMenuItem("Open", func() {
					fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
						if err == nil && reader == nil {
							return
						}
						if err != nil {
							dialog.ShowError(err, dg.MainWindow)
							return
						}
						dg.fileOpened(reader)
					}, dg.MainWindow)
					fd.Show()
				}),
				// A quit item will be appended to our first menu
			),
			fyne.NewMenu("Theme",
				fyne.NewMenuItem("Light", func() {
					dg.App.Settings().SetTheme(theme.LightTheme())
					dg.App.Preferences().SetString("theme", "light")
				}),
				fyne.NewMenuItem("Dark", func() {
					dg.App.Settings().SetTheme(theme.DarkTheme())
					dg.App.Preferences().SetString("theme", "dark")
				}),
			),
			fyne.NewMenu("Help",
				fyne.NewMenuItem("About", func() {
					dialog.ShowInformation("About", "deen is a DEcoding/ENcoding application that processes arbitrary input data with a wide range of plugins.", dg.MainWindow)
				}),
			)))
}

// Populate the DeenGUI.PluginList field
func (dg *DeenGUI) loadPluginList() {
	dg.Plugins = []string{}
	pluginList := widget.NewAccordionContainer()
	var pluginGroup *widget.AccordionItem
	for _, c := range plugins.PluginCategories {
		filteredPlugins := plugins.GetForCategory(c, false)
		var groupList *widget.Box
		groupList = widget.NewVBox()
		for _, p := range filteredPlugins {
			pluginName := p
			groupList.Append(widget.NewButton(p, func() {
				dg.RunPlugin(pluginName)
			}))
		}
		allPlugins := plugins.GetForCategory(c, true)
		for _, p := range allPlugins {
			pluginName := p
			dg.Plugins = append(dg.Plugins, pluginName)
		}

		pluginGroup = widget.NewAccordionItem(c, groupList)
		pluginList.Append(pluginGroup)
	}
	dg.PluginList = widget.NewScrollContainer(pluginList)
	// Ensure that the scroll container is wide enough
	dg.PluginList.SetMinSize(fyne.NewSize(pluginList.MinSize().Width, 0))
	dg.PluginList.Refresh()
	return
}

func (dg *DeenGUI) addCustomShortcuts() {
	f2 := desktop.CustomShortcut{KeyName: fyne.KeyF2}
	dg.MainWindow.Canvas().AddShortcut(&f2, func(shortcut fyne.Shortcut) {
		// Show fuzzy search
		dg.showPluginSearch()
	})
}

// Reprocess all encoder widgets and update the GUI elements.
func (dg *DeenGUI) updateGUI() {
	log.Println("[DEBUG] Updating GUI")
	// We should only start processing
	// when at least the root widget
	// has a plugin set.
	if dg.Encoders[0].Plugin != nil {
		// We have to process all
		// encoders before creating
		// the GUI layouts.
		dg.processChain()
	}

	dg.updateEncoderWidgets()
	dg.EncoderWidgetsScroll = widget.NewScrollContainer(dg.EncoderWidgets)
	dg.newMainLayout()

	dg.MainWindow.SetContent(dg.MainLayout)
	dg.EncoderWidgetsScroll.ScrollToBottom()
	// Always set focus to the newest encoder.
	dg.SetEncoderFocus(len(dg.Encoders) - 1)
	return
}

func (dg *DeenGUI) updateEncoderWidgets() {
	dg.EncoderWidgets = widget.NewVBox()
	dg.HistoryList = widget.NewGroup("History")
	var historyName string
	for _, e := range dg.Encoders {
		dg.EncoderWidgets.Append(e.createLayout())
		if e.Plugin != nil {
			if e.Plugin.Unprocess {
				historyName = "." + e.Plugin.Name
			} else {
				historyName = e.Plugin.Name
			}
			dg.HistoryList.Append(widget.NewLabel(historyName))
		}
	}
}

func (dg *DeenGUI) fileOpened(f fyne.URIReadCloser) {
	input, err := ioutil.ReadAll(f)
	if err != nil {
		dialog.ShowError(err, dg.MainWindow)
	}
	dg.Encoders[0].SetContent(input)
}
