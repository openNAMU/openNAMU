function do_preview(name) {
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
            p_data.innerHTML = JSON.parse(xhr.responseText)['data'];

            xhr_2.onreadystatechange = function() {
                if(xhr_2.readyState === 4 && xhr_2.status === 200) {
                    markup = JSON.parse(xhr_2.responseText)['markup'];

                    if(markup === 'markdown') {
                        render_markdown();
                    }
                }
            }
        }
    }
}