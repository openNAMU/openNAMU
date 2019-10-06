function do_open_foot(name) {
    var found_include = name.match(/^(include_(?:[0-9]+)\-)/);
    if(found_include) {
        var o_data = document.getElementById(found_include[1] + 'c' + name.replace(/^(?:include_(?:[0-9]+)\-)/, ''));
    } else {
        var o_data = document.getElementById('c' + name);
    }
    
    var g_data = document.getElementById(name.replace(/\.([0-9]+)$/, ''));

    if(o_data.innerHTML === '') {
        o_data.style.display = 'block';
        o_data.innerHTML += '' + 
            '<div class="foot_in">' + 
                '<a onclick="do_open_foot(\'' + name + '\')" href="#' + name + '">(Go)</a> ' + g_data.innerHTML + 
            '</div> <a class="foot_close" onclick="do_open_foot(\'' + name + '\')" href="javascript:void(0);">(X)</a>' + 
        '';
    } else {
        o_data.style.display = 'none';
        o_data.innerHTML = '';
    }
}