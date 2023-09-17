function opennamu_do_editor_preview() {
    const input = document.querySelector('#opennamu_edit_textarea') as HTMLInputElement | null;
    if(input !== null) {
        fetch("/api/w/test/doc_tool/preview", {
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