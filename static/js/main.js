// allow error responses to swap as we are using this as a signal that
// a form was submitted with bad data and want to rerender with the
// error as a flash message
document.body.addEventListener("htmx:beforeSwap", function(evt) {
	// unset current flash message
	htmx.find("#flash").innerHTML = "";

	// set new flash message upon error
	if (evt.detail.isError) {
		evt.detail.shouldSwap = true;
		evt.detail.target = htmx.find("#flash");
	}
});

// disable save button before each request
document.body.addEventListener("htmx:beforeRequest", function(evt) {
	htmx.find("#save-button").disabled = true;
});

// enable save button after successful searches
document.body.addEventListener("htmx:afterRequest", function(evt) {
	let path = evt.detail.requestConfig.path;
	if (evt.detail.successful && path == "/search") {
		htmx.find("#save-button").disabled = false;
	}
});
