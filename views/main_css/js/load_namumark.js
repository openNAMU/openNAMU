function get_link_state(data, i = 0) { 
    var get_class = document.getElementsByClassName(data + 'link_finder')[i];
    if(get_class) {
        var xhr = new XMLHttpRequest();
        xhr.open("GET", get_class.href.replace('/w/', '/api/w/').replace(/#([^#]*)/, '') + "?exist=1");
        xhr.send();

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(JSON.parse(this.responseText)['exist'] !== '1') {
                    document.getElementsByClassName(data + 'link_finder')[i].id = "not_thing";
                } else {
                    document.getElementsByClassName(data + 'link_finder')[i].id = "";
                }
            }
        }

        get_link_state(data, i + 1);
    }
}

function load_image_link(data) {
    data.innerHTML = '' +
        '<img   style="' + data.getAttribute('under_style') + '" ' + 
                'alt="' + data.getAttribute('under_alt') + '" ' + 
                'src="' + data.getAttribute('under_src') + '">' +
    '';
}

function get_file_state(data, i = 0) {       
    var get_class = document.getElementsByClassName(data + 'file_finder')[i];
    if(get_class) {            
        if(get_class.getAttribute('under_href') === 'out_link') {
            if(
                document.cookie.match(main_css_regex_data('main_css_image_set')) &&
                document.cookie.match(main_css_regex_data('main_css_image_set'))[1] === '1'
            ) {
                document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                    '<a href="' + get_class.getAttribute('under_src') + '" ' +
                        'title="' + get_class.getAttribute('under_src') + '">' + 
                        '(External image link)' + 
                    '</a>' +
                '';
            } else if(
                document.cookie.match(main_css_regex_data('main_css_image_set')) &&
                document.cookie.match(main_css_regex_data('main_css_image_set'))[1] === '2'
            ) {
                document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                    '<a href="javascript:void(0);" ' +
                        'onclick="load_image_link(this); this.onclick = \'\';" ' + 
                        'under_style="' + get_class.getAttribute('under_style') + '" ' +
                        'under_alt="' + get_class.getAttribute('under_alt') + '" ' +
                        'under_src="' + get_class.getAttribute('under_src') + '" ' +
                        'title="' + get_class.getAttribute('under_src') + '">' + 
                        '(External image load)' + 
                    '</a>' +
                '';
            } else {
                document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                    '<img   style="' + get_class.getAttribute('under_style') + '" ' + 
                            'alt="' + get_class.getAttribute('under_alt') + '" ' + 
                            'src="' + get_class.getAttribute('under_src') + '">' +
                '';
            }
        } else {
            var xhr = new XMLHttpRequest();
            xhr.open("GET", get_class.getAttribute('under_src').replace('/image/', '/api/image/'));
            xhr.send();
            
            xhr.onreadystatechange = function() {
                if(this.readyState === 4 && this.status === 200) {
                    if(JSON.parse(this.responseText)['exist'] !== '1') {
                        document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                            '<a href="' + get_class.getAttribute('under_href') + '" ' + 
                                'id="not_thing">' +
                                '(' + get_class.getAttribute('under_alt') + ')' +
                            '</a>' +
                        '';
                    } else {
                        if(
                            document.cookie.match(main_css_regex_data('main_css_image_set')) &&
                            document.cookie.match(main_css_regex_data('main_css_image_set'))[1] === '1'
                        ) {
                            document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                                '<a href="' + get_class.getAttribute('under_src') + '">' +
                                    '(' + get_class.getAttribute('under_alt') + ')' +
                                '</a>' +
                            '';
                        } else if(
                            document.cookie.match(main_css_regex_data('main_css_image_set')) &&
                            document.cookie.match(main_css_regex_data('main_css_image_set'))[1] === '2'
                        ) {
                            document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                                '<a href="javascript:void(0);" ' +
                                    'onclick="load_image_link(this); this.onclick = \'\';" ' + 
                                    'under_style="' + get_class.getAttribute('under_style') + '" ' +
                                    'under_alt="' + get_class.getAttribute('under_alt') + '" ' +
                                    'under_src="' + get_class.getAttribute('under_src') + '">' + 
                                    '(' + get_class.getAttribute('under_alt') + ' load)' +
                                '</a>' +
                            '';
                        } else {
                            document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                                '<img   style="' + get_class.getAttribute('under_style') + '" ' + 
                                        'alt="' + get_class.getAttribute('under_alt') + '" ' + 
                                        'src="' + get_class.getAttribute('under_src') + '">' +
                            '';
                        }
                    }
                }
            }
        }

        get_file_state(data, i + 1);
    }
}

function load_include(title, name, p_data) {
    var change = '';
    for(key in p_data) {
        change += '@' + p_data[key][0].replace('&', '<amp>') + '@,' + p_data[key][1].replace(',', '<comma>').replace('&', '<amp>') + ','
    }
    
    var url = "/api/w/" + encodeURI(title) + "?include=" + name + "&change=" + change;

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url);
    xhr.send();

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            if(this.responseText === "{}\n") {
                document.getElementById(name).innerHTML = "";
                document.getElementsByClassName(name)[0].id = "not_thing";
            } else {
                var o_p_data = JSON.parse(this.responseText);
                document.getElementById(name).innerHTML = o_p_data['data'];
                eval(o_p_data['js_data']);
            }
        }
    }
}

function page_count() {
    var url = "/api/title_index";

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url);
    xhr.send();

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var i = 0;
            while(document.getElementsByClassName('all_page_count')[i]) {
                document.getElementsByClassName('all_page_count')[i].innerHTML = JSON.parse(this.responseText)['count'];
                
                i += 1;
            }
        }
    }
}

function do_open_folding(data, element = '') {
    var fol = document.getElementById(data);
    if(fol.style.display === '' || (fol.style.display === 'inline-block' || fol.style.display === 'block')) {
        document.getElementById(data).style.display = 'none';
    } else {
        document.getElementById(data).style.display = 'block';
    }
    
    if(element != '') {
        var fol_data = element.innerHTML;
        if(fol_data != '(-)') {
            element.innerHTML = '(-)';
        } else {
            element.innerHTML = '(+)';
        }
    }
}

function do_open_foot(name, num = 0) {
    var found_include = name.match(/^(include_(?:[0-9]+)\-)/);
    if(found_include) {
        var include_name = name.replace(/^(?:include_(?:[0-9]+)\-)/, '');
        var front_data = found_include[1];
    } else {
        var include_name = name;
        var front_data = '';
    }

    if(
        document.cookie.match(main_css_regex_data('main_css_footnote_set')) &&
        document.cookie.match(main_css_regex_data('main_css_footnote_set'))[1] === '1'
    ) {
        if(num === 1) {
            document.getElementById(front_data + 'r' + include_name).focus();
        } else {
            var get_data = document.getElementById(front_data + include_name).innerHTML;
            var org_data = document.getElementById(front_data + 'd' + include_name).innerHTML;
            if(org_data === '') {
                document.getElementById(front_data + 'd' + include_name).innerHTML = '' +
                    '<a href="#' + front_data + 'c' + include_name + '">(Go)</a> ' + get_data +
                '';
                document.getElementById(front_data + 'd' + include_name).className = 'spead_footnote';
            } else {
                document.getElementById(front_data + 'd' + include_name).innerHTML = '';
                document.getElementById(front_data + 'd' + include_name).className = '';
            }
        }
    } else {
        document.getElementById(front_data + 'r' + include_name).style.color = 'red';
        document.getElementById(front_data + 'c' + include_name).style.color = (num === 1 ? 'inherit' : 'red');
        document.getElementById(front_data + (num === 1 ? 'r' : 'c') + include_name).focus();
    }
}