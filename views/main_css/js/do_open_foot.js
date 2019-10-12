function do_open_foot(name, num = 0) {
    var found_include = name.match(/^(include_(?:[0-9]+)\-)/);
    if(found_include) {
        var include_name = name.replace(/^(?:include_(?:[0-9]+)\-)/, '');
        if(num === 1) {
            document.getElementById(found_include[1] + 'c' + include_name).style.color = 'inherit';
        } else {
            document.getElementById(found_include[1] + 'c' + include_name).style.color = 'red';
        }
    } else {
        var d_color = document.getElementById('c' + name).style.color;
        if(num === 1) {
            document.getElementById('c' + name).style.color = 'inherit';
        } else {
            document.getElementById('c' + name).style.color = 'red';
        }
    }

    
}