function topic_main_load(topic_num, s_num) {
    var o_data = document.getElementById('main_topic');
    if(s_num) {
        var url = "/api/thread/" + String(topic_num) + "?render=1&num=" + s_num;
    } else {
        var url = "/api/thread/" + String(topic_num) + "?render=1";
    }
    var url_2 = "/api/markup";
    var n_data = "";
    var num = 1;

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    var xhr_2 = new XMLHttpRequest();
    xhr_2.open("GET", url_2, true);
    xhr_2.send(null);

    xhr.onreadystatechange = function() {
        if(xhr.readyState === 4 && xhr.status === 200) {
            t_data = JSON.parse(xhr.responseText);
            for(var key in t_data) {
                n_data += t_data[key]['data'];
                num = key;
            }

            o_data.innerHTML = n_data;
            if(!s_num) {
                topic_plus_load(topic_num, String(Number(num) + 1));
            }

            xhr_2.onreadystatechange = function() {
                if(xhr_2.readyState === 4 && xhr_2.status === 200) {
                    markup = JSON.parse(xhr_2.responseText)['markup'];

                    if(markup === 'markdown') {
                        render_markdown();
                    } else {
                        for(var key in t_data) {
                            render_html('topic_' + String(key) + '-');
                        }
                    }
                }
            }
        }
    }

}