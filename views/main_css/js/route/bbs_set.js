function opennamu_bbs_set_post() {
    let post_data = new FormData();
    post_data.append('data', '');

    fetch('/api/v2/bbs/set', {
        method : 'put',
        body : post_data,
    }).then(function(res) {
        return res.json();
    }).then(function(data) {

    });
}

function opennamu_bbs_set() {
    let lang_data = new FormData();
    lang_data.append('data', 'title_start_document title_end_document title_include_document move document_name');

    fetch('/api/lang', {
        method : 'post',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];
    
        document.getElementById('opennamu_bbs_set').innerHTML = '' +
            '<button onclick="opennamu_bbs_set_post();"></button>' +
        '';
    });
}