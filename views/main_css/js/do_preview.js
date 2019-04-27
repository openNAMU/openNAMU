function do_preview(name) {
    var o_data = document.getElementById('content');
    var p_data = document.getElementById('see_preview');

    var s_data = new FormData();
    s_data.append('data', o_data.value);

    var url = "/api/w/" + name;
    
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send(s_data);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            p_data.innerHTML = JSON.parse(this.responseText)['data'];
        }
    }
}