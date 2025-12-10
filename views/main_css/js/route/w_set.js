"use strict";

function w_set_reset() {
    let lang_data = new FormData();
    lang_data.append('data', 'reset end authority_error');

    fetch('/api/v2/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];
        
        let check = confirm(lang['reset']);
        if(check === true) {
            const url = window.location.pathname;
            const url_split = url.split('/');
    
            let doc_name = url_split.slice(2, undefined).join('/');

            fetch('/api/v2/set_reset/' + doc_name).then(function(res) {
                return res.json();
            }).then(function(data) {
                if(data) {
                    if(data["response"] === "require auth") {
                        alert(lang["authority_error"]);
                    } else {
                        alert(lang['end']);
                    }

                    history.go(0);
                }
            });
        }
    });
}