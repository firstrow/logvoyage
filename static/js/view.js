// Helps to view log records
$(function() {
	$("body").on("click", "a.view", function(e) {
		e.preventDefault();
		$.getJSON($(this).attr("href"), function(data) {
			$(".modal-body").JSONView(data);
			$("#viewRecordModal").modal();
		}).fail(function() {
			$(".modal-body").html("Error: Record not found or wrong JSON structure.");
			$("#viewRecordModal").modal();
		});
	});
});