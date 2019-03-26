function open_foot(name) {
    var g_data = document.getElementById(name);
    var o_data = document.getElementById('c' + name);

    if(o_data.innerHTML === '') {
        o_data.innerHTML += '<sup><a onclick="open_foot(\'' + name + '\')" href="#' + name + '">(Go)</a></sup> ' + g_data.innerHTML;
    } else {
        o_data.innerHTML = '';
    }
}

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

            xhr.open("GET", url);
            xhr.send(null);

            xhr.onreadystatechange = function() {
                if(this.readyState == XMLHttpRequest.DONE && xhr.status == 200 && xhr.responseText != "{}\n") {
                    console.log(xhr.responseText);
                    console.log(url);

                    doc_data.innerText += '(New)\n\n';

                    clearInterval(test);
                }
            }
        }, 3000)
    }, 4000);
}

function folding(num) { 
    var fol = document.getElementById('folding_' + num); 
    if(fol.style.display == 'inline-block' || fol.style.display == 'block') { 
        fol.style.display = 'none';
    } else {
        fol.style.display = 'block'; 
    } 
}