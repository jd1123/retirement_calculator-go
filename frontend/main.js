var HighCharts = require('highcharts-browserify');

var beautify = require('js-beautify').js_beautify;

$(document).ready(function () {
  var loading = false;
  $('#gocalc').click(function () {
    if (loading) {
      alert('please wait asshole');
    } else {
      loading = true;
      $('pre#json').text('please wait');
      $.getJSON('/incomes/', function (data) {
        loading = false;
        var output = beautify(JSON.stringify(data), { indent_size: 2 });
        $('pre#json').text(output);
        require('./histo')(data);
      });
    }
  });
});
