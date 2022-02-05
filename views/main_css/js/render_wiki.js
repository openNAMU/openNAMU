function get_link_state(data) {
    let data_exter_link = '0';
    if(document.cookie.match(main_css_regex_data('main_css_exter_link'))) {
        data_exter_link = document.cookie.match(main_css_regex_data('main_css_exter_link'))[1];
    }
    
    var link_list = [];
    var link_list_2 = {}
    for(var i = 0; document.getElementsByClassName(data + 'link_finder')[i]; i++) {
        var data_class = document.getElementsByClassName(data + 'link_finder')[i];
        console.log(data_class.href)
        if(
            data_class.id !== 'out_link' && 
            data_class.id !== 'inside' && 
            data_class.id !== 'in_doc_link'
        ) {            
            link_list.push(data_class.title);
            
            if(!link_list_2[data_class.title]) {
                link_list_2[data_class.title] = [i];
            } else {
                link_list_2[data_class.title].push(i);
            }
        } else if(
            data_exter_link === '1' && 
            (
                data_class.id === 'out_link' ||
                data_class.id === 'inside'
            )
        ) {
            document.getElementsByClassName(data + 'link_finder')[i].target = '_self';
        }
    }
    
    var data_form = new FormData();
    data_form.append('title_list', JSON.stringify(link_list));
    
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/w/test/doc_tool/exist");
    xhr.send(data_form);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var data_xhr = JSON.parse(this.responseText);
            for(var key in link_list_2) {
                if(!data_xhr[key]) {
                    for(var key_2 in link_list_2[key]) {
                        document.getElementsByClassName(data + 'link_finder')[link_list_2[key][key_2]].id = "not_thing";
                    }
                }
            }
        }
    }
}

function get_heading_name() {
    let heading_name = document.getElementsByClassName('render_heading_text');
    for(let i = 0; i < heading_name.length; i++) {
        heading_name[i].id = heading_name[i].innerText.replace(/^([0-9]+\.)+ /, '').replace(/ ✎ ⊖$/, '');
    }
}

function load_image_link(data) {
    data.innerHTML = '' +
        '<img   style="' + data.getAttribute('under_style') + '" ' + 
                'alt="' + data.getAttribute('under_alt') + '" ' + 
                'src="' + data.getAttribute('under_src') + '">' +
    '';
}

function get_file_state_extermal(data, data_exter) {
    if(document.cookie.match(main_css_regex_data('main_css_image_set'))) {
        var data_image_set = document.cookie.match(main_css_regex_data('main_css_image_set'))[1];
    } else {
        var data_image_set = '0';
    }
    
    var data_class = document.getElementsByClassName(data + 'file_finder');
    for(var key in data_exter) {
        var key = data_exter[key];
        
        if(data_image_set === '1') {
            document.getElementsByClassName(data + 'file_finder')[key].innerHTML = '' +
                '<a href="' + data_class[key].getAttribute('under_src') + '" ' +
                    'title="' + data_class[key].getAttribute('under_src') + '">' + 
                    '(External image link)' + 
                '</a>' +
            '';
        } else if(data_image_set === '2') {
            document.getElementsByClassName(data + 'file_finder')[key].innerHTML = '' +
                '<a href="javascript:void(0);" ' +
                    'onclick="load_image_link(this); this.onclick = \'\';" ' + 
                    'under_style="' + data_class[key].getAttribute('under_style') + '" ' +
                    'under_alt="' + data_class[key].getAttribute('under_alt') + '" ' +
                    'under_src="' + data_class[key].getAttribute('under_src') + '" ' +
                    'title="' + data_class[key].getAttribute('under_src') + '">' + 
                    '(External image load)' + 
                '</a>' +
            '';
        } else {
            document.getElementsByClassName(data + 'file_finder')[key].innerHTML = '' +
                '<img   style="' + data_class[key].getAttribute('under_style') + '" ' + 
                        'alt="' + data_class[key].getAttribute('under_alt') + '" ' + 
                        'src="' + data_class[key].getAttribute('under_src') + '">' +
            '';
        }
    }
}

function get_file_state_intermal(data, data_inter) {
    var data_dict = {};
    var data_list = [];
    for(var key in data_inter) {
        var data_class = document.getElementsByClassName(data + 'file_finder')[key];
    
        var file_org = data_class.getAttribute('under_alt');
        var file_type = file_org.split('.');
        var file_name = file_type.slice(0, file_type.length - 1).join('.');
        file_type = file_type[file_type.length - 1];
        
        if(!data_dict[file_org]) {
            data_dict[file_org] = {};
        }
        
        data_dict[file_org]['file_name'] = file_name;
        data_dict[file_org]['file_type'] = file_type;
        
        data_list.push(file_name);
        
        if(!data_dict[file_org]['list']) {
            data_dict[file_org]['list'] = [key];
        } else {
            data_dict[file_org]['list'].push(key);
        }
    }
    
    if(document.cookie.match(main_css_regex_data('main_css_image_set'))) {
        var data_image_set = document.cookie.match(main_css_regex_data('main_css_image_set'))[1];
    } else {
        var data_image_set = '0';
    }
    
    var data_form = new FormData();
    data_form.append('title_list', JSON.stringify(data_list));

    var xhr = new XMLHttpRequest();
    xhr.open("POST", '/api/sha224/test');
    xhr.send(data_form);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var file_sha224 = JSON.parse(this.responseText);
            var data_list_2 = [];
            for(var key in data_dict) {
                data_dict[key]['file_sha224'] = file_sha224[data_dict[key]['file_name']]

                data_list_2.push(data_dict[key]['file_sha224'] + '.' + data_dict[key]['file_type'])
            }
            
            var data_form_2 = new FormData();
            data_form_2.append('title_list', JSON.stringify(data_list_2));
            
            var xhr_2 = new XMLHttpRequest();
            xhr_2.open("POST", '/api/image/test');
            xhr_2.send(data_form_2);

            xhr_2.onreadystatechange = function() {
                if(this.readyState === 4 && this.status === 200) {
                    var file_data = JSON.parse(this.responseText);
                    var data_class = document.getElementsByClassName(data + 'file_finder');
                    for(var key_3 in data_dict) {
                        if(!file_data[data_dict[key_3]['file_sha224'] + '.' + data_dict[key_3]['file_type']]) {
                            for(var key_4 in data_dict[key_3]['list']) {
                                var key_4 = data_dict[key_3]['list'][key_4];
                                document.getElementsByClassName(data + 'file_finder')[key_4].innerHTML = '' +
                                    '<a href="' + data_class[key_4].getAttribute('under_href') + '" ' + 
                                        'id="not_thing">' +
                                        '(' + data_class[key_4].getAttribute('under_alt') + ')' +
                                    '</a>' +
                                '';
                            }
                        } else {
                            if(data_image_set === '1') {
                                for(var key_4 in data_dict[key_3]['list']) {
                                    var key_4 = data_dict[key_3]['list'][key_4];
                                    document.getElementsByClassName(data + 'file_finder')[key_4].innerHTML = '' +
                                        '<a href="/image/' + data_dict[key_3]['file_sha224'] + '.' + data_dict[key_3]['file_type'] + '">' +
                                            '(' + data_class[key_4].getAttribute('under_alt') + ')' +
                                        '</a>' +
                                    '';
                                }
                            } else if(data_image_set === '2') {
                                for(var key_4 in data_dict[key_3]['list']) {
                                    var key_4 = data_dict[key_3]['list'][key_4];
                                    document.getElementsByClassName(data + 'file_finder')[key_4].innerHTML = '' +
                                        '<a href="javascript:void(0);" ' +
                                            'onclick="load_image_link(this); this.onclick = \'\';" ' + 
                                            'under_style="' + data_class[key_4].getAttribute('under_style') + '" ' +
                                            'under_alt="' + data_class[key_4].getAttribute('under_alt') + '" ' +
                                            'under_src="/image/' + data_dict[key_3]['file_sha224'] + '.' + data_dict[key_3]['file_type'] + '">' + 
                                            '(' + data_class[key_4].getAttribute('under_alt') + ' load)' +
                                        '</a>' +
                                    '';
                                }
                            } else {
                                for(var key_4 in data_dict[key_3]['list']) {
                                    var key_4 = data_dict[key_3]['list'][key_4];
                                    document.getElementsByClassName(data + 'file_finder')[key_4].innerHTML = '' +
                                        '<img   style="' + data_class[key_4].getAttribute('under_style') + '" ' + 
                                                'alt="' + data_class[key_4].getAttribute('under_alt') + '" ' + 
                                                'src="/image/' + data_dict[key_3]['file_sha224'] + '.' + data_dict[key_3]['file_type'] + '">' +
                                        '' +
                                    '';
                                }
                            }
                        }   
                    }
                }
            }
        }
    }
}

function get_file_state(data, i = 0) {
    var data_exter = [];
    var data_inter = [];
    for(var i = 0; document.getElementsByClassName(data + 'file_finder')[i]; i++) {
        var data_class = document.getElementsByClassName(data + 'file_finder')[i];
        if(data_class.getAttribute('under_href') === 'out_link') {
            data_exter.push(i);
        } else {
            data_inter.push(i);
        }
    }
    
    get_file_state_extermal(data, data_exter);
    get_file_state_intermal(data, data_inter);
}

function load_include(name_doc, name_ob, data_include, name_org = '') {
    var data_form = new FormData();
    data_form.append('include_list', JSON.stringify(data_include));
    data_form.append('name_include', name_ob);
    data_form.append('name_org', name_org);
    
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/w/" + encodeURI(name_doc) + "/doc_tool/include");
    xhr.send(data_form);

    document.getElementsByClassName(name_ob)[0].href = "/w/" + do_url_change(name_doc);
    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            if(this.responseText === "{}\n") {
                document.getElementById(name_ob).innerHTML = "";
                document.getElementsByClassName(name_ob)[0].id = "not_thing";
            } else {
                var data_load = JSON.parse(this.responseText);
                document.getElementById(name_ob).innerHTML = data_load['data'];
                eval(data_load['js_data']);
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
        if(fol_data != '⊖') {
            element.innerHTML = '⊖';
        } else {
            element.innerHTML = '⊕';
        }
    }
}

function do_open_foot(front_data, name, num = 0) {    
    if(
        document.cookie.match(main_css_regex_data('main_css_footnote_set')) &&
        document.cookie.match(main_css_regex_data('main_css_footnote_set'))[1] === '1'
    ) {
        if(num === 1) {
            document.getElementById(front_data + 'r' + name).focus();
        } else {
            var get_data = document.getElementById(front_data + name).innerHTML;
            var org_data = document.getElementById(front_data + 'd' + name).innerHTML;
            if(org_data === '') {
                document.getElementById(front_data + 'd' + name).innerHTML = '' +
                    '<a href="#' + front_data + 'c' + name + '">(Go)</a> ' + get_data +
                '';
                document.getElementById(front_data + 'd' + name).className = 'spead_footnote';
            } else {
                document.getElementById(front_data + 'd' + name).innerHTML = '';
                document.getElementById(front_data + 'd' + name).className = '';
            }
        }
    } else {
        document.getElementById(front_data + 'r' + name).style.color = 'red';
        document.getElementById(front_data + 'c' + name).style.color = (num === 1 ? 'inherit' : 'red');
        document.getElementById(front_data + (num === 1 ? 'r' : 'c') + name).focus();
    }
}