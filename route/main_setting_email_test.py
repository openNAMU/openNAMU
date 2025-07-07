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

            return easy_minify(flask.render_template(await skin_check(conn),
                imp = [await get_lang("email_test"), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = data,
                menu = [["setting/external", await get_lang('return')]]
            ))
        else:
            return easy_minify(flask.render_template(await skin_check(conn),
                imp = [await get_lang("email_test"), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input name="title" placeholder="''' + await get_lang("title") + '''">
                        <hr class="main_hr">
                        <input name="email" placeholder="''' + await get_lang("email") + '''">
                        <hr class="main_hr">
                        <textarea  name="data" class="opennamu_textarea_500" placeholder="''' + await get_lang("content") + '''"></textarea>
                        <hr class="main_hr">
                        <button type="submit">''' + await get_lang("send") + '''</button>
                    </form>
                ''',
                menu = [["setting/external", await get_lang('return')]]
            ))