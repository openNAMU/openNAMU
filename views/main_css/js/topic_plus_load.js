function topic_plus_load(name, sub, num) {
    var test = setInterval(function() {
        var url = "/api/topic/" + encodeURI(name) + "/sub/" + encodeURI(sub) + "?num=" + num + "&render=1";
        var p_data = document.getElementById("plus_topic");
        var n_data = '';
        var n_num = 1;

        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200 && this.responseText !== '{}\n') {                
                t_data = JSON.parse(this.responseText);
                for(key in t_data) {
                    n_data += t_data[key]['data'];
                    n_num = key;
                }
                
                p_data.innerHTML += n_data;

                // https://programmingsummaries.tistory.com/379
                var options = {
                    body: '#' + n_num
                }
               
                var notification = new Notification("openNAMU", options);
                
                setTimeout(function () {
                    notification.close();
                }, 5000);

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

                topic_plus_load(name, sub, String(Number(num) + 1));
                clearInterval(test);
            }
        }
    }, 2000);
}