function do_onmark_text_render(data) {
    data = data.replace(/'''((?:(?!''').)+)'''/g, '<b>$1</b>');
    data = data.replace(/''((?:(?!'').)+)''/g, '<i>$1</i>');
    data = data.replace(/__((?:(?!__).)+)__/g, '<u>$1</u>');
    data = data.replace(/\^\^((?:(?!\^\^).)+)\^\^/g, '<sup>$1</sup>');
    data = data.replace(/,,((?:(?!,,).)+),,/g, '<sub>$1</sub>');
    data = data.replace(/--((?:(?!--).)+)--/g, '<s>$1</s>');
    data = data.replace(/~~((?:(?!~~).)+)~~/g, '<s>$1</s>');
    
    return data;
}

function do_onmark_heading_render(data) {
    var heading_re = /<br>(={1,6}) ?([^=]+) ?={1,6}<br>/;
    var heading_level_all = [0, 0, 0, 0, 0, 0];
    while(data.match(heading_re)) {
        console.log(data.match(heading_re));
        data = data.replace(heading_re, function(x, heading_level, heading_data) {
            console.log(heading_level);
            heading_level = heading_level.length;
            heading_level_all[heading_level - 1] += 1;

            var i = 6;
            while(i > heading_level - 1) {
                heading_level_all[i] = 0;

                i -= 1;
            }

            console.log(heading_level_all);
            heading_level = String(heading_level);
            heading_level_string = heading_level_all.join('.');
            console.log(heading_level_string);
            
            return '<h' + heading_level + '>' + heading_data + '</h' + heading_level + '><br>';
        });
    }
    
    console.log(data.match(heading_re));
    
    return data;
}

function do_onmark_render(id_name) {
    var data = document.getElementById(id_name).innerHTML;
    data = '<br>' + data.replace(/\n/g, '<br>') + '<br>';
    
    data = do_onmark_text_render(data);
    data = do_onmark_heading_render(data);
    
    document.getElementById(id_name).innerHTML = data;
}