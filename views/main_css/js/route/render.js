"use strict";

function opennamu_heading_folding(data, element = '') {
    let fol = document.getElementById(data);
    if(fol.style.display === '' || fol.style.display === 'inline-block' || fol.style.display === 'block') {
        document.getElementById(data).style.display = 'none';
    } else {
        document.getElementById(data).style.display = 'block';
    }
    
    if(element !== '') {
        console.log(element.innerHTML);
        if(element.innerHTML !== '⊖') {
            element.innerHTML = '⊖';
        } else {
            element.innerHTML = '⊕';
        }
    }
}

function opennamu_render_html(name = '') {
    console.log(name);
    if(document.getElementById(name)) {
        let data = document.getElementById(name).innerHTML;

        let src_list = ['www.youtube.com', 'www.google.com', 'play-tv.kakao.com'];
        let t_data = [
            'b', 'i', 's', 'del', 'strong', 'bold', 'em', 'sub', 'sup', 
            'div', 'span', 
            'a',
            'iframe'
        ];
        for(let key in t_data) {
            let patt = new RegExp(
                '&lt;' + t_data[key] + '( (?:(?:(?!&gt;).)+))?&gt;((?:(?!&lt;\/' + t_data[key] + '&gt;).)*)&lt;\/' + t_data[key] + '&gt;',
                'ig'
            );
            
            data = data.replace(patt, function(full, in_data, in_data_2) {
                if(['b', 'i', 's', 'del', 'strong', 'bold', 'em', 'sub', 'sup'].includes(t_data[key])) {
                    return '<' + t_data[key] + '>' + in_data_2 + '</' + t_data[key] + '>'
                } else if(t_data[key] === 'div' || t_data[key] === 'span') {
                    let style_data = in_data.match(/ style=['"]([^'"]*)['"]/);
                    if(style_data) {
                        style_data = style_data[1].replace(/position/ig, '');
                    } else {
                        style_data = '';
                    }

                    return '<' + t_data[key] + ' style="' + style_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                } else if(t_data[key] === 'a') {
                    let link_data = in_data.match(/ href=['"]([^'"]*)['"]/);
                    if(link_data) {
                        link_data = link_data[1].replace(/^javascript:/ig, '');
                    } else {
                        link_data = '';
                    }

                    return '<' + t_data[key] + ' id="out_link" href="' + link_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                } else if(t_data[key] === 'iframe') {
                    let src_data = in_data.match(/ src=['"]([^'"]*)['"]/);
                    if(src_data) {
                        src_data = src_data[1];

                        let src_check = src_data.match(/^http(?:s)?:\/\/([^/]+)/);
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

                    let width_data = in_data.match(/ width=['"]([^'"]*)['"]/);
                    if(width_data) {
                        width_data = width_data[1];
                    } else {
                        width_data = '';
                    }

                    let height_data = in_data.match(/ height=['"]([^'"]*)['"]/);
                    if(height_data) {
                        height_data = height_data[1];
                    } else {
                        height_data = '';
                    }

                    return '<' + t_data[key] + ' src="' + src_data + '" width="' + width_data + '" height="' + height_data + '" allowfullscreen frameborder="0">' + in_data_2 + '</' + t_data[key] + '>';
                } else {
                    let src_data = in_data.match(/ src=['"]([^'"]*)['"]/);
                    if(src_data) {
                        src_data = src_data[1];
                    } else {
                        src_data = '';
                    }

                    let width_data = in_data.match(/ width=['"]([^'"]*)['"]/);
                    if(width_data) {
                        width_data = width_data[1];
                    } else {
                        width_data = '';
                    }

                    let height_data = in_data.match(/ height=['"]([^'"]*)['"]/);
                    if(height_data) {
                        height_data = height_data[1];
                    } else {
                        height_data = '';
                    }

                    return '<' + t_data[key] + ' controls src="' + src_data + '" width="' + width_data + '" height="' + height_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                }
            });
        }

        document.getElementById(name).innerHTML = data;
    }
}