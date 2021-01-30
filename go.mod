module github.com/rcssggb/ggb-single

go 1.15

// replace github.com/rcssggb/ggb-lib => ../ggb-lib
replace github.com/aunum/goro => ../goro

require (
	github.com/Arafatk/glot v0.0.0-20180312013246-79d5219000f0
	github.com/aunum/goro v0.0.0-20200405193922-6843b7790935
	github.com/rcssggb/ggb-lib v0.2.9
	gorgonia.org/gorgonia v0.9.16
	gorgonia.org/tensor v0.9.19
)
