// Log view popup logic
$(function() {
	$("body").on("click", "a.view", function(e) {
		e.preventDefault();
		var el = this;
		$("#recordViewLabel").html($(this).data("type"));
		$("#recordViewDateTime").html($(this).data("datetime"));
		$("#viewRecordModal .btn-danger").unbind("click").click(function() {
			if (confirm("Are you sure want to delete this event?")) {
				$.ajax({
					url: $(el).attr("href"),
					type: 'DELETE',
					success: function() {
						$(".modal .close").click();
						$(el).parents("tr").css("opacity", "0.2");
					},
					error: function() {
						alert("Error: Record not deleted.")
					}
				});
			} else {
				e.preventDefault();
			}
		});
		$.getJSON($(this).attr("href"), function(data) {
			$(".modal-body").JSONView(data);
			$("#viewRecordModal").modal();
		}).fail(function() {
			$(".modal-body").html("Error: Record not found or wrong JSON structure.");
			$("#viewRecordModal").modal();
		});
	});
});