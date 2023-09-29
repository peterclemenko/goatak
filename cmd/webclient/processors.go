package main

import (
	"encoding/json"
	"github.com/kdudkov/goatak/pkg/cot"
	"github.com/kdudkov/goatak/pkg/model"
)

func (app *App) InitMessageProcessors() {
	app.eventProcessors["t-x-d-d"] = app.removeItemProcessor
	// ping
	app.eventProcessors["t-x-c-t"] = app.justLogProcessor
	// pong
	app.eventProcessors["t-x-c-t-r"] = app.justLogProcessor
	// chat
	app.eventProcessors["b-t-f"] = app.chatProcessor
	app.eventProcessors["b-t-f-"] = app.logInterestingProcessor
	app.eventProcessors["a-"] = app.aProcessor
	app.eventProcessors["b-"] = app.bProcessor
	app.eventProcessors["u-"] = app.logInterestingProcessor
	// video feed
	app.eventProcessors["b-i-v"] = app.logInterestingProcessor
	// photo
	app.eventProcessors["b-f-t-r"] = app.logInterestingProcessor
	// b-r-f-h-c casevac
	app.eventProcessors["b-r-f-h-c"] = app.logInterestingProcessor

	// u-rb-a Range & Bearing – Line
	// u-r-b-c-c R&b - Circle
	// u-d-c-c Drawing Shapes – Circle
	// u-d-r Drawing Shapes – Rectangle
	// u-d-f Drawing Shapes - Free Form
	// u-d-c-e Drawing Shapes – Ellipse
}

func (app *App) GetProcessor(t string) (string, EventProcessor) {
	var found string
	for k, v := range app.eventProcessors {
		if k == t {
			return k, v
		}
		if cot.MatchPattern(t, k) && len(k) > len(found) {
			found = k
		}
	}

	if found != "" {
		return found, app.eventProcessors[found]
	}

	return "", nil
}

func (app *App) justLogProcessor(msg *cot.CotMessage) {
	app.Logger.Debugf("%s %s", msg.GetType(), msg.GetUid())
}

func (app *App) logInterestingProcessor(msg *cot.CotMessage) {
	b, err := json.Marshal(msg.TakMessage)
	if err == nil {
		app.Logger.Info(string(b))
		app.Logger.Info(msg.TakMessage.GetCotEvent().GetDetail().GetXmlDetail())
	}
}

func (app *App) removeItemProcessor(msg *cot.CotMessage) {
	// t-x-d-d
	if msg.Detail != nil && msg.Detail.Has("link") {
		uid := msg.Detail.GetFirst("link").GetAttr("uid")
		typ := msg.Detail.GetFirst("link").GetAttr("type")
		if uid == "" {
			app.Logger.Warnf("invalid remove message: %s", msg.Detail)
			return
		}
		if v := app.items.Get(uid); v != nil {
			switch v.GetClass() {
			case model.CONTACT:
				app.Logger.Debugf("remove %s by message", uid)
				v.SetOffline()
				app.processChange(v)
				return
			case model.UNIT, model.POINT:
				app.Logger.Debugf("remove unit/point %s type %s by message", uid, typ)
				//app.units.Delete(uid)
				return
			}
		}
	}
}

func (app *App) chatProcessor(msg *cot.CotMessage) {
	c := model.MsgToChat(msg)
	if c == nil {
		app.Logger.Errorf("invalid chat message %s", msg.TakMessage)
		return
	}
	if c.From == "" {
		c.From = app.items.GetCallsign(c.FromUid)
	}
	app.Logger.Infof("%s", c)
	app.messages.Add(c)
}

func (app *App) aProcessor(msg *cot.CotMessage) {
	if msg.GetUid() != app.uid {
		app.ProcessItem(msg)
	}
}

func (app *App) bProcessor(msg *cot.CotMessage) {
	if uid, _ := msg.GetParent(); uid != app.uid {

		name := "point"

		switch msg.GetType() {
		case "b-i":
			name = "document"
		case "b-i-x-i":
			name = "photo"
		case "b-m-p-w":
			name = "checkpoint"
		case "b-m-p-w-GOTO":
			name = "flag"
		case "b-m-p-a":
			name = "aimpoint"
		case "b-m-p-c-ip":
			name = "initial point"
		case "b-m-p-c-cp":
			name = "contact point"
		case "b-m-p-c-z":
			name = "black triangle"
		case "b-m-p-s-p-op":
			name = "Observer point"
		case "b-m-p-s-p-loc":
			name = "sensor"
		case "b-m-p-s-p-i":
			name = "target point"
		}
		app.Logger.Debugf("%s %s (%s) %s", name, msg.GetUid(), msg.GetCallsign(), msg.GetType())
		app.ProcessItem(msg)
	}
}
