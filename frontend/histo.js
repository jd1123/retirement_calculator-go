var _ = require('lodash');

module.exports = function (data) {
  return new Highcharts.Chart({
    chart: {
      renderTo: 'chart',
      type: 'column'
    },

    xAxis: {
      categories: _.map(_.pluck(data['Bins'], 'Min'), function (number) {
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
      data: _.pluck(data["Bins"], 'Weight')
    }]
  })
};
