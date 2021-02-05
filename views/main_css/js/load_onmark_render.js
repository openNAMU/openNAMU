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
    while(1) {        
        var heading_data = data.match(heading_re);
        if(!heading_data) {
            break;
        }
          
        var heading_level = heading_data[1].length;
        heading_level_all[heading_level - 1] += 1;

        var i = 6;
        while(i > heading_level - 1) {
            heading_level_all[i] = 0;

            i -= 1;
        }

        heading_level = String(heading_level);
        var heading_level_string = '';
        i = 0;
        while(i < 6) {
            if(heading_level_all[i] !== 0) {
                heading_level_string += String(heading_level_all[i]) + '.';
            }

            i += 1;
        }

        data = data.replace(heading_re, '<h' + heading_level + '>' + heading_level_string + ' ' + heading_data[2] + '</h' + heading_level + '><br>');
    }
    
    data = data.replace(/(<\/h[0-9]>)<br>/g, '$1');
    
    return data;
}

function do_onmark_render(id_name) {
    var data = document.getElementById(id_name).innerHTML;
    data = '<br>' + data.replace(/\n/g, '<br>') + '<br>';
    
    data = do_onmark_text_render(data);
    data = do_onmark_heading_render(data);
    
    data = data.replace(/^(<br>| )+/, '');
    data = data.replace(/(<br>| )+$/, '');
    
    document.getElementById(id_name).innerHTML = data;
}