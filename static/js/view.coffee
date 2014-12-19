# Log view popup logic
$ ->
  $("body").on "click", "a.view", (e) ->
    e.preventDefault()
    el = this
    $("#recordViewLabel").html $(this).data("type")
    $("#recordViewDateTime").html $(this).data("datetime")
    $("#viewRecordModal .btn-danger").unbind("click").click ->
      if confirm("Are you sure want to delete this event?")
        $.ajax
          url: $(el).attr("href")
          type: "DELETE"
          success: ->
            $(".modal .close").click()
            $(el).parents("tr").css "opacity", "0.2"
          error: ->
            alert "Error: Record not deleted."
      else
        e.preventDefault()

    $.getJSON($(this).attr("href"), (data) ->
      $(".modal-body").JSONView data
      $("#viewRecordModal").modal()
    ).fail ->
      $(".modal-body").html "Error: Record not found or wrong JSON structure."
      $("#viewRecordModal").modal()
