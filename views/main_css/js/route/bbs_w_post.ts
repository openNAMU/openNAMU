function opennamu_change_comment(get_id : string): void {
    const input = document.querySelector('#opennamu_comment_select') as HTMLInputElement | null;
    if(input !== null) {
        input.value = get_id;
        document.getElementById('opennamu_edit_textarea')?.focus();
    }
}

function opennamu_return_comment(): void {
    const input = document.querySelector('#opennamu_comment_select') as HTMLInputElement | null;
    if(input !== null) {
        document.getElementById(input.value)?.focus();
    }
}