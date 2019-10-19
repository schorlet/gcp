
function doUpload() {
	media.canvas.toBlob(processBlob, "image/png");
}

function processBlob(blob) {
	if (blob.size > 5242880) {
		displayError(`Photo size ${blob.size} bytes, is too big`);
		return;
	}

	getUploadURL(blob)
		.then(uploadPhoto(blob))
		.then(detectText)
		.then(displayText)
		.catch(displayError)
		.then(hideOverlay);
}

function getUploadURL(blob) {
	const contentType = blob.type == "" ? "image/*" : blob.type;

	const headers = new Headers({
		"X-Content-Type": contentType,
		"X-Content-Length": blob.size
	});

	const options = {
		method: 'GET',
		headers: headers
	};

	showOverlay("Getting upload URL...");
	return fetch("/uploadURL", options)
		.then(validateResponse)
		.then(readResponseText);
}

function uploadPhoto(blob) {
	const contentType = blob.type == "" ? "image/*" : blob.type;

	const headers = new Headers({
		"Content-Type": contentType,
		"Content-Length": blob.size,
		"x-goog-content-length-range": "0,5242880",
		"x-goog-if-generation-match": "0",
		"x-goog-storage-class": "STANDARD"
	});

	const options = {
		method: 'PUT',
		headers: headers,
		body: blob
	};

	return function(url) {
		showOverlay("Uploading photo...");
		return fetch(url, options)
			.then(validateResponse)
			.then(function() {
				return url;
			});
	}
}

function detectText(url) {
	const u = new URL(url);
	const arr = u.pathname.split("/");
	const path = arr.pop();
	const retry = [5, 3, 3, 2, 2];

	const fn = function(resolve, reject) {
		fetch("/detectText/"+path)
		.then(validateResponse)
		.then(readResponseJSON)
		.then(resolve)
		.catch(function(err) {
			const delay = retry.pop();
			if (delay) {
				setTimeout(fn, delay*1000, resolve, reject);
			} else {
				reject(new Error("max retry exceeded"));
			}
		});
	};

	showOverlay("Detecting text...");
	return new Promise(fn);
}

function displayText(dt) {
	const parent = panels.text.querySelector("#text_step");
	// remove children
	while (parent.firstChild) {
		parent.removeChild(parent.firstChild);
	}
	if (dt && dt.annotations && dt.annotations.length > 0) {
		// add annotations
		const count = dt.annotations.length;
		let i = count > 1 ? 1 : 0;
		for (; i < count; i++) {
			const span = document.createElement("span");
			span.style.padding = "1px 1px";
			span.style.borderWidth = 1;
			span.style.borderStyle = "solid";
			span.style.borderColor = selectColor(i, count);
			span.appendChild(document.createTextNode(dt.annotations[i].description));
			parent.appendChild(span);
		}
		drawBoundaries(dt);
	} else {
		// no annotation found
		const p = document.createElement("p");
		p.appendChild(document.createTextNode("No text detected!"));
		parent.appendChild(p);
	}
}

function displayError(err) {
	const parent = panels.text.querySelector("#text_step");
	// remove children
	while (parent.firstChild) {
		parent.removeChild(parent.firstChild);
	}
	// add error message
	const p = document.createElement("p");
	p.appendChild(document.createTextNode(err));
	parent.appendChild(p);
}

function validateResponse(response) {
	if (!response.ok) {
		throw Error(response.statusText);
	}
	return response;
}

function readResponseJSON(response) {
	const contentType = response.headers.get('content-type');
	if(contentType && contentType.includes('application/json')) {
		return response.json();
	}
	throw Error("not a json response");
}

function readResponseText(response) {
	// console.log(Array.from(response.headers.entries()));
	return response.text();
}

function logResult(result) {
	console.log(result);
}

function logError(err) {
	console.error("Upload error", err);
}
