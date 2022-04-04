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
        
        this.render_part_id_add = render_part_id_add;
        this.render_part_id_after = render_part_id_after;
    }
    
    // Func Part
    do_func_parser_to_text(data, parser_type = 'parser') {
        let parser_data_temp = this.parser_data_temp;
        let parser_match;
        if(parser_type === 'nowiki') {
            parser_match = /<(\/?nowiki[0-9]+Span)>/;
        } else {
            parser_match = /<(\/?render[0-9]+Span)>/;
        }
        
        while(data.match(parser_match)) {
            data = data.replace(parser_match, function(match, x1) {
                return parser_data_temp[x1];
            });
        }
        
        return data;
    }
    
    do_func_xss_encode(data) {
        data = data.replace(/"/g, '&quot;');
        data = data.replace(/</g, '&lt;');
        data = data.replace(/</g, '&gt;');
        
        return data;
    }
    
    do_func_xss_decode(data) {
        data = data.replace(/&#x27;/g, '\'');
        data = data.replace(/&quot;/g, '"');
        data = data.replace(/&lt;/g, '<');
        data = data.replace(/&gt;/g, '<');
        
        return data;
    }
    
    do_func_url_encode(data) {
        return encodeURIComponent(data);
    }
    
    // Render Part
    do_part_text() {
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        this.doc_data = this.doc_data.replace(/~~((?:(?!~~).)+)~~/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<s>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</s>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/\*\*((?:(?!\*\*).)+)\*\*/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<b>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</b>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/__((?:(?!__).)+)__/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<b>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</b>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/\*([^*]+)\*/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<i>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</i>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/_([^_]+)_/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<i>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</i>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/&lt;ins&gt;((?:(?!&lt;ins&gt;|&lt;\/ins&gt;).)+)&lt;\/ins&gt;/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<u>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</u>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/&lt;sub&gt;((?:(?!&lt;sub&gt;|&lt;\/sub&gt;).)+)&lt;\/sub&gt;/g, function(match, x1) {
            parser_count += 1;
            let parser_count_str = String(parser_count);
            
            parser_data_temp['render' + parser_count_str + 'Span'] = '<sub>';
            parser_data_temp['/render' + parser_count_str + 'Span'] = '</sub>';
            
            return '<render' + parser_count_str + 'Span>' + x1 + '</render' + parser_count_str + 'Span>';
        });
        
        this.doc_data = this.doc_data.replace(/&lt;sup&gt;((?:(?!&lt;sup&gt;|&lt;\/sup&gt;).)+)&lt;\/sup&gt;/g, function(match, x1) {
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
        let toc_list = [0, 0, 0, 0, 0, 0];
        let toc_regex = /\n(#{1,6})([^\n]+)\n/;
        while(this.doc_data.match(toc_regex)) {
            this.doc_data = this.doc_data.replace(toc_regex, function(match, x1, x2) {
                let toc_level = x1.length - 1;
                let toc_level_str = String(toc_level + 1);

                toc_list[toc_level] += 1;
                for(let for_a = toc_level + 1; for_a < 6; for_a++) {
                    toc_list[for_a] = 0;
                }

                let toc_list_str = '';
                for(let for_a = 0; for_a < 6; for_a++) {
                    if(toc_list[for_a] !== 0) {
                        toc_list_str += String(toc_list[for_a]) + '.'
                    }
                }

                return '\n<brEnd><h' + toc_level_str + '>' + toc_list_str + x2 + '</h' + toc_level_str + '><brStart>\n';
            });
        }
        
        this.parser_data_temp_other['toc'] = toc_data;
    }
    
    do_part_image() {
        let render_main = this;
        let render_part_id_add = this.render_part_id_add;
        
        let parser_count = this.parser_count['parser'];
        let parser_data_temp = this.parser_data_temp;
        
        this.doc_data = this.doc_data.replace(/!\[([^\[\]]*)\]\(([^\(\)]*)\)/g, function(match, x1, x2) {
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
        
        this.doc_data = this.doc_data.replace(/\[([^\[\]]*)\]\(([^\(\)]*)\)/g, function(match, x1, x2) {
            if(x1 === '' && x2 === '') {
                return '<linkBlink>';
            } else {
                if(x2 !== '' && x2.match(/^https?:\/\//)) {
                    parser_count += 1;
                    let parser_count_str = String(parser_count);
                    
                    let link_main = render_main.do_func_xss_encode(x2);;
                    let link_sub;
                    if(x1 === '') {
                        link_sub = x2;
                    } else {
                        link_sub = x1;
                    }
                    
                    parser_data_temp['render' + parser_count_str + 'Span'] = '<a href="' + link_main + '">';
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
                    
                    link_title = render_main.do_func_xss_encode(link_main);

                    link_main = render_main.do_func_xss_decode(link_main);
                    link_main = render_main.do_func_url_encode(link_main);

                    parser_data_temp['render' + parser_count_str + 'Span'] = '<a class="' + render_part_id_add + 'opennamuLink" title="' + link_title + '" href="/w/' + link_main + '">';
                    parser_data_temp['/render' + parser_count_str + 'Span'] = '</a>';

                    return '<render' + parser_count_str + 'Span>' + link_sub + '</render' + parser_count_str + 'Span>';
                }
            }
        });
        
        this.doc_data = this.doc_data.replace(/<linkBlink>/g, '[]()');
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
    }
    
    do_part_footnote_list() {
        if(this.parser_data_temp_other['footnote'] !== '') {
            let footnote = '';
            footnote += '<ul id="footnote<underBar>data">';
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
        
        this.doc_data = this.doc_data.replace(/\[([^\[\(<>]+)\(((?:(?!\(|\)\]|<|>).)+)\)\]/g, function(match, x1, x2) {
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
                link_main = render_main.do_func_url_encode(link_main);
                
                parser_data_temp_other['category'] += '<a class="' + render_part_id_add + 'opennamuLink" title="' + link_title + '" href="/w/' + link_main + '">';
                parser_data_temp_other['category'] += category_data;
                parser_data_temp_other['category'] += '</a>';
                parser_data_temp_other['category'] += ' | ';
                
                return '';
            } else if(x1 === 'toc') {                
                return parser_data_temp_other['toc'];
            }
            else {
                return '';
            }
        });
        
        this.parser_count['parser'] = parser_count;
        this.parser_data_temp = parser_data_temp;
        this.parser_data_temp_other = parser_data_temp_other;
    }
    
    do_part_final() {
        this.doc_data = this.doc_data.replace(/<brStart>\n?/g, '');
        this.doc_data = this.doc_data.replace(/\n?<brEnd>/g, '');
        
        this.doc_data = this.doc_data.replace(/<underBar>/g, '_');
        
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
        this.do_part_heading();
        this.do_part_footnote();
        this.do_part_macro();
        this.do_part_image();
        this.do_part_link();
        this.do_part_text();
        this.do_part_final();
        
        this.doc_data = this.do_func_parser_to_text(this.doc_data);
        document.getElementById(this.render_part_id_add + this.render_part_id_after).innerHTML = this.doc_data;
        for(let x1 in this.parser_data_js) {
            eval(this.parser_data_js[x1]);
        }
        
        new opennamu_render_wiki(
            render_part_id_add = this.render_part_id_add
        ).do_main();
    }
}