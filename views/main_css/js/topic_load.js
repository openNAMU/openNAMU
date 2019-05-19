function topic_load(name, sub, num) {
    var test = setInterval(function() {
        var url = "/api/topic/" + name + "/sub/" + sub + "?num=" + num;
        var doc_data = document.getElementById("plus");

        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(this.responseText) {
                    doc_data.innerHTML += '<hr class="main_hr">(New)<hr class="main_hr">';

                    clearInterval(test);
                }
            }
        }
    }, 2000);
}