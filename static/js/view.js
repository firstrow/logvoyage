// Helps to view log records
$(function() {
	$("body").on("click", "a.view", function(e) {
		e.preventDefault();
		$("#recordViewLabel").html($(this).data("type"));
		$("#recordViewDateTime").html($(this).data("datetime"));
		$.getJSON($(this).attr("href"), function(data) {
			$(".modal-body").JSONView(data);
			$("#viewRecordModal").modal();
		}).fail(function() {
			$(".modal-body").html("Error: Record not found or wrong JSON structure.");
			$("#viewRecordModal").modal();
		});
	});
});