package romanizerhandlers

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"

	"github.com/keypeearr/romanizer/src/utils"
	"github.com/keypeearr/romanizer/src/views/pages"
	"github.com/keypeearr/romanizer/src/views/props"
)

var romans = []struct {
	value int
	digit string
}{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

var alphas = map[byte]int{'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000}

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func DisplayRomanizer(ctx fiber.Ctx) error {
	props := props.RomanizerPageProps{
		Main: props.MainProps{
			Title: "Romanizer",
		},
		Alpha: props.RomanizerInputProps{
			Name:  "alpha",
			Class: "romanizer-alpha-input",
			Value: "",
			Attrs: templ.Attributes{
				"placeholder":   "9",
				"minlength":     "1",
				"maxlength":     "4",
				"hx-post":       "/api/v1/romanizer/alphaToRoman",
				"hx-target":     ".romanizer-roman-input",
				"hx-target-400": ".romanizer-error",
				"hx-on":         "input changed",
				"hx-swap":       "outerHTML",
			},
		},
		Roman: props.RomanizerInputProps{
			Name:  "roman",
			Class: "romanizer-roman-input",
			Value: "",
			Attrs: templ.Attributes{
				"placeholder":   "IX",
				"minlength":     "1",
				"maxlength":     "16",
				"hx-post":       "/api/v1/romanizer/romanToAlpha",
				"hx-target":     ".romanizer-alpha-input",
				"hx-target-400": ".romanizer-error",
				"hx-on":         "input changed",
				"hx-swap":       "outerHTML",
			},
		},
	}

	return utils.Render(ctx, pages.Romanizer(props))
}

func DisplayRomanValue(ctx fiber.Ctx) error {
	props := props.RomanizerPageProps{
		Roman: props.RomanizerInputProps{
			Name:  "roman",
			Class: "romanizer-roman-input",
			Value: "",
			Attrs: templ.Attributes{
				"placeholder":   "IX",
				"minlength":     "1",
				"maxlength":     "16",
				"hx-post":       "/api/v1/romanizer/romanToAlpha",
				"hx-target":     ".romanizer-alpha-input",
				"hx-target-400": ".romanizer-error",
				"hx-on":         "input changed",
				"hx-swap":       "outerHTML",
			},
		},
	}

	type RequestBody struct {
		Value int `form:"alpha"`
	}

	body := new(RequestBody)
	if err := ctx.Bind().Body(body); err != nil {
		ctx.Status(fiber.ErrBadRequest.Code)
		return utils.Render(ctx, pages.RomanizerError(err))
	}

	if body.Value < 1 {
		ctx.Status(fiber.ErrBadRequest.Code)
		return utils.Render(ctx, pages.RomanizerError(errors.New("value must not be less than 1")))
	}

	if body.Value > 3999 {
		ctx.Status(fiber.ErrBadRequest.Code)
		return utils.Render(
			ctx,
			pages.RomanizerError(errors.New("value must not be greater than 3999")),
		)
	}

	// Reference:
	// https://freshman.tech/snippets/go/roman-numerals/
	temp := body.Value
	var output strings.Builder
	for _, roman := range romans {
		if temp == 0 {
			break
		}
		for temp >= roman.value {
			output.WriteString(roman.digit)
			temp -= roman.value
		}
	}

	props.Roman.Value = output.String()
	return utils.Render(ctx, pages.RomanizerInput(props.Roman))
}

func DisplayAlphaValue(ctx fiber.Ctx) error {
	props := props.RomanizerPageProps{
		Alpha: props.RomanizerInputProps{
			Name:  "alpha",
			Class: "romanizer-alpha-input",
			Value: "",
			Attrs: templ.Attributes{
				"placeholder":   "9",
				"minlength":     "1",
				"maxlength":     "4",
				"hx-post":       "/api/v1/romanizer/alphaToRoman",
				"hx-target":     ".romanizer-roman-input",
				"hx-target-400": ".romanizer-error",
				"hx-on":         "input changed",
				"hx-swap":       "outerHTML",
			},
		},
	}

	type RequestBody struct {
		Value string `form:"roman"`
	}

	body := new(RequestBody)
	if err := ctx.Bind().Body(body); err != nil {
		ctx.Status(fiber.ErrBadRequest.Code)
		return utils.Render(ctx, pages.RomanizerError(err))
	}

	if body.Value == "" {
		return utils.Render(ctx, pages.RomanizerInput(props.Alpha))
	}

	// Reference:
	// https://stackoverflow.com/questions/38554353/how-to-check-if-a-string-only-contains-alphabetic-characters-in-go
	if !IsLetter(body.Value) {
		ctx.Status(fiber.ErrBadRequest.Code)
		return utils.Render(
			ctx,
			pages.RomanizerError(errors.New("It is not a valid input (numbers)")),
		)
	}

	sum := 0
	prevValue := 0
	for i := range len(body.Value) {
		currentValue, ok := alphas[body.Value[i]]
		if !ok {
			ctx.Status(fiber.ErrBadRequest.Code)
			return utils.Render(
				ctx,
				pages.RomanizerError(errors.New("It is not a valid input (letters)")),
			)
		}
		if currentValue > prevValue {
			sum += currentValue - 2*prevValue
		} else {
			sum += currentValue
		}
		prevValue = currentValue

	}

	props.Alpha.Value = strconv.Itoa(sum)
	return utils.Render(ctx, pages.RomanizerInput(props.Alpha))
}
