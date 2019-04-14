function open_foot(name) {
    var o_data = document.getElementById('c' + name);
    
    name = name.replace(/\.([0-9]+)$/, '');
    var g_data = document.getElementById(name);

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

function folding(num) { 
    var fol = document.getElementById('folding_' + num); 
    if(fol.style.display === 'inline-block' || fol.style.display === 'block') { 
        fol.style.display = 'none';
    } else {
        fol.style.display = 'block'; 
    } 
}

function do_preview(name) {
    var o_data = document.getElementById('content');
    var p_data = document.getElementById('see_preview');

    var s_data = new FormData();
    s_data.append('data', o_data.value);

    var url = "/api/w/" + name;
    
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send(s_data);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            p_data.innerHTML = JSON.parse(this.responseText)['data'];
        }
    }
}