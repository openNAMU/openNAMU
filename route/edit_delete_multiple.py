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
            return await render_template(
                await get_lang('many_delete'),
                '''
                    <form method="post">
                        <textarea class="opennamu_textarea_500" placeholder="''' + await get_lang('many_delete_help') + '''" name="content"></textarea>
                        <hr class="main_hr">
                        <input class="__ON_INPUT__" placeholder="''' + await get_lang('why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        ''' + await captcha_get(conn) + await ip_warning(conn) + get_edit_text_bottom_check_box(conn) + get_edit_text_bottom(conn, 'edit')  + '''
                        <button type="submit">''' + await get_lang('delete') + '''</button>
                    </form>
                ''',
                0,
                [['manager/1', await get_lang('return')]]
            )
