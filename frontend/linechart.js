var HighCharts = require('highcharts-browserify')
  , _ = require('lodash')

module.exports = function (entries) {
  this.render = function() {
    return new Highcharts.Chart({
      chart: {
        renderTo: 'linechart',
        type: 'line'
      },

      xAxis: {
        categories: _.map(entries, function(entry) {
          return entry['Year'].substring(0, 4)
        })
      },

      plotOptions: {
        column: {
          groupPadding: 0,
          pointPadding: 0,
          borderWidth: 0
        }
      },

      series: [{
        name: 'Total Balance',
        data: _.map(entries, function(entry) {
          return parseFloat(entry['EOY_total_balance'].toFixed(2))
        })
      }]
    })
  }

  this.show = function() {
    $('#linechart').show()
    $(window).resize();
    $("html, body").animate({ scrollTop: $(document).height() }, "slow");
  }
};
