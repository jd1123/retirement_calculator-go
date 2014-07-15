var _ = require('lodash');

module.exports = function (data) {
  return new Highcharts.Chart({
    chart: {
      renderTo: 'chart',
      type: 'column'
    },

    xAxis: {
      categories: _.pluck(data['Bins'], 'Max')
    },

    plotOptions: {
      column: {
        groupPadding: 0,
        pointPadding: 0,
        borderWidth: 0
      }
    },

    series: [{
      data: _.pluck(data["Bins"], 'Weight')
    }]
  })
};
