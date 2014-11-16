// Helps to view log records
$(function() {
	$("a.view").click(function(event) {
		event.preventDefault();
		$.getJSON($(this).attr("href"), function(data) {
			$(".modal-body").JSONView(data);
			$("#viewRecordModal").modal();
		}).fail(function() {
			$(".modal-body").html("Error: Record not found or wrong JSON structure.");
			$("#viewRecordModal").modal();
		});
	});
});