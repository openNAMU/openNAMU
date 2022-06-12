"use strict";

class opennamu_render_markdown {
    // Init Part
    constructor(
        render_part_id,
        render_part_id_after,
        render_part_id_add,
        doc_name
    ) {
        this.doc_data = document.getElementById(render_part_id_add + render_part_id).innerHTML;
            
        this.doc_data = this.doc_data.replace(/_/g, '<uBar>');

        this.doc_data = this.doc_data.replace(/&amp;/g, '&');
        this.doc_data = '<brStart>\n' + this.doc_data + '\n<brEnd>';
        this.doc_data.replace(/\r/g, '');

        this.doc_name = doc_name;
            
        this.parser_data_temp = {};
        this.parser_data_temp_other = {};
        this.parser_data_temp_other['toc'] = '';
        this.parser_data_temp_other['footnote'] = '';
        this.parser_data_temp_other['category'] = '';

        this.parser_data_js = [];
        
        this.parser_count = {};
        this.parser_count['parser'] = 0;
        this.parser_count['nowiki'] = 0;
        
        this.render_part_id = render_part_id;
        this.render_part_id_add = render_part_id_add;
        this.render_part_id_after = render_part_id_after;
    }
    
    // Func Part
    do_func_parser_to_text(data, parser_type = 'parser') {
        let parser_data_temp = this.parser_data_temp;
        let parser_match;
        if(parser_type === 'nowiki' || parser_type === 'nowikiLink' || parser_type === 'nowikiEnd') {
            parser_match = /<(\/?nowiki[0-9]+Span)>/;
        } else {
            parser_match = /<(\/?render[0-9]+Span)>/;
        }
        
        while(data.match(parser_match)) {
            data = data.replace(parser_match, function(match, x1) {
                if(parser_type === 'nowikiEnd') {
                    return parser_data_temp[x1 + 'End'];
                } else if(parser_type === 'nowikiLink') {
                    let nowiki_data = parser_data_temp[x1];
                    nowiki_data = nowiki_data.replace(/\\(.)/g, '$1');
                    
                    return nowiki_data;
                } else {
                    return parser_data_temp[x1];
                }
            });
        }
        
        return data;
    }
    
    do_func_xss_encode(data) {
        data = data.replace(/'/g, '&#x27;');
        data = data.replace(/"/g, '&quot;');
        data = data.replace(/</g, '&lt;');
        data = data.replace(/</g, '&gt;');
        data = data.replace(/&/g, '&amp;');
        
        return data;
    }
    
    do_func_xss_decode(data) {
        data = data.replace(/&#x27;/g, '\'');
        data = data.replace(/&quot;/g, '"');
        data = data.replace(/&lt;/g, '<');
        data = data.replace(/&gt;/g, '>');
        data = data.replace(/&amp;/g, '&');
        
        return data;
    }
    
    // Render Part
    do_part_text() {
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        this.doc_data = this.doc_data.replace(/~~((?:(?!~~|\n).)+)~~/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<s>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</s>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/\*\*((?:(?!\*\*|\n).)+)\*\*/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<b>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</b>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/<uBar><uBar>((?:(?!<uBar><uBar>|\n).)+)<uBar><uBar>/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<b>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</b>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/\*([^*\n]+)\*/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<i>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</i>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/<uBar>(((?!<uBar>).)+)<uBar>/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<i>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</i>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/&lt;ins&gt;((?:(?!&lt;ins&gt;|&lt;\/ins&gt;|\n).)+)&lt;\/ins&gt;/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<u>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</u>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/&lt;sub&gt;((?:(?!&lt;sub&gt;|&lt;\/sub&gt;|\n).)+)&lt;\/sub&gt;/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<sub>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</sub>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/&lt;sup&gt;((?:(?!&lt;sup&gt;|&lt;\/sup&gt;|\n).)+)&lt;\/sup&gt;/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<sup>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</sup>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
    }
    
    do_part_heading() {
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        let toc_data = '';
        
        let heading_n = 0;
        let heading_list = [0, 0, 0, 0, 0, 0];
        let heading_regex = /\n(#{1,6})([^#][^\n]*)\n/;
        while(this.doc_data.match(heading_regex)) {
            this.doc_data = this.doc_data.replace(heading_regex, function(match, x1, x2) {
                let heading_level = x1.length - 1;
                let heading_level_str = String(heading_level + 1);

                heading_list[heading_level] += 1;
                for(let for_a = heading_level + 1; for_a < 6; for_a++) {
                    heading_list[for_a] = 0;
                }

                let heading_list_str = '';
                for(let for_a = 0; for_a < 6; for_a++) {
                    if(heading_list[for_a] !== 0) {
                        heading_list_str += String(heading_list[for_a]) + '.'
                    }
                }
                
                let heading_list_str_2 = heading_list_str.replace(/\.$/, '');
                
                heading_n += 1;
                let heading_n_str = String(heading_n);
                
                toc_data += '' +
                    '<a href="#opennamu_heading_' + heading_list_str_2 + '">' + heading_list_str + '</a> ' +
                    '<span id="opennamu_TOC_content_' + heading_n_str + '"></span>' +
                    '<br>' +
                ''
                
                let heading_data = x2;
                heading_data = heading_data.replace(/^ /, '');

                return '' + 
                    '\n<brEnd>' + 
                    '<h' + heading_level_str + ' id="opennamu_heading_' + heading_list_str_2 + '">' + 
                        '<a href="#opennamu_TOC">' + heading_list_str + '</a> ' + 
                        '<span id="opennamu_heading_content_' + heading_n_str + '">' + heading_data + '</span>' + 
                    '</h' + heading_level_str + '>' +
                    '<brStart>\n' +
                '';
            });
        }
        
        this.parser_data_temp_other['toc'] = toc_data;
    }
    
    do_part_image() {
        let render_main = this;
        let render_part_id_add = this.render_part_id_add;
        
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        this.doc_data = this.doc_data.replace(/!\[([^\[\]\n]*)\]\(([^\(\)\n]*)\)/g, function(match, x1, x2) {
            if(x1 === '' && x2 === '') {
                return '<imageBlink>';
            } else {
                if(x2 !== '' && x2.match(/^https?:\/\//)) {
                    parser_count += 1;
                    let parser_count_str = String(parser_count);
                    
                    let image_src = render_main.do_func_xss_encode(x2);
                    let image_alt;
                    if(x1 !== '') {
                        image_alt = render_main.do_func_xss_encode(x1);
                    } else {
                        image_alt = image_src;
                    }
                    
                    parser_data_temp['render' + parser_count_str + 'Span'] = '<img alt="' + image_alt + '" src="' + image_src + '">';
                    parser_data_temp['/render' + parser_count_str + 'Span'] = '';

                    return '<render' + parser_count_str + 'Span></render' + parser_count_str + 'Span>';
                } else {
                    parser_count += 1;
                    let parser_count_str = String(parser_count);

                    parser_data_temp['render' + parser_count_str + 'Span'] = '<img>';
                    parser_data_temp['/render' + parser_count_str + 'Span'] = '';

                    return '<render' + parser_count_str + 'Span></render' + parser_count_str + 'Span>';
                }
            }
        });
        
        this.doc_data = this.doc_data.replace(/<imageBlink>/g, '![]()');
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
    }
    
    do_part_link() {
        let render_main = this;
        let render_part_id_add = this.render_part_id_add;
        
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        this.doc_data = this.doc_data.replace(/\[([^\[\]\n]*)\]\(([^\(\)\n]*)\)/g, function(match, x1, x2) {
            if(x1 === '' && x2 === '') {
                return '<linkBlink>';
            } else {
                if(x2 !== '' && x2.match(/^https?:\/\//)) {
                    parser_count += 1;
                    let parser_count_str = String(parser_count);
                    
                    let link_main = render_main.do_func_parser_to_text(x2, 'nowikiLink');
                    link_main = render_main.do_func_xss_encode(link_main);
                    
                    let link_sub;
                    if(x1 === '') {
                        link_sub = x2;
                    } else {
                        link_sub = x1;
                    }
                    
                    parser_data_temp['render' + parser_count_str + 'Span'] = '<a class="opennamu_link_out" href="' + link_main + '">';
                    parser_data_temp['/render' + parser_count_str + 'Span'] = '</a>';

                    return '<render' + parser_count_str + 'Span>' + link_sub + '</render' + parser_count_str + 'Span>';
                } else {
                    parser_count += 1;
                    let parser_count_str = String(parser_count);

                    let link_main;
                    let link_sub;
                    let link_title;
                    if(x2 === '') {
                        link_main = x1;
                        link_sub = x1;
                    } else if(x1 === '') {
                        link_main = x2;
                        link_sub = x2;
                    } else {
                        link_main = x2;
                        link_sub = x1;
                    }
                    
                    link_main = render_main.do_func_parser_to_text(link_main, 'nowikiLink');
                    link_main = render_main.do_func_xss_encode(link_main);
                    
                    link_title = link_main;

                    link_main = render_main.do_func_xss_decode(link_main);
                    link_main = opennamu_do_url_encode(link_main);

                    parser_data_temp['render' + parser_count_str + 'Span'] = '<a class="' + render_part_id_add + 'opennamu_link" title="' + link_title + '" href="/w/' + link_main + '">';
                    parser_data_temp['/render' + parser_count_str + 'Span'] = '</a>';

                    return '<render' + parser_count_str + 'Span>' + link_sub + '</render' + parser_count_str + 'Span>';
                }
            }
        });
        
        this.doc_data = this.doc_data.replace(/<linkBlink>/g, '[]()');
        
        this.doc_data = this.doc_data.replace(/&lt;(https?:\/\/(?:(?:(?!&lt;|&gt;).)+))&gt;/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            let link_main = render_main.do_func_parser_to_text(x1, 'nowikiLink');
            link_main = render_main.do_func_xss_encode(link_main);
            
            let link_sub = x1;

            parser_data_temp['render' + parser_count_str + 'Span'] = '<a class="opennamu_link_out" href="' + link_main + '">';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</a>';

            return '<render' + parser_count_str + 'Span>' + link_sub + '</render' + parser_count_str + 'Span>';
        });
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
    }
    
    do_part_footnote_list() {
        if(this.parser_data_temp_other['footnote'] !== '') {
            let footnote = '';
            footnote += '<ul id="footnote_data">';
            footnote += this.parser_data_temp_other['footnote'];
            footnote += '</ul>';
            
            this.parser_data_temp_other['footnote'] = '';
            
            return footnote;
        } else {
            return '';
        }
    }
    
    do_part_footnote() {
        let render_main = this;
        
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        let parser_data_temp_other = this.parser_data_temp_other;
        
        let footnote_n = 0;
        let footnote_name_all = {};
        this.doc_data = this.doc_data.replace(/(?:\[\^((?:(?!\[\^| ).)*)(?: ((?:(?!\[\^|\]).)+))?\]|(\[fnote\(\)]))/g, function(match, x1, x2, x3) {
            if(x3 === undefined) {
                if(x1 === '' && x2 === undefined) {
                    return '<footnoteBlink>';
                } else {
                    footnote_n += 1;

                    let footnote_name;
                    let footnote_id = String(footnote_n);
                    if(x1 === '') {
                        footnote_name = footnote_id;
                    } else {
                        footnote_name = x1;
                    }
                    
                    let footnote_content;
                    if(x2 === undefined) {
                        if(footnote_name_all[footnote_name]) {
                            footnote_content = footnote_name_all[footnote_name];
                        } else {
                            footnote_content = '';
                            footnote_name_all[footnote_name] = footnote_content;
                        }
                    } else {
                        footnote_content = x2;
                        footnote_name_all[footnote_name] = footnote_content;
                    }
                    
                    let footnote_list = '';
                    footnote_list += '<li>';
                    footnote_list += '<a id="opennamuFnGo' + footnote_id + '" href="#opennamuFnIn' + footnote_id + '">'
                    footnote_list += '(' + footnote_name + ')';
                    footnote_list += '</a>';
                    footnote_list += ' ';
                    footnote_list += footnote_content;
                    footnote_list += '</li>';
                    
                    parser_data_temp_other['footnote'] += footnote_list;

                    parser_count += 1;
                    let parser_count_str = String(parser_count);
                    
                    let footnote_data = '';
                    footnote_data += '<sup>'
                    footnote_data += '<a id="opennamuFnIn' + footnote_id + '" href="#opennamuFnGo' + footnote_id + '">';
                    footnote_data += '(' + footnote_name + ')';
                    footnote_data += '</a>';
                    footnote_data += '</sup>';

                    parser_data_temp['render' + parser_count_str + 'Span'] = footnote_data;
                    parser_data_temp['/render' + parser_count_str + 'Span'] = '';

                    return '<render' + parser_count_str + 'Span></render' + parser_count_str + 'Span>';
                }
            } else {
                return render_main.do_part_footnote_list();
            }
        });
        
        this.doc_data = this.doc_data.replace(/<footnoteBlink>/g, '[^]');
        this.doc_data += this.do_part_footnote_list();
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
        this.parser_data_temp_other = parser_data_temp_other;
    }
    
    do_part_macro() {
        let render_main = this;
        let render_part_id_add = this.render_part_id_add;
        
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        let parser_data_temp_other = this.parser_data_temp_other;
        
        this.doc_data = this.doc_data.replace(/\[([^\[\(<>\n]+)\(((?:(?!\(|\)\]|<|>|\n).)*)\)\]/g, function(match, x1, x2) {
            if(x1 === 'anchor') {
                parser_count += 1;
                let parser_count_str = String(parser_count);
                
                let anchor_data = render_main.do_func_parser_to_text(x2, 'nowiki');
                anchor_data = render_main.do_func_xss_encode(anchor_data);

                parser_data_temp['render' + parser_count_str + 'Span'] = '<span id="' + anchor_data + '">';
                parser_data_temp['/render' + parser_count_str + 'Span'] = '</span>';

                return '<render' + parser_count_str + 'Span></render' + parser_count_str + 'Span>' 
            } else if(x1 === 'category') {
                let category = x2.split(',');

                let category_data = render_main.do_func_parser_to_text(category[0], 'nowiki');
                
                let link_main = 'category:' + category_data;
                let link_title = render_main.do_func_xss_encode(link_main);
                
                link_main = render_main.do_func_xss_decode(link_main);
                link_main = opennamu_do_url_encode(link_main);
                
                parser_data_temp_other['category'] += '<a class="' + render_part_id_add + 'opennamu_link" title="' + link_title + '" href="/w/' + link_main + '">';
                parser_data_temp_other['category'] += category_data;
                parser_data_temp_other['category'] += '</a>';
                parser_data_temp_other['category'] += ' | ';
                
                return '';
            } else if(x1 === 'toc') {
                if(parser_data_temp_other['toc'] !== '') {
                    return '' +
                        '<div id="opennamu_TOC" class="opennamu_TOC">' + 
                            '<span class="opennamu_TOC_title">' +
                                'TOC' +
                            '</span>' + 
                            '<br>' + 
                            '<br>' + 
                            parser_data_temp_other['toc'] + 
                        '</div>' +
                    '';
                } else {
                    return '';
                }
            }
            else {
                return match;
            }
        });
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
        this.parser_data_temp_other = parser_data_temp_other;
    }
    
    do_part_nowiki() {
        let render_main = this;
        
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;

        this.doc_data = this.doc_data.replace(/\\\n/g, '\n');
        this.doc_data = this.doc_data.replace(/\\(&#x27;|&quot;|&lt;|&gt;|&amp;|.)/g, function(match, x1) {
            let nowiki_data = x1;
            
            parser_count += 1;
            let parser_count_str = String(parser_count);

            parser_data_temp['nowiki' + parser_count_str + 'Span'] = match;
            parser_data_temp['/nowiki' + parser_count_str + 'Span'] = '';
            
            parser_data_temp['nowiki' + parser_count_str + 'SpanEnd'] = '<code>' + nowiki_data;
            parser_data_temp['/nowiki' + parser_count_str + 'SpanEnd'] = '</code>';
            
            return '<nowiki' + parser_count_str + 'Span>' + '</nowiki' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/```((?:(?:(?!```).)|\n)+)```/g, function(match, x1) {
            let nowiki_data = render_main.do_func_parser_to_text(x1, 'nowiki');
            
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            if(nowiki_data.match(/\n/)) {
                parser_data_temp['nowiki' + parser_count_str + 'Span'] = match;
                parser_data_temp['/nowiki' + parser_count_str + 'Span'] = '';
                
                parser_data_temp['nowiki' + parser_count_str + 'SpanEnd'] = '<pre>' + nowiki_data;
                parser_data_temp['/nowiki' + parser_count_str + 'SpanEnd'] = '</pre>';
            } else {
                parser_data_temp['nowiki' + parser_count_str + 'Span'] = match;
                parser_data_temp['/nowiki' + parser_count_str + 'Span'] = '';
                
                parser_data_temp['nowiki' + parser_count_str + 'SpanEnd'] = '<code>' + nowiki_data;
                parser_data_temp['/nowiki' + parser_count_str + 'SpanEnd'] = '</code>';
            }
            
            return '<nowiki' + parser_count_str + 'Span>' + '</nowiki' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/`([^`\n]+)`/g, function(match, x1) {
            let nowiki_data = render_main.do_func_parser_to_text(x1, 'nowiki');
            
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['nowiki' + parser_count_str + 'Span'] = match;
            parser_data_temp['/nowiki' + parser_count_str + 'Span'] = '';

            parser_data_temp['nowiki' + parser_count_str + 'SpanEnd'] = '<code>' + nowiki_data;
            parser_data_temp['/nowiki' + parser_count_str + 'SpanEnd'] = '</code>';
            
            return '<nowiki' + parser_count_str + 'Span>' + '</nowiki' + parser_count_str + 'Span>';
        });
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
    }
    
    do_part_horizon() {        
        let horizone_regex = /\n((?:\*|\* |-|- ){3,})\n/;
        while(this.doc_data.match(horizone_regex)) {
            this.doc_data = this.doc_data.replace(horizone_regex, function(match, x1, x2) {
                return '\n<brEnd><hr><brStart>\n';
            });
        }
    }
    
    do_part_blockquote() {
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        let blockquote_regex = /((?:\n&gt;(?:[^\n]+))+)/;
        while(this.doc_data.match(blockquote_regex)) {
            this.doc_data = this.doc_data.replace(blockquote_regex, function(match, x1) {
                let blockquote_data = '<brStart>' + x1;
                blockquote_data = blockquote_data.replace(/\n&gt;	?/g, '\n');

                parser_count += 1;
                let parser_count_str = String(parser_count);

                parser_data_temp['render' + parser_count_str + 'Span'] = '<blockquote>';
                parser_data_temp['/render' + parser_count_str + 'Span'] = '</blockquote>';

                return '\n<brEnd><render' + parser_count_str + 'Span>' + blockquote_data + '</render' + parser_count_str + 'Span>';
            });
        }
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
    }
    
    do_part_list() {
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        /*
        // 여기 파트 개선 필요
        this.doc_data = this.doc_data.replace(/((?:\n(?: )*\* ?(?:[^\n]+))+)/g, function(match, x1) {
            let list_data = x1;
            let list_depth = -1;
            
            list_data = list_data.replace(/\n( )*\* ?([^\n]+)/g, function(match, x1, x2) {
                return '';
            });

            parser_count += 1;
            let parser_count_str = String(parser_count);

            parser_data_temp['render' + parser_count_str + 'Span'] = '<ul>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</ul>';

            return '\n<brEnd><render' + parser_count_str + 'Span>' + list_data + '</render' + parser_count_str + 'Span>';
        });
        */
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
    }
    
    do_part_final() {
        this.doc_data = this.doc_data.replace(/<brStart>\n?/g, '');
        this.doc_data = this.doc_data.replace(/\n?<brEnd>/g, '');
        
        this.doc_data = this.doc_data.replace(/<uBar>/g, '_');
        
        this.doc_data = this.doc_data.replace(/\n/g, '<br>');
        
        if(this.parser_data_temp_other['category'] !== '') {
            let category = this.parser_data_temp_other['category'];
            category = category.replace(/ \| $/, '');
            
            this.doc_data += '<hr class="main_hr">';
            this.doc_data += '<div id="cate_all"><div id="cate">';
            this.doc_data += 'Category : ';
            this.doc_data += category;
            this.doc_data += '</div></div>';
        }
    }
    
    // Main Part
    do_main() {
        this.do_part_nowiki();
        this.do_part_list();
        this.do_part_blockquote();
        this.do_part_heading();
        this.do_part_horizon();
        this.do_part_footnote();
        this.do_part_macro();
        this.do_part_image();
        this.do_part_link();
        this.do_part_text();
        this.do_part_final();
        
        this.doc_data = this.do_func_parser_to_text(this.doc_data);
        this.doc_data = this.do_func_parser_to_text(this.doc_data, 'nowikiEnd');
        
        document.getElementById(this.render_part_id_add + this.render_part_id_after).innerHTML = this.doc_data;
        
        document.getElementById(this.render_part_id_add + this.render_part_id).style.display = "none";
        document.getElementById(this.render_part_id_add + this.render_part_id_after).style.display = "";
        
        for(let x1 in this.parser_data_js) {
            eval(this.parser_data_js[x1]);
        }
        
        new opennamu_render_wiki(
            render_part_id_add = this.render_part_id_add
        ).do_main();
    }
}