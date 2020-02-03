function topic_plus_load(topic_num, num) {
    var test = setInterval(function() {
        var url = "/api/thread/" + String(topic_num) + "?num=" + num + "&render=1";
        var p_data = document.getElementById("plus_topic");
        var n_data = '';
        var n_num = 1;
        var url_2 = "/api/markup";

        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.send(null);

        var xhr_2 = new XMLHttpRequest();
        xhr_2.open("GET", url_2, true);
        xhr_2.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200 && this.responseText !== '{}\n') {
                var t_data = JSON.parse(this.responseText);
                var t_plus_data = '';

                for(key in t_data) {
                    n_data += t_data[key]['data'];
                    n_num = key;

                    t_plus_data += t_data[key]['plus_data'].replace(/<script>/g, '').replace(/<\/script>/g, '');
                }

                p_data.innerHTML += n_data;
                eval(t_plus_data);

                // https://programmingsummaries.tistory.com/379
                var options = {
                    body: '#' + n_num
                }

                var notification = new Notification("openNAMU", options);

                setTimeout(function () {
                    notification.close();
                }, 5000);

                topic_plus_load(topic_num, String(Number(num) + 1));
                clearInterval(test);
            }
        }
    }, 2000);
}