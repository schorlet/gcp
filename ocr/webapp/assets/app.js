// ♥ thank you MDN

document.addEventListener("DOMContentLoaded", initMedia);
document.addEventListener("DOMContentLoaded", initApp);

const panels = {
	media: null,
	text: null,
	overlay: null
}

const steps = {
	video: null,
	file: null,
	preview: null,
	overlay: null
};

function initApp() {
	panels.media = document.getElementById("media_panel");
	panels.text = document.getElementById("text_panel");
	panels.overlay = document.getElementById("overlay_panel");

	steps.video = document.getElementById("video_step");
	steps.file = document.getElementById("file_step");
	steps.preview = document.getElementById("preview_step");
	steps.overlay = document.getElementById("overlay_step");

	const reset = document.getElementById("reset");
	reset.addEventListener("click", resetSteps, false);
	reset.addEventListener("touch", resetSteps, false);

	panels.media.addEventListener("preview_ready", showPreview, false);
	panels.media.addEventListener("preview_ready", doUpload, false);

	startVideo()
	.then(function() {
		showVideo();
	})
	.catch(function(err) {
		console.error("Enable to capture environment video", err);
		showFile();
	});
}

function resetSteps(ev) {
	ev.preventDefault();
	if (media.streaming) {
		showVideo()
	} else {
		showFile();
	}
}

function showVideo() {
	show(steps.video);
	hide(steps.file);
	hide(steps.preview);
}
function showFile() {
	hide(steps.video);
	show(steps.file);
	hide(steps.preview);
}
function showPreview() {
	hide(steps.video);
	hide(steps.file);
	show(steps.preview);
}

function show(step) {
	step.removeAttribute("hidden");
}
function hide(step) {
	step.setAttribute("hidden", "");
}

function showOverlay(message) {
	panels.text.style.opacity = 0.6;
	panels.media.style.opacity = 0.6;
	panels.overlay.removeAttribute("hidden");
	steps.overlay.innerText = `⚙ ${message}`;
}
function hideOverlay() {
	panels.text.style.opacity = 1;
	panels.media.style.opacity = 1;
	panels.overlay.setAttribute("hidden", "");
}
