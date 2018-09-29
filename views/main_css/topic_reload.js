function topic_load(name, sub) {
    function addZero(i) {
        if(i < 10) {
            i = "0" + i;
        }
        
        return i;
    }

    setInterval(
        function() {
            var d = new Date();
            d.setSeconds(d.getSeconds() - 3);
            
            var date = d.getFullYear() + '-' + addZero(d.getMonth() + 1) + '-' + d.getDate() + ' ' + addZero(d.getHours()) + ':' + addZero(d.getMinutes()) + ':' + addZero(d.getSeconds());
            var url = "/api/topic/" + name + "/sub/" + sub + "?time=" + date;
            var xhr = new XMLHttpRequest();

            doc_data = document.getElementById("plus");

            test = '';

            xhr.onreadystatechange = function() {
                if(xhr.status == 200) {
                    if(xhr.responseText != "{}\n" && test != xhr.responseText) {
                        test = xhr.responseText;
                        doc_data.innerText += xhr.responseText + '\n';
                    }
                }
            }

            xhr.open("GET", url);
            xhr.send();
        }
    , 3000);
}