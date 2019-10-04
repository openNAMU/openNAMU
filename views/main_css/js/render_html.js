function render_html() {
    data = document.getElementById('render_contect').innerHTML;

    t_data = ['b', 'i', 's', 'del']
    for(var key in t_data) {
        var patt = new RegExp('&lt;' + t_data[key] + '&gt;((?:(?!&lt;\/' + t_data[key] + '&gt;).)*)&lt;\/' + t_data[key] + '&gt;', 'ig');
        data = data.replace(patt, '<' + t_data[key] + '>$1</' + t_data[key] + '>');
    }
    
    // document.getElementById('render_contect').innerHTML = data;
}