function topic_main_load(name, sub, s_num = null) {
    var o_data = document.getElementById('main_topic');
    if(s_num) {
        var url = "/api/topic/" + name + "/sub/" + sub + "?render=1&num=" + s_num;
    } else {
        var url = "/api/topic/" + name + "/sub/" + sub + "?render=1";
    }
    var n_data = "";
    var num = 1;
    
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            t_data = JSON.parse(this.responseText);
            for(key in t_data) {
                n_data += t_data[key]['data'];
                num = key;
            }
            
            o_data.innerHTML = n_data;
            if(!s_num) {
                topic_plus_load(name, sub, String(Number(num) + 1));
            }
        }
    }
    
}