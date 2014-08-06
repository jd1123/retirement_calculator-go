var beautify = require('js-beautify').js_beautify;

$(document).ready(function () {
  var sessionId = Math.random().toString(16).substring(2)
  sessionStorage.setItem("SessionID", sessionId)
  var loading = false;
  $('#gocalc').click(function () {
    if (loading) {
      alert('please wait asshole');
    } else {
      loading = true;
      $('pre#json').text('please wait');
      var payload = { SessionId: sessionId };
      $('#inputs input').each(function(i, e) {
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
        $('pre#json').text(output);
        require('./histo')(data["Bins"]);
        $('#mainchart').show();
        $('#mainchart .highcharts-tracker rect').on('click', function(){
          $.ajax({
            type: "GET",
            url: '/paths/',
            beforeSend: function(xhr) {
              xhr.setRequestHeader('X-Session-Id', sessionId)
			  xhr.setRequestHeader('X-Percentile-Req', 0.25)
            },
            success: function (data) {
              console.log(data)
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
