const socket = new WebSocket("ws://localhost:8080/ws")

socket.onopen = (event) => {
    console.log("WebSocket connection opened");
}

socket.onmessage = (event) => {
    const output = document.getElementById('output');
    output.innerHTML += `<p>Received: ${event.data}</p>`;
}

function sendMessage() {
    //const messageInput = document.getElementById("messageInput");
    const nameInput = document.getElementById("playerName");
    const scoreInput = document.getElementById("playerScore");
    const timeInput = document.getElementById("playerTime");

    //const message = messageInput.value;
    /* const name = nameInput.value;
    const score = scoreInput.value;
    const time = timeInput.value; */

    const messageToSend = JSON.stringify({ name: nameInput.value, score: +scoreInput.value, time: +timeInput.value })
    socket.send(messageToSend);

    //messageInput.value = "";
    nameInput.value = "";
    scoreInput.value = "";
    timeInput.value = "";
}