function render_html(name = '') {
    var num = 0;
    while(1) {
        num += 1;

        if(document.getElementById(name + '_' + String(num))) {
            data = document.getElementById(name + '_' + String(num)).innerHTML;

            var src_list = ['www.youtube.com', 'www.google.com', 'play-tv.kakao.com'];
            var t_data = [
                'b', 'i', 's', 'del', 'strong', 'bold', 'em', 'sub', 'sup', 
                'div', 'span', 
                'a',
                'iframe'
            ];
            for(var key in t_data) {
                patt = new RegExp(
                    '&lt;' + t_data[key] + '( (?:(?:(?!&gt;).)+))?&gt;((?:(?!&lt;\/' + t_data[key] + '&gt;).)*)&lt;\/' + t_data[key] + '&gt;',
                    'ig'
                );
                
                data = data.replace(patt, function(full, in_data, in_data_2) {
                    if(['b', 'i', 's', 'del', 'strong', 'bold', 'em', 'sub', 'sup'].includes(t_data[key])) {
                        return '<' + t_data[key] + '>' + in_data_2 + '</' + t_data[key] + '>'
                    } else if(t_data[key] === 'div' || t_data[key] === 'span') {
                        var style_data = in_data.match(/ style=['"]([^'"]*)['"]/);
                        if(style_data) {
                            style_data = style_data[1].replace(/position/ig, '');
                        } else {
                            style_data = '';
                        }

                        return '<' + t_data[key] + ' style="' + style_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                    } else if(t_data[key] === 'a') {
                        var link_data = in_data.match(/ href=['"]([^'"]*)['"]/);
                        if(link_data) {
                            link_data = link_data[1].replace(/^javascript:/ig, '');
                        } else {
                            link_data = '';
                        }

                        return '<' + t_data[key] + ' id="out_link" href="' + link_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                    } else if(t_data[key] === 'iframe') {
                        var src_data = in_data.match(/ src=['"]([^'"]*)['"]/);
                        if(src_data) {
                            src_data = src_data[1];

                            var src_check = src_data.match(/^http(?:s)?:\/\/([^/]+)/);
                            if(src_check) { 
                                if(!src_list.includes(src_check[1])) {
                                    src_data = '';
                                }
                            } else {
                                src_data = '';
                            }
                        } else {
                            src_data = '';
                        }

                        var width_data = in_data.match(/ width=['"]([^'"]*)['"]/);
                        if(width_data) {
                            width_data = width_data[1];
                        } else {
                            width_data = '';
                        }

                        var height_data = in_data.match(/ height=['"]([^'"]*)['"]/);
                        if(height_data) {
                            height_data = height_data[1];
                        } else {
                            height_data = '';
                        }

                        return '<' + t_data[key] + ' src="' + src_data + '" width="' + width_data + '" height="' + height_data + '" allowfullscreen frameborder="0">' + in_data_2 + '</' + t_data[key] + '>';
                    } else {
                        var src_data = in_data.match(/ src=['"]([^'"]*)['"]/);
                        if(src_data) {
                            src_data = src_data[1];
                        } else {
                            src_data = '';
                        }

                        var width_data = in_data.match(/ width=['"]([^'"]*)['"]/);
                        if(width_data) {
                            width_data = width_data[1];
                        } else {
                            width_data = '';
                        }

                        var height_data = in_data.match(/ height=['"]([^'"]*)['"]/);
                        if(height_data) {
                            height_data = height_data[1];
                        } else {
                            height_data = '';
                        }

                        return '<' + t_data[key] + ' controls src="' + src_data + '" width="' + width_data + '" height="' + height_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                    }
                });
            }

            document.getElementById(name + '_' + String(num)).innerHTML = data;
        } else {
            break;
        }
    }
}