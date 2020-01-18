function render_namumark(target) {
    function get_today() {
        var today_data = new Date();

        return '' +
            String(today_data.getFullYear()) + '-' + 
            String(today_data.getMonth() + 1) + '-' + 
            String(today_data.getDate()) + ' ' + 
            (today_data.getHours() < 10 ? '0' + String(today_data.getHours()) : String(today_data.getHours())) + ':' + 
            (today_data.getMinutes() < 10 ? '0' + String(today_data.getMinutes()) : String(today_data.getMinutes())) + ':' + 
            (today_data.getSeconds() < 10 ? '0' + String(today_data.getSeconds()) : String(today_data.getSeconds())) +
        '';
    }

    function get_link_state(link_data) {            
        var xhr = new XMLHttpRequest();
        xhr.open("GET", "/api/w/" + encodeURIComponent(link_data[0]) + "?exist=1", true);
        xhr.send(null);
        
        xhr.onreadystatechange = function() {
            var i = 0;
            while(1) {
                if(document.getElementsByClassName(link_data[1])[i]) {
                    if(this.readyState === 4 && this.status === 200) {
                        if(JSON.parse(this.responseText)['exist'] !== '1') {
                            document.getElementsByClassName(link_data[1])[i].id = "not_thing";
                        } else {
                            document.getElementsByClassName(link_data[1])[i].id = "";
                        }
                    } else {
                        document.getElementsByClassName(link_data[1])[i].id = "not_thing";
                    }

                    i += 1;
                } else {
                    break;
                }
            }
        }
    }

    function get_file_state(file_data) {            
        var file_part = file_data[0].match(/^([^.]+)\.(.+)$/);
        if(file_part) {
            var file_name = file_part[1];
            var file_type = '.' + file_part[2];
        } else {
            var file_name = file_data;
            var file_type = '';
        }

        var xhr = new XMLHttpRequest();
        xhr.open("GET", "/api/sha224/" + encodeURIComponent(file_name), true);
        xhr.send(null);
        
        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                var xhr = new XMLHttpRequest();
                xhr.open("GET", "/api/w/file:" + encodeURIComponent(file_data[0]) + "?exist=1", true);
                xhr.send(null);

                var img_src = JSON.parse(this.responseText)['data'];
                
                xhr.onreadystatechange = function() {
                    if(this.readyState === 4 && this.status === 200 && JSON.parse(this.responseText)['exist'] === '1') {
                        document.getElementById(file_data[1]).innerHTML = '' +
                            '<img style="' + file_data[2] + '" src="/image/' + img_src + file_type + '">' +
                        '';
                    } else {
                        document.getElementById(file_data[1]).innerHTML = '' +
                            '<a href="/upload?name=' + encodeURIComponent(file_name) + '" id="not_thing">' + file_data[0] + '</a>' +
                        '';
                    }
                }
            }
        }
    }

    function divi_link(link_data) {
        var link_part = link_data.match(/^([^|]+)\|(.+)$/);
        if(link_part) {
            return [link_part[2], link_part[1]]
        } else {
            return [link_data, link_data]
        }
    }

    function table_analysis(main_data, cel_data, start_cel, num = 0) {
        var table_class = 'class="'
        
        var div_style = 'style="'
        var table_style = 'style="'
        var cel_style = 'style="'
        var row_style = 'style="'
        
        var row = ''
        var cel = ''

        var table_state_get = main_data.match(/&lt;table ?width=((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            if(main_data.match('^[0-9]+$', table_state_get[1])) {
                div_style += 'width: ' + table_state_get[1] + 'px;';
            } else {
                div_style += 'width: ' + table_state_get[1] + ';';
            }

            table_style += 'width: 100%;';
        }

        table_state_get = main_data.match(/&lt;table ?height=((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            if(main_data.match(/^[0-9]+$/, table_state_get[1])) {
                table_style += 'height: ' + table_state_get[1] + 'px;';
            } else {
                table_style += 'height: ' + table_state_get[1] + ';';
            }
        }

        table_state_get = main_data.match(/&lt;table ?align=((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            if(table_state_get[1] == 'right') {
                div_style += 'float: right;';
            } else if(table_state_get[1] == 'center') {
                table_style += 'margin: auto;';
            }
        }

        table_state_get = main_data.match(/&lt;table ?textalign=((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            num = 1

            if(table_state_get[1] == 'right') {
                table_style += 'text-align: right;';
            } else if(table_state_get[1] == 'center') {
                table_style += 'text-align: center;';
            }
        }

        table_state_get = main_data.match(/&lt;row ?textalign=((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            if(table_state_get[1] == 'right') {
                row_style += 'text-align: right;';
            } else if(table_state_get[1] == 'center') {
                row_style += 'text-align: center;';
            } else {
                row_style += 'text-align: left;';
            }
        }

        table_state_get = main_data.match(/&lt;-((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            cel = 'colspan="' + table_state_get[1] + '"';
        } else {
            cel = 'colspan="' + String(Math.round(start_cel.length / 2)) + '"';
        }

        table_state_get = main_data.match(/&lt;\|((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            row = 'rowspan="' + table_state_get[1] + '"';
        }

        table_state_get = main_data.match(/&lt;rowbgcolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)(?:,(#(?:[0-9a-f-A-F]{3}){1,2}|\w+))?&gt;/);
        if(table_state_get) {
            row_style += 'background: ' + table_state_get[1] + ';';
        }

        table_state_get = main_data.match(/&lt;rowcolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)(?:,(#(?:[0-9a-f-A-F]{3}){1,2}|\w+))?&gt;/);
        if(table_state_get) {
            row_style += 'color: ' + table_state_get[1] + ';';
        }

        table_state_get = main_data.match(/&lt;table ?bordercolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)(?:,(#(?:[0-9a-f-A-F]{3}){1,2}|\w+))?&gt;/);
        if(table_state_get) {
            table_style += 'border: ' + table_state_get[1] + ' 2px solid;';
        }

        table_state_get = main_data.match(/&lt;table ?bgcolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)(?:,(#(?:[0-9a-f-A-F]{3}){1,2}|\w+))?&gt;/);
        if(table_state_get) {
            table_style += 'background: ' + table_state_get[1] + ';';
        }

        table_state_get = main_data.match(/&lt;table ?color=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)(?:,(#(?:[0-9a-f-A-F]{3}){1,2}|\w+))?&gt;/);
        if(table_state_get) {
            table_style += 'color: ' + table_state_get[1] + ';';
        }

        table_state_get = main_data.match(/&lt;(?:bgcolor=)?(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)(?:,(#(?:[0-9a-f-A-F]{3}){1,2}|\w+))?&gt;/);
        if(table_state_get) {
            cel_style += 'background: ' + table_state_get[1] + ';';
        }

        table_state_get = main_data.match(/&lt;color=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)(?:,(#(?:[0-9a-f-A-F]{3}){1,2}|\w+))?&gt;/);
        if(table_state_get) {
            cel_style += 'color: ' + table_state_get[1] + ';';
        }

        table_state_get = main_data.match(/&lt;width=((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            if(table_state_get[1].match(/^[0-9]+$/)) {
                cel_style += 'width: ' + table_state_get[1] + 'px;';
            } else {
                cel_style += 'width: ' + table_state_get[1] + ';';
            }
        }

        table_state_get = main_data.match(/&lt;height=((?:(?!&gt;).)*)&gt;/);
        if(table_state_get) {
            if(table_state_get[1].match(/^[0-9]+$/)) {
                cel_style += 'height: ' + table_state_get[1] + 'px;';
            } else {
                cel_style += 'height: ' + table_state_get[1] + ';';
            }
        }

        var text_right = main_data.match(/&lt;\)&gt;/);
        var text_center = main_data.match(/&lt;:&gt;/);
        var text_left = main_data.match(/&lt;\(&gt;/);
        if(text_right) {
            cel_style += 'text-align: right;';
        } else if(text_center) {
            cel_style += 'text-align: center;';
        } else if(text_left) {
            cel_style += 'text-align: left;';
        } else if(num == 0) {
            if(cel_data.match(/^ /) && cel_data.match(/ $/)) {
                cel_style += 'text-align: center;';
            } else if(cel_data.match(/^ /)) {
                cel_style += 'text-align: right;';
            } else if(cel_data.match(/ $/)) {
                cel_style += 'text-align: left;';
            }
        }

        table_state_get = main_data.match(/&lt;table ?class=((?:(?!&gt;).)+)&gt;/);
        if(table_state_get) {
            table_class += table_state_get[1];
        }

        div_style += '"';
        table_style += '"';
        cel_style += '"';
        row_style += '"';

        table_class += '"';

        return [table_style, row_style, cel_style, row, cel, table_class, num, div_style]
    }

    function table_render(data) {
        var table_num = 0;
        while(1) {
            var table_data = data.match(/\n((?:(?:(?:(?:\|\|)+(?:(?:(?!\|\|).(?:\n)*)*))+)\|\|(?:\n)?)+)/);
            if(table_data) {
                table_data = table_data[1];
                
                var get_table_data = table_data.match(/^((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*((?:(?!\|\|).\n*)*)/);
                if(get_table_data) {
                    table_return_data = table_analysis(get_table_data[2], get_table_data[3], get_table_data[1]);
                    table_num = table_return_data[6];
    
                    table_data = table_data.replace(
                        /^((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*/,
                        '\n' + 
                        '<div class="table_safe" ' + table_return_data[7] + '>' + 
                            '<table ' + table_return_data[5] + ' ' + table_return_data[0] + '>' + 
                                '<tr ' + table_return_data[1] + '>' + 
                                    '<td ' + table_return_data[2] + ' ' + table_return_data[3] + ' ' + table_return_data[4] + '>'
                    );
                }
    
                table_data = table_data.replace(/\|\|\n?$/, '</td></tr></table></div>');
    
                while(1) {
                    get_table_data = table_data.match(/\|\|\n((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*((?:(?!\|\||<\/td>).\n*)*)/);
                    if(get_table_data) {
                        table_return_data = table_analysis(get_table_data[2], get_table_data[3], get_table_data[1], table_num);
    
                        table_data = table_data.replace(
                            /\|\|\n((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*/,
                            '</td></tr><tr ' + table_return_data[1] + '><td ' + table_return_data[2] + ' ' + table_return_data[3] + ' ' + table_return_data[4] + '>'
                        );
                    } else {
                        break;
                    }
                }
    
                while(1) {
                    get_table_data = table_data.match(/((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*((?:(?:(?!\|\||<\/td>).)|\n)*\n*)/);
                    if(get_table_data) {
                        table_return_data = table_analysis(get_table_data[2], get_table_data[3].replace(/\n/g, ' '), get_table_data[1], table_num);
    
                        table_data = table_data.replace(
                            /((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*/,
                            '</td><td ' + table_return_data[2] + ' ' + table_return_data[3] + ' ' + table_return_data[4] + '>'
                        );
                    } else {
                        break;
                    }
                }
    
                data = data.replace(/\n((?:(?:(?:(?:\|\|)+(?:(?:(?!\|\|).(?:\n)*)*))+)\|\|(?:\n)?)+)/, table_data);
            } else {
                break;
            }
        }

        return data;
    }

    var data = '\n' + document.getElementById(target).innerHTML + '\n';
    var title = window.location.pathname.replace(/^\/w\//, '');

    var math_list = [];
    var math_num = 0;
    data = data.replace(/\[math\(((?:(?!\)]).)+)\)]/ig, function(all, in_data) {
        var math_data = in_data.replace(/&amp;/g, "&").replace(/&lt;/g, "<").replace(/&gt;/g, ">").replace(/&quot;/g, "\"").replace(/&#039;/g, "'");

        math_num += 1;
        math_list.push(['math_' + String(math_num), math_data]);

        return '<span id="math_' + String(math_num) + '"></span>';
    });

    var mid_num = 0;
    var mid_stack = 0;
    var mid_list = [];
    var html_num = 0;
    var fol_num = 0;
    var mid_regex = /(?:{{{(?:((?:(?! |{{{|}}}|&lt;).)*) ?)|(?:}}}))/;
    while(1) {
        var all_mid_data = data.match(mid_regex);
        if(all_mid_data) {
            var all = all_mid_data[0];
            var in_data = all_mid_data[1];

            if(all === '}}}') {
                if(mid_stack > 0) {
                    mid_stack -= 1;
                }

                if(mid_stack > 0) {
                    data = data.replace(mid_regex, '</mid>');
                } else {
                    if(mid_num > 0) {
                        mid_num -= 1;
                    }

                    if(!mid_list[mid_num]) {
                        var return_data = '';
                    } else if(mid_list[mid_num] === 'pre') {
                        var return_data = '</code></pre>';
                    } else if(mid_list[mid_num] === 'div_2') {
                        var return_data = '</div_1></div>';
                    } else {
                        var return_data = '</' + mid_list[mid_num] + '>';
                    }

                    if(return_data !== '') {
                        mid_list.splice(mid_num, 1);

                        data = data.replace(mid_regex, return_data);
                    } else {
                        data = data.replace(mid_regex, '</mid>');
                    }
                }
            } else {
                if(mid_stack > 0) {
                    mid_stack += 1;

                    data = data.replace(mid_regex, all.replace('{{{', '<mid>'));
                } else {
                    mid_num += 1;

                    if(in_data.match(/^(#|@|\+|\-)/) && !in_data.match(/^(#|@|\+|\-){2}|(#|@|\+|\-)\\\\/)) {                    
                        if(in_data.match(/^((#|@)([0-9a-f-A-F]{3}){1,2})/)) {
                            mid_list.push('span');

                            if(in_data.match(/^#/)) {
                                data = data.replace(mid_regex, '<span style="color: ' + in_data + ';">');
                            } else {
                                data = data.replace(mid_regex, '<span style="background: ' + in_data + ';">');
                            }
                        } else if(in_data.match(/^((#|@)(\w+))/)) {
                            mid_list.push('span');

                            if(in_data.match(/^#/)) {
                                data = data.replace(mid_regex, '<span style="color: ' + in_data.replace(/^#/, '') + ';">');
                            } else {
                                data = data.replace(mid_regex, '<span style="background: ' + in_data.replace(/^@/, '') + ';">');
                            }
                        } else if(in_data.match(/^(\+|-)([1-5])/)) {
                            mid_list.push('span');

                            var font_size_data = in_data.match(/^(\+|-)([1-5])/);
                            if(font_size_data[1] == '+') {
                                font_size_data = String(Number(font_size_data[2]) * 20 + 100);
                            } else {
                                font_size_data = String(100 - Number(font_size_data[2]) * 10);
                            }

                            data = data.replace(mid_regex, '<span style="font-size: ' + font_size_data + '%;">');
                        } else if(in_data.match(/#!wiki/i)) {
                            mid_list.push('div_1');


                            var div_style_data = data.match(/{{{#!wiki(?: style=([^\n]+))?\n?/i);
                            if(!div_style_data[1]) {
                                div_style_data = '';
                            } else {
                                div_style_data = div_style_data[1];
                                console.log(div_style_data);

                                div_style_data = div_style_data.replace(/('|")/g, '');
                                div_style_data = div_style_data.replace(/position/ig, '');
                            }

                            data = data.replace(/{{{#!wiki(?: style=([^\n]+))?\n?/i, '<div id="wiki_div" style="' + div_style_data + '">');
                        } else if(in_data.match(/#!syntax/i)) {
                            mid_list.push('pre');
                            mid_stack += 1;

                            var syntax_data = data.match(/{{{#!syntax(?: ([^\n]+))?\n?/i);
                            if(!syntax_data[1]) {
                                syntax_data = 'python';
                            } else {
                                syntax_data = syntax_data[1];
                            }
                            
                            data = data.replace(/{{{#!syntax(?: ([^\n]+))?\n?/i, '<pre id="syntax"><code class="' + syntax_data + '">');
                        } else if(in_data.match(/#!folding/i)) {
                            mid_list.push('div_2');

                            var fol_data = data.match(/{{{#!folding(?: ([^\n]+))?\n?/i);
                            if(!fol_data[1]) {
                                fol_data = 'TEST';
                            } else {
                                fol_data = fol_data[1];
                            }
                            
                            fol_num += 1;
                            
                            data = data.replace(/{{{#!folding( ([^\n]+))?\n?/i, '' +
                                fol_data +
                                '<div style="display: inline-block;">' + 
                                    '<a href="javascript:void(0);" onclick="do_open_folding(\'folding_' + String(fol_num) + '\', this);">' +
                                        '[+]' +
                                    '</a>' +
                                '</div>' +
                                '<div id="folding_' + String(fol_num) + '" style="display: none;">' +
                                    '<div id="wiki_div" style="">' +
                            '');
                        } else if(in_data.match(/#!html/i)) {
                            mid_list.push('span');
                            html_num += 1;

                            data = data.replace(mid_regex, '<span id="html_render_contect_' + String(html_num) + '">');
                        } else {
                            mid_list.push('code');
                            mid_stack += 1;
        
                            data = data.replace(mid_regex, '<code>' + in_data);
                        }
                    } else {
                        mid_list.push('code');
                        mid_stack += 1;

                        data = data.replace(mid_regex, '<code>' + in_data);
                    }
                }
            }
        } else {
            break;
        }
    }

    data = data.replace(/<mid>/g, '{{{');
    data = data.replace(/<\/mid>/g, '}}}');
    data = data.replace(/<\/div> *\n/ig, '</div>');

    var nowiki_num = 0;
    var nowiki_list = {};
    data = data.replace(/<code( (?:class="(?:[^"]+)"))?>(\n*((?:(?!<\/code>).)+\n*)+)<\/code>/g, function(all, class_data, in_data) {
        nowiki_num += 1;
        nowiki_list['nowiki_' + String(nowiki_num)] = in_data;

        if(class_data) {
            return '<code' + class_data + '><span id="nowiki_' + String(nowiki_num) + '"></span></code>';
        } else {
            return '<span id="nowiki_' + String(nowiki_num) + '"></span>';
        }
    });

    data = data.replace(/\r\n/g, '\n');
    data = data.replace(/&amp;/g, '&');

    data = data.replace(/\n(?: +)\|\|/g, '\n||');
    data = data.replace(/\|\|(?: +)\n/g, '||\n');

    data = data.replace(/\n##(?:(?:(?!\n).)+)/g, '');

    while(1) {
        wiki_table_data = data.match(/<div id="wiki_div"( (?:[^>]*))>((?:(?:\n| )*)(?:(?!<div id="wiki_div"|<\/div_1>).\n*)+)<\/div_1>/i);
        if(wiki_table_data) {
            if(wiki_table_data[2].match(/\|\|/)) {
                var end_table_render = table_render('\n' + wiki_table_data[2] + '\n').replace(/^\n/, '').replace(/\n$/, '');
            } else {
                var end_table_render = wiki_table_data[2];
            }

            data = data.replace(
                /<div id="wiki_div"( (?:[^>]*))>((?:(?:\n| )*)(?:(?!<div id="wiki_div"|<\/div_1>).\n*)+)<\/div_1>/i, 
                '<div' + wiki_table_data[1] + '>' + end_table_render + '</div>'
            );
        } else {
            break;
        }
    }
 
    data = data.replace(/<\/td>/g, '</td_1>');

    data = data.replace(/~~((?:(?!~~).)+)~~/g, '<s>$1</s>');
    data = data.replace(/--((?:(?!--).)+)--/g, '<s>$1</s>');
    data = data.replace(/__((?:(?!__).)+)__/g, '<u>$1</u>');
    data = data.replace(/'''((?:(?!''').)+)'''/g, '<b>$1</b>');
    data = data.replace(/''((?:(?!'').)+)''/g, '<i>$1</i>');
    data = data.replace(/\^\^((?:(?!\^\^).)+)\^\^/g, '<sup>$1</sup>');
    data = data.replace(/,,((?:(?!,,).)+),,/g, '<sub>$1</sub>');

    var toc_array = [0, 0, 0, 0, 0, 0];
    var before_data = 0;
    var edit_number = 0;
    var toc_data = '<div id="toc"><span id="toc_title">TOC</span>\n\n'
    data = data.replace(/\n(={1,6}) ?([^\n]+) (?:={1,6})/g, function(all, num, in_data) {
        num = num.length;
        edit_number += 1;
        
        if(before_data > num) {
            var i = num;
            while(1) {
                if(i == 6) {
                    break;
                }

                toc_array[i] = 0;
                i += 1;
            }
        }

        before_data = num;
        toc_array[num - 1] += 1;
        num = String(num);
        var toc_num = (toc_array.join('.') + '.').replace(/0\./g, '');
        if(!toc_num.match(/\./)) {
            toc_num += '0.';
        }

        toc_data += '' + 
            '<span style="margin-left: ' + String(10 * (toc_num.length / 2) - 10) + 'px;">' + 
                '<a href="#s-' + toc_num.replace(/\.$/, '') + '">' + toc_num + '</a> ' + in_data + 
            '</span>' +
            '\n' + 
        '';

        return '' +
            '\n' +
            '<h' + num + ' id="s-' + toc_num.replace(/\.$/, '') + '">' +
                '<a href="#toc">' + toc_num + '</a> ' + in_data +
                '<span style="font-size: 12px">' +
                    '<a href="/edit/' + title + '?section=' + String(edit_number) + '">(Edit)</a>' +
                '</span>' +
            '</h' + num + '>' +
        '';
    });

    toc_data += '</div>';
    data = data.replace(/<\/h([0-9])>\n/g, '</h$1>');

    while(1) {
        if(data.match(/(\n(?:&gt; ?(?:[^\n]+)?\n?)+)/)) {
            data = data.replace(/(\n(?:&gt; ?(?:[^\n]+)?\n?)+)/, function(all, in_data) {
                var new_in_data = in_data;
                new_in_data = new_in_data.replace(/^\n&gt; ?/, '');
                new_in_data = new_in_data.replace(/\n&gt; ?/g, '\n');
                new_in_data = new_in_data.replace(/\n$/, '');

                return '\n<blockquote>' + new_in_data + '</blockquote>\n';
            });
        } else {
            break;
        }
    }

    while(1) {
        if(data.match(/\n-{4,9}\n/)) {
            data = data.replace(/\n-{4,9}\n/, function() {
                return '\n<hr>\n';
            });
        } else {
            break;
        }
    }

    data = data.replace(/(\n +\* ?(?:(?:(?!\|\|).)+))\|\|/g, '$1\n ||');

    data = data.replace(/\n( {1,})\* ([^\n]+)/g, function(all, margin_data, in_data) {
        return '<li style="margin-left: ' + String(margin_data.length * 20) + 'px;">' + in_data + '</li>'
    });
    data = data.replace(/\|\|<li/g, '||\n<li');

    data = data.replace(/\n( {1,})/g, function(all, margin_data) {
        return '\n<span style="margin-left: ' + String(margin_data.length * 10) + 'px"></span>'
    });

    data = table_render(data);

    var link_list = [];
    var file_list = [];
    var link_num = 0;
    var file_num = 0;
    var category = ''
    while(1) {
        if(data.match(/\[\[((?:(?!\[\[|]]).)+)]]/)) {
            data = data.replace(/\[\[((?:(?!\[\[|]]).)+)]]/, function(all, in_data) {
                if(in_data.match(/^(?:category|분류):/i)) {
                    var back_data = in_data.replace(/^(?:category|분류):/i, '');
                    var front_data = back_data;
                    back_data = 'category:' + back_data.replace(/#blur$/, '');
                    
                    if(front_data.match(/#blur$/)) {
                        front_data = '#blur';
                    }

                    link_list.push([back_data, 'link_' + String(link_num)]);
                    link_num += 1;

                    if(category === '') {
                        category += '<div id="cate_all"><hr><div id="cate">Category : '
                    }

                    category += '<a class="link_' + String(link_num - 1) + '" href="' + encodeURIComponent(back_data) + '">' + front_data + '</a> | ';

                    return '';
                } else if(in_data.match(/^(?:file|파일):/i)) {
                    if(in_data.match(/\|/)) {
                        var file_name = in_data.replace(/^(?:file|파일):/i, '');
                        file_name = file_name.match(/^([^|]+)/)[1];
                    } else {
                        var file_name = in_data.replace(/^(?:file|파일):/i, '');
                    }

                    var file_style = '';
                    
                    var file_state = in_data.match(/\|width=([^|]+)/);
                    if(file_state) {
                        file_style += 'width:' + file_state[1];
                    }

                    file_state = in_data.match(/\|height=([^|]+)/);
                    if(file_state) {
                        file_style += 'height:' + file_state[1];
                    }

                    file_list.push([file_name, 'file_' + String(file_num), file_style]);
                    file_num += 1;
                    
                    return '<span id="file_' + String(file_num - 1) + '"></span>';
                } else if(in_data.match(/^http(?:s)?:\/\//i)) {
                    var link_part = divi_link(in_data);
                    
                    var front_data = link_part[0];
                    var back_data = link_part[1];

                    return '<a id="out_link" href="' + back_data + '">' + front_data + '</a>'; 
                } else {
                    var link_part = divi_link(in_data.replace(/^:/, ''));
                    
                    var front_data = link_part[0];
                    var back_data = link_part[1];

                    link_list.push([back_data, 'link_' + String(link_num)]);
                    link_num += 1;

                    return '<a class="link_' + String(link_num - 1) + '" href="/w/' + encodeURIComponent(back_data) + '">' + front_data + '</a>'; 
                }
            });
        } else {
            break;
        }
    }

    if(category !== '') {
        category = category.replace(/ \| $/, '') + '</div></div>'
    }

    data = data.replace(/\[([^(\]]+)\(((?:(?!\)]).)+)\)]/g, function(all, name, in_data) {
        if(name.match(/^youtube|kakaotv|nicovideo$/i)) {
            var video_code = in_data.match(/^([^,]+)/);
            if(video_code) {
                video_code = video_code[1];
            } else {
                video_code = 'test';
            }

            if(name === 'youtube') {
                var video_src = 'https://www.youtube.com/embed/' + video_code
            } else if(name === 'kakaotv') {
                var video_src = 'https://tv.kakao.com/embed/player/cliplink/' + video_code +'?service=kakao_tv'
            } else {
                var video_src = 'https://embed.nicovideo.jp/watch/' + video_code
            }

            var width_data = in_data.match(/, *width=([^,]+)/);
            if(width_data) {
                width_data = width_data[1];
            } else {
                width_data = '560';
            }

            var height_data = in_data.match(/, *height=([^,]+)/);
            if(height_data) {
                height_data = height_data[1];
            } else {
                height_data = '315';
            }

            return '' +
                '<iframe ' +
                    'width="' + width_data + '" ' +
                    'height="' + height_data + '" ' +
                    'src="' + video_src + '" ' +
                    'allowfullscreen>' +
                '</iframe>' +
            '';
        } else if(name.match(/^ruby$/i)) {
            var main_text = in_data.match(/^([^,]+)/);
            if(main_text) {
                main_text = main_text[1];
            } else {
                main_text = 'test';
            }

            var ruby_text = in_data.match(/, *ruby=([^,]+)/);
            if(ruby_text) {
                ruby_text = ruby_text[1];
            } else {
                ruby_text = 'test';
            }

            var color_text = in_data.match(/, *color=([^,]+)/);
            if(color_text) {
                color_text = 'color:' + color_text[1];
            } else {
                color_text = '';              
            }

            return '' +
                '<ruby>' +
                    main_text +
                    '<rp>(</rp>' +
                    '<rt>' +
                        '<span style="' + color_text + '">' + ruby_text + '</span>' +
                    '</rt>' +
                    '<rp>)</rp>' +
                '</ruby>' +
            '';
        } else if(name.match(/^anchor$/i)) {
            return '<span id="' + in_data + '"></span>';
        } else {
            return all;
        }
    });

    data = data.replace(/\[([^\]]+)\]/g, function(all, name) {
        if(name.match(/^br$/i)) {
            return '\n'
        } else if(name.match(/^목차|tableofcontents$/i)) {
            return toc_data;
        } else if(name.match(/^date|datetime$/i)) {
            return get_today();
        } else {
            return all;
        }
    });

    var ref_num = 0;
    var ref_data = '';
    var name_ref_data = {};
    while(1) {
        if(data.match(/(?:\[\*([^ \]]*)(?: ((?:(?!\[\*|\]).)+))?\]|\[(?:각주|footnote)])/)) {
            data = data.replace(/(?:\[\*([^ \]]*)(?: ((?:(?!\[\*|\]).)+))?\]|\[(?:각주|footnote)])/, function(all, name_data, in_data) {
                if(ref_num === 0) {
                    ref_data += '<hr><ul id="footnote_data">';
                }

                if(all.match(/^\[(?:각주|footnote)]$/i)) {
                    var new_ref_data = ref_data;
                    ref_data = '<hr><ul id="footnote_data">';
                    
                    return new_ref_data + '</ul>';
                } else {
                    ref_num += 1;
                    var fn_num = String(ref_num);

                    if(name_data) {
                        if(in_data) {
                            var fn_data = in_data;
                            var fn_name = name_data;
                            name_ref_data[name_data] = fn_data;
                        } else {
                            var fn_name = name_data;
                            if(name_ref_data[name_data]) {
                                var fn_data = name_ref_data[name_data];
                            } else {
                                var fn_data = '';
                            }
                        }
                    } else {
                        var fn_name = fn_num;
                        var fn_data = in_data;
                    }

                    ref_data += '' +
                        '<li>' +
                            '<a id="cfn-' + fn_num + '" ' +
                                'href="#rfn-' + fn_num + '" ' +
                                'onclick="do_open_foot(\'fn-' + fn_num + '\', 1);">' +
                                '(' + fn_name + ')' +
                            '</a> ' + fn_data +
                        '</li>' +
                    ''

                    if(name_data) {
                        fn_name = name_data;
                    } else {
                        fn_name = fn_num;
                    }

                    return '' +
                        '<sup>' +
                            '<a href="#fn-' + fn_num + '" ' +
                                'id="rfn-' + fn_num + '" ' +
                                'onclick="do_open_foot(\'fn-' + fn_num + '\');">' +
                                '(' + fn_name + ')' +
                            '</a>' +
                        '</sup>' +
                    '';
                }
            });
        } else {
            break;
        }
    }

    if(ref_data !== '') {
        data += ref_data + '</ul>';
    }

    var i = 1;
    while(1) {
        if(nowiki_list['nowiki_' + String(i)]) {
            data = data.replace('<span id="nowiki_' + String(i) + '"></span>', '<code>' + nowiki_list['nowiki_' + String(i)] + '</code>');

            i += 1;
        } else {
            break;
        }
    }

    data = data.replace(/<\/td_1>/g, '</td>');

    data = data.replace(/^(\n| )+/g, '');
    data = data.replace(/(\n| )+$/g, '');
    data = data.replace(/\n/g, '<br>');

    data = data.replace(/&amp;/g, '&');
    data += category;

    document.getElementById(target).innerHTML = data;

    i = 0;
    while(1) {
        if(math_list[i]) {
            try {
                katex.render(math_list[i][1], document.getElementById(math_list[i][0]));
            } catch {
                try {
                    document.getElementById(math_list[i][0]).innerHTML = '<span style="color: red;">' + math_list[i][1] + '</span>';
                } catch {}
            }

            i += 1;
        } else {
            break;
        }
    }

    i = 0;
    while(1) {
        if(link_list[i]) {
            get_link_state(link_list[i]);

            i += 1;
        } else {
            break;
        }
    }

    i = 0;
    while(1) {
        if(file_list[i]) {
            get_file_state(file_list[i]);

            i += 1;
        } else {
            break;
        }
    }

    console.log('render end')
    hljs.initHighlightingOnLoad();
    render_html("html_render_contect");    

    // v0.0.8
    // 완성 직전
}