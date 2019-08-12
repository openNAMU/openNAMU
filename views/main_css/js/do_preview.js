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
            g_data = JSON.parse(this.responseText)['data'];
            p_data.innerHTML = g_data;
            
            while(1) {
                m_data = g_data.match(/<script>((?:(?!<\/script>).)+)<\/script>/);
                if(m_data) {
                    eval(m_data[1]);
                    
                    g_data = g_data.replace(/<script>((?:(?!<\/script>).)+)<\/script>/, '', 1);
                } else {
                    break;
                }
            }
        }
    }
}