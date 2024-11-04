package table

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

const (
	SquareDownHorizontal = "┬"
	SquareRightVertial   = "├"
	SquareLeftVertial    = "┤"
	SquareUpHorizontal   = "┴"
	SquareCross          = "┼"

	SquareTopLeft     = "┌"
	SquareTopRight    = "┐"
	SquareBottomLeft  = "└"
	SquareBottomRight = "┘"

	RoundedTopLeft     = "╭"
	RoundedTopRight    = "╮"
	RoundedBottomLeft  = "╰"
	RoundedBottomRight = "╯"

	SeperatorHorizontal = "═"
	SeperatorLeft       = "╞"
	SeperatorRight      = "╡"
	SeperatorSeperator  = "╪"

	HorizontalLine = "─"
	VerticalLine   = "│"
)

const sampleTable = `
┌───────┬─────┬─────┬────────┐
│Hallo  │Was  │ Geht│ Hiiiiii│
├───────┼─────┼─────┼────────┤
│  sdfa	│asdf │asdf │  asdf  │
└───────┴─────┴─────┴────────┘
`
const bubbleteaTable = `
┌────────────────────────────┐
│Rank    City      County    │
│────────────────────────────│
│                            │
│                            │
│                            │
│                            │
└────────────────────────────┘
`

type Table struct {
	headers        []*Header
	rows           [][]string
	row            int
	col            int
	minLengths     []int
	padding        int
	title          string
	seperators     []int
	roundedCorners bool
}

func NewTable(headers ...*Header) *Table {
	minLen := make([]int, len(headers))

	for i, h := range headers {
		minLen[i] = utf8.RuneCountInString(h.text)
	}

	return &Table{
		headers:    headers,
		minLengths: minLen,
		padding:    1,
	}
}

type Header struct {
	text string

	headingCentered bool
	rowCentered     bool
}

func NewHeader(text string, centered ...bool) *Header {
	var (
		headingCentered = false
		rowCentered     = false
	)

	switch len(centered) {
	case 0:
		break
	case 1:
		headingCentered = centered[0]
		rowCentered = centered[0]
	default:
		headingCentered = centered[0]
		rowCentered = centered[1]
	}

	return &Header{
		text:            text,
		headingCentered: headingCentered,
		rowCentered:     rowCentered,
	}
}

func (h *Header) RowCentered() *Header {
	h.rowCentered = true
	return h
}

func (h *Header) HeadingCentered() *Header {
	h.headingCentered = true
	return h
}

func (t *Table) AddSeperator() {
	t.seperators = append(t.seperators, len(t.rows)-1)
}

func (t *Table) AddRow(row []string) {
	lenH, lenR := len(t.headers), len(row)

	if lenR > lenH {
		t.rows = append(t.rows, row[:lenH])
		return
	}

	if lenR < lenH {
		for range lenH - lenR {
			row = append(row, "")
		}
	}

	for i, r := range row {
		length := utf8.RuneCountInString(r)
		if t.minLengths[i] < length {
			t.minLengths[i] = length
		}
	}

	t.rows = append(t.rows, row)
}

func (t Table) createTitle() string {
	b := new(strings.Builder)

	var (
		topLeft  = SquareTopLeft
		topRight = SquareTopRight
	)

	if t.roundedCorners {
		topLeft = RoundedTopLeft
		topRight = RoundedTopRight
	}

	b.WriteString(topLeft)
	b.WriteString(strings.Repeat(HorizontalLine, len(t.title)+2))
	b.WriteString(topRight + "\n")

	b.WriteString(VerticalLine)
	b.WriteString(fmt.Sprintf(" %s ", t.title))
	b.WriteString(VerticalLine)

	return b.String() + "\n"
}

func (t Table) String() string {
	b := new(strings.Builder)

	if t.title != "" {
		b.WriteString(t.createTitle())
	}

	b.WriteString(t.createTopLine())
	b.WriteString(t.createHeading())

	if len(t.rows) > 0 {
		b.WriteString(t.createHeadingSeperator())
	}

	for i := range t.rows {
		b.WriteString(t.createRow(i))
		for _, s := range t.seperators {
			if i != s {
				continue
			}
			b.WriteString(SeperatorLeft)
			for i, m := range t.minLengths {
				b.WriteString(strings.Repeat(SeperatorHorizontal, m+2*t.padding))
				if i == len(t.headers)-1 {
					b.WriteString(SeperatorRight)
				} else {
					b.WriteString(SeperatorSeperator)
				}
			}
			b.WriteString("\n")
		}
	}

	b.WriteString(t.createBottomLine())

	return b.String()
}

func (t Table) createRow(idx int) string {
	centerText := func(s string, length int) string {
		l := utf8.RuneCountInString(s)
		space := strings.Repeat(" ", (length-l)/2)

		return fmt.Sprintf("%s%s%s", space, s, strings.Repeat(" ", length-l-len(space)))
	}

	b := new(strings.Builder)

	for i, l := range t.minLengths {
		var (
			item   = t.rows[idx][i]
			length = utf8.RuneCountInString(item)
		)
		b.WriteString(VerticalLine)
		b.WriteString(strings.Repeat(" ", t.padding))
		if t.headers[i].rowCentered {
			b.WriteString(centerText(item, l))
		} else {
			b.WriteString(item)
			b.WriteString(strings.Repeat(" ", l-length))
		}
		b.WriteString(strings.Repeat(" ", t.padding))
	}

	b.WriteString(VerticalLine)

	return b.String() + "\n"
}

func (t *Table) WithTitle(title string) *Table {
	t.title = title
	return t
}

func (t Table) createHeadingSeperator() string {
	b := new(strings.Builder)
	b.WriteString(SquareRightVertial)

	for i, l := range t.minLengths {
		b.WriteString(strings.Repeat(HorizontalLine, l+2*t.padding))
		if i != len(t.headers)-1 {
			b.WriteString(SquareCross)
			continue
		}
		b.WriteString(SquareLeftVertial)
	}

	return b.String() + "\n"
}

func (t Table) createBottomLine() string {
	b := new(strings.Builder)

	var (
		bottomLeft  = SquareBottomLeft
		bottomRight = SquareBottomRight
	)

	if t.roundedCorners {
		bottomLeft = RoundedBottomLeft
		bottomRight = RoundedBottomRight
	}

	b.WriteString(bottomLeft)
	for i, l := range t.minLengths {
		b.WriteString(strings.Repeat(HorizontalLine, l+2*t.padding))
		if i != len(t.headers)-1 {
			b.WriteString(SquareUpHorizontal)
			continue
		}
		b.WriteString(bottomRight)
	}

	return b.String()
}

func (t Table) createHeading() string {
	b := new(strings.Builder)

	for i, h := range t.headers {
		b.WriteString(VerticalLine)

		if h.headingCentered {
			totalLength := t.minLengths[i] + 2*t.padding
			padding := float64(totalLength-len(h.text)) / 2.0
			b.WriteString(strings.Repeat(" ", int(math.Floor(padding))))
			b.WriteString(h.text)
			b.WriteString(strings.Repeat(" ", int(math.Ceil(padding))))
		} else {
			b.WriteString(strings.Repeat(" ", t.padding))
			b.WriteString(h.text)
			b.WriteString(strings.Repeat(" ", t.minLengths[i]-utf8.RuneCountInString(h.text)))
			b.WriteString(strings.Repeat(" ", t.padding))
		}
		if i == len(t.headers)-1 {
			b.WriteString(VerticalLine)
		}
	}

	return b.String() + "\n"
}

// HACK: Works but really messy
func (t Table) createTopLine() string {
	b := new(strings.Builder)

	var (
		topLeft  = SquareTopLeft
		topRight = SquareTopRight
	)

	if t.roundedCorners {
		topLeft = RoundedTopLeft
		topRight = RoundedTopRight
	}

	if t.title == "" {

		b.WriteString(topLeft)
		for i := range t.headers {
			b.WriteString(strings.Repeat(HorizontalLine, t.minLengths[i]+2*t.padding))
			if i != len(t.headers)-1 {
				b.WriteString(SquareDownHorizontal)
				continue
			}
			b.WriteString(topRight)
		}
	} else {
		totalLength := 0

		for _, l := range t.minLengths {
			totalLength += l + 2*t.padding
		}
		totalLength += len(t.headers) + 1

		titleDown := len(t.title) + 2
		cur := 0
		nextSep := t.minLengths[cur] + 2*t.padding

		b.WriteString(SquareRightVertial)
		for i := range totalLength - 2 {
			if i == titleDown && i == nextSep {
				b.WriteString(SquareCross)
				cur++
				nextSep += t.minLengths[cur] + 2*t.padding + 1
				continue
			}

			switch i {
			case titleDown:
				b.WriteString(SquareUpHorizontal)
			case nextSep:
				b.WriteString(SquareDownHorizontal)
				if cur+1 < len(t.minLengths) {
					cur++
				}
				nextSep += t.minLengths[cur] + 2*t.padding + 1
			default:
				b.WriteString(HorizontalLine)
			}
		}
		b.WriteString(topRight)
	}

	return b.String() + "\n"
}

func (t *Table) WithRoundedCorners() *Table {
	t.roundedCorners = true
	return t
}
