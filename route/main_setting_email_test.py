from .tool.func import *

from .go_api_func_email import api_func_email

async def main_setting_email_test():
    with get_db_connect() as conn:
        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
        if flask.request.method == 'POST':
            render_data = await api_func_email()
            if render_data["response"] == "ok":
                data = await get_lang("ok")
            else:
                data = await get_lang("error")

            return await render_template(
                await get_lang("email_test"),
                data,
                0,
                [["setting/external", await get_lang('return')]]
            )
        else:
            return await render_template(
                await get_lang("email_test"),
                '''
                    <form method="post">
                        <input class="__ON_INPUT__" name="title" placeholder="''' + await get_lang("title") + '''">
                        <hr class="main_hr">
                        <input class="__ON_INPUT__" name="email" placeholder="''' + await get_lang("email") + '''">
                        <hr class="main_hr">
                        <textarea  name="data" class="opennamu_textarea_500 __ON_TEXTAREA__" placeholder="''' + await get_lang("content") + '''"></textarea>
                        <hr class="main_hr">
                        <button type="submit">''' + await get_lang("send") + '''</button>
                    </form>
                ''',
                0,
                [["setting/external", await get_lang('return')]]
            )
