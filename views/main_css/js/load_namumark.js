function get_link_state(data, i = 0) { 
    if(document.getElementsByClassName(data + 'link_finder')[i]) {
        var link_data = document.getElementsByClassName(data + 'link_finder')[i];

        var xhr = new XMLHttpRequest();
        xhr.open("GET", link_data.href.replace('/w/', '/api/w/').replace(/#([^#]*)/, '') + "?exist=1", true);
        xhr.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(JSON.parse(this.responseText)['exist'] !== '1') {
                    document.getElementsByClassName(data + 'link_finder')[i].id = "not_thing";
                } else {
                    document.getElementsByClassName(data + 'link_finder')[i].id = "";
                }

                get_link_state(data, i + 1);
            }
        }
    }
}

function get_file_state(data, i = 0) {       
    if(document.getElementsByClassName(data + 'file_finder_1')[i]) {
        var file_data = document.getElementsByClassName(data + 'file_finder_1')[i];

        var xhr = new XMLHttpRequest();
        xhr.open("GET", file_data.src.replace('/image/', '/api/image/'), true);
        xhr.send(null);
        
        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(JSON.parse(this.responseText)['exist'] !== '1') {
                    document.getElementsByClassName(data + 'file_finder_1')[i].style = "display: none;";
                } else {
                    document.getElementsByClassName(data + 'file_finder_2')[i].innerHTML = "";
                }
            
                get_file_state(data, i + 1);
            }
        }
    }
}

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