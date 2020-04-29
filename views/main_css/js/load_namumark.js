function get_link_state(data, i = 0) { 
    var get_class = document.getElementsByClassName(data + 'link_finder')[i];
    if(get_class) {
        get_link_state(data, i + 1);

        var xhr = new XMLHttpRequest();
        xhr.open(
            "GET", 
            get_class.href.replace('/w/', '/api/w/').replace(/#([^#]*)/, '') + "?exist=1", 
            true
        );
        xhr.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(JSON.parse(this.responseText)['exist'] !== '1') {
                    document.getElementsByClassName(data + 'link_finder')[i].id = "not_thing";
                } else {
                    document.getElementsByClassName(data + 'link_finder')[i].id = "";
                }
            }
        }
    }
}

function get_file_state(data, i = 0) {       
    var get_class = document.getElementsByClassName(data + 'file_finder')[i];
    if(get_class) {        
        get_file_state(data, i + 1);
    
        if(get_class.getAttribute('under_href') === 'out_link') {
            if(
                document.cookie.match(main_css_regex_data('main_css_image_set')) &&
                document.cookie.match(main_css_regex_data('main_css_image_set'))[1] === '0'
            ) {
                document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                    '<img   style="' + get_class.getAttribute('under_style') + '" ' + 
                            'alt="' + get_class.getAttribute('under_alt') + '" ' + 
                            'src="' + get_class.getAttribute('under_src') + '">' +
                '';
            } else {
                document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                    '<a href="' + get_class.getAttribute('under_src') + '">(' +
                        get_class.getAttribute('under_src') +
                    ')</a>' +
                '';
            }
        } else {
            var xhr = new XMLHttpRequest();
            xhr.open(
                "GET", 
                get_class.getAttribute('under_src').replace('/image/', '/api/image/'), 
                true
            );
            xhr.send(null);
            
            xhr.onreadystatechange = function() {
                if(this.readyState === 4 && this.status === 200) {
                    if(JSON.parse(this.responseText)['exist'] !== '1') {
                        document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                            '<a href="' + get_class.getAttribute('under_href') + '" ' + 
                                'id="not_thing">' +
                                get_class.getAttribute('under_alt') +
                            '</a>' +
                        '';
                    } else {
                        if(
                            document.cookie.match(main_css_regex_data('main_css_image_set')) &&
                            document.cookie.match(main_css_regex_data('main_css_image_set'))[1] === '0'
                        ) {
                            document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                                '<img   style="' + get_class.getAttribute('under_style') + '" ' + 
                                        'alt="' + get_class.getAttribute('under_alt') + '" ' + 
                                        'src="' + get_class.getAttribute('under_src') + '">' +
                            '';
                        } else {
                            document.getElementsByClassName(data + 'file_finder')[i].innerHTML = '' +
                                '<a href="' + get_class.getAttribute('under_src') + '">(' +
                                    get_class.getAttribute('under_alt') +
                                ')</a>' +
                            '';
                        }
                    }
                }
            }
        }
    }
}

function load_include(title, name, p_data) {
    var change = '';
    for(key in p_data) {
        change += '@' + p_data[key][0].replace('&', '<amp>') + '@,' + p_data[key][1].replace(',', '<comma>').replace('&', '<amp>') + ','
    }
    
    var url = "/api/w/" + encodeURI(title) + "?include=" + name + "&change=" + change;

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

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
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var i = 0;
            while(1) {
                if(document.getElementsByClassName('all_page_count')[i]) {
                    document.getElementsByClassName('all_page_count')[i].innerHTML = JSON.parse(this.responseText)['count'];
                    i += 1;
                } else {
                    break;
                }
            }
        }
    }
}

function not_from_exist() {
    window.addEventListener('DOMContentLoaded', function() {
        if(document.getElementById('go_redirect_link')) {
            var r_link = document.getElementById('go_redirect_link').href;
            if(r_link.match(/#([^#]+)$/)) {
                var s_link = '#' + r_link.match(/#([^#]+)$/)[1];
                r_link = r_link.replace(/#([^#]+)$/, '');
            } else {
                var s_link = '';
            }

            window.location.href = r_link + '?from=' + location.pathname.replace(/^\/w\//, '') + s_link;
        }
    });
}

function do_open_folding(data, element) {
    var fol = document.getElementById(data);
    if(fol.style.display === '' || (fol.style.display === 'inline-block' || fol.style.display === 'block')) {
        document.getElementById(data).style.display = 'none';
        element.innerHTML = '(+)'
    } else {
        document.getElementById(data).style.display = 'block';
        element.innerHTML = '(-)'
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