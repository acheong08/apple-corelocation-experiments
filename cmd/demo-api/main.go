package main

import (
	_ "embed"
	"github.com/acheong08/apple-corelocation-experiments/lib"
	"log"

	"github.com/a-h/templ"
	"github.com/acheong08/clir"

	"github.com/labstack/echo/v4"
)

//go:embed main.js
var mainJs []byte

func init() {
	log.SetFlags(log.Lshortfile)
}

type gps struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func main() {
	lat := 51.51493459648336
	long := -3.1548554460964624
	var china bool
	cli := clir.NewCli("demo", "Interactive user interface to demonstrate the functionality of Apple's Geolocation services", "v0.0.1")

	cli.WithFlags(
		clir.Float64Flag("lat", "default latitude", &lat),
		clir.Float64Flag("long", "default longitude", &long),
		clir.BoolFlag("china", "use the Chinese API", &china),
	)

	cli.Action(func() error {
		e := echo.New()
		e.GET("/", func(c echo.Context) error {
			return Render(c, 200, Index(lat, long, china))
		})
		e.GET("/main.js", func(c echo.Context) error {
			c.Response().Header().Set("content-type", "application/javascript")
			// Set status code
			c.Response().WriteHeader(200)
			_, err := c.Response().Write(mainJs)
			return err
		})
		e.POST("/gps", func(c echo.Context) error {
			var g gps
			if err := c.Bind(&g); err != nil {
				return c.String(400, "Bad Request")
			}
			if g.Lat < -90 || g.Lat > 90 || g.Long < -180 || g.Long > 180 || g.Lat == 0 || g.Long == 0 {
				return c.String(400, "Bad Request")
			}

			var options []lib.Modifier = make([]lib.Modifier, 0)
			if china {
				options = append(options, lib.Options.WithRegion(lib.Options.China))
			}

			points, err := lib.SearchProximity(g.Lat, g.Long, 20, options...)
			if err != nil {
				log.Println(err)
				return c.String(404, "did not find any points nearby")
			}

			return c.JSON(200, map[string]any{
				"closest": points[0],
				"points":  points[1:],
			})
		})
		e.Logger.Fatal(e.Start("127.0.0.1:1974"))
		return nil
	})
	if err := cli.Run(); err != nil {
		panic(err)
	}
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx.Response().Writer.WriteHeader(statusCode)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}
