// 좀 더 개선 필요
// use strict 적용 필요 (eval 동작에 문제 있음)
function opennamu_do_thread_make(topic_num, type_do = 'top', some = '', where = 'top_topic') {
    let url = '';
    if(type_do === 'top') {
        url = "/api/thread/" + topic_num + "/top";
    } else if(type_do === 'main') {
        url = "/api/thread/" + topic_num;
    } else {
        url = "/api/thread/" + topic_num + some;
    }

    let xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            let data_t = JSON.parse(this.responseText);
            
            let start = 0;
            let key_v = '/normal/1';
            let admin = '';
            let ip_first = '';
            
            let data_all = '';
            let data_all_js = '';
            
            let count = 0;
            for(let key in data_t) {
                let data_a = '';
                
                if(start === 0) {
                    admin = data_t['data_main']['admin'];
                    ip_first = data_t['data_main']['ip_first'];
                
                    start = 1;
                }
                
                if(key === 'data_main') {
                    continue;
                }
                
                key_v = '/normal/' + String(Number(key) + 1);
                
                let color_b = '';
                let color_t = '';
                
                let ip = data_t[key]['ip_pas'];
                let ip_o = data_t[key]['ip'];
                let blind = data_t[key]['blind'];
                let data_i_pas = data_t[key]['data_pas'][0];
                
                if(data_i_pas === '') {
                    data_i_pas = '<br>';
                } else {
                    data_i_pas = data_i_pas.replace(
                        /&lt;topic_a&gt;((?:(?!&lt;\/topic_a&gt;).)+)&lt;\/topic_a&gt;/g,
                        '<a href="$1">$1</a>'
                    );
                    data_i_pas = data_i_pas.replace(
                        /&lt;topic_call&gt;@((?:(?!&lt;\/topic_call&gt;).)+)&lt;\/topic_call&gt;/g,
                        '<a href="/w/user:$1">@$1</a>', 
                    );
                }
                
                if(blind === 'O') {
                    color_b = 'opennamu_comment_blind';
                } else {
                    color_b = 'opennamu_comment_blind_not';
                }
                
                if(blind === 'O') {
                    ip += ' <a href="/admin_log?search=blind%20(code%20' + topic_num + '#' + key + '">(B)</a>';
                    
                    if(admin === '1') {
                        ip += ' <a href="javascript:opennamu_do_open_comment(\'' + key + '\');">(O)</a>';
                    }
                }
                
                if(admin === '1' || blind !== 'O') {
                    ip += ' <a href="/thread/' + topic_num + '/comment/' + key + '/tool">(T)</a>';
                }
                
                if(type_do === 'top') {
                    color_t = 'opennamu_comment_color_red';
                } else if(blind === '1') {
                    color_t = 'opennamu_comment_color_blue';
                } else if(blind === 'O') {
                    color_t = 'opennamu_comment_color_gray';
                } else if(ip_o === ip_first) {
                    color_t = 'opennamu_comment_color_green';
                } else {
                    color_t = 'opennamu_comment_color_default';
                }
                
                data_a += '' + 
                    '<table class="opennamu_comment">' + 
                        '<tr>' + 
                            '<td class="' + color_t + '">' + 
                                '<a href="javascript:void(0);" id="' + key + '">#' + key + '</a> ' + 
                                ip + 
                                '<span style="float: right;">' + data_t[key]['date'] + '</span>' + 
                            '</td>' + 
                        '</tr>' + 
                        '<tr>' + 
                            '<td class="' + color_b + '" id="opennamu_comment_data_' + key + '">' + 
                                '<div class="opennamu_comment_scroll">' + data_i_pas + '</div>' + 
                            '</td>' + 
                        '</tr>' +
                    '</table>' + 
                    '<hr class="main_hr">' + 
                ''

                document.getElementById(where).innerHTML += data_a;
                
                count += 1;
                data_all_js += data_t[key]['data_pas'][1] + '\n';
                
                if(count > 100) {
                    eval(data_all_js);
                    
                    count = 0;
                    data_all_js = '';
                }
            }
            
            eval(data_all_js);
            
            opennamu_do_ip_parser();
            
            if(type_do === 'top') {
                opennamu_do_thread_make(topic_num, 'main', '', 'main_topic');
            } else if(type_do === 'main') {
                let data_url_v = window.location.hash.replace(/^#/, '');
                if(data_url_v !== '') {
                    if(document.getElementById(data_url_v)) {
                        document.getElementById(data_url_v).focus();
                    }
                }
                
                opennamu_do_thread_make(topic_num, 're', key_v, where);
            } else if(type_do === 're') {
                setTimeout(function() {
                    if(start === 0) {
                        opennamu_do_thread_make(topic_num, 're', some, where);
                    } else {
                        opennamu_do_thread_make(topic_num, 're', key_v, where);
                    }
                }, 2000);
            }
        }
    }
}

function opennamu_do_open_comment(key) {
    let elementState = document.getElementById('opennamu_comment_data_' + key).style.display;
    if(!elementState || elementState === 'none') {
        document.getElementById('opennamu_comment_data_' + key).style.display = 'block';
    } else {
        document.getElementById('opennamu_comment_data_' + key).style.display = 'none';
    }
}