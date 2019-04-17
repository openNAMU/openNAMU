function topic_load(name, sub) {
    function addZero(i) {
        if(i < 10) {
            i = "0" + i;
        }
        
        return i;
    }

    setTimeout(function() {
        var test = setInterval(function() {
            var d = new Date();
            d.setSeconds(d.getSeconds() - 3);
            
            var date = d.getFullYear() + '-' + addZero(d.getMonth() + 1) + '-' + addZero(d.getDate());
            date += ' ' + addZero(d.getHours()) + ':' + addZero(d.getMinutes()) + ':' + addZero(d.getSeconds());

            var url = "/api/topic/" + name + "/sub/" + sub + "?time=" + date;
            var xhr = new XMLHttpRequest();
            var doc_data = document.getElementById("plus");

            xhr.open("GET", url, true);
            xhr.send(null);

            xhr.onreadystatechange = function() {
                if(this.readyState === 4 && this.status === 200 && this.responseText !== "{}\n") {
                    doc_data.innerText += '(New)\n\n';

                    clearInterval(test);
                }
            }
        }, 1000)
    }, 4000);
}