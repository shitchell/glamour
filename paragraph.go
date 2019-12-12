package gold

import (
	"bytes"
	"io"
	"strings"

	"github.com/muesli/reflow"
)

type ParagraphElement struct {
}

func (e *ParagraphElement) Render(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack
	var rules ElementStyle

	rules = ctx.style[Paragraph]
	_, _ = w.Write([]byte("\n"))
	be := BlockElement{
		Block: &bytes.Buffer{},
		Style: cascadeStyle(bs.Current().Style, rules, true),
	}
	bs.Push(be)

	renderText(w, bs.Parent().Style, rules.Prefix)
	renderText(w, bs.Current().Style, rules.StyledPrefix)
	return nil
}

func (e *ParagraphElement) Finish(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack
	rules := bs.Current().Style

	keepNewlines := false

	mw := NewMarginWriter(ctx, w, rules)
	if len(strings.TrimSpace(bs.Current().Block.String())) > 0 {
		flow := reflow.NewReflow(int(bs.Width(ctx)))
		flow.KeepNewlines = keepNewlines
		_, _ = flow.Write(bs.Current().Block.Bytes())
		flow.Close()

		_, err := mw.Write(flow.Bytes())
		if err != nil {
			return err
		}
		_, _ = mw.Write([]byte("\n"))
	}

	renderText(w, bs.Current().Style, rules.StyledSuffix)
	renderText(w, bs.Parent().Style, rules.Suffix)

	bs.Current().Block.Reset()
	bs.Pop()
	return nil
}
