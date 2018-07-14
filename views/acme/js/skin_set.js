$(document).ready(function() {
    if(window.location.pathname == "/skin_set") {
        const title = document.getElementById("fix_title");
        const data = document.getElementById("fix_data");

        get_title = "Skin Setting";
        get_data =  ' \
                        <form> \
                            <input type="checkbox" name="chk_info" value="dark"> Dark Mode \
                            <hr> \
                            <button>Save</button> \
                        </form> \
                    ';

        document.title = document.title.replace(/.*(\- .*)$/, get_title + " $1");

        title.innerHTML = get_title;
        data.innerHTML = get_data;
    } else {
    }
});