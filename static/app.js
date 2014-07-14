$(document).ready(function () {
  var loading = false;
  $('#gocalc').click(function () {
    if (loading) {
      alert('please wait asshole');
    } else {
      loading = true;
      $('pre#json').text('please wait');
      $.getJSON('/recalc/', function (data) {
        loading = false;
        // data.age
        $('pre#json').text('too big.. but it\'s here');
      });
    }
  });
});
