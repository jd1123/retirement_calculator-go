var HighCharts = require('highcharts-browserify')
  , _ = require('lodash');

module.exports = function (raw_bins) {

  var bins = _.filter(raw_bins, function (bin) {
    return bin['Weight'] > 0.001;
  });

  return new Highcharts.Chart({
    chart: {
      renderTo: 'mainchart',
      type: 'column'
    },

    xAxis: {
      categories: _.map(_.pluck(bins, 'Min'), function (number) {
        return parseInt(number);
      }),
      labels: { enabled: false }
    },

    plotOptions: {
      column: {
        groupPadding: 0,
        pointPadding: 0,
        borderWidth: 0
      }
    },

    series: [{
      data: _.pluck(bins, 'Weight')
    }]
  })
};
