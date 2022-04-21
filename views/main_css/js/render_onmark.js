// 중괄호 문법 정리
// Tool
function do_url_change(data) {
    return encodeURIComponent(data);
}

function do_nowiki_change(data, data_nowiki, type = 'normal') {
    return data.replace(/<span id="((?:[^"]*)(?:nowiki_(?:[^"]+)))"><\/span>/g, function(x, x_1) {
        if(type === 'normal') {
            return data_nowiki[x_1];
        } else {
            return '\\' + data_nowiki[x_1];
        }
    });
}

function do_link_change(data, data_nowiki, no_change) {
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

function do_darkmode_split(data) {
    if(
        document.cookie.match(regex_data('main_css_darkmode')) &&
        document.cookie.match(regex_data('main_css_darkmode'))[1] === '1'
    ) {
        let data_split = data.split(',');
        if(data_split.length > 1) {
            return data.split(',')[1];
        } else {
            return data.split(',')[0];
        }
    } else {
        return data.split(',')[0];
    }
}

function do_js_safe_change(data, br_on = 1) {
    data = data.replace(/\\/g, '\\\\');
    data = data.replace(/"/g, '\\"');
    if(br_on === 1) {
        data = data.replace(/\n/g, '<br>');
    } else {
        data = data.replace(/\n/g, '\\n');
    }
    
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

function do_all_try(data) {
    return '' +
        'try {\n'
            + data + 
        '} catch {}\n' +
    ''
}

function do_px_add(data) {
    if(data) {
        return data.match(/^[0-9]+$/) ? (data + 'px') : data;
    } else {
        return '';
    }
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
    data = data.replace(/&quot;/g, '"');
    
    return data;
}

function do_html_escape(data) {
    data = data.replace(/</g, '&lt;');
    data = data.replace(/>/g, '&gt;');
    data = data.replace(/&/g, '&amp;');
    data = data.replace(/"/g, '&quot;');
    
    return data
}

function do_end_br_replace(data) {
    data = data.replace(/(\n| )+$/, '\n');
    
    return data;
}

// Sub
function do_onmark_text_render(data) {    
    data = data.replace(/'''((?:(?!''').)+)'''/g, '<b>$1</b>');
    data = data.replace(/''((?:(?!'').)+)''/g, '<i>$1</i>');
    data = data.replace(/__((?:(?!__).)+)__/g, '<u>$1</u>');
    
    data = data.replace(/\^\^\^((?:(?!\^\^\^).)+)\^\^\^/g, '<sup>$1</sup>');
    data = data.replace(/\^\^((?:(?!\^\^).)+)\^\^/g, '<sup>$1</sup>');
    
    data = data.replace(/,,,((?:(?!,,,).)+),,,/g, '<sub>$1</sub>');
    data = data.replace(/,,((?:(?!,,).)+),,/g, '<sub>$1</sub>');
    
    data = data.replace(/--((?:(?!--).)+)--/g, '<s>$1</s>');
    data = data.replace(/~~((?:(?!~~).)+)~~/g, '<s>$1</s>');
    
    return data;
}

function do_onmark_set_toc_name(toc_name) {
    let toc_data = document.getElementById("heading_text_" + toc_name).innerText;
    
    for(let for_a = 0; document.getElementsByClassName("toc_text_" + toc_name)[for_a]; for_a++) {
        document.getElementsByClassName("toc_text_" + toc_name)[for_a].innerHTML = toc_data;
    }
}

function do_onmark_heading_render(
    data, 
    data_js, 
    name_doc, 
    name_include
) {
    var heading_re = /\n(={1,6})(#)? ?([^\n]+) ?#?={1,6}\n/;
    var heading_level_all = [0, 0, 0, 0, 0, 0];
    var toc_data = '';
    var toc_n = 0;
    while(1) {        
        toc_n += 1;
        
        var heading_data = data.match(heading_re);
        if(!heading_data) {
            break;
        }
        
        if(toc_data === '') {
            toc_data += '<div id="toc"><div id="toc_title">TOC</div>\n';
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
        
        var heading_data_text = heading_data[3].replace(/=+$/, '');
        heading_data_text = heading_data_text.replace(/#$/, '');
        heading_data_text = heading_data_text.replace(/ $/, '');
        
        toc_data += '' +
            '<span style="margin-left: ' + String((heading_level_string.match(/\./g).length - 1) * 10) + 'px;">' +
                '<a href="#s-' + heading_level_string_no_end + '">' + 
                    heading_level_string + 
                '</a> ' + 
                '<span class="toc_text_' + heading_level_string_no_end + '"></span>' +
            '</span>' +
            '\n' +
        ''
        data_js += 'do_onmark_set_toc_name("' + heading_level_string_no_end + '");\n';
        
        data = data.replace(heading_re, 
            '\n' +
            (toc_n === 1 ? '' : '</div>') +
            '<h' + heading_level + ' class="render_heading_text">' + 
                '<a href="#toc" id="s-' + heading_level_string_no_end + '">' + heading_level_string + '</a> ' + 
                '<span id="heading_text_' + heading_level_string_no_end + '">' +
                    heading_data_text + 
                '</span> ' +
                '<a id="edit_load_' + String(toc_n) + '" ' +
                    'style="font-size: 70%;"' +
                    'href="/edit/' + do_url_change(name_doc) + '/doc_section/' + String(toc_n) + '">✎</a> ' +
                '<a href="javascript:void(0);" ' +
                    'onclick="javascript:do_open_folding(\'' + name_include + 'in_data_' + String(toc_n) + '\', this);"' +
                    'style="font-size: 70%;">' + (heading_data[2] ? '⊕' : '⊖') + '</a>' +
            '</h' + heading_level + '>' +
            '<div   id="' + name_include + 'in_data_' + String(toc_n) + '" ' +
                    'style="display: ' + (heading_data[2] ? 'none' : 'block') + ';">' +
            '<end_point>\n'
       );
    }
    
    if(toc_data !== '') {
        toc_data += '</div>';

        data = do_end_br_replace(data) + '</div>';
    }
    
    var toc_auto_add = data.match(/\[(?:목차|toc)\(no\)\]/);
    var toc_re = /\[(?:toc|목차)\]/g;
    if(toc_auto_add) {
        data = data.replace(/\[(?:목차|toc)\(no\)\]/g, '');
    } else {
        if(name_include === '' && !data.match(toc_re)) {
            data = data.replace(/(<h[1-6] (?:[^>]+)>)/, '<div id="auto_toc">' + toc_data + '</div>$1');
        }
    }
    
    data = data.replace(toc_re, toc_data);
    
    return [data, data_js];
}

function do_onmark_link_render(data, data_js, name_doc, name_include, data_nowiki, data_wiki_set) {
    var num_link = 0;
    
    var category_data = '';
    
    var category_re = /^(분류|category):/i;
    let inter_re = /^inter:([^:]+):/i;
    let out_link_re = /^http(s)?:\/\//i;
    var file_re = /^(파일|file|외부|out):/i;
    
    var link_re = /\[\[(((?!\[\[|\]\]).)+)\]\]/;
    
    while(data.match(link_re)) {
        data = data.replace(link_re, function(x, x_1) {
            var link_split = x_1.split('|');
            var link_real = link_split[0];
            var link_out = link_split[1] ? link_split[1] : link_split[0];
            var link_out_2 = link_split[1] ? link_split[1] : '';

            num_link += 1;
            var num_link_str = String(num_link - 1);
            if(link_real.match(/<|>/)) {
                return '<link_s>' + x_1 + '<link_e>';
            } else if(link_real.match(file_re)) {
                var file_load_type = link_real.match(file_re)[1];
                var file_name = link_real.replace(file_re, '');

                if(file_load_type === '파일' || file_load_type === 'file') {
                    var file_type = file_name.split('.');
                    file_name = file_type.slice(0, file_type.length - 1).join('.');
                    file_type = file_type[file_type.length - 1];

                    var file_src = do_url_change(do_xss_change(file_name)) + '.' + do_html_escape(file_type);
                    var file_alt = do_html_escape(file_name + '.' + file_type);
                    var file_exist = 1;
                } else {
                    var file_src = do_html_escape(file_name);
                    var file_alt = do_html_escape(file_name);
                    var file_exist = 0;
                }

                var file_style = '';
                var file_bgcolor = '';
                var file_align = '';

                var file_set = link_out_2.split('&amp;');
                for(let i = 0; file_set[i]; i++) {
                    var file_set_name = file_set[i].split('=');
                    var file_set_data = file_set_name[1];
                    file_set_name = file_set_name[0];
                    
                    if(file_set_data) {
                        if(file_set_name === 'width') {
                            file_style += 'width:' + do_px_add(file_set_data) + ';';
                        } else if(file_set_name === 'height') {
                            file_style += 'height:' + do_px_add(file_set_data) + ';';
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
                    }
                }

                return '' +
                    '<span style="' + file_align + '">' + 
                        '<span  style="' + file_bgcolor + '" ' +
                                'class="' + name_include + 'file_finder" ' +
                                'under_style="' + file_style + '" ' +
                                'under_alt="' + file_alt + '" ' +
                                'under_src="' + file_src + '" ' +
                                'under_href="' + (file_exist === 0 ? "out_link" : '/upload?name=' + file_src.replace(/\.[^.]+$/, '')) + '">' +
                        '</span>' + 
                    '</span>' +
                ''
            } else if(link_real.match(category_re)) {
                var category_link = link_real.replace(category_re, '');

                data_js += '' +
                    'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].href = ' + 
                        '"/w/category:' + do_url_change(category_link) + '";' +
                    '\n' +
                '';
                data_js += '' +
                    'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].title = ' + 
                        '"' + do_js_safe_change(do_xss_change('category:' + category_link)) + '";' +
                    '\n' +
                '';

                category_data += '' +
                    '<a class="' + name_include + 'link_finder" ' +
                        'name="' + name_include + 'set_link_' + num_link_str + '" ' +
                        'href="" ' +
                        'title="">' +
                        category_link +
                    '</a> | ' +
                ''

                return '';
            } else if(link_real.match(out_link_re)) {
                var i = 0;
                while(i < 2) {
                    if(i === 0) {
                        var var_link_type = 'href';
                    } else {
                        var var_link_type = 'title';
                    }

                    data_js += '' +
                        'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].' + var_link_type + ' = ' + 
                            '"' + do_js_safe_change(do_xss_change(link_real)) + '";' +
                        '\n' +
                    '';

                    i += 1;
                }

                return  '<a id="out_link" ' +
                            'class="' + name_include + 'link_finder" ' +
                            'target="_blank" ' +
                            'name="' + name_include + 'set_link_' + num_link_str + '" ' + 
                            'title=""' +
                            'href="">' + link_out + '</a>';
            } else if(link_real.match(inter_re)) {
                let data_inter = link_real.match(inter_re);
                
                let data_inter_link = '';
                let data_inter_logo = '';
                if(data_inter) {
                    if(link_real === link_out) {
                        link_real = link_real.replace(
                            inter_re, 
                            ''
                        );
                        
                        link_out = link_real;
                    } else {
                        link_real = link_real.replace(
                            inter_re, 
                            ''
                        );
                    }
                    
                    var data_inter_var = do_link_change(link_real, data_nowiki, 1);
                    var data_inter_link_main = data_inter_var[0];
                    var data_inter_link_sub = data_inter_var[1];
                        
                    let data_inter_get = data_wiki_set['inter_wiki'][data_inter[1]];
                    if(data_inter_get) {
                        data_inter_link = data_inter_get['link'];
                        if(data_inter_get['logo'] !== '') {
                            data_inter_logo = data_inter_get['logo'];
                            data_inter_logo = data_inter_logo.replace(/&lt;/g, '<').replace(/&gt;/g, '>');
                        } else {
                            data_inter_logo = data_inter[1] + ':';
                        }
                    } else {
                        return '';
                    }
                    
                    data_js += '' +
                        'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].title = ' + 
                        '"' + do_js_safe_change(do_xss_change(data_inter[1] + ':' + link_real)) + '";' +
                            '\n' +
                    '';
                    data_js += '' +
                        'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].href = ' + 
                        '"' + data_inter_link + do_url_change(data_inter_link_main) + data_inter_link_sub + '";' +
                            '\n' +
                    '';

                    return  '<a id="inside" ' +
                                'class="' + name_include + 'link_finder" ' +
                                'target="_blank" ' +
                                'name="' + name_include + 'set_link_' + num_link_str + '" ' + 
                                'title=""' +
                                'href="">' + data_inter_logo + link_out + '</a>'; 
                } else {
                    return '';
                }
            } else {
                if(link_real.match(/^\//)) {
                    link_real = name_doc + link_real;
                } else if(link_real.match(/^\.\.\//)) {
                    link_real = link_real.replace(/^\.\.\//, '');
                    link_real = name_doc.replace(/\/[^/]+$/, '') + (link_real !== '' ? '/' + link_real : '');
                }
                
                var link_data_var = do_link_change(link_real, data_nowiki, 0);
                var link_main = link_data_var[0];
                var link_sub = link_data_var[1];
                
                let link_id = "real_normal_link"

                var i = 0;
                while(i < 2) {
                    if(i === 0) {
                        var var_link_type = 'href';
                        if(link_main === '') {
                            link_id = "in_doc_link"
                            var var_link_data = link_sub;
                        } else {
                            var var_link_data = '/w/' + do_url_change(link_main) + link_sub;
                        }
                    } else {
                        var var_link_type = 'title';
                        var var_link_data = do_js_safe_change(link_main);
                    }

                    data_js += '' +
                        'document.getElementsByName("' + name_include + 'set_link_' + num_link_str + '")[0].' + var_link_type + ' = ' + 
                            '"' + var_link_data + '";' +
                        '\n' +
                    '';

                    i += 1;
                }

                return  '<a class="' + name_include + 'link_finder" ' +
                            'id="' + link_id + '"' +
                            'name="' + name_include + 'set_link_' + num_link_str + '" ' +
                            'title="" ' +
                            'href="">' + link_out + '</a>';
            }
        });
    }
    
    data = data.replace(/<link_s>/, '[[');
    data = data.replace(/<link_e>/, ']]');

    if(category_data !== '') {
        if(name_include === '') {
            category_data = '<div id="cate_all"><div id="cate">Category : ' + category_data;
        } else {
            category_data = '<div style="display: none;" id="cate_all"><div id="cate">Category : ' + category_data;
        }
        
        category_data = category_data.replace(/\| $/, '') + '</div></div>';
    }
    
    return [data, data_js, category_data];
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
                    '<a href="javascript:do_open_foot(\'' + name_include + '\', \'fn-' + String(i) + '\', 1);" ' +
                        'id="' + name_include + 'cfn-' + String(i) + '">' +
                        '(' + footnote_name + ')' +
                    '</a> <span id="' + name_include + 'fn-' + String(i) + '">' + footnote_line_data + '</span>' +
                '</li>' +
            '';
            data = data.replace(footnote_re, '' +
                '<sup>' +
                    '<a href="javascript:do_open_foot(\'' + name_include + '\', \'fn-' + String(i) + '\', 0);" ' +
                        'id="' + name_include + 'rfn-' + String(i) + '">' +
                        '(' + footnote_name + ')' +
                    '</a>' +
                '</sup><span id="' + name_include + 'dfn-' + String(i) + '"></span>' +
           '');
            
            i += 1;
        } else {
            if(footnote_end_data !== '') {
                footnote_end_data = '<ul id="footnote_data">' + footnote_end_data + '</ul>';   
            }
            
            data = data.replace(footnote_re, footnote_end_data);    
            footnote_end_data = '';
        }
    }
    
    if(footnote_end_data !== '') {
        data = do_end_br_replace(data) + '<ul id="footnote_data">' + footnote_end_data + '</ul>';
    }
    
    return data;
}

function do_onmark_macro_render(data, data_js) {
    data = data.replace(/\[([^[\](]+)\(((?:(?!\)\]).)+)\)\]/g, function(x, x_1, x_2) {
        x_1 = x_1.toLowerCase();
        if(x_1 === 'youtube' || x_1 === 'kakaotv' || x_1 === 'nicovideo' || x_1 === 'navertv' || x_1 === 'vimeo') {
            var video_code = x_2.match(/^([^,]+)/);
            video_code = video_code ? video_code[1] : '';
            
            var video_width = x_2.match(/,(?: *)width=([0-9]+)/);
            video_width = video_width ? (video_width[1] + 'px') : '640px';
            
            var video_height = x_2.match(/,(?: *)height=([0-9]+)/);
            video_height = video_height ? (video_height[1] + 'px') : '360px';
            
            if(x_1 === 'youtube') {
                var video_start = x_2.match(/,(?: *)start=([0-9]+)/);
                video_start = video_start ? ('?start=' + video_start[1]) : '';
                
                video_code = video_code.replace(/^https:\/\/www\.youtube\.com\/watch\?v=/, '');
                video_code = video_code.replace(/^https:\/\/youtu\.be\//, '');
                
                var video_src = 'https://www.youtube.com/embed/' + video_code + video_start
            } else if(x_1 === 'kakaotv') {
                video_code = video_code.replace(/^https:\/\/tv\.kakao\.com\/channel\/9262\/cliplink\//, '');
                video_code = video_code.replace(/^http:\/\/tv\.kakao\.com\/v\//, '');
                
                var video_src = 'https://tv.kakao.com/embed/player/cliplink/' + video_code +'?service=kakao_tv'
            } else if(x_1 === 'nicovideo') {
                var video_src = 'https://embed.nicovideo.jp/watch/' + video_code
            } else if(x_1 === 'navertv') {
                var video_src = 'https://tv.naver.com/embed/' + video_code
            } else {
		var video_src = 'https://player.vimeo.com/video/' + video_code
	    }
            
            return '<iframe style="width: ' + video_width + '; height: ' + video_height + ';" src="' + video_src + '" frameborder="0" allowfullscreen></iframe>';
        } else if(x_1 === 'anchor') {
            return '<span id="' + x_2 + '"></span>';
        } else if(x_1 === 'ruby') {
            let ruby_main_data = x_2.match(/^([^,]+)/);
            if(ruby_main_data) {
                ruby_main_data = ruby_main_data[1];
            } else {
                ruby_main_data = 'Test';
            }
            
            let ruby_sub_data = x_2.match(/,(?: *)ruby=([^,]+)/);
            if(ruby_sub_data) {
                ruby_sub_data = ruby_sub_data[1];
            } else {
                ruby_sub_data = 'Test';
            }
            
            return '<ruby>' + ruby_main_data + '<rp>(</rp><rt>' + ruby_sub_data + '</rt><rp>)</rp></ruby>';
        } else if(x_1 === 'dday') {
            var date_old = new Date(x_2);
            var date_now = new Date(do_return_date());
            
            var date_end = Math.floor((date_now - date_old) / (24 * 60 * 60 * 1000));
            
            return (date_end > 0 ? '+' : '') + date_end;
        } else if(x_1 === 'age') {
            var date_old = new Date(x_2);
            var date_now = new Date(do_return_date());
            
            var date_end = Math.floor((date_now - date_old) / (365 * 24 * 60 * 60 * 1000));
            
            return date_end > 0 ? date_end : '';            
        } else if(x_1 === 'pagecount') {
            return '0';
        } else {
            return '<macro_start>' + x_1 + '(' + x_2 + ')<macro_end>';
        }
    });
    
    var pagecount_n = 0;
    data = data.replace(/\[([^[*()\]]+)\]/g, function(x, x_1) {
        x_1 = x_1.toLowerCase();
        if(x_1 === 'date' || x_1 === 'datetime') {
            return do_return_date();
        } else if(x_1 === 'clearfix') {
            return '<div style="clear:both"></div>';
        } else if(x_1 === 'br') { 
            return '<br>';
        } else if(x_1 === 'pagecount') {
            if(pagecount_n === 0) {
                pagecount_n += 1;
                
                data_js += 'page_count();\n';
            }
            
            return '<span class="all_page_count"></span>';
        } else {
            return '<macro_start>' + x_1 + '<macro_end>';
        }
    });
    
    data = data.replace(/<macro_start>/g, '[');
    data = data.replace(/<macro_end>/g, ']');
    
    return [data, data_js];
}

function do_onmark_middle_render(data, data_js, name_include, data_nowiki, name_doc) {
    // 이것도 nowiki 처럼 가야할 듯
    
    var middle_re = /{{{((?:(?!{{{|}}}).)+)}}}/s;
    
    var html_n = 0;
    var syntax_n = 0;
    var nowiki_n = 0;
    var folding_n = 0;
    
    while(data.match(middle_re)) {        
        data = data.replace(middle_re, function(x, x_1) {
            var middle_data_before = x_1.match(/^({+)/);
            middle_data_before = middle_data_before ? middle_data_before[1] : '';
            var middle_data_x_1 = x_1.replace(/^({+)/, '');
            
            var middle_data = middle_data_x_1.match(/^([^ ]+) /);
            middle_data = middle_data ? middle_data[1] : '';
            if(middle_data) {
                var middle_data_all = middle_data_x_1.replace(/^([^ ]+) /, '');    
            }

            var middle_type = middle_data.match(
                /^(?:(?:(#|@)((?:[0-9a-f-A-F]{3}){1,2}|\w+))(?:,(?:(#|@)((?:[0-9a-f-A-F]{3}){1,2}|\w+)))?|(\+|-)([1-5])|#!(html|wiki|syntax|folding|html))$/i
            );
            if(middle_type) {
                if(middle_data_x_1[middle_data_x_1.length - 1] === '\\') {
                    return middle_data_before + '{{{' + middle_data_x_1 + '<mid_e>';
                } else if(middle_type[1]) {
                    let data_color = middle_type[2];
                    if(middle_type[3]) {
                        data_color = do_darkmode_split(middle_type[2] + ',' + middle_type[4])    
                    }
                    
                    let data_sharp = '';
                    if(data_color.match(/^(?:[0-9a-f-A-F]{3}){1,2}$/)) {
                        data_sharp = '#';
                    }
                    
                    if(middle_type[1] === '@') {
                        return middle_data_before + '<span style="background: ' + data_sharp + data_color + '">' + middle_data_all + '</span>';
                    } else {
                        return middle_data_before + '<span style="color: ' + data_sharp + data_color + '">' + middle_data_all + '</span>';
                    }
                } else if(middle_type[5]) {
                    if(middle_type[5] === '+') {
                        return middle_data_before + '<span style="font-size: ' + String(100 + (Number(middle_type[6]) * 20)) + '%">' + middle_data_all + '</span>';
                    } else {
                        return middle_data_before + '<span style="font-size: ' + String(100 - (Number(middle_type[6]) * 10)) + '%">' + middle_data_all + '</span>';
                    }
                } else if(middle_type[7]) {
                    var middle_type_sub = middle_type[7].toLowerCase();
                    if(middle_type_sub === 'html') {
                        html_n += 1;

                        data_nowiki[name_include + 'nowiki_html_' + String(html_n)] = middle_data_all;
                        data_js += do_data_try_insert(
                            name_include + 'nowiki_html_' + String(html_n),
                            do_js_safe_change(middle_data_all)
                        );

                        return middle_data_before + '<span id="' + name_include + 'nowiki_html_' + String(html_n) + '"></span>';
                    } else if(middle_type_sub === 'wiki') {
                        var middle_wiki_re = /^(?:[^ ]+)(?: style=['"]([^\n'"]*)['"])?[^\n]*\n?/;
                        var middle_wiki = middle_data_x_1.match(middle_wiki_re);
                        middle_wiki = middle_wiki[1] ? middle_wiki[1] : '';
                        middle_wiki = middle_wiki.replace(/position/, '');

                        middle_data_all = middle_data_x_1.replace(middle_wiki_re, '');
                        
                        return middle_data_before + '' +
                            '<wiki_s style="' + middle_wiki + '">' +
                                '<end_point>\n' + 
                                middle_data_all + 
                                '\n<start_point>' +
                            '<wiki_e>' +
                        '';
                    } else if(middle_type_sub === 'folding') {
                        folding_n += 1;
                        
                        var middle_folding_re = /^(?:[^ ]+)(?: ([^\n]*))?\n?/;
                        var middle_folding = middle_data_x_1.match(middle_folding_re);
                        middle_folding = middle_folding ? middle_folding[1] : 'open';
                        
                        middle_data_all = middle_data_x_1.replace(middle_folding_re, '');

                        data_js += do_data_try_insert('get_' + name_include + 'folding_' + String(folding_n), do_js_safe_change(middle_folding));
                        return middle_data_before + 
                            '<div>' +
                                '<b>' + 
                                    '<a href="javascript:do_open_folding(\'' + name_include + 'folding_' + String(folding_n) + '\');" ' + 
                                        'id="get_' + name_include + 'folding_' + String(folding_n) + '">' + 
                                    '</a>' + 
                                '</b>' + 
                                '<div id="' + name_include + 'folding_' + String(folding_n) + '" style="display: none;">' +
                                    '<wiki_s style="">' +
                                        '\n' +
                                        middle_data_all +
                                        '\n<start_point>' +
                                    '<wiki_e>' +
                                '</div>' +
                            '</div>' +
                        ''
                    } else if(middle_type_sub === 'syntax') {
                        syntax_n += 1;

                        let middle_syntax_re = /^(?:[^ ]+) ([^\n]+)\n?/;
                        var middle_syntax = middle_data_x_1.match(middle_syntax_re);
                        middle_syntax = middle_syntax ? middle_syntax[1] : 'python';

                        middle_data_all = middle_data_x_1.replace(middle_syntax_re, '');

                        data_nowiki[name_include + 'nowiki_syntax_' + String(syntax_n)] = middle_data_all;
                        data_js += do_data_try_insert(
                            name_include + 'nowiki_syntax_' + String(syntax_n),
                            do_js_safe_change(do_xss_change(middle_data_all), 0)
                        );

                        return middle_data_before +
                            '<pre id="syntax">' +
                                '<code  id="' + name_include + 'nowiki_syntax_' + String(syntax_n) + '" ' +
                                        'class="' + middle_syntax + '"></code>' +
                            '</pre>' +
                        ''
                    }
                }
            }
            
            // 최대한 노력해봐야함
            nowiki_n += 1;
            
            data_nowiki[name_include + 'nowiki_' + String(nowiki_n)] = middle_data_x_1;
            data_js = do_data_try_insert(
                name_include + 'nowiki_' + String(nowiki_n),
                do_js_safe_change(middle_data_x_1)
            ) + data_js;
            
            return middle_data_before + '<span id="' + name_include + 'nowiki_' + String(nowiki_n) + '"></span>';
        });
    }
    
    data = data.replace(/<mid_s>/g, '{{{');
    data = data.replace(/<mid_e>/g, '}}}');
    
    if(syntax_n > 0) {
        data_js += 'hljs.highlightAll();\n';
    }
    
    return [data, data_js, data_nowiki];
}

function do_onmark_last_render(data, name_include, data_category) {       
    // middle_render 마지막 처리
    data = data.replace(/<wiki_s_[0-9] /g, '<div ');
    data = data.replace(/<wiki_e_[0-9]>/g, '</div>');
    
    // heading_render 마지막 처리
    data = data.replace(/\n?<start_point>/g, '');
    data = data.replace(/<end_point>\n?/g, '');
    
    // br 마지막 처리
    data = data.replace(/^(\n| )+/, '');
    data = do_end_br_replace(data);
    data = data.replace(/\n/g, '<br>');
    
    data += data_category;
    
    return data;
}

function do_onmark_include_render(data, data_js, name_include, data_nowiki) {
    var include_re = /\[include\(((?:(?!\)\]).)+)\)\]/i;
    
    if(name_include === '') {
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
    } else {
        while(1) {
            var include_data = data.match(include_re);
            if(!include_data) {
                break;
            }
            
            data = data.replace(include_re, '');
        }
    }
    
    return [data, data_js];
}

function do_onmark_nowiki_before_render(data, data_js, name_include, data_nowiki) {   
    var num_nowiki = 0;
    data = data.replace(/\\(&gt;|&lt;|.)/g, function(x, x_1) {
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
    var align_auto = 1;
    
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
            var table_option_data = data_option_var[1].replace(/"/g, '');
            if(table_option_name === 'tablebgcolor') {
                // table
               data_option_all['table'] += 'background:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'tablewidth') {
                data_option_all['table'] += 'width:' + do_px_add(table_option_data) + ';';
            } else if(table_option_name === 'tableheight') {
                data_option_all['table'] += 'height:' + do_px_add(table_option_data) + ';';
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
                data_option_all['table'] += 'color:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'tablebordercolor') {
                data_option_all['table'] += 'border:2px solid ' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'rowbgcolor') {
                // tr
                data_option_all['tr'] += 'background:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'rowtextalign') {
                data_option_all['tr'] += 'text-align:' + table_option_data + ';';
            } else if(table_option_name === 'rowcolor') {
                data_option_all['tr'] += 'color:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'colcolor') {
                // col
                data_option_all['col'] += 'color:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'colbgcolor') {
                data_option_all['col'] += 'background:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'bgcolor') {
                // td
                data_option_all['td'] += 'background:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'color') {
                data_option_all['td'] += 'color:' + do_darkmode_split(table_option_data) + ';';
            } else if(table_option_name === 'width') {
                data_option_all['td'] += 'width:' + do_px_add(table_option_data) + ';';
            } else if(table_option_name === 'height') {
                data_option_all['td'] += 'height:' + do_px_add(table_option_data) + ';';
            } else {
                no_option = '<lt>' + data_option + '<gt>';
            }
        } else {
            if(data_option.match(/^-[0-9]+$/)) {
                // span
                data_option_all['colspan'] = data_option.replace('-', '');
            } else if(data_option.match(/^(\^|v)?\|[0-9]+$/)) {
                if(data_option[0] === '^') {
                    data_option_all['td'] += 'vertical-align: top;';
                } else if(data_option[0] === 'v') {
                    data_option_all['td'] += 'vertical-align: bottom;';
                }
                
                data_option_all['rowspan'] = data_option.replace(/[^0-9]+/g, '');
            } else if(
                data_option === '(' ||
                data_option === ':' ||
                data_option === ')'
            ) {
                align_auto = 0;
                
                // align
                if(data_option === '(') {
                    data_option_all['td'] += 'text-align: left;';
                } else if(data_option === ':') {
                    data_option_all['td'] += 'text-align: center;';
                } else {
                    data_option_all['td'] += 'text-align: right;';
                }
                
                
            } else {
                var table_option_data = data_option.replace(/"/g, '')
                table_option_data = table_option_data.match(/^((?:(?:#(?:[a-zA-Z0-9]{3}){1,2})|\w+)(?:,(?:(?:#(?:[a-zA-Z0-9]{3}){1,2})|\w+))?)/);
                if(table_option_data) {
                    data_option_all['td'] += 'background:' + do_darkmode_split(table_option_data[1]) + ';';
                } else {
                    no_option = '<lt>' + data_option + '<gt>';
                }
            }
        }
        
        data = data.replace(table_option_re, no_option);
    }
    
    data = data.replace('<lt>', '&lt;');
    data = data.replace('<gt>', '&gt;');
    data_option_all['data'] = data;
    
    if(align_auto === 1) {
        if(
            data_option_all['data'][0] === ' ' &&
            data_option_all['data'][data_option_all['data'].length - 1] === ' '
        ) {
            data_option_all['td'] += 'text-align:center;';
        } else if(data_option_all['data'][0] === ' ') {
            data_option_all['td'] += 'text-align:right;';
        }
    }
    
    data_option_all['data'] = data_option_all['data'].replace(/^ +| +$/g, '');
    
    return data_option_all;
}

function do_onmark_table_render_main(data) {
    var table_re = /\n((?:(?:(?:(?:\|\|)+)|(?:\|[^|]+\|(?:\|\|)*))\n?(?:(?:(?!\|\|).)+))(?:(?:\|\||\|\|\n|(?:\|\|)+(?!\n)(?:(?:(?!\|\|).)+)\n*)*)\|\|)\n/gs;
    data = data.replace(table_re, function(x, x_1) {
        var table_cel_re = /((?:\|\|)+)((?:(?!\|\|).)*)/gs;
        var table_data = '';
        var table_data_org = x_1;
        
        var table_col = 0;
        var table_col_data = {};
        var table_col_count = {};
        
        let table_caption_re = /^\|([^|]+)\|/;
        let table_caption = '';
        let table_caption_get = table_data_org.match(table_caption_re);
        if(table_caption_get) {
            table_caption = '<caption>' + table_caption_get[1] + '</caption>';
            table_data_org = table_data_org.replace(table_caption_re, '||');
        }
            
        table_data_org = table_data_org.replace(table_cel_re, function(x, x_1, x_2) {
            if(!table_col_data[table_col]) {
                table_col_data[table_col] = '';
            }
            
            if(!table_col_count[table_col]) {
                table_col_count[table_col] = 0;
            }
            
            if(table_col_count[table_col] !== 0) {
                table_col_count[table_col] -= 1;
                table_col += 1;
            }

            var table_data_option = do_onmark_table_render_sub(x_2, table_col_data[table_col]);
            table_col_data[table_col] = table_data_option['col'];
            if(table_data_option['colspan'] === '') {
                table_data_option['colspan'] = String(x_1.length / 2);
            }

            if(table_data === '') {
                table_data += '' + 
                    '<div class="table_safe" style="' + table_data_option['div'] + '">' +
                        '<table style="' + table_data_option['table'] + '">' +
                            table_caption +
                '';
            }

            if(x_1 === '||' && (x_2 === '\n' || x_2 === '')) {
                table_data += '</tr>';
                table_col = 0;
            } else if(x_2 === '\n' || x_2 === '') {
                table_data += '</tr><tr></tr>';
                table_col = 0;
            } else {
                if(table_col === 0) {
                    table_data += '' + 
                        '<tr style="' + table_data_option['tr'] + '">' +
                    ''
                }

                if(table_data_option['rowspan'] !== '') {
                    table_col_count[table_col] += Number(table_data_option['rowspan']) - 1;
                }
                
                table_data += '' +
                    '<td     colspan="' + table_data_option['colspan'] + '" ' +
                            'rowspan="' + table_data_option['rowspan'] + '" ' +
                            'style="' + table_data_option['col'] + table_data_option['td'] + '">' + 
                        table_data_option['data'] + 
                    '</td>' +
                '';
                table_col += 1;
            }

            return '';
        });
        
        if(table_col === 0) {
            table_data += '</table></div>';
        } else {
            table_data += '</tr></table></div>';
        }
        
        return '\n' + table_data + '\n';
    });
    
    return data;
}

function do_onmark_table_render(data) {
    data = data.replace(/\n +\|\|/g, '\n||');
    
    var wiki_re = /<wiki_s ([^>]+)>((?:(?!<wiki_s |<wiki_e>).)+)<wiki_e>/s;
    while(1) {
        if(!data.match(wiki_re)) {
            break;
        }
        
        data = data.replace(wiki_re, function(x, x_1, x_2) {
            return '<wiki_s_2 ' + x_1 + '>' + do_onmark_table_render_main(x_2) + '<wiki_e_2>';
        });
    }
    
    data = do_onmark_table_render_main(data);
    
    return data;
}

function do_onmark_list_sub_render(data) {
    var quote_re = /\n((?:(?:(?:&gt;)+) ?(?:(?:(?!\n).)*)\n)+)/;
    var quote_short_re = /((?:&gt;)+) ?((?:(?!\n).)*)\n/g;
    var quote_leng = 1;
    while(1) {
        var quote_data = data.match(quote_re);
        if(!quote_data) {
            break;
        }
        
        var quote_end_data = quote_data[1].replace(quote_short_re, function(x, x_1, x_2) {
            var quote_leng_now = (x_1.length / 4);
            var quote_data_part_1 = '';
            var quote_data_part_2 = '';
            if(quote_leng < quote_leng_now) {
                var quote_data_part_1 = ('<blockquote><end_point>\n'.repeat(quote_leng_now - quote_leng));
            } else if(quote_leng > quote_leng_now) {
                var quote_data_part_1 = ('\n<start_point></blockquote>'.repeat(quote_leng - quote_leng_now));
            }
            
            quote_leng = quote_leng_now;
            
            return '' +
                quote_data_part_1 + 
                x_2 + '\n'
                quote_data_part_2 +
            '';
        });

        data = data.replace(quote_re, '' +
            '\n' +
            '<blockquote>' + 
                '<end_point>\n' +
                quote_end_data + 
                '\n<start_point>' +
            '</blockquote>' +
            '<end_point>\n' +
       '');
    }
    
    var list_re = /\n((?:(?:(?: )*)\* ?(?:(?:(?!\n).)+)\n)+)/;
    var list_short_re = /((?: )*)\* ?((?:(?!\n).)+)\n/g;
    while(1) {
        var list_data = data.match(list_re);
        if(!list_data) {
            break;
        }
        
        var list_end_data = list_data[1].replace(list_short_re, function(x, x_1, x_2) {
            var list_leng = x_1.length;
            
            list_leng = list_leng > 0 ? list_leng : 1;
            
            return '<li style="margin-left: ' + String(list_leng * 20) + 'px">' + x_2 + '</li>';
        });

        data = data.replace(list_re, '\n<ul>' + list_end_data + '</ul><end_point>\n');
    }
    
    return data;
}

function do_onmark_list_render(data) {
    var wiki_re = /<wiki_s_2 ([^>]+)>((?:(?!<wiki_s_2 |<wiki_e_2>).)+)<wiki_e_2>/s;
    while(1) {
        if(!data.match(wiki_re)) {
            break;
        }

        data = data.replace(wiki_re, function(x, x_1, x_2) {
            return '<wiki_s_3 ' + x_1 + '>' + x_2.replace(/\n/g, '<t_br>') + '<wiki_e_3>';
        });
    }
    
    data = do_onmark_list_sub_render(data);
    
    var wiki_re = /<wiki_s_3 ([^>]+)>((?:(?!<wiki_s_3 |<wiki_e_3>).)+)<wiki_e_3>/s;
    while(1) {
        if(!data.match(wiki_re)) {
            break;
        }

        data = data.replace(wiki_re, function(x, x_1, x_2) {
            return '<wiki_s_4 ' + x_1 + '>' + x_2.replace(/<t_br>/g, '\n') + '<wiki_e_4>';
        });
    }
    
    data = do_onmark_list_sub_render(data);
    
    return data;
}

function do_onmark_math_render(data, data_js, name_include, data_nowiki) {
    data = data.replace(/&lt;math&gt;((?:(?!&lt;\/math&gt;).)+)&lt;\/math&gt;/g, '[math($1)]');
    
    var i = 0;
    data = data.replace(/\[math\((((?!\)]).)+)\)]/g, function(x, x_1) {
        i += 1;
        
        data_js += do_math_try_insert(
            name_include + 'math_' + String(i), 
            do_js_safe_change(do_xss_change(do_nowiki_change(x_1, data_nowiki, 'math')))
        );
        
        return '<span id="' + name_include + 'math_' + String(i) + '"></span>';
    });
    
    return [data, data_js];
}

function do_onmark_hr_render(data) {
    var hr_re = /\n-{4,9}\n/;
    while(1) {
        if(!data.match(hr_re)) {
            break;
        }
        
        data = data.replace(hr_re, '\n<hr><end_point>\n');
    }
    
    return data;
}

function do_onmark_redirect_render(data, data_js, name_doc) {
    var redirect_re = /^\n#(?:redirect|넘겨주기) ([^\n]+)/i;
    var data_redirect = data.match(redirect_re);
    if(data_redirect) {
        var link_data_var = do_link_change(data_redirect[1], {}, 1);
        var link_main = link_data_var[0];
        var link_sub = link_data_var[1];
        
        // 임시 조치
        if(
            name_include == '' &&
            window.location.search === '' &&
            window.location.pathname.match(/\/w\//) &&
            !window.location.pathname.match(/\/doc_from\//)
        ) {
            window.location.href = '/w/' + do_url_change(link_main) + '/doc_from/' + do_url_change(name_doc) + link_sub;
        }
        
        return [
            '/w/' + do_url_change(link_main) + 
            '/doc_from/' + do_url_change(name_doc) + 
            link_sub,
            data_js, 
            1
        ];
    } else {
        return [data, data_js, 0];
    }
}

function do_onmark_remark_render(data) {
    data = data.replace(/\n##([^\n]+)/g, '');
    
    return data;
}

// Main
// var 쓰인 부분 전부 let으로 변경하기 (호이스팅 혼용 방지)
// 중첩 함수 구조로 개편하기
function do_onmark_render(
    test_mode = 'test', 
    name_id = '', 
    name_include = '', 
    name_doc = '', 
    doc_data = ''
) {
    let data_wiki_set = {};
	if(test_mode === 'normal') {
        var data = '\n' + 
            document.getElementById(name_id + '_load').innerHTML.replace(/\r/g, '') + 
        '\n';
    } else if(test_mode === 'manual') { 
        var data = '\n' + 
            doc_data.replace(/\r/g, '') + 
        '\n';
    } else {
    	var data = '\n' + (
``
        ).replace(/\r/g, '') + '\n';
    }
    var data_js = '';
    var data_backlink = [];
    var data_nowiki = {};
        
    name_doc = do_xss_change(name_doc);
    console.log(name_doc);

    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/api/setting/inter_wiki");
    xhr.send();

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            data_wiki_set = JSON.parse(this.responseText);
            let data_wiki_set_inter_wiki = { "inter_wiki" : {}};
            if(data_wiki_set["inter_wiki"]) {
                for(let i = 0; i < data_wiki_set["inter_wiki"].length; i++) {
                    data_wiki_set_inter_wiki["inter_wiki"][
                        data_wiki_set["inter_wiki"][i][0]
                    ] = {
                        "link" : data_wiki_set["inter_wiki"][i][1],
                        "logo" : data_wiki_set["inter_wiki"][i][2]
                    }
                }
            }
            
            let data_var = do_onmark_redirect_render(
                data, 
                data_js, 
                name_doc, 
                name_include
            );
            data = data_var[0];
            data_js = data_var[1];
            let passing = data_var[2];
            
            if(passing === 1) {
                if(test_mode === 'normal') {
                    document.getElementById(name_id).innerHTML = data + '<script>' + data_js + '</script>';
                    eval(data_js);
                } else if(test_mode === 'manual') {
                    return [data, data_js];
                } else {
                	console.log([data, data_js]);
                }
                
                return 0;
            }
            
            data = do_onmark_remark_render(data);
            
            data_var = do_onmark_nowiki_before_render(
                data, 
                data_js, 
                name_include, 
                data_nowiki
            );
            data = data_var[0];
            data_js = data_var[1];
            data_nowiki = data_var[2];
            
            data_var = do_onmark_math_render(
                data, 
                data_js, 
                name_include, 
                data_nowiki
            );
            data = data_var[0];
            data_js = data_var[1];
        
            data_var = do_onmark_include_render(
                data, 
                data_js, 
                name_include, 
                data_nowiki
            );
            data = data_var[0];
            data_js = data_var[1];
        
            data_var = do_onmark_middle_render(
                data, 
                data_js, 
                name_include, 
                data_nowiki, 
                name_doc
            );
            data = data_var[0];
            data_js = data_var[1];
            data_nowiki = data_var[2];
        
            data = do_onmark_text_render(data);
            data_var = do_onmark_heading_render(
                data, 
                data_js, 
                name_doc, 
                name_include
            );
            data = data_var[0];
            data_js = data_var[1];
            
            data = do_onmark_table_render(data);
        
            data_var = do_onmark_link_render(
                data, 
                data_js, 
                name_doc, 
                name_include,
                data_nowiki,
                data_wiki_set_inter_wiki
            );
            data = data_var[0];
            data_js = data_var[1];
            var data_category = data_var[2];
        
            data_var = do_onmark_macro_render(data, data_js);
            data = data_var[0];
            data_js = data_var[1];
            
            data = do_onmark_list_render(data);
            data = do_onmark_hr_render(data);
            data = do_onmark_footnote_render(data, name_include);
            data = do_onmark_last_render(
                data, 
                name_include, 
                data_category
            );
            
            data_js += '' + 
                'do_heading_move();\n' + 
                'get_link_state("' + name_include + '");\n' + 
                'get_file_state("' + name_include + '");\n' + 
        		'get_heading_name();\n' +
                'render_html("' + name_include + 'nowiki_html");\n' +
            ''
            
            if(test_mode === 'normal') {
                document.getElementById(name_id).innerHTML = data + '<script>' + data_js + '</script>';
                
                document.getElementById(name_id).style.display = "";
                document.getElementById(name_id + '_load').style.display = "none";
                
                eval(data_js);
            } else if(test_mode === 'manual') {
                return [data, data_js];
            } else {
            	console.log([data, data_js]);
            }
        }
    }
}
