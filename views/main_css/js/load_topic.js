function topic_list_load(topic_num, s_num, where) {
    var url = "/api/thread/" + String(topic_num) + "?render=1&num=" + String(s_num);
    var n_data = "";

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var t_data = JSON.parse(this.responseText);
            var t_plus_data = '';
            
            for(key in t_data) {
                n_data += t_data[key]['data'];
                t_plus_data += t_data[key]['plus_data'];
            }

            document.getElementById(where).innerHTML = n_data;
            eval(t_plus_data);
        }
    }
}

function topic_plus_load(topic_num, num) {
    var test = setInterval(function() {
        var url = "/api/thread/" + String(topic_num) + "?num=" + num + "&render=1";
        var n_data = '';
        var n_num = 1;

        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200 && this.responseText !== '{}\n') {
                var t_data = JSON.parse(this.responseText);
                var t_plus_data = '';

                for(key in t_data) {
                    n_data += t_data[key]['data'];
                    n_num = key;

                    t_plus_data += t_data[key]['plus_data'];
                }

                document.getElementById("plus_topic").innerHTML += n_data;
                eval(t_plus_data);

                topic_plus_load(topic_num, String(Number(num) + 1));
                clearInterval(test);
            }
        }
    }, 5000);
}

function topic_main_load(topic_num, s_num) {
    if(s_num) {
        var url = "/api/thread/" + String(topic_num) + "?render=1&num=" + s_num;
    } else {
        var url = "/api/thread/" + String(topic_num) + "?render=1";
    }
    var n_data = "";
    var num = 1;

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var t_data = JSON.parse(this.responseText);
            var t_plus_data = '';

            for(var key in t_data) {
                n_data += t_data[key]['data'];
                num = key;

                t_plus_data += t_data[key]['plus_data'];
            }

            document.getElementById('main_topic').innerHTML = n_data;
            eval(t_plus_data);

            if(window.location.hash) {
                document.getElementById(window.location.hash.replace(/^#/, '')).focus();
            }
            
            if(!s_num) {
                topic_plus_load(topic_num, String(Number(num) + 1));
            }
        }
    }
}

function topic_top_load(topic_num) {
    var url = "/api/thread/" + String(topic_num) + "?top=1&render=1";
    var n_data = "";
    var num = 1;

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var t_data = JSON.parse(this.responseText);
            var t_plus_data = '';

            for(var key in t_data) {
                n_data += t_data[key]['data'];
                num = key;

                t_plus_data += t_data[key]['plus_data'];
            }

            document.getElementById('top_topic').innerHTML = n_data;
            eval(t_plus_data);

            topic_main_load(topic_num, null);
        }
    }
}

function new_topic_load(topic_num, con = 1, type_do = 0, some = '', where = 'top_topic') {
    if(type_do === 0) {
        var url = "/api/thread/" + topic_num + "?top=1";
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
                
                var color_b = '';
                var color_t = '';
                
                var ip = data_t[key]['ip_pas'];
                var ip_o = data_t[key]['ip'];
                var blind = data_t[key]['blind'];
                var data_i = data_t[key]['data'];
                
                if(data_i === '') {
                    data_i = '<br>';
                }
                
                if(blind === 'O') {
                    color_b = 'toron_color_not';
                } else {
                    color_b = 'toron_color';
                }
                
                if(blind === 'O') {
                    ip += ' <a href="/admin_log?search=blind%20(code%20' + topic_num + '#' + i[0] + '">(B)</a>';
                    
                    if(admin === '1') {
                        ip += ' <a href="/thread/' + topic_num + '/raw/' + key + '">(R)</a>';
                    }
                }
                
                if(admin === '1' || blind !== 'O') {
                    ip += ' <a href="/thread/' + topic_num + '/admin/' + key + '">(T)</a>';
                }
                
                if(type_do === 0) {
                    color_t = 'toron_color_red';
                } else if(blind === '1') {
                    color_t = 'toron_color_blue';
                } else if(ip_o === ip_first) {
                    color_t = 'toron_color_green';
                } else {
                    color_t = 'toron_color';
                }
                
                console.log(data_t[key])
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
                                '<div id="topic_scroll">' + data_t[key]['data_pas'][0] + '</div>' + 
                            '</td>' + 
                        '</tr>' +
                    '</table>' + 
                    '<hr class="main_hr">' + 
                ''

                document.getElementById(where).innerHTML += data_a;
                eval(data_t[key]['data_pas'][1]);
            }
            
            if(con === 1) {
                new_topic_load(topic_num, 0, type_do + 1, '', where);
            }
        }
    }
}