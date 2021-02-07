// Tool
function do_url_change(data) {
    return encodeURIComponent(data);
}

// Main
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
    var toc_data = '<div id="toc"><div id="toc_title">TOC</div><br>';
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
        
        var heading_level_string_no_end = heading_level_string.replace(/\.$/, '');

        toc_data += '' +
            '<span style="margin-left: ' + String((heading_level_string.match(/\./g).length - 1) * 10) + 'px;">' +
                '<a href="#s-' + heading_level_string_no_end + '">' + 
                    heading_level_string + ' ' +
                '</a>' + heading_data[2] +
            '</span>' +
            '<br>' +
        ''
        data = data.replace(heading_re, 
            '<h' + heading_level + ' id="s-' + heading_level_string_no_end + '">' + 
                '<a href="#toc">' + heading_level_string + '</a> ' + heading_data[2] + 
            '</h' + heading_level + '>' +
            '<br>'
       );
    }
    
    data = data.replace(/(<\/h[0-9]>)<br>/g, '$1');
    data = data.replace(/\[(?:toc|목차)\]/g, toc_data + '</div>');
    
    return data;
}

function do_onmark_link_render(data, data_js, name_doc, name_include) {
    var link_num = 0;
    data = data.replace(/\[\[(((?!\]\]).)+)\]\]/g, function(x, x_1) {
        var link_split = x_1.split('|');
        var link_real = link_split[0];
        var link_out = link_split[1] ? link_split[1] : link_split[0];
        
        link_num += 1;
        var link_num_str = String(link_num - 1);
        
        if(link_real.match(/^http(s)?:\/\//)) {
            var i = 0;
            while(i < 2) {
                if(i === 0) {
                    var var_link_type = 'href';
                } else {
                    var var_link_type = 'title';
                }
                
                data_js += '' +
                    'document.getElementsByName("' + name_include + 'set_link_' + link_num_str + '")[0].' + var_link_type + ' = ' + 
                        '"' + link_real.replace(/"/g, '\\"') + '";' +
                    '\n' +
                '';
                
                i += 1;
            }
            
            return  '<a id="out_link" ' +
                        'name="' + name_include + 'set_link_' + link_num_str + '" ' + 
                        'title=""' +
                        'href="">' + link_out + '</a>';
        } else {
            var i = 0;
            while(i < 2) {
                if(i === 0) {
                    var var_link_type = 'href';
                    var var_link_data = '/w/' + do_url_change(link_real);
                } else {
                    var var_link_type = 'title';
                    var var_link_data = link_real.replace(/"/g, '\\"');
                }
                
                data_js += '' +
                    'document.getElementsByName("' + name_include + 'set_link_' + link_num_str + '")[0].' + var_link_type + ' = ' + 
                        '"' + var_link_data + '";' +
                    '\n' +
                '';
                
                i += 1;
            }
            
            return  '<a class="' + name_include + 'link_finder" ' +
                        'name="' + name_include + 'set_link_' + link_num_str + '" ' +
                        'title="" ' +
                        'href="">' + link_out + '</a>';
        }
    });
    
    return [data, data_js];
}

function do_onmark_render(name_id, name_include = '', name_doc = '') {
    var data = document.getElementById(name_id).innerHTML;
    var data_js = '';
    data = '<br>' + data.replace(/\n/g, '<br>') + '<br>';
    
    data = do_onmark_text_render(data);
    data = do_onmark_heading_render(data);
    
    var var_data = do_onmark_link_render(data, data_js, name_doc, name_include);
    data = var_data[0];
    data_js = var_data[1];
    
    data = data.replace(/^(<br>| )+/, '');
    data = data.replace(/(<br>| )+$/, '');
    
    data_js += '' + 
        'get_link_state("' + name_include + '");\n' + 
        'get_file_state("' + name_include + '");\n' + 
    ''
    data_js = 'render_html("' + name_include + 'render_contect");\n' + data_js
    
    document.getElementById(name_id).innerHTML = data;
    eval(data_js);
}