// Tool
function do_url_change(data) {
    return encodeURIComponent(data);
}

function do_nowiki_change(data, data_nowiki) {
    return data.replace(/<span id="((?:.*)(?:nowiki_(?:[^"]+)))"><\/span>/, function(x, x_1) {
        return data_nowiki[x_1];
    });
}

function do_link_change(data, data_nowiki, no_change, data_nowiki) {
    data = data.replace(/^:/, '');
    
    if(no_change === 0) {
        data = data.replace(/^사용자:/, 'user:');
        data = data.replace(/^분류:/, 'category:');
        data = data.replace(/^파일:/, 'file:');
    }
    
    var data_var = data.split('#');
    var link_main = data.replace(/#(.*)$/, '');
    var link_sub = data_var.length !== 1 ? ('#' + data_var[data_var.length - 1]) : '';
    
    link_main = do_nowiki_change(link_main, data_nowiki);
    link_main = do_xss_change(link_main);
    
    return [link_main, link_sub];
}

function do_js_safe_change(data) {
    data = data.replace(/\\/g, '\\\\');
    data = data.replace(/"/g, '\\"');
    
    return data;
}

function do_math_try_insert(name_ob, data) {
    return '' + 
        'try {\n' + 
            'katex.render("' + data + '", document.getElementById(\"' + name_ob + '\"));\n' + 
        '} catch {\n' + 
            'document.getElementById(\"' + name_ob + '\").innerHTML = "<span style=\'color: red;\'>' + data + '</span>";\n' + 
        '}\n' + 
    ''
}

function do_data_try_insert(name_ob, data) {
    return '' +
        'if(document.getElementById("' + name_ob + '")) {\n' + 
            'document.getElementById("' + name_ob + '").innerHTML = "' + data + '";\n' + 
        '}\n' +
    ''
}

function do_return_date() {
    var today_data = new Date();

    return '' +
        String(today_data.getFullYear()) + '-' + 
        ((today_data.getMonth() + 1) < 10 ? '0' : '') + String(today_data.getMonth() + 1) + '-' + 
        (today_data.getDate() < 10 ? '0' : '') + String(today_data.getDate()) + ' ' + 
        (today_data.getHours() < 10 ? '0' : '') + String(today_data.getHours()) + ':' + 
        (today_data.getMinutes() < 10 ? '0' : '') + String(today_data.getMinutes()) + ':' + 
        (today_data.getSeconds() < 10 ? '0' : '') + String(today_data.getSeconds()) +
    '';
}

function do_xss_change(data) {
    data = data.replace(/&lt;/g, '<');
    data = data.replace(/&gt;/g, '>');
    data = data.replace(/&amp;/g, '&');
    
    return data;
}

// Sub
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
    var heading_re = /\n(={1,6}) ?([^=]+) ?={1,6}\n/;
    var heading_level_all = [0, 0, 0, 0, 0, 0];
    var toc_data = '<div id="toc"><div id="toc_title">TOC</div>\n';
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
            '\n' +
        ''
        data = data.replace(heading_re, 
            '<h' + heading_level + ' id="s-' + heading_level_string_no_end + '">' + 
                '<a href="#toc">' + heading_level_string + '</a> ' + heading_data[2] + 
            '</h' + heading_level + '>' +
            '\n'
       );
    }
    
    toc_data += '</div>';
    
    data = data.replace(/\[(?:toc|목차)\]/g, toc_data);
    
    var toc_auto_add = data.match(/\[(?:목차|toc)\(no\)\]/);
    if(toc_auto_add) {
        data = data.replace(/\[(?:목차|toc)\(no\)\]/g, '');
    } else {
        data = data.replace(/(<h[1-6] (?:[^>]+)>)/, toc_data + '$1');
    }
    
    return data;
}

function do_onmark_link_render(data, data_js, name_doc, name_include, data_nowiki) {
    var num_link = 0;
    var category_data = '';
    var category_re = /^(분류|category):/i;
    var file_re = /^(파일|file|외부|out):/i;
    data = data.replace(/\[\[(((?!\]\]).)+)\]\]/g, function(x, x_1) {
        var link_split = x_1.split('|');
        var link_real = link_split[0];
        var link_out = link_split[1] ? link_split[1] : link_split[0];
        var link_out_2 = link_split[1] ? link_split[1] : '';
        
        num_link += 1;
        var num_link_str = String(num_link - 1);
        if(link_real.match(file_re)) {
            var file_load_type = link_real.match(file_re)[1];
            var file_name = link_real.replace(file_re, '');
                 
            if(file_load_type === '파일' || file_load_type === 'file') {
                var file_type = file_name.split('.');
                file_name = file_type.slice(0, file_type.length - 1).join('.');
                file_type = file_type[file_type.length - 1];
                
                var file_src = do_url_change(file_name) + '.' + file_type;       
                var file_alt = file_name + '.' + file_type;
                var file_exist = 1;
            } else {
                var file_src = file_name;
                var file_alt = file_name;
                var file_exist = 0;
            }
            
            var file_style = '';
            var file_bgcolor = '';
            var file_align = '';
            
            var file_set = link_out_2.split('&amp;');
            var i = 0;
            while(file_set[i]) {
                var file_set_name = file_set[i].split('=');
                var file_set_data = file_set_name[1];
                file_set_name = file_set_name[0];
                
                if(file_set_name === 'width') {
                    file_style += 'width:' + file_set_data + ';';
                } else if(file_set_name === 'height') {
                    file_style += 'height:' + file_set_data + ';';
                } else if(file_set_name === 'bgcolor') {
                    file_bgcolor += 'background:' + file_set_data + ';';
                } else if(file_set_name === 'alt') {
                    file_alt += file_set_data;
                } else if(file_set_name === 'align') {
                    if(file_set_data === 'center') {
                        file_align = 'display: block; text-align: center;';
                    } else {
                        file_align = 'float: ' + file_set_data + ';';
                    }
                } 
                
                i += 1;
            }
            
            return '' +
                '<span style="' + file_align + '">' + 
                    '<span  style="' + file_bgcolor + '" ' +
                            'class="' + name_include + 'file_finder" ' +
                            'under_style="' + file_style + '" ' +
                            'under_alt="' + file_alt + '" ' +
                            'under_src="' + file_src + '" ' +
                            'under_href="' + (file_exist === 0 ? "out_link" : '/upload?name=' + do_url_change(file_name)) + '">' +
                    '</span>' + 
                '</span>' +
            ''
        } else if(link_real.match(category_re)) {
            var category_link = link_real.replace(category_re, '');
            
            category_data = (category_data === '' ? '<div id="cate_all"><div id="cate">Category : ' : category_data);
            category_data += '' +
                '<a class="' + name_include + 'link_finder" ' +
                    'href="/w/category:' + do_url_change(category_link) + '">' +
                    category_link +
                '</a> | ' +
            ''
            
            return '';
        } else if(link_real.match(/^http(s)?:\/\//)) {
            var i = 0;
            while(i < 2) {
                if(i === 0) {
                    var var_link_type = 'href';
                } else {
                    var var_link_type = 'title';
                }
                
                data_js += '' +
                    'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].' + var_link_type + ' = ' + 
                        '"' + do_js_safe_change(link_real) + '";' +
                    '\n' +
                '';
                
                i += 1;
            }
            
            return  '<a id="out_link" ' +
                        'name="' + name_include + 'set_link_' + num_link_str + '" ' + 
                        'title=""' +
                        'href="">' + link_out + '</a>';
        } else {
            var i = 0;
            while(i < 2) {
                if(i === 0) {
                    var link_data_var = do_link_change(link_real, data_nowiki, 0, data_nowiki);
                    var link_main = link_data_var[0];
                    var link_sub = link_data_var[1];
                    
                    var var_link_type = 'href';
                    var var_link_data = '/w/' + do_url_change(link_main) + link_sub;
                } else {
                    var var_link_type = 'title';
                    var var_link_data = do_js_safe_change(link_main) + link_sub;
                }
                
                data_js += '' +
                    'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].' + var_link_type + ' = ' + 
                        '"' + var_link_data + '";' +
                    '\n' +
                '';
                
                i += 1;
            }
            
            return  '<a class="' + name_include + 'link_finder" ' +
                        'name="' + name_include + 'set_link_' + num_link_str + '" ' +
                        'title="" ' +
                        'href="">' + link_out + '</a>';
        }
    });
    
    data += (category_data === '' ? '' : (category_data.replace(/\| $/, '') + '</div></div>'));
    
    return [data, data_js];
}

function do_onmark_footnote_render(data, name_include) {
    var footnote_end_data = '';
    var footnote_all_data = {};
    var footnote_re = /(?:\[\*([^ \]]*)(?: ((?:(?!\]).)+))?\]|\[(footnote|각주)\])/;
    var i = 1;
    while(1) {
        var footnote_data = data.match(footnote_re);
        if(!footnote_data) {
            break;
        }
        
        if(!footnote_data[3]) {
            if(!footnote_data[2]) {
                var footnote_line_data = '';
            } else {
                var footnote_line_data = footnote_data[2];
            }
            
            if(!footnote_data[1]) {
                var footnote_name = String(i);
            } else {
                var footnote_name = footnote_data[1];
            }
            
            if(!footnote_all_data[footnote_name]) {
                footnote_all_data[footnote_name] = footnote_line_data;
            }

            footnote_line_data = footnote_all_data[footnote_name];
            
            footnote_end_data += '' +
                '<li>' +
                    '<a href="javascript:do_open_foot(\'' + name_include + 'fn-' + String(i) + '\', 1);" ' +
                        'id="' + name_include + 'cfn-' + String(i) + '">' +
                        '(' + footnote_name + ')' +
                    '</a> <span id="' + name_include + 'fn-' + String(i) + '">' + footnote_line_data + '</span>' +
                '</li>' +
            '';
            data = data.replace(footnote_re, '' +
                '<sup>' +
                    '<a href="javascript:do_open_foot(\'' + name_include + 'fn-' + String(i) + '\', 0);" ' +
                        'id="' + name_include + 'rfn-' + String(i) + '">' +
                        '(' + footnote_name + ')' +
                    '</a>' +
                '</sup><span id="' + name_include + 'dfn-' + String(i) + '"></span>' +
           '');
            
            i += 1;
        } else {
            if(footnote_end_data !== '') {
                data = data.replace(footnote_re, '<ul id="footnote_data">' + footnote_end_data + '</ul>');    
            }
            
            footnote_end_data = '';
        }
    }
    
    if(footnote_end_data !== '') {
        data += '<ul id="footnote_data">' + footnote_end_data + '</ul>';
    }
    
    return data;
}

function do_onmark_macro_render(data) {
    data = data.replace(/\[([^[\](]+)\(((?:(?!\)\]).)+)\)\]/g, function(x, x_1, x_2) {
        x_1 = x_1.toLowerCase();
        if(x_1 === 'youtube' || x_1 === 'kakaotv' || x_1 === 'nicovideo') {
            var video_code = x_2.match(/^([^,]+)/);
            video_code = video_code ? video_code[1] : '';
            
            var video_width = x_2.match(/,(?: *)width=([0-9]+)/);
            video_width = video_width ? (video_width[1] + 'px') : '640px';
            
            var video_height = x_2.match(/,(?: *)height=([0-9]+)/);
            video_height = video_height ? (video_height[1] + 'px') : '360px';
            
            if(x_1 === 'youtube') {
                var video_start = x_2.match(/,(?: *)start=([0-9]+)/);
                video_start = video_start ? ('?' + video_start[1]) : '';
                
                video_code = video_code.replace(/^https:\/\/www\.youtube\.com\/watch\?v=/, '');
                video_code = video_code.replace(/^https:\/\/youtu\.be\//, '');
                
                var video_src = 'https://www.youtube.com/embed/' + video_code + video_start
            } else if(x_1 === 'kakaotv') {
                video_code = video_code.replace(/^https:\/\/tv\.kakao\.com\/channel\/9262\/cliplink\//, '');
                video_code = video_code.replace(/^http:\/\/tv\.kakao\.com\/v\//, '');
                
                var video_src = 'https://tv.kakao.com/embed/player/cliplink/' + video_code +'?service=kakao_tv'
            } else {
                var video_src = 'https://embed.nicovideo.jp/watch/' + video_code
            }
            
            return '<iframe style="width: ' + video_width + '; height: ' + video_height + ';" src="' + video_src + '" frameborder="0" allowfullscreen></iframe>';
        } else if(x_1 === 'anchor') {
            return '<span id="' + x_2 + '"></span>';
        } else {
            return '<macro_start>' + x_1 + '(' + x_2 + ')<macro_end>';
        }
    });
    
    data = data.replace(/\[([^[*()\]]+)\]/g, function(x, x_1) {
        x_1 = x_1.toLowerCase();
        if(x_1 === 'date') {
            return do_return_date();
        } else if(x_1 === 'clearfix') {
            return '<div style="clear:both"></div>';
        } else if(x_1 === 'br') { 
            return '<br>';
        } else {
            return '<macro_start>' + x_1 + '<macro_end>';
        }
    });
    
    data = data.replace(/<macro_start>/g, '[');
    data = data.replace(/<macro_end>/g, ']');
    
    return data;
}

function do_onmark_middle_render(data, data_js, name_include, data_nowiki, name_doc) {
    var middle_stack = [];
    var middle_re = /(?:{{{([^{} ]*)|(}}}))/;
    
    var syntax_on = 0;
    
    var html_n = 0;
    
    while(1) {
        var middle_data = data.match(middle_re);
        if(!middle_data) {
            break;
        }
        
        if(middle_data[2]) {
            if(middle_stack.length === 0) {
                data = data.replace(middle_re, '<middle_end>');   
            } else {
                data = data.replace(middle_re, middle_stack[middle_stack.length - 1]);    
                middle_stack.pop();
            }
        } else {
            if(middle_stack.includes('</code>')) {
                data = data.replace(middle_re, '<middle_start>' + middle_data[1]);
                middle_stack.push('<middle_end>');
            } else {
                if(middle_data[1].match(/^(?:(#(?:[0-9a-f-A-F]{3}){1,2})|#([a-zA-Z]+))/)) {
                    var color = middle_data[1].match(/^(?:(#(?:[0-9a-f-A-F]{3}){1,2})|#([a-zA-Z]+))/);
                    color = color[1] ? color[1] : color[2];

                    data = data.replace(middle_re, '<span style="color: ' + color + ';">');
                    middle_stack.push('</span>');
                } else if(middle_data[1].match(/^(\+|-)([1-5])/)) {
                    var font = middle_data[1].match(/^(\+|-)([1-5])/);
                    if(font[1] === '+') {
                        var font_size = String(100 + (20 * Number(font[2]))) + '%';
                    } else {
                        var font_size = String(100 - (10 * Number(font[2]))) + '%';
                    }

                    data = data.replace(middle_re, '<span style="font-size: ' + font_size + ';">');
                    middle_stack.push('</span>');
                } else if(middle_data[1] === '#!wiki') {
                    var wiki_re = /{{{#!wiki(?: style=["']([^"']*)["']\n)?/;

                    var wiki = data.match(wiki_re);
                    var wiki_style = wiki[1] ? wiki[1] : '';

                    data = data.replace(wiki_re, '<wiki_start style="' + wiki_style + '">');  
                    middle_stack.push('<wiki_end>');
                } else if(middle_data[1] === '#!html') {
                    html_n += 1;

                    data = data.replace(middle_re, '<span id="' + name_include + 'render_contect_' + String(html_n) + '">');
                    middle_stack.push('</span>');
                } else if(middle_data[1] === '#!folding') {
                    data = data.replace(middle_re, '<div>');
                    middle_stack.push('</div>');
                } else {
                    data = data.replace(middle_re, '<nowiki_start>' + middle_data[1]);
                    middle_stack.push('<nowiki_end>');
                }
            }
        }
    }
    
    while(middle_stack.length !== 0) {
        data += middle_stack[middle_stack.length - 1];
        middle_stack.pop();
    }
    
    data = data.replace(/\n<div_wiki_end>/g, '<div_wiki_end>');
    
    data = data.replace(/<middle_start>/g, '{{{');
    data = data.replace(/<middle_end>/g, '}}}');
    
    var code_re = /<nowiki_start>(\n*(?:(?:(?!<nowiki_start>|<nowiki_end>).)+\n*)+)<nowiki_end>/;
    var code_n = 0;
    while(1) {
        code_n += 1;
        
        var code_data = data.match(code_re);
        if(!code_data) {
            break;
        }
        
        data_nowiki[name_include + 'nowiki_mid_' + String(code_n)] = code_data[1];
        data_js += do_data_try_insert(name_include + 'nowiki_mid_' + String(code_n), do_js_safe_change(code_data[1]));
        data = data.replace(code_re, '<span id="' + name_include + 'nowiki_mid_' + String(code_n) + '"></span>');
    }
    
    var wiki_re = /<wiki_start (?:[^>]+)>(\n*(?:(?:(?!<wiki_start (?:[^>]+)>|<wiki_end>).)+\n*)+)<wiki_end>/;
    var wiki_n = 0;
    while(1) {
        wiki_n += 1;
        
        var wiki_data = data.match(wiki_re);
        if(!wiki_data) {
            break;
        }
        
        wiki_data = do_nowiki_change(wiki_data[1]);
        wiki_data = do_onmark_render('manual', '', name_include + 'wiki_' + String(wiki_n) + '_', name_doc, wiki_data);
        
        data_js += wiki_data[1];
        data = data.replace(wiki_re, wiki_data[0]);
    }
    
    return [data, data_js, data_nowiki];
}

function do_onmark_last_render(data) {
    // middle_render 마지막 처리
    data = data.replace(/<wiki_start /g, '<div ');
    data = data.replace(/<wiki_end>/g, '</div>');
    
    // heading_render 마지막 처리
    data = data.replace(/(<\/h[0-9]>)\n/g, '$1');
    
    // list_render 마지막 처리
    data = data.replace(/(<\/ul>)\n/g, '$1');
    
    // br 마지막 처리
    data = data.replace(/^(\n| )+/, '');
    data = data.replace(/(\n| )+$/, '');
    data = data.replace(/\n/g, '<br>');
    
    return data;
}

function do_onmark_include_render(data, data_js, name_include, data_nowiki) {
    var include_re = /\[include\(((?:(?!\)\]).)+)\)\]/;
    var i = 0;
    while(1) {
        i += 1;
        
        var include_data = data.match(include_re);
        if(!include_data) {
            break;
        }
        
        var include_name = do_nowiki_change(
            include_data[1].match(/^([^,]+)/)[1],
            data_nowiki
        );
        var include_add_re = /, *([^=]+)=((?:(?:(?!\)]|,).)+)+)/;
        var include_add_data = []
        var include_data = include_data[1];
        while(1) {
            var include_add = include_data.match(include_add_re);
            if(!include_add) {
                break;
            }
            
            include_add_data.push([
                include_add[1], 
                do_nowiki_change(include_add[2], data_nowiki)
            ]);
            include_data = include_data.replace(include_add_re, '');
        }
        
        data = data.replace(include_re,
            '<a id="' + name_include + 'include_link" class="include_' + String(i) + '" href="">(' + include_name + ')</a>' +
            '<div id="' + name_include + 'include_' + String(i) + '"></div>'
        );
        
        data_js += 'load_include("' + do_js_safe_change(include_name) + '", "' + name_include + 'include_' + String(i) + '", ' + JSON.stringify(include_add_data) + ');\n'
    }
    
    return [data, data_js];
}

function do_onmark_nowiki_before_render(data, data_js, name_include, data_nowiki) {
    var num_nowiki = 0;
    data = data.replace(/\\(.)/g, function(x, x_1) {
        num_nowiki += 1;
        data_nowiki[name_include + 'nowiki_one_' + String(num_nowiki)] = x_1;
        data_js += do_data_try_insert(name_include + 'nowiki_one_' + String(num_nowiki), do_js_safe_change(x_1));
        return '<span id="' + name_include + 'nowiki_one_' + String(num_nowiki) + '"></span>';
    });
    
    return [data, data_js, data_nowiki, num_nowiki];
}

function do_onmark_table_render_sub(data, data_col) {
    var data_option_all = {
        "div" : "",
        "table" : "",
        "tr" : "",
        "td" : "",
        "col" : data_col,
        "colspan" : "",
        "rowspan" : "",
        "data" : ""
    };
    
    var table_option_re = /&lt;((?:(?!&lt;|&gt;).)+)&gt;/;
    while(1) {
        var no_option = '';
        var data_option = data.match(table_option_re);
        if(!data_option) {
            break;
        }
        
        data_option = data_option[1];
        var data_option_var = data_option.split('=');
        if(data_option_var.length === 2) {
            var table_option_name = data_option_var[0].replace(/ /g, '');
            var table_option_data = data_option_var[1].replace(/[^a-zA-Z0-9]/g, '');
            if(table_option_name === 'tablebgcolor') {
                // table
               data_option_all['table'] += 'background:' + table_option_data + ';';
            } else if(table_option_name === 'tablewidth') {
                data_option_all['table'] += 'width:' + table_option_data + ';';
            } else if(table_option_name === 'tableheight') {
                data_option_all['table'] += 'height:' + table_option_data + ';';
            } else if(table_option_name === 'tablealign') {
                if(table_option_data === 'right') {
                    data_option_all['div'] += 'float:right;';
                } else if(table_option_data === 'center') {
                    data_option_all['div'] += 'margin:auto;';
                    data_option_all['table'] += 'margin:auto;';
                }
            } else if(table_option_name === 'tabletextalign') {
                data_option_all['table'] += 'text-align:' + table_option_data + ';';
            } else if(table_option_name === 'tablecolor') {
                data_option_all['table'] += 'color:' + table_option_data + ';';
            } else if(table_option_name === 'tablebordercolor') {
                data_option_all['table'] += 'border:2px solid ' + table_option_data + ';';
            } else if(table_option_name === 'rowbgcolor') {
                // tr
                data_option_all['tr'] += 'background:' + table_option_data + ';';
            } else if(table_option_name === 'rowtextalign') {
                data_option_all['tr'] += 'text-align:' + table_option_data + ';';
            } else if(table_option_name === 'rowcolor') {
                data_option_all['tr'] += 'color:' + table_option_data + ';';
            } else if(table_option_name === 'colcolor') {
                // col
                data_option_all['col'] += 'color:' + table_option_data + ';';
            } else if(table_option_name === 'colbgcolor') {
                data_option_all['col'] += 'background:' + table_option_data + ';';
            } else if(table_option_name === 'bgcolor') {
                // td
                data_option_all['td'] += 'background:' + table_option_data + ';';
            } else if(table_option_name === 'color') {
                data_option_all['td'] += 'color:' + table_option_data + ';';
            } else if(table_option_name === 'width') {
                data_option_all['td'] += 'width:' + table_option_data + ';';
            } else if(table_option_name === 'height') {
                data_option_all['td'] += 'height:' + table_option_data + ';';
            } else {
                no_option = '<lt>' + data_option + '<gt>';
            }
        } else {
            if(data_option.match(/^-[0-9]+$/)) {
                // span
                data_option_all['colspan'] = data_option.replace('-', '');
            } else if(data_option.match(/^\|[0-9]+$/)) {
                data_option_all['rowspan'] = data_option.replace('|', '');
            } else if(data_option === '(') {
                // align
                data_option_all['td'] += 'text-align:right;';
            } else if(data_option === ':') {
                data_option_all['td'] += 'text-align:center;';
            } else if(data_option === ')') {
                data_option_all['td'] += 'text-align:left;';
            } else {
                no_option = '<lt>' + data_option + '<gt>';
            }
        }
        
        data = data.replace(table_option_re, no_option);
    }
    
    data = data.replace('<lt>', '&lt;');
    data = data.replace('<gt>', '&gt;');
    data_option_all['data'] = data;
    
    return data_option_all;
}

function do_onmark_table_render(data) {
    var table_before_re = /\|\|((?:(?:(?:(?!\|\|).)+)\n)(?:(?:(?:(?!\|\|).)+)\n*)*)\|\|/;
    while(1) {
        if(data.match(table_before_re)) {
            data = data.replace(table_before_re, function(x, x_1) {
                return '||' + x_1.replace(/\n/g, '<t_br>') + '||';
            })
        } else {
            break;
        }
    }
    
    var table_re = /\n((?:(?:(?:\|\|)+)(?:(?:(?:(?!\|\|).)+)|\n))+)\|\|\n/g;
    var table_in_re = /(?:((?:\|\|)+)((?:(?:(?!\|\|).)+)|\n))/;
    data = data.replace(table_re, function(x, x_1) {
        var table_data_real = '';
        var table_data_var = x_1;
        var table_new_line = 0;
        var table_col = 0;
        var table_data_col = {};
        while(1) {
            var table_data_in = table_data_var.match(table_in_re);
            if(!table_data_in) {
                break;
            }
            
            if(!table_data_col[table_col]) {
                table_data_col[table_col] = '';
            }

            var table_data_option = do_onmark_table_render_sub(
                table_data_in[2],
                table_data_col[table_col]
            );
            table_data_col[table_col] = table_data_option['col'];
            if(table_data_option['colspan'] === "") {
                table_data_option['colspan'] = String(table_data_in[1].length / 2);
            }
            
            if(table_data_real === '') {
                table_data_real += '' + 
                    '<div style="' + table_data_option['div'] + '">' +
                        '<table style="' + table_data_option['table'] + '">' +
                            '<tr style="' + table_data_option['tr'] + '">' +
                '';
            }

            if(table_data_in[2] === '\n') {
                table_data_real += '</tr>';
                table_new_line = 1;
                table_col = 0;
            } else {
                if(table_new_line === 1) {
                    table_data_real += '<tr style="' + table_data_option['tr'] + '">';
                    table_new_line = 0;
                }
                
                table_data_real += '' +
                    '<td     colspan="' + table_data_option['colspan'] + '" ' +
                            'rowspan="' + table_data_option['rowspan'] + '" ' +
                            'style="' + table_data_option['col'] + table_data_option['td'] + '">' + 
                        table_data_option['data'] + 
                    '</td>' +
                '';
                table_col += 1;
            }

            table_data_var = table_data_var.replace(table_in_re, '');
        }
        table_data_real += '</tr></table></div>';
        
        return table_data_real;
    });
    
    data = data.replace(/<t_br>/g, '\n');
    
    return data;
}

function do_onmark_list_render(data) {
    var list_re = /\n((?:(?:(?: )+)\* (?:(?:(?!\n).)+)\n)+)/;
    var list_short_re = /((?: )+)\* ((?:(?!\n).)+)\n/g;
    while(1) {
        var list_data = data.match(list_re);
        if(!list_data) {
            break;
        }
        
        var list_end_data = '<ul>' + list_data[1].replace(list_short_re, function(x, x_1, x_2) {
            return '<li style="margin-left: ' + String(x_1.length * 20) + 'px;">' + x_2 + '</li>';
        }) + '</ul>';

        data = data.replace(list_re, '\n' + list_end_data + '\n');
    }
    
    return data;
}

function do_onmark_math_render(data, data_js, name_include) {
    data = data.replace(/<math>((?:(?!<\/math>).)+)<\/math>/g, '[math($1)]');
    
    var i = 0;
    data = data.replace(/\[math\((((?!\)]).)+)\)]/g, function(x, x_1) {
        i += 1;
        
        data_js += do_math_try_insert(name_include + 'math_' + String(i), do_js_safe_change(do_xss_change(x_1)));
        return '<span id="' + name_include + 'math_' + String(i) + '"></span>';
    });
    
    return [data, data_js];
}

// Main
function do_onmark_render(test_mode = 'test', name_id = '', name_include = '', name_doc = '', doc_data = '') {    
	if(test_mode === 'normal') {
        var data = '\n' + document.getElementById(name_id).innerHTML.replace(/\r/g, '') + '\n';
    } else if(test_mode === 'manual') { 
        var data = '\n' + doc_data.replace(/\r/g, '') + '\n';
    } else {
    	var data = '\n' + (
`||||<tablebgcolor=red><bgcolor=red> test ||
|| A || B ||
||<bgcolor=red> test ||<bgcolor=red> test || asdf
asdf || asdf`
        ).replace(/\r/g, '') + '\n';
    }
    var data_js = '';
    var data_backlink = [];
    var data_nowiki = {};
    
    var data_var = do_onmark_math_render(data, data_js, name_include);
    data = data_var[0];
    data_js = data_var[1];
    
    data_var = do_onmark_nowiki_before_render(data, data_js, name_include, data_nowiki);
    data = data_var[0];
    data_js = data_var[1];
    data_nowiki = data_var[2];
    
    data_var = do_onmark_include_render(data, data_js, name_include, data_nowiki);
    data = data_var[0];
    data_js = data_var[1];
    
    data_var = do_onmark_middle_render(data, data_js, name_include, data_nowiki, name_doc);
    data = data_var[0];
    data_js = data_var[1];
    data_nowiki = data_var[2];
    
    data = do_onmark_text_render(data);
    data = do_onmark_heading_render(data);
    data = do_onmark_table_render(data);
    
    data_var = do_onmark_link_render(data, data_js, name_doc, name_include, data_nowiki);
    data = data_var[0];
    data_js = data_var[1];
    
    data = do_onmark_macro_render(data);
    data = do_onmark_list_render(data);
    data = do_onmark_footnote_render(data, name_include);
    data = do_onmark_last_render(data, name_include);
    
    data_js += '' + 
        'get_link_state("' + name_include + '");\n' + 
        'get_file_state("' + name_include + '");\n' + 
    ''
    data_js = 'render_html("' + name_include + 'render_contect");\n' + data_js
    
    if(test_mode === 'normal') {
        document.getElementById(name_id).innerHTML = data + '<script>' + data_js + '</script>';
        eval(data_js);
    } else if(test_mode === 'manual') {
        return [data, data_js];
    } else {
    	console.log([data, data_js]);
    }
}

do_onmark_render();