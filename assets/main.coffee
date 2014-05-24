localTime = (tstr) ->
    moment.utc(tstr, "YY-MM-DD H:m:s").local()

duration = (tstr) ->
    localTime(tstr).fromNow()

time = (tstr) ->
    localTime(tstr).format "MMM D, hh:m A"

$('#content').on('mouseenter', '.entry', ->
    $e = $(@).find('.time')
    $e.html time $e.data 'time'
).on('mouseleave', '.entry', ->
    $e = $(@).find('.time')
    $e.html duration $e.data 'time'
)

markupify_single = (entry) ->
    entry = $(entry).html().split(" - ")
    return "" if !entry[0].length

    "<div class='entry'>"+
    "<div class='time' data-time='#{entry[0]}'>#{duration(entry[0])}</div>" +
    "<div class='author'>#{entry[1]}</div>" +
    "<div class='msg'>#{entry[2..]}</div>" + "</div>"

markupify = (resp) ->
    markupify_single(i) for i in resp

$ ->
    $.get "/logs?page=1", (data) ->
        console.log markupify $(data).find 'p'
        $('#content').append markupify $(data).find 'p'
        console.log('started')


        $('#content').infinitescroll {
            navSelector  	: "#next:last"
            nextSelector 	: "a#next:last"
            itemSelector 	: "#content p"
            appendCallback  : false
            debug		 	: true
            dataType	 	: 'html'
            # maxPage         : 3,
            finishedMsg: 'we are done here',
            path: (index) ->
                console.log "curr_page", index
                "/logs?page=#{index}"
            behavior		: 'reverse'
        }, (resp, opts) ->
            console.log resp

            $curr_top = $(opts.itemSelector).children().first()

            $('#content').prepend markupify resp

            prev_height = 0
            $curr_top.parent().prevAll().each -> prev_height += $(@).outerHeight(true);

            $("html, body").scrollTop prev_height
