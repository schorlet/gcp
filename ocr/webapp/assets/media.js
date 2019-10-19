
const media = {
	streaming: false,
	maxWidth: 640,
	video: null,
	canvas: null,
	file: null,
	dropzone: null,
	photo: null
}

function initMedia() {
	media.video = document.getElementById("video");
    media.canvas = document.getElementById("canvas");
	media.file = document.getElementById("file");
	media.dropzone = document.getElementById("drop_zone");
    media.photo = document.getElementById("photo");

	media.video.addEventListener("canplay", canplay, false);
	media.file.addEventListener("change", fileChange, false);
	media.photo.addEventListener("error", clearCanvas, false);

	media.dropzone.addEventListener("dragenter", dragEnter, false);
	media.dropzone.addEventListener("dragover", dragOver, false);
	media.dropzone.addEventListener("drop", dropFile, false);
	media.dropzone.addEventListener("click", dropClick, false);
	media.dropzone.addEventListener("touch", dropClick, false);
}

function startVideo() {
	const options = {
		audio: false,
		video: {
			width: { min: 640 },
			height: { min: 640 },
			facingMode: { exact: "environment" }
		}
	};
	// facingMode: { exact: "environment" }

	return supportedConstraints()
	.then(function() {
		return navigator.mediaDevices.getUserMedia(options);
	})
	.then(function(stream) {
		media.video.srcObject = stream;
		return media.video.play();
	})
	.then(function() {
		const button = document.getElementById("take_photo");
		button.addEventListener("click", takePhoto, false);
		button.addEventListener("touch", takePhoto, false);
		button.removeAttribute("disabled");
	});
}

function supportedConstraints() {
	return new Promise(function(resolve, reject) {
		// TODO: support for <video> and <canvas>

		if (navigator.mediaDevices === undefined ||
			navigator.mediaDevices.getUserMedia === undefined) {
			reject(Error("Browser is too old"));

		} else if (navigator.mediaDevices.getSupportedConstraints === undefined ||
			navigator.mediaDevices.getSupportedConstraints().facingMode === undefined) {
			reject(Error("No facingMode constraint"));

		} else if (!("srcObject" in media.video)) {
			reject(Error("Browser is too old"));

		} else {
			resolve();
		}
	})
}

function canplay() {
	media.streaming = true;
}

function takePhoto(ev) {
	// prevent the click from being handled more than once.
	ev.preventDefault();

	if (!media.streaming) {
		return clearCanvas();
	}

	drawCanvas(media.video);
}

function dragEnter(ev) {
	ev.stopPropagation();
	ev.preventDefault();
}
function dragOver(ev) {
	ev.stopPropagation();
	ev.preventDefault();
}
function dropFile(ev) {
	ev.stopPropagation();
	ev.preventDefault();

	const fileList = ev.dataTransfer.files;
	handleFiles(fileList);
}
function dropClick(ev) {
	media.file.click();
}

function fileChange() {
	const fileList = media.file.files;
	handleFiles(fileList)
}

function handleFiles(fileList) {
	if (fileList.length == 0) {
		return clearCanvas();
	}

	const data = window.URL.createObjectURL(fileList.item(0));
	media.photo.setAttribute("src", data);

	const onload = function() {
		drawCanvas(this);
		window.URL.revokeObjectURL(this.src);
	}
	media.photo.addEventListener("load", onload, {once: true});
}

function drawCanvas(src) {
	const width = src.videoWidth || src.naturalWidth;
	const height = src.videoHeight || src.naturalHeight;

	media.canvas.width = media.maxWidth;
	media.canvas.height = height / (width/media.maxWidth);

	const context = media.canvas.getContext("2d");
	context.drawImage(src, 0, 0, media.canvas.width, media.canvas.height);

	panels.media.dispatchEvent(new Event('preview_ready'));
}

function clearCanvas() {
	const context = media.canvas.getContext("2d");
	context.clearRect(0, 0, media.canvas.width, media.canvas.height);
}

function drawBoundaries(dt) {
	const context = media.canvas.getContext("2d");
	context.lineWidth = 2;
	context.lineJoin = "round";

	const count = dt.annotations.length;
	let i = count > 1 ? 1 : 0;
	for (; i < count; i++) {
		const vertices = dt.annotations[i].vertices;
		context.beginPath();
		context.moveTo(vertices[0].x, vertices[0].y);
		context.lineTo(vertices[1].x, vertices[1].y);
		context.lineTo(vertices[2].x, vertices[2].y);
		context.lineTo(vertices[3].x, vertices[3].y);
		context.closePath();
		context.strokeStyle = selectColor(i, count);
		context.stroke();
	}
}

function selectColor(num, colors) {
    if (colors < 1) colors = 1;
    return "hsl(" + (num * (360 / colors) % 360) + ",100%,50%)";
}
