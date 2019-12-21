function render_html(name = '') {
    var num = 0;
    while(1) {
        num += 1

        if(document.getElementById(name + '_' + String(num))) {
            data = document.getElementById(name + '_' + String(num)).innerHTML;

            var t_data = ['b', 'i', 's', 'del', 'strong', 'bold', 'em', 'sub', 'sup']
            for(var key in t_data) {
                var patt = new RegExp('&lt;' + t_data[key] + '&gt;((?:(?!&lt;\/' + t_data[key] + '&gt;).)*)&lt;\/' + t_data[key] + '&gt;', 'ig');
                data = data.replace(patt, '<' + t_data[key] + '>$1</' + t_data[key] + '>');
            }

            var src_list = {
                'www.youtube.com' : '1',
                'www.google.com' : '1'
            }
            data = data.replace(/&lt;iframe( (?:(?:(?!&gt;).)+))&gt;&lt;\/iframe&gt;/ig, function(full, in_data) {
                var src_data = in_data.match(/ src=['"]https:\/\/([^/'"]+)(?:[^'"]+)['"](?: |$)/);
                if(src_data) {
                    if(src_list[src_data[1]]) {
                        return '<iframe' + in_data + '></iframe>';
                    } else {
                        return full;
                    }
                }
            });

            t_data = ['div', 'span']
            for(var key in t_data) {
                patt = new RegExp('&lt;' + t_data[key] + ' style=["\']((?:(?!["\']).)+)["\']&gt;((?:(?!&lt;\/' + t_data[key] + '&gt;).)*)&lt;\/' + t_data[key] + '&gt;', 'ig');
                data = data.replace(patt, function(full, in_data, in_data_2) {
                    return '<' + t_data[key] + ' style="' + in_data.replace(/position/ig, '') + '">' + in_data_2 + '</' + t_data[key] + '>'
                });
            }

            data = data.replace(/&lt;a href=["\']((?:(?!["\']).)+)["\']&gt;((?:(?!&lt;\/a&gt;).)*)&lt;\/a&gt;/ig, function(full, in_data, in_data_2) {
                return '<a id="out_link" href="' + in_data.replace(/^javascript/ig, '') + '">' + in_data_2 + '</a>'
            });

            document.getElementById(name + '_' + String(num)).innerHTML = data;
        } else {
            break;
        }
    }
}