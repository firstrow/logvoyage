// Initialize components
$(function() {
	jQuery('#time_start, #time_stop').datetimepicker();
	$("#pagination").center();
});

$(window).resize(function(){
	$("#pagination").center();
});

jQuery.fn.center = function () {
	var left = Math.max(0, (($(window).width() - $(this).outerWidth()) / 2) + $(window).scrollLeft());
    this.css("left", left + "px");
    return this;
}