function new_topic_load(topic_num, type_do = 'top', some = '', where = 'top_topic') {
    if(type_do === 'top') {
        var url = "/api/thread/" + topic_num + "/top";
    } else if(type_do === 'main') {
        var url = "/api/thread/" + topic_num;
    } else {
        var url = "/api/thread/" + topic_num + some;
    }

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var data_t = JSON.parse(this.responseText);
            var start = 0;
            var key_v = '/normal/1';
            
            for(var key in data_t) {
                var data_a = '';
                if(start === 0) {
                    var admin = data_t['data_main']['admin'];
                    var ip_first = data_t['data_main']['ip_first'];
                
                    start = 1;
                }
                
                if(key === 'data_main') {
                    continue;
                }
                
                key_v = '/normal/' + String(Number(key) + 1);
                
                var color_b = '';
                var color_t = '';
                
                var ip = data_t[key]['ip_pas'];
                var ip_o = data_t[key]['ip'];
                var blind = data_t[key]['blind'];
                var data_i_pas = data_t[key]['data_pas'][0];
                
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
                    color_b = 'toron_color_not';
                } else {
                    color_b = 'toron_color';
                }
                
                if(blind === 'O') {
                    ip += ' <a href="/admin_log?search=blind%20(code%20' + topic_num + '#' + key + '">(B)</a>';
                    
                    if(admin === '1') {
                        ip += ' <a href="/thread/' + topic_num + '/raw/' + key + '">(R)</a>';
                    }
                }
                
                if(admin === '1' || blind !== 'O') {
                    ip += ' <a href="/thread/' + topic_num + '/admin/' + key + '">(T)</a>';
                }
                
                if(type_do === 'top') {
                    color_t = 'toron_color_red';
                } else if(blind === '1') {
                    color_t = 'toron_color_blue';
                } else if(ip_o === ip_first) {
                    color_t = 'toron_color_green';
                } else {
                    color_t = 'toron_color_normal';
                }
                
                data_a += '' + 
                    '<table id="toron">' + 
                        '<tr>' + 
                            '<td id="' + color_t + '">' + 
                                '<a href="javascript:void(0);" id="' + key + '">#' + key + '</a> ' + 
                                ip + 
                                '<span style="float: right;">' + data_t[key]['date'] + '</span>' + 
                            '</td>' + 
                        '</tr>' + 
                        '<tr>' + 
                            '<td id="' + color_b + '">' + 
                                '<div id="topic_scroll">' + data_i_pas + '</div>' + 
                            '</td>' + 
                        '</tr>' +
                    '</table>' + 
                    '<hr class="main_hr">' + 
                ''

                document.getElementById(where).innerHTML += data_a;
                eval(data_t[key]['data_pas'][1]);
            }
            
            if(type_do === 'top') {
                new_topic_load(topic_num, 'main', '', 'main_topic');
            } else if(type_do === 'main') {
                data_url_v = window.location.href.split('#');
                if(data_url_v.length !== 0) {
                    if(document.getElementById(data_url_v[1])) {
                        document.getElementById(data_url_v[1]).focus();
                    }
                }
                
                new_topic_load(topic_num, 're', key_v, where);
            } else if(type_do === 're') {
                setTimeout(function() {
                    if(start === 0) {
                        new_topic_load(topic_num, 're', some, where);
                    } else {
                        new_topic_load(topic_num, 're', key_v, where);
                    }
                }, 2000);
            }
        }
    }
}