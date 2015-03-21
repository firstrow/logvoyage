$ ->
  # On submit log search form send ajax request
  # Enable ajax search only on start page.
  # If user is on other page will work regular GET request
  p = window.location.pathname
  if p != "/" and p.substring(0, 9) != "/project/"
    return

  $("#searchForm").attr 'action', p
  $("#searchForm").submit (e) ->
    e.preventDefault()
    $.ajax
      type: "GET"
      url: $(this).attr("action")
      data: $(this).serialize()
      success: (data) ->
        $("#logTableContainer").html data
      complete: ->
        # Search is really fast, we should add delay
        setTimeout (->
          Ladda.stopAll()
          $("html, body").animate(scrollTop: 0, "fast")
        ), 300

  $("body").on "click", "#pagination a", (e) ->
    e.preventDefault()
    $("#logTableContainer").load $(this).attr("href"), ->
      $("html, body").animate
        scrollTop: 0
      , "fast"
