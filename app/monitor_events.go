package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/moncho/dry/appui"
	"github.com/moncho/dry/ui"
)

type monitorScreenEventHandler struct {
	baseEventHandler
	widget *appui.Monitor
}

func (h *monitorScreenEventHandler) handle(event *tcell.EventKey, f func(eventHandler)) {
	handled := false
	cursor := h.dry.screen.Cursor()
	switch event.Key() {
	case tcell.KeyF1:
		handled = true
		h.widget.Sort()
		h.widget.OnEvent(nil)
	case tcell.KeyUp: //cursor up
		handled = true
		cursor.ScrollCursorUp()
		h.widget.OnEvent(nil)
	case tcell.KeyDown: // cursor down
		handled = true
		cursor.ScrollCursorDown()
		h.widget.OnEvent(nil)
	case tcell.KeyEnter: //Container menu
		showMenu := func(id string) error {
			h.widget.Unmount()
			h.dry.screen.Cursor().Reset()
			widgets.ContainerMenu.ForContainer(id)
			widgets.ContainerMenu.OnUnmount = func() error {
				h.dry.screen.Cursor().Reset()
				h.dry.changeView(Monitor)
				f(h)
				return refreshScreen()
			}
			h.dry.changeView(ContainerMenu)
			f(viewsToHandlers[ContainerMenu])
			return refreshScreen()
		}
		if err := h.widget.OnEvent(showMenu); err != nil {
			h.dry.message(err.Error())
		}
	}
	if !handled {
		switch event.Rune() {
		case 'g': //Cursor to the top
			handled = true
			cursor.Reset()
			h.widget.OnEvent(nil)

		case 'G': //Cursor to the bottom
			handled = true
			cursor.Bottom()
			h.widget.OnEvent(nil)
		case 's': // Set the delay between updates to <delay> seconds.
			//widget is mounted on render, dont Mount here
			h.widget.Unmount()
			prompt := appui.NewPrompt(
				h.dry.screen.Dimensions(),
				"Set the delay between updates (in milliseconds)")
			widgets.add(prompt)
			forwarder := newEventForwarder()
			f(forwarder)
			h.dry.changeView(NoView)
			refreshScreen()
			go func() {
				defer h.dry.changeView(Monitor)
				defer f(h)
				events := ui.EventSource{
					Events: forwarder.events(),
					EventHandledCallback: func(e *tcell.EventKey) error {
						return refreshScreen()
					},
				}
				prompt.OnFocus(events)
				input, cancel := prompt.Text()
				widgets.remove(prompt)
				if cancel {
					return
				}
				refreshRate, err := toInt(input)
				if err != nil {
					h.dry.message(
						fmt.Sprintf("Error setting refresh rate: %s", err.Error()))
					return
				}
				h.widget.RefreshRate(refreshRate)
			}()
		}
	}
	if !handled {
		h.baseEventHandler.handle(event, func(eh eventHandler) {
			h.widget.Unmount()
			f(eh)
		})
	}
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return -1, errors.New("be nice, a number is expected")
	}
	if i < 0 {
		return -1, errors.New("negative values are not allowed")
	}
	return i, nil
}
