package retcalc

type Portfolio struct {
	Mean  float64
	Stdev float64
}

var HIGHRISKPORTFOLIO Portfolio = Portfolio{0.07, 0.15}
var MEDIUMRISKPORTFOLIO Portfolio = Portfolio{0.055, 0.12}
var LOWRISKPORTFOLIO Portfolio = Portfolio{0.04, 0.09}
var BLANKPORTFOLIO Portfolio = Portfolio{0.0, 0.0}
