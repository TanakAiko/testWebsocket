const socket = new WebSocket("ws://localhost:8080/ws")

socket.onopen = (event) => {
    console.log("WebSocket connection opened");
}

socket.onmessage = (event) => {
    const output = document.getElementById('output');
    output.innerHTML += `<p>Received: ${event.data}</p>`;
}

function sendMessage() {
    const name = document.getElementById("playerName");
    const score = document.getElementById("playerScore");
    const time = document.getElementById("playerTime");

    const message = JSON.stringify({ name: name.value, score: +score.value, time: +time.value })
    socket.send(message);
    
    name.value = "";
    score.value = "";
    time.value = "";
}
