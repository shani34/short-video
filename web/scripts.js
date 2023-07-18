// JavaScript code for video recording functionality
let stream;
let mediaRecorder;
const recordedChunks = [];

// Access the user's camera and microphone
navigator.mediaDevices.getUserMedia({ video: true, audio: true })
    .then(function (streamObj) {
        stream = streamObj;
        const videoElement = document.getElementById('videoElement');
        videoElement.srcObject = stream;
    })
    .catch(function (error) {
        console.error('Error accessing the camera and microphone:', error);
    });

// Record the video
document.getElementById('recordButton').addEventListener('click', function () {
    recordedChunks.length = 0;
    mediaRecorder = new MediaRecorder(stream);

    mediaRecorder.ondataavailable = function (event) {
        if (event.data.size > 0) {
            recordedChunks.push(event.data);
        }
    };

    mediaRecorder.start();
});

// Stop the recording
document.getElementById('stopButton').addEventListener('click', function () {
    mediaRecorder.stop();
});

// Share the video
document.getElementById('shareButton').addEventListener('click', function () {
    const blob = new Blob(recordedChunks, { type: 'video/webm' });
    const formData = new FormData();
    formData.append('video', blob, 'myvideo.webm');

    // Send the video to the backend for processing and storage
    fetch('/upload', { method: 'POST', body: formData })
        .then(function (response) {
            if (response.ok) {
                alert('Video shared successfully!');
            } else {
                alert('Error sharing the video.');
            }
        })
        .catch(function (error) {
            console.error('Error sharing the video:', error);
        });
});
