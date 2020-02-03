function get_link_state(i = 0) {       
    if(document.getElementsByClassName('link_finder')[i]) {
        var link_data = document.getElementsByClassName('link_finder')[i];

        var xhr = new XMLHttpRequest();
        xhr.open("GET", link_data.href.replace('/w/', '/api/w/') + "?exist=1", true);
        xhr.send(null);
        
        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(JSON.parse(this.responseText)['exist'] !== '1') {
                    document.getElementsByClassName('link_finder')[i].id = "not_thing";
                } else {
                    document.getElementsByClassName('link_finder')[i].id = "";
                }

                get_link_state(i + 1);
            }
        }
    }
}

function get_file_state(i = 0) {       
    if(document.getElementsByClassName('file_finder_1')[i]) {
        var file_data = document.getElementsByClassName('file_finder_1')[i];

        var xhr = new XMLHttpRequest();
        xhr.open("GET", file_data.src.replace('/image/', '/api/image/'), true);
        xhr.send(null);
        
        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                if(JSON.parse(this.responseText)['exist'] !== '1') {
                    document.getElementsByClassName('file_finder_1')[i].style = "display: none;";
                } else {
                    document.getElementsByClassName('file_finder_2')[i].innerHTML = "";
                }
            
                get_file_state(i + 1);
            }
        }
    }
}