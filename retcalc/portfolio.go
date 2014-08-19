package retcalc

import "math"

/* 	Portfolio objects
Will let the user decide which type of portfolio he/she will
use for their retirement portfolio. Will be based on historical
data. Will also let the user make a portfolio
*/

type Portfolio struct {
	Mean  float64
	Stdev float64
}

// Basic returns and volatility parameters
// FIXME: check these against history
var STOCK_RETURNS = 0.07
var STOCK_VOLATILITY = 0.15
var BOND_RETURNS = 0.04
var BOND_VOLATILITY = 0.07
var CORRELATION = 0.0

var HIGHRISKPORTFOLIO Portfolio = Portfolio{0.07, 0.15}    // All stocks
var MEDIUMRISKPORTFOLIO Portfolio = Portfolio{0.055, 0.12} // 80% stocks 20% bond
var LOWRISKPORTFOLIO Portfolio = Portfolio{0.04, 0.09}     // 50% stocks 50% bonds
var BLANKPORTFOLIO Portfolio = Portfolio{0.0, 0.0}         // the blank portfolio

func NewCustomPortfolio(expectedReturn, assetVolatility float64) Portfolio {
	return Portfolio{expectedReturn, assetVolatility}
}

func CustomAssetMix(stocks, bonds float64) Portfolio {
	if stocks+bonds != 1.0 {
		dem := stocks + bonds
		stocks /= dem
		bonds /= dem
	}
	expReturn := STOCK_RETURNS*stocks + BOND_RETURNS*bonds
	expVol := math.Sqrt(stocks*math.Pow(STOCK_VOLATILITY, 2) + bonds*math.Pow(BOND_VOLATILITY, 2) +
		2*CORRELATION*STOCK_VOLATILITY*BOND_VOLATILITY*bonds*stocks)
	return Portfolio{expReturn, expVol}
}
