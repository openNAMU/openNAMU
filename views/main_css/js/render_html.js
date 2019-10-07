function render_html(name = '') {
    if(name === '') {
        data = document.getElementById('render_contect').innerHTML;
    } else {
        data = document.getElementById(name).innerHTML;
    }

    t_data = ['b', 'i', 's', 'del']
    for(var key in t_data) {
        var patt = new RegExp('&lt;' + t_data[key] + '&gt;((?:(?!&lt;\/' + t_data[key] + '&gt;).)*)&lt;\/' + t_data[key] + '&gt;', 'ig');
        data = data.replace(patt, '<' + t_data[key] + '>$1</' + t_data[key] + '>');
    }
    
    src_list = {
        'www.youtube.com' : '1'
    }
    data = data.replace(/&lt;iframe( (?:(?:(?!&gt;).)+))&gt;&lt;\/iframe&gt;/g, function(full, in_data) {
        src_data = data.match(/ src=['"]https:\/\/([^/'"]+)(?:[^'"]+)['"](?: |$)/);
        if(src_data) {
            if(src_list[src_data[1]]) {
                return '<iframe' + in_data + '></iframe>';
            } else {
                return full;
            }
        }
    });
    
    if(name === '') {
        document.getElementById('render_contect').innerHTML = data;
    } else {
        document.getElementById(name).innerHTML = data;
    }
}