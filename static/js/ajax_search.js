// On submit log search form send ajax request
$(function() {
	// Enable ajax search only on start page.
	// If user is on other page will work regular GET request
	if (window.location.pathname != "/") {
		return;
	}
	
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
					$("html, body").animate({ scrollTop: 0 }, "fast");
				}, 300);
			}
		});
	});

	$("body").on("click", "#pagination a", function(e) {
		e.preventDefault();
		$("#logTableContainer").load($(this).attr("href"), function(){
			 $("html, body").animate({ scrollTop: 0 }, "fast");
		});
	});
});