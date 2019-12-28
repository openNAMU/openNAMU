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
        xhr.open("GET", "/api/w/" + encodeURI(link_data[0]) + "?exist=1", true);
        xhr.send(null);
        
        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(JSON.parse(this.responseText)['exist'] !== '1') {
                    document.getElementById(link_data[1]).className = "not_thing";
                } else {
                    document.getElementById(link_data[1]).className = "";
                }
            } else {
                document.getElementById(link_data[1]).className = "not_thing";
            }
        }
    }

    var data = '\n' + document.getElementById(target).innerHTML + '\n';
    var nowiki_list = [];

    var mid_num = 0;
    var mid_stack = 0;
    var mid_list = [];
    var html_number = 0;
    data = data.replace(/(?:{{{(?:((?:(?! |{{{|}}}|&lt;).)*) ?)|(}}}))/g, function(all, in_data) {
        if(all === '}}}') {
            if(!mid_list[mid_num]) {
                return all;
            } else {
                if(mid_stack > 0) {
                    mid_stack -= 1;
                }

                if(mid_stack > 0) {
                    return all;
                } else {
                    if(mid_num > 0) {
                        mid_num -= 1;
                    }

                    if(mid_list[mid_num] === 'pre') {
                        var return_data = '</code></pre>';
                    } else {
                        var return_data = '</' + mid_list[mid_num] + '>';
                    }

                    mid_list.splice(mid_num, 1);

                    return return_data;
                }
            }
        } else {
            if(mid_stack > 0) {
                mid_stack += 1;

                return all;
            } else {
                mid_num += 1;

                if(in_data.match(/^(#|@|\+|\-)/) && !in_data.match(/^(#|@|\+|\-){2}|(#|@|\+|\-)\\\\/)) {                    
                    if(in_data.match(/^((#|@)([0-9a-f-A-F]{3}){1,2})/)) {
                        mid_list.push('span');

                        if(in_data.match(/^#/)) {
                            return '<span style="color: ' + in_data + ';">';
                        } else {
                            return '<span style="background: ' + in_data + ';">';
                        }
                    } else if(in_data.match(/^((#|@)(\w+))/)) {
                        mid_list.push('span');

                        if(in_data.match(/^#/)) {
                            return '<span style="color: ' + in_data.replace(/^#/, '') + ';">';
                        } else {
                            return '<span style="background: ' + in_data.replace(/^@/, '') + ';">';
                        }
                    } else if(in_data.match(/^(\+|-)([1-5])/)) {
                        mid_list.push('span');

                        var font_size_data = in_data.match(/^(\+|-)([1-5])/);
                        if(font_size_data[1] == '+') {
                            font_size_data = String(Number(font_size_data[2]) * 20 + 100);
                        } else {
                            font_size_data = String(100 - Number(font_size_data[2]) * 10);
                        }

                        return '<span style="font-size: ' + font_size_data + '%;">'
                    } else if(in_data.match(/#!wiki/i)) {
                        mid_list.push('div');

                        if(data.match(/{{{#!wiki style=((?:(?!\n).)+) *\n/i)) {
                            return '<div id="wiki_div_before">';
                        } else {
                            return '<div id="wiki_div">'
                        }
                    } else if(in_data.match(/#!syntax/i)) {
                        mid_list.push('pre');
                        mid_stack += 1;
                        
                        return '<pre><code id="syntax_before">';
                    } else if(in_data.match(/#!folding/i)) {
                        mid_list.push('div');
                        
                        return '<div id="folding_before">';
                    } else if(in_data.match(/#!html/i)) {
                        mid_list.push('span');
                        html_number += 1;

                        return '<span id="html_render_contect_' + String(html_number) + '">';
                    } else {
                        mid_list.push('code');
                        mid_stack += 1;
    
                        return '<code>' + in_data;
                    }
                } else {
                    mid_list.push('code');
                    mid_stack += 1;

                    return '<code>' + in_data;
                }
            }
        }
    });

    console.log(mid_stack);
    console.log(mid_num);
    console.log(mid_list);

    data = data.replace(/<\/div> *\n/ig, '</div>');

    data = data.replace(/<div id="folding_before">((?:(?!\n).)+) *\n/ig, function(all, in_data) {
        return in_data + ' [+]<div id="folding">';
    });
    
    data = data.replace(/<pre><code id="syntax_before">((?:(?!\n).)+) *\n/ig, function(all, in_data) {
        return '<pre><code>';
    });

    data = data.replace(/<div id="wiki_div_before">style=((?:(?!\n).)+) *\n/ig, function(all, in_data) {
        return '<div style=' + in_data.replace(/&quot;/g, "\"").replace(/&#039;/g, "'") + ' id="wiki_div">';
    });

    var nowiki_num = 0;
    data = data.replace(/<code>((?:(?!<\/code>).)+)<\/code>/gm, function(all, in_data) {
        nowiki_num += 1;
        nowiki_list.push(['nowiki_' + String(nowiki_num), in_data]);

        return '<span id="nowiki_' + String(nowiki_num) + '"></span>';
    });

    var math_list = [];
    var math_num = 0;
    data = data.replace(/\[math\(((?:(?!\)]).)+)\)]/ig, function(all, in_data) {
        var math_data = in_data.replace(/&amp;/g, "&").replace(/&lt;/g, "<").replace(/&gt;/g, ">").replace(/&quot;/g, "\"").replace(/&#039;/g, "'");

        math_num += 1;
        math_list.push(['math_' + String(math_num), math_data]);

        return '<span id="math_' + String(math_num) + '"></span>';
    });

    data = data.replace(/~~((?:(?!~~).)+)~~/g, '<s>$1</s>');
    data = data.replace(/--((?:(?!--).)+)--/g, '<s>$1</s>');
    data = data.replace(/__((?:(?!__).)+)__/g, '<u>$1</u>');
    data = data.replace(/'''((?:(?!''').)+)'''/g, '<b>$1</b>');
    data = data.replace(/''((?:(?!'').)+)''/g, '<i>$1</i>');
    data = data.replace(/\^\^((?:(?!\^\^).)+)\^\^/g, '<sup>$1</sup>');
    data = data.replace(/,,((?:(?!,,).)+),,/g, '<sub>$1</sub>');

    data = data.replace(/\n( {1,})\* ([^\n]+)/g, function(all, margin_data, in_data) {
        return '<li style="margin-left: ' + String(margin_data.length * 20) + 'px;">' + in_data + '</li>'
    });

    data = data.replace(/\n( {1,})/g, function(all, margin_data) {
        return '\n<span style="margin-left: ' + String(margin_data.length * 10) + 'px"></span>'
    });

    var link_list = [];
    var link_num = 0;
    var category = '<div id="cate_all"><hr><div id="cate">Category : '
    data = data.replace(/\[\[((?:(?!]]).)+)]]/g, function(all, in_data) {
        if(in_data.match(/^(?:[^|]+)\|(?:.+)$/)) {
            return in_data.replace(/^([^|]+)\|([^|]+)$/, function(all, front_data, back_name) {
                link_list.push([front_data, 'link_' + String(link_num)]);
                link_num += 1;

                return '<a id="link_' + String(link_num - 1) + '" href="/w/' + encodeURI(front_data) + '">' + back_name + '</a>'; 
            });   
        } else {
            if(in_data.match(/^(?:category|분류):/i)) {
                category += '<a href="' + encodeURI(in_data) + '">' + in_data + '</a> | ';

                return '';
            } else if(in_data.match(/^(?:file|파일):/i)) {

            } else {
                link_list.push([in_data, 'link_' + String(link_num)]);
                link_num += 1;

                return '<a id="link_' + String(link_num - 1) + '" href="/w/' + encodeURI(in_data) + '">' + in_data + '</a>'; 
            }
        }
    });

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

    var toc_array = [0, 0, 0, 0, 0, 0];
    var before_data = 0;
    var toc_data = '<div id="toc"><span id="toc_title">TOC</span>\n\n'
    data = data.replace(/\n(={1,6}) ?([^\n]+) (?:={1,6})/g, function(all, num, in_data) {
        num = num.length;
        
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

        return '\n<h' + num + ' id="s-' + toc_num.replace(/\.$/, '') + '"><a href="#toc">' + toc_num + '</a> ' + in_data + '</h' + num + '>';
    });

    toc_data += '</div>';
    data = data.replace(/<\/h([0-9])>\n/g, '</h$1>');

    data = data.replace(/\[([^\]]+)\]/g, function(all, name) {
        if(name.match(/^br$/i)) {
            return '\n'
        } else if(name.match(/^목차$/i)) {
            return toc_data;
        } else if(name.match(/^date|datetime$/i)) {
            return get_today();
        } else {
            return all;
        }
    });

    var ref_num = 0;
    var ref_data = '<hr><ul id="footnote_data">';
    var name_ref_data = {};
    data = data.replace(/(?:\[\*([^ \]]*)(?: ([^\]]+))?\]|\[(?:각주|footnote)])/g, function(all, name_data, in_data) {
        if(all.match(/^\[(?:각주|footnote)]$/i)) {
            var new_ref_data = ref_data;
            ref_data = '<hr><ul id="footnote_data">';
            
            return new_ref_data + '</ul>';
        } else {
            ref_num += 1;
            if(name_data) {
                if(in_data) {
                    name_ref_data[name_data] = in_data;

                    ref_data += '' +
                        '<li>' +
                            '<a id="fn-' + name_data + '" href="#rfn-' + String(ref_num) + '">(' + name_data + ')</a> ' + in_data + ''
                        '</li>' +
                    ''    
                } else {
                    ref_data += '' +
                        '<li>' +
                            '<a href="#rfn-' + String(ref_num) + '">(' + name_data + ')</a>' +
                        '</li>' +
                    ''
                }
            } else {
                ref_data += '' +
                    '<li>' +
                        '<a id="fn-' + String(ref_num) + '" href="#rfn-' + String(ref_num) + '">(' + String(ref_num) + ')</a> ' + in_data + ''
                    '</li>' +
                ''
            }

            if(name_data) {
                return '' +
                    '<sup>' +
                        '<a href="#fn-' + name_data + '" id="rfn-' + String(ref_num) + '" title="' + name_ref_data[name_data].replace(/<([^>]*)>/g, '') + '">' +
                            '(' + name_data + ')' +
                        '</a>' +
                    '</sup>' +
                '';
            } else {
                return '' +
                '<sup>' +
                    '<a href="#fn-' + String(ref_num) + '" id="rfn-' + String(ref_num) + '" title="' + in_data.replace(/<([^>]*)>/g, '') + '">' +
                        '(' + String(ref_num) + ')' +
                    '</a>' +
                '</sup>' +
            '';
            }
        }
    });

    if(ref_data !== '<hr><ul id="footnote_data">') {
        data += ref_data + '</ul>';
    }

    i = 0;
    while(1) {
        if(nowiki_list[i]) {
            data = data.replace('<span id="' + nowiki_list[i][0] + '"></span>', '<code>' + nowiki_list[i][1] + '</code>');

            i += 1;
        } else {
            break;
        }
    }

    data = data.replace(/^(\n| )+/g, '');
    data = data.replace(/(\n| )+$/g, '');
    data = data.replace(/\n/g, '<br>');

    data = data.replace(/&amp;/g, '&');

    document.getElementById(target).innerHTML = data;

    var i = 0;
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

    var i = 0;
    while(1) {
        if(link_list[i]) {
            get_link_state(link_list[i]);

            i += 1;
        } else {
            break;
        }
    }
    render_html("html_render_contect");    
}