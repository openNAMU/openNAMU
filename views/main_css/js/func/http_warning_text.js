"use strict";

function opennamu_do_warning_text() {
    if(window.location.protocol !== 'https:') {
        let http_warning_text = document.getElementById('opennamu_http_warning_text_lang').innerHTML;
        
        document.getElementById('opennamu_http_warning_text').innerHTML = http_warning_text;
        document.getElementById('opennamu_http_warning_text').style.margin = "10px 0px 0px 0px";
    }
}