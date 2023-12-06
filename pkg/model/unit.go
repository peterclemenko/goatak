package model

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kdudkov/goatak/pkg/cot"
)

const (
	StaleContactDelete = time.Hour * 12
	POINT              = "point"
	UNIT               = "unit"
	CONTACT            = "contact"
	MaxTrackPoints     = 5000
)

type Item struct {
	mx             sync.RWMutex
	uid            string
	cottype        string
	class          string
	callsign       string
	staleTime      time.Time
	startTime      time.Time
	sendTime       time.Time
	lastSeen       time.Time
	online         bool
	local          bool
	send           bool
	parentCallsign string
	parentUID      string
	color          string
	icon           string
	track          []*Pos
	msg            *cot.CotMessage
}

func (i *Item) String() string {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return fmt.Sprintf("%s: %s %s %s", i.class, i.uid, i.cottype, i.callsign)
}

func (i *Item) GetClass() string {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.class
}

func (i *Item) GetCotType() string {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.cottype
}

func (i *Item) GetMsg() *cot.CotMessage {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.msg
}

func (i *Item) GetUID() string {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.uid
}

func (i *Item) GetCallsign() string {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.callsign
}

func (i *Item) GetLastSeen() time.Time {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.lastSeen
}

func (i *Item) GetStartTime() time.Time {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.startTime
}

func (i *Item) IsOld() bool {
	i.mx.RLock()
	defer i.mx.RUnlock()

	switch i.class {
	case CONTACT:
		return (!i.online) && i.lastSeen.Add(StaleContactDelete).Before(time.Now())
	default:
		return i.staleTime.Before(time.Now())
	}
}

func (i *Item) IsOnline() bool {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.online
}

func (i *Item) SetOffline() {
	i.mx.Lock()
	defer i.mx.Unlock()
	i.online = false
}

func (i *Item) SetOnline() {
	i.mx.Lock()
	defer i.mx.Unlock()
	i.online = true
	i.lastSeen = time.Now()
}

func (i *Item) SetLocal(local, send bool) {
	i.mx.Lock()
	defer i.mx.Unlock()
	i.local = local
	i.send = send
}

func (i *Item) IsSend() bool {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.send
}

func GetClass(msg *cot.CotMessage) string {
	if msg == nil {
		return ""
	}

	t := msg.GetType()

	switch {
	case strings.HasPrefix(t, "a-"):
		if msg.IsContact() {
			return CONTACT
		} else {
			return UNIT
		}
	case strings.HasPrefix(t, "b-"):
		return POINT
	}

	return ""
}

func FromMsg(msg *cot.CotMessage) *Item {
	cls := GetClass(msg)

	if cls == "" {
		return nil
	}

	parent, parentCs := msg.GetParent()

	i := &Item{
		mx:             sync.RWMutex{},
		uid:            msg.GetUid(),
		cottype:        msg.GetType(),
		class:          cls,
		callsign:       msg.GetCallsign(),
		staleTime:      msg.GetStaleTime(),
		startTime:      msg.GetStartTime(),
		sendTime:       msg.GetSendTime(),
		lastSeen:       time.Now(),
		online:         true,
		local:          false,
		send:           false,
		parentCallsign: parentCs,
		parentUID:      parent,
		color:          msg.Detail.GetFirst("color").GetAttr("argb"),
		icon:           msg.Detail.GetFirst("usericon").GetAttr("iconsetpath"),
		track:          nil,
		msg:            msg,
	}

	if i.class == UNIT || i.class == CONTACT {
		if msg.GetLat() != 0 || msg.GetLon() != 0 {
			pos := &Pos{
				time:  msg.GetSendTime(),
				lat:   msg.GetLat(),
				lon:   msg.GetLon(),
				speed: msg.TakMessage.GetCotEvent().GetDetail().GetTrack().GetSpeed(),
				mx:    sync.RWMutex{},
			}

			i.track = []*Pos{pos}
		}
	}

	return i
}

func FromMsgLocal(msg *cot.CotMessage, send bool) *Item {
	i := FromMsg(msg)
	i.local = true
	i.send = send

	return i
}

func (i *Item) GetLanLon() (float64, float64) {
	return i.msg.GetLat(), i.msg.GetLon()
}

func (i *Item) Update(msg *cot.CotMessage) {
	if msg == nil {
		i.SetOnline()

		return
	}

	i.mx.Lock()
	defer i.mx.Unlock()

	i.class = GetClass(msg)
	i.cottype = msg.GetType()
	i.callsign = msg.GetCallsign()
	i.staleTime = msg.GetStaleTime()
	i.startTime = msg.GetStartTime()
	i.sendTime = msg.GetSendTime()
	i.msg = msg
	i.lastSeen = time.Now()

	i.parentUID, i.parentCallsign = msg.GetParent()

	if c := msg.Detail.GetFirst("color"); c != nil {
		i.color = c.GetAttr("argb")
	}

	i.icon = msg.Detail.GetFirst("usericon").GetAttr("iconsetpath")

	if i.class == UNIT || i.class == CONTACT {
		i.online = true

		if msg.GetLat() != 0 || msg.GetLon() != 0 {
			pos := &Pos{
				time:  msg.GetSendTime(),
				lat:   msg.GetLat(),
				lon:   msg.GetLon(),
				speed: msg.TakMessage.GetCotEvent().GetDetail().GetTrack().GetSpeed(),
				mx:    sync.RWMutex{},
			}

			i.track = append(i.track, pos)
			if len(i.track) > MaxTrackPoints {
				i.track = i.track[len(i.track)-MaxTrackPoints:]
			}
		}
	}
}

func (i *Item) GetTrack() []*Pos {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.track
}

func (i *Item) UpdateFromWeb(w *WebUnit, m *cot.CotMessage) {
	if w == nil {
		return
	}

	i.mx.Lock()
	defer i.mx.Unlock()

	i.class = w.Category
	i.cottype = w.Type
	i.callsign = w.Callsign
	i.staleTime = w.StaleTime
	i.startTime = w.StartTime
	i.sendTime = w.SendTime
	i.lastSeen = time.Now()
	i.parentUID = w.ParentUID
	i.parentCallsign = w.ParentCallsign
	i.icon = w.Icon
	i.color = w.Color
	i.local = w.Local
	i.send = w.Send

	i.msg = m
}
