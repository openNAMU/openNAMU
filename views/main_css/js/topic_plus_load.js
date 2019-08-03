function topic_plus_load(name, sub, num) {
    var test = setInterval(function() {
        var url = "/api/topic/" + name + "/sub/" + sub + "?num=" + num;
        var doc_data = document.getElementById("plus_topic");

        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(this.responseText) {
                    doc_data.innerHTML += '<hr class="main_hr">(New)<hr class="main_hr">';

                    // https://programmingsummaries.tistory.com/379
                    var options = {
                        body: 'New topic'
                    }
                   
                    var notification = new Notification("openNAMU", options);
                    
                    setTimeout(function () {
                        notification.close();
                    }, 5000);

                    clearInterval(test);
                }
            }
        }
    }, 2000);
}