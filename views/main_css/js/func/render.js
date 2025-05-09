"use strict";

// https://css-tricks.com/how-to-animate-the-details-element/
class Accordion {
    constructor(el) {
        this.el = el;
        this.summary = el.querySelector('summary');
        this.content = el.querySelector('.opennamu_folding');
    
        this.animation = null;
        this.isClosing = false;
        this.isExpanding = false;
        this.summary.addEventListener('click', (e) => this.onClick(e));
    }
  
    onClick(e) {
        e.preventDefault();
        this.el.style.overflow = 'hidden';
        if(this.isClosing || !this.el.open) {
            this.open();
        } else if(this.isExpanding || this.el.open) {
            this.shrink();
        }
    }
  
    shrink() {
        this.isClosing = true;
        
        const startHeight = `${this.el.offsetHeight}px`;
        const endHeight = `${this.summary.offsetHeight}px`;
        
        if(this.animation) {
            this.animation.cancel();
        }
        
        this.animation = this.el.animate({
            height: [startHeight, endHeight]
        }, {
            duration: 200,
            easing: 'ease-out'
        });
        
        this.animation.onfinish = () => this.onAnimationFinish(false);
        this.animation.oncancel = () => this.isClosing = false;
    }
  
    open() {
        this.el.style.height = `${this.el.offsetHeight}px`;
        this.el.open = true;
        window.requestAnimationFrame(() => this.expand());
    }
  
    expand() {
        this.isExpanding = true;
        const startHeight = `${this.el.offsetHeight}px`;
        const endHeight = `${this.summary.offsetHeight + this.content.offsetHeight}px`;
        
        if(this.animation) {
            this.animation.cancel();
        }
        
        this.animation = this.el.animate({
            height: [startHeight, endHeight]
        }, {
            duration: 200,
            easing: 'ease-out'
        });
        this.animation.onfinish = () => this.onAnimationFinish(true);
        this.animation.oncancel = () => this.isExpanding = false;
    }
  
    onAnimationFinish(open) {
        this.el.open = open;
        this.animation = null;
        this.isClosing = false;
        this.isExpanding = false;
        this.el.style.height = this.el.style.overflow = '';
    }
}

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
                '&lt;' + t_data[key] + '( (?:(?:(?!&gt;).)+))?&gt;((?:(?!&lt;\/' + t_data[key] + '&gt;).)*)(&lt;\/' + t_data[key] + '&gt;|$)',
                'ig'
            );
            
            data = data.replace(patt, function(full, in_data, in_data_2) {
                if(['b', 'i', 's', 'del', 'strong', 'bold', 'em', 'sub', 'sup'].includes(t_data[key])) {
                    return '<' + t_data[key] + '>' + in_data_2 + '</' + t_data[key] + '>'
                } else if(t_data[key] === 'div' || t_data[key] === 'span') {
                    let style_data = in_data?.match(/ style=['"]([^'"]*)['"]/);
                    if(style_data) {
                        style_data = style_data[1];

                        const find_regex = [
                            / *box-shadow *: *(([^,;]*)(,|;|$)){10,}/i,
                            / *url\([^()]*\)/i,
                            / *linear-gradient\((([^(),]+)(,|\))){10,}/i,
                            / *position *: */i
                        ];
                        for(let regex of find_regex) {
                            if(regex.test(style_data)) {
                                style_data = "";
                                break;
                            }
                        }
                    } else {
                        style_data = '';
                    }

                    return '<' + t_data[key] + ' style="' + style_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                } else if(t_data[key] === 'a') {
                    let link_data = in_data?.match(/ href=['"]([^'"]*)['"]/);
                    if(link_data) {
                        link_data = link_data[1].replace(/^javascript:/ig, '');
                    } else {
                        link_data = '';
                    }

                    return '<' + t_data[key] + ' class="opennamu_link_out" href="' + link_data + '">' + in_data_2 + '</' + t_data[key] + '>';
                } else if(t_data[key] === 'iframe') {
                    let src_data = in_data?.match(/ src=['"]([^'"]*)['"]/);
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

                    let width_data = in_data?.match(/ width=['"]([^'"]*)['"]/);
                    if(width_data) {
                        width_data = width_data[1];
                    } else {
                        width_data = '';
                    }

                    let height_data = in_data?.match(/ height=['"]([^'"]*)['"]/);
                    if(height_data) {
                        height_data = height_data[1];
                    } else {
                        height_data = '';
                    }

                    return '<' + t_data[key] + ' src="' + src_data + '" width="' + width_data + '" height="' + height_data + '" allowfullscreen frameborder="0">' + in_data_2 + '</' + t_data[key] + '>';
                } else {
                    let src_data = in_data?.match(/ src=['"]([^'"]*)['"]/);
                    if(src_data) {
                        src_data = src_data[1];
                    } else {
                        src_data = '';
                    }

                    let width_data = in_data?.match(/ width=['"]([^'"]*)['"]/);
                    if(width_data) {
                        width_data = width_data[1];
                    } else {
                        width_data = '';
                    }

                    let height_data = in_data?.match(/ height=['"]([^'"]*)['"]/);
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

function opennamu_do_footnote_popover(set_name, load_name, sub_obj = undefined, do_type = 'open') {
    if(document.getElementById(set_name + '_load')) {
        if(do_type === 'open') {
            if(sub_obj !== undefined) {
                document.getElementById(set_name + '_load').innerHTML = document.getElementById(sub_obj).innerHTML; 
            } else {
                document.getElementById(set_name).title = '';
                document.getElementById(set_name + '_load').innerHTML = '<a href="#' + load_name + '">(Go)</a> ' + document.getElementById(load_name + '_title').innerHTML;   
            }
            document.getElementById(set_name + '_load').style.display = "inline-block";
            document.getElementById(set_name + '_load').count = 0;

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
            if(document.getElementById(set_name + '_load').count === 1) {
                document.getElementById(set_name + '_load').style.display = "none";
            } else {
                document.getElementById(set_name + '_load').count = 1;
            }
        }
    }
}

function opennamu_do_category_spread() {
    if(document.getElementsByClassName('opennamu_render_complete')) {
        document.getElementsByClassName('opennamu_render_complete')[0].innerHTML = '' +
            '<style>.opennamu_main .opennamu_category_button { display: none; } .opennamu_main .opennamu_category { white-space: pre-wrap; overflow-x: unset; text-overflow: unset; }</style>' +
        '' + document.getElementsByClassName('opennamu_render_complete')[0].innerHTML;
    }
}

function opennamu_do_toc() {
    let data = document.getElementById('opennamu_render_complete');
    let h_tag = data.querySelectorAll("h1, h2, h3, h4, h5, h6");
    let toc_count = [0, 0, 0, 0, 0, 0];
    let toc_html = '';

    for(let for_a = 0; for_a < h_tag.length; for_a++) {
        let tag = h_tag[for_a].tagName.toLowerCase();
        tag = tag.replace('h', '');
        tag = Number(tag) - 1;

        for(let for_b = tag + 1; for_b < 6; for_b++) {
            toc_count[for_b] = 0;
        }

        toc_count[tag] += 1;

        let toc_string = '';
        let add_on = false;
        for(let for_b = 5; for_b >= 0; for_b--) {
            if(add_on === false && toc_count[for_b] != 0) {
                add_on = true;
            }

            if(add_on === true) {
                toc_string = String(toc_count[for_b]) + '.' + toc_string;
            }
        }

        toc_string = toc_string.replace(/^(0\.)+/, '');

        let toc_string_sub =  toc_string.replace(/\.$/, '');
        let toc_margin = '<span style="margin-left: 10px;"></span>'.repeat(toc_string_sub.split('.').length - 1);

        toc_html += toc_margin + '<a href="#s-' + toc_string_sub + '">' + toc_string + '</a> ' + h_tag[for_a].innerHTML + '<br>';
        h_tag[for_a].innerHTML = '<a id="s-' + toc_string_sub + '" href="#toc">' + toc_string + '</a> ' + h_tag[for_a].innerHTML;
    }

    data.innerHTML = data.innerHTML.replace(/(<h[1-6]>)/, '<div class="opennamu_toc"></div>$1');
    data.innerHTML = data.innerHTML.replace(/<div class="opennamu_toc"><\/div>/g, function(match) {
        return '<div class="opennamu_TOC" id="toc"><div class="opennamu_TOC_title">TOC</div><br>' + toc_html + '</div>';
    });
}
