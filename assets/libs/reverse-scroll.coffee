log = (args...) ->
    console.log(args...) if console?.log?

$ = jQuery

log 'behaviour: reverse-scroll initialized'    

reverse_scroll =
    _neartop: -> @options.binder.scrollTop() <= 0
    
    scroll_reverse: ->
        s = @options.state

        if s.isDuringAjax or s.isInvalidPage or s.isDone or s.isDestroyed or s.isPaused
            return

        unless @loadingsetup
            @loadingsetup = true
            @options.loading.start = =>
                console.group('loading')
                @.beginAjax @options
            @options.loading.finished = ->
                console.groupEnd()

        return unless @._neartop()

        @.retrieve();

    _loadcallback_reverse: (box, data, url) ->
        opts = @options
        callback = opts.callback # GLOBAL OBJECT FOR CALLBACK
        result = if opts.state.isDone
            'done'
        else
            if opts.appendCallback then 'append' else 'no-append'

        log 'operation: ', result
        
        switch result
          when 'done' then @._showdonemsg()
          when 'no-append'
            if opts.dataType is 'html'
              data = $("<div>#{data}</div>").find(opts.itemSelector)
          when 'append'
            children = box.children()

            return @._error('end') if children.length is 0

            log 'box', box, children

            # use a documentFragment because it works when content is going into a table or UL
            frag = document.createDocumentFragment()
            frag.appendChild(box[0].firstChild) while box[0].firstChild

            content = $(opts.contentSelector)[0]
            $curr_top = $(content).children().first()
            $(content).prepend frag
            
            # previously, we would pass in the new DOM element as context for the callback
            # however we're now using a documentfragment, which doesn't have parents or children,
            # so the context is the contentContainer guy, and we pass in an array
            # of the elements collected as the first argument.
            data = children.get()

            # debugger
         
        # loadingEnd function
        opts.loading.finished.call content, opts
    
        if opts.appendCallback 
            prev_height = 0
            $curr_top.prevAll().each -> prev_height += $(@).outerHeight(true);

            $("html, body").scrollTop prev_height

        # once the call is done, we can allow it again.
        opts.state.isDuringAjax = false
        callback this, data, url
        @._prefill() if opts.prefill   

        
$.extend $.infinitescroll.prototype, reverse_scroll
