# _renderItem = (data) ->
#     "<p>Title :" + data.title + "</p>"

# $('#content').infinitescroll {
# 		navSelector  	: "a#next:last"
# 		nextSelector 	: "a#next:last"
# 		itemSelector 	: "#content p"
# 		debug		 	: true
# 		dataType	 	: 'json'
# 		appendCallback	: false
#     }, ( response ) ->
#         jsonData = response.results
#         $theCntr = $("#content");
#         newElements = "";
#         for d in jsonData
#             item = $ _renderItem d
#             $theCntr.append item


# LOGS = "http://#{window.location.hostname}:8004"

timeify = (tstr) ->
    moment.utc(tstr, "YY-MM-DD H:m:s").local().format("MMM D, hh:m A")

markupify_single = (entry) ->
    entry = $(entry).html().split(" - ")

    "<p class='entry'>"+
    "<span class='time'>#{timeify(entry[0])}</span>" +
    "<span class='author'>#{entry[1]}:</span>" +
    "<span class='msg'>#{entry[2..]}</span>" + "</p>"

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
            

