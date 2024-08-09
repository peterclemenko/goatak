package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"path/filepath"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"

	"github.com/aofei/air"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kdudkov/goatak/cmd/goatak_server/tak_ws"
	"github.com/kdudkov/goatak/internal/client"
	"github.com/kdudkov/goatak/internal/model"
	"github.com/kdudkov/goatak/internal/wshandler"
	"github.com/kdudkov/goatak/pkg/cot"
	"github.com/kdudkov/goatak/staticfiles"
)

func getAdminAPI(app *App, addr string, renderer *staticfiles.Renderer, webtakRoot string) *air.Air {
	adminAPI := air.New()
	adminAPI.Address = addr
	adminAPI.NotFoundHandler = getNotFoundHandler()

	staticfiles.EmbedFiles(adminAPI, "/static")
	adminAPI.GET("/", getIndexHandler(app, renderer))
	adminAPI.GET("/points", getPointsHandler(app, renderer))
	adminAPI.GET("/map", getMapHandler(app, renderer))
	adminAPI.GET("/missions", getMissionsPageHandler(app, renderer))
	adminAPI.GET("/packages", getMPPageHandler(app, renderer))
	adminAPI.GET("/config", getConfigHandler(app))
	adminAPI.GET("/connections", getConnHandler(app))

	adminAPI.GET("/unit", getUnitsHandler(app))
	adminAPI.GET("/unit/:uid/track", getUnitTrackHandler(app))
	adminAPI.DELETE("/unit/:uid", deleteItemHandler(app))
	adminAPI.GET("/message", getMessagesHandler(app))

	adminAPI.GET("/ws", getWsHandler(app))
	adminAPI.GET("/takproto/1", getTakWsHandler(app))
	adminAPI.POST("/cot", getCotPostHandler(app))
	adminAPI.POST("/cot_xml", getCotXMLPostHandler(app))

	adminAPI.GET("/mp", getAllMissionPackagesHandler(app))
	adminAPI.GET("/mp/:uid", getPackageHandler(app))

	if app.missions != nil {
		adminAPI.GET("/mission", getAllMissionHandler(app))
	}

	if webtakRoot != "" {
		adminAPI.FILE("/webtak/", filepath.Join(webtakRoot, "index.html"))
		adminAPI.FILES("/webtak", webtakRoot)
		addMartiRoutes(app, adminAPI)
	}

	adminAPI.GET("/stack", getStackHandler())
	adminAPI.GET("/metrics", getMetricsHandler())

	adminAPI.RendererTemplateLeftDelim = "[["
	adminAPI.RendererTemplateRightDelim = "]]"

	return adminAPI
}

func getIndexHandler(app *App, r *staticfiles.Renderer) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		data := map[string]any{
			"theme": "auto",
			"page":  " dash",
			"js":    []string{"util.js", "main.js"},
		}

		s, err := r.Render(data, "index.html", "menu.html", "header.html")
		if err != nil {
			app.logger.Error("error", "error", err)
			_ = res.WriteString(err.Error())

			return err
		}

		return res.WriteHTML(s)
	}
}

func getPointsHandler(app *App, r *staticfiles.Renderer) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		data := map[string]any{
			"theme": "auto",
			"page":  " points",
			"js":    []string{"util.js", "points.js"},
		}

		s, err := r.Render(data, "points.html", "menu.html", "header.html")
		if err != nil {
			app.logger.Error("error", "error", err)
			_ = res.WriteString(err.Error())

			return err
		}

		return res.WriteHTML(s)
	}
}

func getMapHandler(app *App, r *staticfiles.Renderer) air.Handler {
	return func(_ *air.Request, res *air.Response) error {
		data := map[string]any{
			"theme": "auto",
			"js":    []string{"util.js", "map.js"},
		}

		s, err := r.Render(data, "map.html", "header.html")
		if err != nil {
			app.logger.Error("error", "error", err)
			_ = res.WriteString(err.Error())

			return err
		}

		return res.WriteHTML(s)
	}
}

func getMissionsPageHandler(app *App, r *staticfiles.Renderer) air.Handler {
	return func(_ *air.Request, res *air.Response) error {
		data := map[string]any{
			"theme": "auto",
			"page":  " missions",
			"js":    []string{"missions.js"},
		}

		s, err := r.Render(data, "missions.html", "menu.html", "header.html")
		if err != nil {
			app.logger.Error("error", "error", err)
			_ = res.WriteString(err.Error())

			return err
		}

		return res.WriteHTML(s)
	}
}

func getMPPageHandler(app *App, r *staticfiles.Renderer) air.Handler {
	return func(_ *air.Request, res *air.Response) error {
		data := map[string]any{
			"theme": "auto",
			"page":  " mp",
			"js":    []string{"mp.js"},
		}

		s, err := r.Render(data, "mp.html", "menu.html", "header.html")
		if err != nil {
			app.logger.Error("error", "error", err)
			_ = res.WriteString(err.Error())

			return err
		}

		return res.WriteHTML(s)
	}
}

func getNotFoundHandler() air.Handler {
	return func(_ *air.Request, res *air.Response) error {
		res.Status = http.StatusNotFound

		return errors.New(http.StatusText(res.Status))
	}
}

func getConfigHandler(app *App) air.Handler {
	m := make(map[string]any, 0)
	m["lat"] = app.lat
	m["lon"] = app.lon
	m["zoom"] = app.zoom
	m["version"] = getVersion()

	m["layers"] = getDefaultLayers()

	return func(_ *air.Request, res *air.Response) error {
		return res.WriteJSON(m)
	}
}

func getUnitsHandler(app *App) air.Handler {
	return func(_ *air.Request, res *air.Response) error {
		return res.WriteJSON(getUnits(app))
	}
}

func getMessagesHandler(app *App) air.Handler {
	return func(_ *air.Request, res *air.Response) error {
		return res.WriteJSON(app.messages)
	}
}

func getStackHandler() air.Handler {
	return func(_ *air.Request, res *air.Response) error {
		return pprof.Lookup("goroutine").WriteTo(res.Body, 1)
	}
}

func getMetricsHandler() air.Handler {
	h := promhttp.Handler()

	return func(req *air.Request, res *air.Response) error {
		h.ServeHTTP(res.HTTPResponseWriter(), req.HTTPRequest())

		return nil
	}
}

func getUnitTrackHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		uid := getStringParam(req, "uid")

		item := app.items.Get(uid)
		if item == nil {
			res.Status = http.StatusNotFound

			return nil
		}

		return res.WriteJSON(item.GetTrack())
	}
}

func deleteItemHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		uid := getStringParam(req, "uid")
		app.items.Remove(uid)

		r := make(map[string]any, 0)
		r["units"] = getUnits(app)
		r["messages"] = app.messages

		return res.WriteJSON(r)
	}
}

func getConnHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		conn := make([]*Connection, 0)

		app.ForAllClients(func(ch client.ClientHandler) bool {
			c := &Connection{
				Uids:     ch.GetUids(),
				User:     ch.GetUser().GetLogin(),
				Ver:      ch.GetVersion(),
				Addr:     ch.GetName(),
				Scope:    ch.GetUser().GetScope(),
				LastSeen: ch.GetLastSeen(),
			}
			conn = append(conn, c)

			return true
		})

		sort.Slice(conn, func(i, j int) bool {
			return conn[i].Addr < conn[j].Addr
		})

		return res.WriteJSON(conn)
	}
}

func getCotPostHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		c := new(cot.CotMessage)

		dec := json.NewDecoder(req.Body)

		if err := dec.Decode(c); err != nil {
			app.logger.Error("cot decode error", "error", err)

			return err
		}

		app.NewCotMessage(c)

		return nil
	}
}

func getCotXMLPostHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		scope := getStringParam(req, "scope")
		if scope == "" {
			scope = "test"
		}

		ev := new(cot.Event)

		dec := xml.NewDecoder(req.Body)

		if err := dec.Decode(ev); err != nil {
			app.logger.Error("cot decode error", "error", err)

			return err
		}

		c, err := cot.EventToProto(ev)
		if err != nil {
			app.logger.Error("cot convert error", "error", err)

			return err
		}

		c.Scope = scope
		app.NewCotMessage(c)

		return nil
	}
}

func getAllMissionHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		data := app.missions.GetAllMissionsAdm()

		result := make([]*model.MissionDTO, len(data))

		for i, m := range data {
			result[i] = model.ToMissionDTOAdm(m, app.packageManager)
		}

		return res.WriteJSON(result)
	}
}

func getAllMissionPackagesHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		data := app.packageManager.GetList(nil)

		return res.WriteJSON(data)
	}
}

func getPackageHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		pi := app.packageManager.Get(getStringParam(req, "uid"))

		if pi == nil {
			res.Status = http.StatusNotFound

			return nil
		}

		f, err := app.packageManager.GetFile(pi.Hash)

		if err != nil {
			app.logger.Error("get file error", "error", err)
			return err
		}

		defer f.Close()

		res.Header.Set("Content-Type", pi.MIMEType)

		if !strings.HasPrefix(pi.MIMEType, "image/") {
			fn := pi.Name
			if pi.MIMEType == "application/x-zip-compressed" && !strings.HasSuffix(fn, ".zip") {
				fn += ".zip"
			}
			res.Header.Set("Content-Disposition", "attachment; filename="+fn)
		}

		res.Header.Set("Last-Modified", pi.SubmissionDateTime.UTC().Format(http.TimeFormat))
		res.Header.Set("Content-Length", strconv.Itoa(pi.Size))

		return res.Write(f)
	}
}

func getWsHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		ws, err := res.WebSocket()
		if err != nil {
			return err
		}

		name := uuid.NewString()

		h := wshandler.NewHandler(name, ws)

		app.logger.Debug("ws listener connected")
		app.changeCb.SubscribeNamed(name, h.SendItem)
		app.deleteCb.SubscribeNamed(name, h.DeleteItem)
		h.Listen()
		app.logger.Debug("ws listener disconnected")

		return nil
	}
}

// handler for WebTAK client - sends/receives protobuf COTs
func getTakWsHandler(app *App) air.Handler {
	return func(req *air.Request, res *air.Response) error {
		ws, err := res.WebSocket()
		if err != nil {
			return err
		}

		defer ws.Close()

		app.logger.Info("WS connection from " + req.ClientAddress())
		name := "ws:" + req.ClientAddress()
		w := tak_ws.New(name, nil, ws, app.NewCotMessage)

		app.AddClientHandler(w)
		w.Listen()
		app.logger.Info("ws disconnected")
		app.RemoveClientHandler(w.GetName())

		return nil
	}
}

func getDefaultLayers() []map[string]any {
	return []map[string]any{
		{
			"name":    "OSM",
			"url":     "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
			"maxzoom": 19,
		},
		{
			"name":    "Opentopo.cz",
			"url":     "https://tile-{s}.opentopomap.cz/{z}/{x}/{y}.png",
			"maxzoom": 18,
			"parts":   []string{"a", "b", "c"},
		},
		{
			"name":    "Google Hybrid",
			"url":     "http://mt{s}.google.com/vt/lyrs=y&x={x}&y={y}&z={z}&s=Galileo",
			"maxzoom": 20,
			"parts":   []string{"0", "1", "2", "3"},
		},
		{
			"name":    "Yandex maps",
			"url":     "https://core-renderer-tiles.maps.yandex.net/tiles?l=map&x={x}&y={y}&z={z}&scale=1&lang=ru_RU&projection=web_mercator",
			"maxzoom": 20,
		},
	}
}
