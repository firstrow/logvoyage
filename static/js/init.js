// Initialize components
$(function() {
	jQuery('#time_start, #time_stop').datetimepicker();
	$("#pagination").center();
	// Multiselect
	$("#time").multiselect();
	$("#logType").multiselect({
        enableClickableOptGroups: true,
		buttonWidth: '200px',
		numberDisplayed: 2,
		nonSelectedText: 'All sources'
    });

	$("#time").change(function(){
		if ($(this).val() == 'custom') {
			$(".timebox").show(100);
			//$("#time_start").focus();
		} else {
			$(".timebox").hide(100);
		}
	});
	$("#time").change();

	$('[data-toggle="tooltip"]').tooltip()
});

$(window).resize(function(){
	$("#pagination").center();
});

jQuery.fn.center = function () {
	var left = Math.max(0, (($(window).width() - $(this).outerWidth()) / 2) + $(window).scrollLeft());
    this.css("left", left + "px");
    return this;
}