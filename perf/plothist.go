package perf

import (
	"image/color"
	"io"

	"github.com/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// PlotHist renders a histogram plot to a writer.
// Param format can be anything supported by gonum.org/v1/plot/vg
func PlotHist(series []float64, writer io.Writer, nBins int, title string, format string) (int64, error) {
	norm := false
	return plotHistWithNormOpt(series, norm, writer, nBins, title, format)
}

// PlotHistStdNormDist renders a histogram plot to a writer.
// Normalizes the area under the series to 1 and adds a std normal distribution plot.
// Param format can be anything supported by gonum.org/v1/plot/vg
func PlotHistStdNormDist(series []float64, writer io.Writer, nBins int, title string, format string) (int64, error) {
	norm := true
	return plotHistWithNormOpt(series, norm, writer, nBins, title, format)
}

func plotHistWithNormOpt(series []float64, plotNorm bool, writer io.Writer, nBins int, title string, format string) (int64, error) {
	// Copy series values into plotter format
	vs := make(plotter.Values, len(series))
	for i := range vs {
		vs[i] = series[i]
	}

	// Make a plot and set its title
	p := plot.New()
	p.Title.Text = title

	// Create a histogram of our values
	h, err := plotter.NewHist(vs, nBins)
	if err != nil {
		return 0, err
	}

	if plotNorm {
		// Normalize the area under the histogram to sum to one
		// Required to correctly render the additional probability distribution
		h.Normalize(1)
		p.Add(h)

		// Plot normal distribution function and add to hist
		norm := plotter.NewFunction(distuv.UnitNormal.Prob)
		norm.Color = color.RGBA{R: 255, A: 255}
		norm.Width = vg.Points(2)
		p.Add(norm)
	} else {
		p.Add(h)
	}

	// Render the plot to the writer.
	writeTo, err := p.WriterTo(4*vg.Inch, 4*vg.Inch, format)
	if err != nil {
		return 0, err
	}

	return writeTo.WriteTo(writer)
}
