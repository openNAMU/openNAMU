function load_preview(name) {
    var o_data = document.getElementById('content');
    var p_data = document.getElementById('see_preview');

    var s_data = new FormData();
    s_data.append('data', o_data.value);

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

            p_data.innerHTML = o_p_data['data'];

            js_data = o_p_data['js_data'];
            js_data = js_data.replace(/<script>/g, '');
            js_data = js_data.replace(/<\/script>/g, '\n');

            eval(js_data)
        }
    }
}