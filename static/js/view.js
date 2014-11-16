// Helps to view log records
$(function() {
	$(".view").click(function(event) {
		event.preventDefault();
		$.getJSON($(this).attr("href"), function(data) {
			$(".modal-body").JSONView(data);
			$("#myModal").modal();
		}).fail(function() {
			$(".modal-body").html("Error: Record not found or wrong JSON structure.");
			$("#myModal").modal();
		});
	});
});