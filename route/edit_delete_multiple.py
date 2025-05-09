from .tool.func import *

from .edit_delete import edit_delete

async def edit_delete_multiple():
    with get_db_connect() as conn:
        if await acl_check('', 'acl_auth', '', '') == 1:
            return await re_error(conn, 0)

        if flask.request.method == 'POST':
            send = flask.request.form.get('send', '')
            agree = flask.request.form.get('copyright_agreement', '')
            
            if await do_edit_send_check(conn, send) == 1:
                return await re_error(conn, 37)
            
            if do_edit_text_bottom_check_box_check(conn, agree) == 1:
                return await re_error(conn, 29)
            
            all_title = re.findall(r'([^\n]+)\n', flask.request.form.get('content', '').replace('\r', '') + '\n')
            for name in all_title:
                await edit_delete(name)

            return redirect(conn, '/recent_change')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'many_delete'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <textarea class="opennamu_textarea_500" placeholder="''' + get_lang(conn, 'many_delete_help') + '''" name="content"></textarea>
                        <hr class="main_hr">
                        <input placeholder="''' + get_lang(conn, 'why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        ''' + await captcha_get(conn) + ip_warning(conn) + get_edit_text_bottom_check_box(conn) + get_edit_text_bottom(conn, 'edit')  + '''
                        <button type="submit">''' + get_lang(conn, 'delete') + '''</button>
                    </form>
                ''',
                menu = [['manager/1', get_lang(conn, 'return')]]
            ))