// allow error responses to swap as we are using this as a signal that
// a form was submitted with bad data and want to rerender with the
// error as a flash message
//
// set isError to false to avoid error logging in console
document.body.addEventListener("htmx:beforeSwap", function(evt) {
	if (evt.detail.isError) {
		evt.detail.isError = false;
		evt.detail.shouldSwap = true;
		evt.detail.target = htmx.find("#flash");
	}
});
