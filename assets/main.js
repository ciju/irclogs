// Generated by CoffeeScript 1.6.2
(function() {
  var duration, localTime, markupify, markupify_single, time;

  localTime = function(tstr) {
    return moment.utc(tstr, "YY-MM-DD H:m:s").local();
  };

  duration = function(tstr) {
    return localTime(tstr).fromNow();
  };

  time = function(tstr) {
    return localTime(tstr).format("MMM D, hh:m A");
  };

  $('#content').on('mouseenter', '.entry', function() {
    var $e;

    $e = $(this).find('.time');
    return $e.html(time($e.data('time')));
  }).on('mouseleave', '.entry', function() {
    var $e;

    $e = $(this).find('.time');
    return $e.html(duration($e.data('time')));
  });

  markupify_single = function(entry) {
    entry = $(entry).html().split(" - ");
    return "<p class='entry'>" + ("<span class='time' data-time='" + entry[0] + "'>" + (duration(entry[0])) + "</span>") + ("<span class='author'>" + entry[1] + ":</span>") + ("<span class='msg'>" + entry.slice(2) + "</span>") + "</p>";
  };

  markupify = function(resp) {
    var i, _i, _len, _results;

    _results = [];
    for (_i = 0, _len = resp.length; _i < _len; _i++) {
      i = resp[_i];
      _results.push(markupify_single(i));
    }
    return _results;
  };

  $(function() {
    return $.get("/logs?page=1", function(data) {
      console.log(markupify($(data).find('p')));
      $('#content').append(markupify($(data).find('p')));
      console.log('started');
      return $('#content').infinitescroll({
        navSelector: "#next:last",
        nextSelector: "a#next:last",
        itemSelector: "#content p",
        appendCallback: false,
        debug: true,
        dataType: 'html',
        finishedMsg: 'we are done here',
        path: function(index) {
          console.log("curr_page", index);
          return "/logs?page=" + index;
        },
        behavior: 'reverse'
      }, function(resp, opts) {
        var $curr_top, prev_height;

        console.log(resp);
        $curr_top = $(opts.itemSelector).children().first();
        $('#content').prepend(markupify(resp));
        prev_height = 0;
        $curr_top.parent().prevAll().each(function() {
          return prev_height += $(this).outerHeight(true);
        });
        return $("html, body").scrollTop(prev_height);
      });
    });
  });

}).call(this);
