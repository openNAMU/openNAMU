function load_preview(name) {
    var s_data = new FormData();
    s_data.append('data', document.getElementById('content').value);

    var url = "/api/w/" + name;
    var url_2 = "/api/markup";

    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send(s_data);

    var xhr_2 = new XMLHttpRequest();
    xhr_2.open("GET", url_2, true);
    xhr_2.send(null);

    xhr.onreadystatechange = function() {
        if(xhr.readyState === 4 && xhr.status === 200) {
            var o_p_data = JSON.parse(xhr.responseText);
            document.getElementById('see_preview').innerHTML = o_p_data['data'];
            eval(o_p_data['js_data'])
        }
    }
}

function load_raw_preview(name_1, name_2) {
    document.getElementById(name_2).innerHTML = document.getElementById(name_1).value;
}