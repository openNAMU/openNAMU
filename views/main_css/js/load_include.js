function load_include(title, name, p_data) {
    var o_data = document.getElementById(name);

    var url = "/api/w/" + encodeURI(title) + "?include=1";
    var url_2 = "/api/markup";
    
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    var xhr_2 = new XMLHttpRequest();
    xhr_2.open("GET", url_2, true);
    xhr_2.send(null);

    xhr.onreadystatechange = function() {
        if(xhr.readyState === 4 && xhr.status === 200) {
            var o_p_data = JSON.parse(xhr.responseText);
            var g_data = o_p_data['data'];
            
            for(key in p_data) {
                try {
                    patt = new RegExp('@' + p_data[key][0] + '@');
                    g_data = g_data.replace(patt, p_data[key][1]);
                } catch {}
            }

            o_data.innerHTML = g_data;

            js_data = o_p_data['js_data'];
            js_data = js_data.replace(/<script>/g, '');
            js_data = js_data.replace(/<\/script>/g, '\n');

            eval(js_data)
        }
    }
}