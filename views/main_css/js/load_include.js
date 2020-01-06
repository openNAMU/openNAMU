function load_include(title, name, p_data) {
    var o_data = document.getElementById(name);

    var change = '';
    for(key in p_data) {
        change += '@' + p_data[key][0].replace('&', '<amp>') + '@,' + p_data[key][1].replace(',', '<comma>').replace('&', '<amp>') + ','
    }
    
    var url = "/api/w/" + encodeURI(title) + "?include=" + name + "&change=" + change;

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var o_p_data = JSON.parse(this.responseText);
            var g_data = o_p_data['data'];

            o_data.innerHTML = g_data;

            js_data = o_p_data['js_data'];
            js_data = js_data.replace(/<script>/g, '');
            js_data = js_data.replace(/<\/script>/g, '\n');

            eval(js_data)
        }
    }
}