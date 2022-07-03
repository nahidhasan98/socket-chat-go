console.log("script linked successfully");

$(document).ready(() => {
    $('#message').keyup((event) => {
        if (event.shiftKey && event.keyCode == 13) {
            console.log("shift+enter pressed");
        } else if (event.keyCode == 13) {
            console.log($('#message').val().trim());
            $('#message').val("");
            //submit
        }
    });

    let socket = new WebSocket("ws://" + location.host + "/ws");
    console.log("Attempting Connection...");

    socket.onopen = () => {
        console.log("Successfully Connected");
        socket.send("Hi From the Client!")
    };

    socket.onclose = (event) => {
        console.log("Socket Closed Connection: ", event);
        socket.send("Client Closed!")
    };

    socket.onerror = (error) => {
        console.log("Socket Error: ", error);
    };
});