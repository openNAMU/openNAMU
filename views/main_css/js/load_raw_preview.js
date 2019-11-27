function load_raw_preview(name_1, name_2) {
    var get = document.getElementById(name_1);
    var send = document.getElementById(name_2);

    send.innerHTML = get.value;
}