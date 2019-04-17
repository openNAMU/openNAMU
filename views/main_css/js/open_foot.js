function open_foot(name) {
    var o_data = document.getElementById('c' + name);
    
    name = name.replace(/\.([0-9]+)$/, '');
    var g_data = document.getElementById(name);

    if(o_data.innerHTML === '') {
        o_data.innerHTML += '<sup><a onclick="open_foot(\'' + name + '\')" href="#' + name + '">(Go)</a></sup> ' + g_data.innerHTML;
    } else {
        o_data.innerHTML = '';
    }
}