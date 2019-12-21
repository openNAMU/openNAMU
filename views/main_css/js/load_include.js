function load_include(title, name, p_data) {
    var o_data = document.getElementById(name);

    var url = "/api/w/" + encodeURI(title) + "?include=" + name;

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var o_p_data = JSON.parse(this.responseText);
            var g_data = o_p_data['data'];

            for(key in p_data) {
                try {
                    var patt = new RegExp('@' + p_data[key][0] + '@', 'g');
                    g_data = g_data.replace(patt, p_data[key][1]);
                } catch(e) {
                    console.log(e);
                }
            }

            o_data.innerHTML = g_data;

            js_data = o_p_data['js_data'];
            js_data = js_data.replace(/<script>/g, '');
            js_data = js_data.replace(/<\/script>/g, '\n');

            eval(js_data)
        }
    }
}