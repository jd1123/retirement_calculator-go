var beautify = require('js-beautify').js_beautify
  , Histogram = require('./histogram')
  , LineChart = require('./linechart')

$(document).ready(function () {
  var sessionId = Math.random().toString(16).substring(2)
  sessionStorage.setItem("SessionID", sessionId)
  var loading = false;
  $('#gocalc').click(function () {
    if (loading) {
      alert('please wait asshole');
    } else {
      loading = true;
      var payload = { SessionId: sessionId };
      $('#inputs input #select select').each(function(i, e) {
		console.log("got one");
        var $el = $(e);
        var type = $el.data('type');
        var val = $el.val();
        if ( type === "float") {
          val = parseFloat(val);
        } else if (type === "int") {
          val = parseInt(val);
        }
        payload[$el.data('key')] = val;
      });

      var onSuccess = function (data) {
        loading = false;
        var output = beautify(JSON.stringify(data), { indent_size: 2 });
        var histo = new Histogram(data["Bins"]);
        histo.render();
        histo.show();
        histo.getBars().on('click', function(event){
          ix = $(event.target).index();
          percentile = data['Bins'][ix]['Weight'];
          $.ajax({
            type: "GET",
            url: '/paths/',
            beforeSend: function(xhr) {
              xhr.setRequestHeader('X-Session-Id', sessionId)
              xhr.setRequestHeader('X-Percentile-Req', percentile)
            },
            success: function (data) {
              var lineChart = new LineChart(data['Yearly_entries']);
              lineChart.render()
              lineChart.show()
            },
            dataType: 'json'
          });
        })
      };

      $.ajax({
        type: "POST",
        url: '/input/',
        data: JSON.stringify(payload),
        success: onSuccess,
        dataType: 'json'
      });

    }
  });
});
