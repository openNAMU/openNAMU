"use strict";

function opennamu_do_wiki_access() {
    let password = document.getElementById('wiki_access').value;
    
    document.cookie = 'opennamu_wiki_access=' + encodeURIComponent(password) + '; path=/;';
    
    history.go(0);
}