package pterm_test

import (
	"errors"
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/pterm/pterm"
)

func TestNewRGB(t *testing.T) {
	type args struct {
		r uint8
		g uint8
		b uint8
	}
	tests := []struct {
		name string
		args args
		want pterm.RGB
	}{
		{name: "1", args: args{0, 0, 0}, want: pterm.RGB{0, 0, 0}},
		{name: "3", args: args{255, 255, 255}, want: pterm.RGB{255, 255, 255}},
		{name: "4", args: args{127, 127, 127}, want: pterm.RGB{127, 127, 127}},
		{name: "5", args: args{1, 2, 3}, want: pterm.RGB{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pterm.NewRGB(tt.args.r, tt.args.g, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRGB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRGBFromHEX(t *testing.T) {
	tests := []struct {
		hex  string
		want pterm.RGB
	}{
		{hex: "#ff0009", want: pterm.RGB{R: 255, G: 0, B: 9}},
		{hex: "ff0009", want: pterm.RGB{R: 255, G: 0, B: 9}},
		{hex: "ff00090x", want: pterm.RGB{R: 255, G: 0, B: 9}},
		{hex: "ff00090X", want: pterm.RGB{R: 255, G: 0, B: 9}},
		{hex: "#fba", want: pterm.RGB{R: 255, G: 187, B: 170}},
		{hex: "fba", want: pterm.RGB{R: 255, G: 187, B: 170}},
		{hex: "fba0x", want: pterm.RGB{R: 255, G: 187, B: 170}},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			rgb, err := pterm.NewRGBFromHEX(test.hex)
			testza.AssertEqual(t, test.want, rgb)
			testza.AssertNoError(t, err)
		})
	}
	testsFail := []struct {
		hex  string
		want error
	}{
		{hex: "faba0x", want: pterm.ErrHexCodeIsInvalid},
		{hex: "faba", want: pterm.ErrHexCodeIsInvalid},
		{hex: "#faba", want: pterm.ErrHexCodeIsInvalid},
		{hex: "faba0x", want: pterm.ErrHexCodeIsInvalid},
		{hex: "fax", want: strconv.ErrSyntax},
	}
	for _, test := range testsFail {
		t.Run("", func(t *testing.T) {
			_, err := pterm.NewRGBFromHEX(test.hex)
			testza.AssertTrue(t, errors.Is(err, test.want))
		})
	}
}

func TestRGB_Fade(t *testing.T) {
	type fields struct {
		R uint8
		G uint8
		B uint8
	}
	type args struct {
		min     float32
		max     float32
		current float32
		end     []pterm.RGB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pterm.RGB
	}{
		{name: "Middle", fields: fields{0, 0, 0}, args: args{min: 0, max: 100, current: 50, end: []pterm.RGB{{255, 255, 255}}}, want: pterm.RGB{127, 127, 127}},
		{name: "ZeroToZero", fields: fields{0, 0, 0}, args: args{min: 0, max: 100, current: 50, end: []pterm.RGB{{0, 0, 0}}}, want: pterm.RGB{0, 0, 0}},
		{name: "DifferentValues", fields: fields{0, 1, 2}, args: args{min: 0, max: 100, current: 50, end: []pterm.RGB{{0, 1, 2}}}, want: pterm.RGB{0, 1, 2}},
		{name: "NegativeRangeMiddle", fields: fields{0, 0, 0}, args: args{min: -50, max: 50, current: 0, end: []pterm.RGB{{255, 255, 255}}}, want: pterm.RGB{127, 127, 127}},
		{name: "NegativeRangeMiddleMultipleRGB", fields: fields{0, 0, 0}, args: args{min: -50, max: 50, current: 0, end: []pterm.RGB{{127, 127, 127}, {255, 255, 255}}}, want: pterm.RGB{127, 127, 127}},
		{name: "MiddleMultipleRGB", fields: fields{0, 0, 0}, args: args{min: 0, max: 100, current: 50, end: []pterm.RGB{{127, 127, 127}, {255, 255, 255}}}, want: pterm.RGB{127, 127, 127}},
		{name: "1/4MultipleRGB", fields: fields{0, 0, 0}, args: args{min: 0, max: 100, current: 25, end: []pterm.RGB{{255, 255, 255}, {255, 255, 255}}}, want: pterm.RGB{127, 127, 127}},
		{name: "MiddleMultipleRGBPositiveMin", fields: fields{0, 0, 0}, args: args{min: 10, max: 110, current: 60, end: []pterm.RGB{{127, 127, 127}, {255, 255, 255}}}, want: pterm.RGB{127, 127, 127}},
		{name: "MiddleNoRGB", fields: fields{0, 0, 0}, args: args{min: 10, max: 110, current: 60, end: []pterm.RGB{}}, want: pterm.RGB{0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pterm.RGB{
				R: tt.fields.R,
				G: tt.fields.G,
				B: tt.fields.B,
			}
			if got := p.Fade(tt.args.min, tt.args.max, tt.args.current, tt.args.end...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGB_GetValues(t *testing.T) {
	type fields struct {
		R uint8
		G uint8
		B uint8
	}
	tests := []struct {
		name   string
		fields fields
		wantR  uint8
		wantG  uint8
		wantB  uint8
	}{
		{name: "Zero", fields: fields{R: 0, G: 0, B: 0}, wantR: uint8(0), wantG: uint8(0), wantB: uint8(0)},
		{name: "Max", fields: fields{R: 255, G: 255, B: 255}, wantR: uint8(255), wantG: uint8(255), wantB: uint8(255)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pterm.RGB{
				R: tt.fields.R,
				G: tt.fields.G,
				B: tt.fields.B,
			}
			gotR, gotG, gotB := p.GetValues()
			if gotR != tt.wantR {
				t.Errorf("GetValues() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotG != tt.wantG {
				t.Errorf("GetValues() gotG = %v, want %v", gotG, tt.wantG)
			}
			if gotB != tt.wantB {
				t.Errorf("GetValues() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestRGB_Print(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run(pterm.Sprintf("%v %v %v", rgb.R, rgb.G, rgb.B), func(t *testing.T) {
			testPrintContains(t, func(w io.Writer, a interface{}) {
				p := rgb.Print(a)
				testza.AssertNotNil(t, p)
			})
		})
	}
}

func TestRGB_Printf(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run(pterm.Sprintf("%v %v %v", rgb.R, rgb.G, rgb.B), func(t *testing.T) {
			testPrintfContains(t, func(w io.Writer, format string, a interface{}) {
				p := rgb.Printf(format, a)
				testza.AssertNotNil(t, p)
			})
		})
	}
}

func TestRGB_Printfln(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run(pterm.Sprintfln("%v %v %v", rgb.R, rgb.G, rgb.B), func(t *testing.T) {
			testPrintflnContains(t, func(w io.Writer, format string, a interface{}) {
				p := rgb.Printfln(format, a)
				testza.AssertNotNil(t, p)
			})
		})
	}
}

func TestRGB_Println(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run(pterm.Sprintf("%v %v %v", rgb.R, rgb.G, rgb.B), func(t *testing.T) {
			testPrintlnContains(t, func(w io.Writer, a interface{}) {
				p := rgb.Println(a)
				testza.AssertNotNil(t, p)
			})
		})
	}
}

func TestRGB_Sprint(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run(pterm.Sprintf("%v %v %v", rgb.R, rgb.G, rgb.B), func(t *testing.T) {
			testSprintContains(t, func(a interface{}) string {
				return rgb.Sprint(a)
			})
		})
	}
}

func TestRGB_Sprintf(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run("", func(t *testing.T) {
			testSprintfContains(t, func(format string, a interface{}) string {
				return rgb.Sprintf(format, a)
			})
		})
	}
}

func TestRGB_Sprintfln(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run("", func(t *testing.T) {
			testSprintflnContains(t, func(format string, a interface{}) string {
				return rgb.Sprintfln(format, a)
			})
		})
	}
}

func TestRGB_Sprintln(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run(pterm.Sprintf("%v %v %v", rgb.R, rgb.G, rgb.B), func(t *testing.T) {
			testSprintlnContains(t, func(a interface{}) string {
				return rgb.Sprintln(a)
			})
		})
	}
}

func TestRGB_PrintOnError(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run("PrintOnError", func(t *testing.T) {
			result := captureStdout(func(w io.Writer) {
				rgb.PrintOnError(errors.New("hello world"))
			})
			testza.AssertContains(t, result, "hello world")
		})
	}
}

func TestRGB_PrintIfError_WithoutError(t *testing.T) {
	RGBs := []pterm.RGB{{0, 0, 0}, {127, 127, 127}, {255, 255, 255}}

	for _, rgb := range RGBs {
		t.Run("PrintIfError_WithoutError", func(t *testing.T) {
			result := captureStdout(func(w io.Writer) {
				rgb.PrintOnError(nil)
			})
			testza.AssertZero(t, result)
		})
	}
}
