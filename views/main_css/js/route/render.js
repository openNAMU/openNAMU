"use strict";

function opennamu_heading_folding(data, element = '') {
    let fol = document.getElementById(data);
    if(fol.style.display === '' || fol.style.display === 'inline-block' || fol.style.display === 'block') {
        document.getElementById(data).style.display = 'none';
        document.getElementById(data + '_sub').style.opacity = '0.5';
    } else {
        document.getElementById(data).style.display = 'block';
        document.getElementById(data + '_sub').style.opacity = '1';
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

function opennamu_do_render_html(name = '') {
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

                    return '<' + t_data[key] + ' class="opennamu_link_out" href="' + link_data + '">' + in_data_2 + '</' + t_data[key] + '>';
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

function opennamu_do_footnote_spread(set_name, load_name) {
    if(document.getElementById(set_name + '_load').style.display === 'none') {
        document.getElementById(set_name).title = '';
        document.getElementById(set_name + '_load').innerHTML = '<a href="#' + load_name + '">(Go)</a> ' + document.getElementById(load_name + '_title').innerHTML;
        document.getElementById(set_name + '_load').style.display = "inline-block";
    } else {
        document.getElementById(set_name + '_load').style.display = "none";
    }
}

function opennamu_do_footnote_popover(set_name, load_name) {
    if(document.getElementById(set_name + '_load').style.display === 'none') {
        document.getElementById(set_name).title = '';
        document.getElementById(set_name + '_load').innerHTML = '<a href="#' + load_name + '">(Go)</a> ' + document.getElementById(load_name + '_title').innerHTML;
        document.getElementById(set_name + '_load').style.display = "inline-block";

        let width = document.getElementById(set_name + '_load').clientWidth;
        let screen_width = window.innerWidth;
        let left = document.getElementById(set_name).getBoundingClientRect().left;
        let left_org = document.getElementById(set_name + '_load').getBoundingClientRect().left;
        let top = window.pageYOffset + document.getElementById(set_name).getBoundingClientRect().top;

        document.getElementById(set_name + '_load').style.top = String(top) + "px";
        if(screen_width - (left + width) < 50) {
            if(left > 350) {
                document.getElementById(set_name + '_load').style.left = String(left - 300) + "px";
            } else {
                document.getElementById(set_name + '_load').style.left = "0px";
            }

            left = document.getElementById(set_name + '_load').getBoundingClientRect().left;
            width = document.getElementById(set_name + '_load').clientWidth;
            if(300 > width) {
                document.getElementById(set_name + '_load').style.left = String(left + (300 - width)) + "px";
            } else {
                document.getElementById(set_name + '_load').style.marginTop = "20px";
            }
        }
    } else {
        document.getElementById(set_name + '_load').style.display = "none";
    }
}

function opennamu_do_category_spread() {
    if(document.getElementsByClassName('opennamu_render_complete')) {
        document.getElementsByClassName('opennamu_render_complete')[0].innerHTML = '' +
            '<style>.opennamu_category_button { display: none; } .opennamu_category { white-space: pre-wrap; overflow-x: unset; text-overflow: unset; }</style>' +
        '' + document.getElementsByClassName('opennamu_render_complete')[0].innerHTML
    }
}