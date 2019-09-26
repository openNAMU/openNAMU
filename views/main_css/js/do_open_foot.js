function do_open_foot(name) {
    var o_data = document.getElementById('c' + name);
    
    name = name.replace(/\.([0-9]+)$/, '');
    var g_data = document.getElementById(name);

    if(o_data.innerHTML === '') {
        o_data.style.display = 'block';
        o_data.innerHTML += '<div class="foot_in"><a onclick="do_open_foot(\'' + name + '\')" href="#' + name + '">(Go)</a> ' + g_data.innerHTML + '</div><a class="foot_close" onclick="do_open_foot(\'' + name + '\')" href="javascript:void();">(X)</a>';
    } else {
        o_data.style.display = 'none';
        o_data.innerHTML = '';
    }
}