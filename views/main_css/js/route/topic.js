function opennamu_do_remove_blind_thread() {
    const style = document.querySelector('#opennamu_remove_blind');
    if(style !== null) {
        if(style.innerHTML !== "") {
            style.innerHTML = '';
        } else {
            style.innerHTML = `
                .opennamu_comment_blind_js {
                    display: none;
                }
            `;
        }
    }
}