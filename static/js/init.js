// Initialize components
$(function() {
	jQuery('#time_start, #time_stop').datetimepicker();
	$("#pagination").center();
	$("#time").multiselect();

	$("#time").change(function(){
		if ($(this).val() == 'custom') {
			$(".timebox").show();
		} else {
			$(".timebox").hide();
		}
	});
	$("#time").change();
});

$(window).resize(function(){
	$("#pagination").center();
});

jQuery.fn.center = function () {
	var left = Math.max(0, (($(window).width() - $(this).outerWidth()) / 2) + $(window).scrollLeft());
    this.css("left", left + "px");
    return this;
}