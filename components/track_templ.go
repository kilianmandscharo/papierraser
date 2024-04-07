// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Track() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

// import "strconv"
// import "github.com/kilianmandscharo/papierraser/utils"
// import "fmt"
// import "github.com/kilianmandscharo/papierraser/types"
//
// templ Track(race types.Race) {
// <svg viewBox="0 0 100 100">
//   @Grid(race.Track)
//   @Line(race.Track.Finish[0], race.Track.Finish[1], "green")
//   @Path(race.Track.Inner, "black")
//   @Path(race.Track.Outer, "black")
//   @Players(race.Players, race.Turn)
// </svg>
// }
//
// templ Players(players []types.Player, turn int) {
// for _, player := range players {
// @Player(player.Path[len(player.Path) - 1], "orange")
// if player.Id == turn {
// @PlayerOptions(player.GetOptions())
// }
// }
// }
//
// templ Player(position types.Point, color string) {
// <circle fill={ color } cx={ strconv.Itoa(position.X * 5) } cy={ strconv.Itoa(position.Y * 5) } r="2" />
// }
//
// templ PlayerOptions(options []types.Point) {
// for _, option := range options {
// <circle id={ fmt.Sprintf("option-%d,%d", option.X, option.Y) } class="player-option" fill="purple" cx={
//   strconv.Itoa(option.X * 5) } cy={ strconv.Itoa(option.Y * 5) } r="2" />
// }
// }
//
// templ Grid(track types.Track) {
// for x := 0; x <= track.Width; x++ {
//   <line { utils.GetVerticalLineAttrs(x, track.Height)... } stroke="lightgray" />
// }
// for y := 0; y <= track.Height; y++ {
//   <line { utils.GetHorizontalLineAttrs(y, track.Width)... } stroke="lightgray" />
// }
// }
//
// templ Line(start, end types.Point, strokeColor string) {
// <line x1={ strconv.Itoa(start.X * 5) } y1={ strconv.Itoa(start.Y * 5) } x2={ strconv.Itoa(end.X * 5) } y2={
//   strconv.Itoa(end.Y * 5) } stroke="green" />
// }
//
// templ Path(path types.Path, strokeColor string) {
// <polygon points={ utils.GetPathString(path) } fill="none" stroke={ strokeColor }></polygon>
// }