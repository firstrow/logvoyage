// On submit log search form send ajax request
$(function() {
	$("#searchForm").submit(function(e) {
		e.preventDefault();
		$.ajax({
			type: "GET",
			url: $(this).attr('action'),
			data: $(this).serialize(), // serializes the form's elements.
			success: function(data) {
				$("#logTableContainer").html(data);
			},
			complete: function() {
				// Search are really fast
				// we should add delay
				setTimeout(function() {
					Ladda.stopAll();
				}, 300);
			}
		});
	});
});