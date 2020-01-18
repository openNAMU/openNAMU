function do_open_foot(name, num = 0) {
    var found_include = name.match(/^(include_(?:[0-9]+)\-)/);
    if(found_include) {
        var include_name = name.replace(/^(?:include_(?:[0-9]+)\-)/, '');
        var front_data = found_include[1];
    } else {
        var include_name = name;
        var front_data = '';
    }

    document.getElementById(front_data + 'r' + include_name).style.color = 'red';
    document.getElementById(front_data + 'c' + include_name).style.color = (num === 1 ? 'inherit' : 'red');
}