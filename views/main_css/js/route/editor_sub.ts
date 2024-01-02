declare function opennamu_do_url_encode(data : any) : string;
declare function do_sync_monaco_and_textarea();
declare function do_textarea_to_manaco();

interface Window {
    editor? : any;
}

function opennamu_do_editor_preview() {
    do_sync_monaco_and_textarea();

    const input = document.querySelector('#opennamu_edit_textarea') as HTMLInputElement | null;
    if(input !== null) {
        let doc_name : string = 'test';

        const doc_name_input = document.querySelector('#opennamu_editor_doc_name') as HTMLInputElement | null;
        if(doc_name_input !== null) {
            doc_name = doc_name_input.value;
        }

        fetch("/api/w_tool/preview/" + (opennamu_do_url_encode(doc_name)), {
            method : 'POST',
            headers : { 'Content-Type': 'application/x-www-form-urlencoded' },
            body : new URLSearchParams({
                'data': input.value,
            })
        }).then(function(res) {
            return res.json();
        }).then(function(text) {
            const preview = document.querySelector('#opennamu_preview_area') as HTMLInputElement | null;
            if(preview !== null) {
                preview.innerHTML = text.data;
                eval(text.js_data);
            }
        });
    }
}

function opennamu_do_editor_temp_save() {
    do_sync_monaco_and_textarea();

    const input = document.querySelector('#opennamu_edit_textarea') as HTMLInputElement | null;
    if(input !== null) {
        localStorage.setItem("key", input.value);
    }
}

function opennamu_do_editor_temp_save_load() {
    const data = localStorage.getItem("key");
    if(data !== null) {
        const input = document.querySelector('#opennamu_edit_textarea') as HTMLInputElement | null;
        if(input !== null) {
            input.value = data;
        }
        
        do_textarea_to_manaco();
    }
}