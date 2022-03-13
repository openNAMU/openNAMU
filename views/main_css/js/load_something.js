function load_user_info(name) {
    var url = "/api/user_info/" + encodeURI(name) + "?render=1";
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            document.getElementById('get_user_info').innerHTML += JSON.parse(this.responseText)['data'];
            opennamu_do_ip_parser();
        }
    }
}

function load_ver() {
    var url = "/api/version";
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send();

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            let get_data = JSON.parse(this.responseText);
            document.getElementById('ver_send_2').innerHTML = get_data['version'];
            
            let url_2 = 'https://raw.githubusercontent.com/openNAMU/openNAMU/' + get_data['build'] + '/version.json';
            var xhr_2 = new XMLHttpRequest();
            xhr_2.open("GET", url_2, true);
            xhr_2.send();
            
            xhr_2.onreadystatechange = function() {
                if(this.readyState === 4 && this.status === 200) {
                    document.getElementById('ver_send').innerHTML += JSON.parse(this.responseText)['beta']['r_ver'];
                    document.getElementById('ver_send').style.display = "list-item";
                }
            }
        }
    }
}

function do_skin_ver_check() {
    var url = "/api/skin_info?all=true";
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send();

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var json_data = JSON.parse(this.responseText);
            for(var key in json_data) {                
                document.getElementById('ver_send_3').innerHTML += '<li>' +
                    json_data[key]['name'] + ' : ' + json_data[key]['skin_ver'] +
                    (json_data[key]['lastest_version'] ? ' (' + json_data[key]['lastest_version']['skin_ver'] + ')' : '') +
                '</li>'
            }
        }
    }
}

function do_twofa_check(init = 0) {
    let twofa_option = document.getElementById('twofa_check_input');
    let twofa_option_num = twofa_option.options.selectedIndex;
    let twofa_select_data = twofa_option.options[twofa_option_num].value;
    
    if(twofa_select_data === 'on') {
        document.getElementById('fa_plus_content').style.display = "block";
    } else {
        document.getElementById('fa_plus_content').style.display = "none";
    }
}

function send_render(i = 0) {
    let get_class = document.getElementsByClassName('send_content')[i];
    if(get_class) {
        send_render(i + 1);
        
        let data = get_class.innerHTML;
        if(data === '&lt;br&gt;') {
            document.getElementsByClassName('send_content')[i].innerHTML = '<br>';
        } else {
            data = data.replace(/javascript:/i, '');
            data = data.replace(/&lt;a(?:(?:(?!&gt;).)*)&gt;((?:(?!&lt;\/a&gt;).)+)&lt;\/a&gt;/g, function(x, x_1) {
                x_1_org = x_1.replace('&lt;', '<').replace('&gt;', '>');
                return '<a href="/w/' + encodeURIComponent(x_1_org) + '">' + x_1 + '</a>';
            });
            
            document.getElementsByClassName('send_content')[i].innerHTML = data;
        }
    }
}

function simple_render(name_ele) {
    var skin_set_data = document.getElementById(name_ele).innerHTML;
    
    // 목차 구현
    var toc_all_data = '<div id="toc"><span id="toc_title">TOC</span><br>';
    var split_toc;
    var toc_data;
    i = 1;
    while(1) {
        toc_data = skin_set_data.match(/<h[1-6]>([^<>]+)<\/h[1-6]>/);
        if(toc_data) {
            split_toc = toc_data[1].match(/^([^ ]+)(.+)/);
            toc_all_data += '' + 
                '<br>' +
                '<span style="margin-left: ' + String(((toc_data[1].match(/\./g) || []).length - 1) * 10) + 'px;">' +
                    '<a href="#toc_' + String(i) + '">' + split_toc[1] + '</a>' + split_toc[2] +
                '</span>' +
            '';

            skin_set_data = skin_set_data.replace(
                /<(h[1-6])>([^<>]+)<\/h[1-6]>/, 
                '<$1 id="toc_' + String(i) + '"><a href="#toc">' + split_toc[1] + '</a>' + split_toc[2] + '</$1>'
            );
            i += 1;
        } else {
            break;
        }
    }
    skin_set_data = toc_all_data + '</div>' + skin_set_data;

    // 각주 구현
    var note_list = {};
    var plus_note;
    i = 1;
    while(1) {
        toc_data = skin_set_data.match(/<sup>([^<>]+)<\/sup>/);
        if(toc_data) {
            if(!note_list[toc_data[1]]) {
                note_list[toc_data[1]] = [String(i), 0];
            } else {
                note_list[toc_data[1]][1] += 1;
            }

            if(note_list[toc_data[1]][1] != 0) {
                plus_note = '_' + String(note_list[toc_data[1]][1]);
            } else {
                plus_note = '';
            }
            
            skin_set_data = skin_set_data.replace(
                /<sup>([^<>]+)<\/sup>/, 
                '<sup><a id="note_' + note_list[toc_data[1]][0] + plus_note + '" href="#note_' + note_list[toc_data[1]][0] + '_end">$1</a></sup>'
            );
            i += 1;
        } else {
            break;
        }
    }
    
    document.getElementById(name_ele).innerHTML = skin_set_data;
}