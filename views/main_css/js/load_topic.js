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