function do_open_foot(name, num = 0) {
    var found_include = name.match(/^(include_(?:[0-9]+)\-)/);
    if(found_include) {
        var include_name = name.replace(/^(?:include_(?:[0-9]+)\-)/, '');
        var front_data = found_include[1];
    } else {
        var include_name = name;
        var front_data = '';
    }

    if(
        document.cookie.match(main_css_regex_data('main_css_footnote_set')) &&
        document.cookie.match(main_css_regex_data('main_css_footnote_set'))[1] === '1'
    ) {
        if(num === 1) {
            document.getElementById(front_data + 'r' + include_name).focus();
        } else {
            var get_data = document.getElementById(front_data + include_name).innerHTML;
            var org_data = document.getElementById(front_data + 'd' + include_name).innerHTML;
            if(org_data === '') {
                document.getElementById(front_data + 'd' + include_name).innerHTML = '' +
                    '<a href="#' + front_data + 'c' + include_name + '">(Go)</a> ' + get_data +
                '';
                document.getElementById(front_data + 'd' + include_name).className = 'spead_footnote';
            } else {
                document.getElementById(front_data + 'd' + include_name).innerHTML = '';
                document.getElementById(front_data + 'd' + include_name).className = '';
            }
        }
    } else {
        document.getElementById(front_data + 'r' + include_name).style.color = 'red';
        document.getElementById(front_data + 'c' + include_name).style.color = (num === 1 ? 'inherit' : 'red');
        document.getElementById(front_data + (num === 1 ? 'r' : 'c') + include_name).focus();
    }
}