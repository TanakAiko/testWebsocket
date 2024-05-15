let ws = new WebSocket("ws://localhost:8080/ws");
ws.onopen = function () {
    console.log("Connection is open...");
};

function sendData() {
    let name = document.getElementById("playerName").value;
    let score = document.getElementById("playerScore").value;
    let time = document.getElementById("playerTime").value;

    if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ name: name, score: +score, time: +time }));
    } else {
        console.log("WebSocket is not open. ReadyState: ", ws.readyState);
    }
}

const buttonOk = document.getElementById('submit');
buttonOk.addEventListener("click", (event) => {
    event.preventDefault();
    sendData();
});